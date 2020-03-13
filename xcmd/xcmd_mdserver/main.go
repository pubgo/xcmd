// Command mdserver is an http server serving directory (recursively) with
// markdown (.md) files. If can render automatically built index of those files
// and render them as html pages.
//
// Its main use-case is reading through directory with documentation written in
// markdown format, i.e. local copy of Github wiki.
//
// To access automatically generated index, request "/?index" path, as
// http://localhost:8080/?index.
//
// To create home page available at / either create index.html file or start
// server with -rootindex flag to render automatically generated index.
//
// If started with -github flag, it will render any absolute links to github
// wikis like "https://github.com/user/project/wiki/Page" to relative ones like
// "Page.md".
//
// To apply custom styling provide css file with -css flag. By default, this
// file is read on server start and then embedded into code of every page,
// making them self-sufficient. If you instead wish to link stylesheet, provide
// absolute root-related path to css file located under the same path as your
// markdown files (-dir flag) and enable -csslink flag. This will link
// stylesheet into head section of page with href being value of -css flag.
//
//
// Note that table of contents generating javascript is a modified version of
// code found at https://github.com/matthewkastor/html-table-of-contents which
// is licensed under GNU GENERAL PUBLIC LICENSE Version 3.
package xcmd_mdserver

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/artyom/autoflags"
	"github.com/artyom/httpgzip"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/browser"
	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

func main() {
	args := runArgs{Dir: ".", Addr: "localhost:8080"}
	autoflags.Parse(&args)
	if err := run(args); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

type runArgs struct {
	Dir     string `flag:"dir,directory with markdown (.md) files"`
	Addr    string `flag:"addr,address to listen"`
	Open    bool   `flag:"open,open index page in default browser on start"`
	Ghub    bool   `flag:"github,rewrite github wiki links to local when rendering"`
	Grep    bool   `flag:"search,enable substring search"`
	Idx     bool   `flag:"rootindex,render autogenerated index at / in addition to /?index"`
	CSS     string `flag:"css,path to custom CSS file (embedded into page unless run with -csslink)"`
	LinkCSS bool   `flag:"csslink,treat -css argument as local href inside <link rel=stylesheet>"`
	HLJS    bool   `flag:"hljs,syntax-highlight code blocks with defined language using highlight.js"`
}

func run(args runArgs) error {
	h := &mdHandler{
		dir:        args.Dir,
		fileServer: http.FileServer(http.Dir(args.Dir)),
		githubWiki: args.Ghub,
		withSearch: args.Grep,
		rootIndex:  args.Idx,
		hljs:       args.HLJS,
		linkStyle:  args.LinkCSS,
		style:      style,
	}
	if args.CSS != "" {
		switch {
		case args.LinkCSS:
			if !path.IsAbs(args.CSS) {
				return fmt.Errorf("with -csslink set, -css must be an absolute / separated path, but %q is not", args.CSS)
			}
			h.style = args.CSS
			reportIfMissing(filepath.Join(args.Dir, filepath.FromSlash(args.CSS)))
		default:
			b, err := ioutil.ReadFile(args.CSS)
			if err != nil {
				return err
			}
			h.style = string(b)
		}
	}
	if !args.LinkCSS {
		sum := sha256.Sum256([]byte(h.style))
		h.styleHash = "sha256-" + base64.StdEncoding.EncodeToString(sum[:])
	}
	srv := http.Server{
		Addr:        args.Addr,
		Handler:     httpgzip.New(h),
		ReadTimeout: time.Second,
	}
	if args.Open {
		go func() {
			time.Sleep(100 * time.Millisecond)
			browser.OpenURL("http://" + args.Addr + "/?index")
		}()
	}
	return srv.ListenAndServe()
}

type mdHandler struct {
	dir        string
	fileServer http.Handler // initialized as http.FileServer(http.Dir(dir))
	githubWiki bool
	withSearch bool
	rootIndex  bool
	hljs       bool
	linkStyle  bool
	style      string
	styleHash  string // sha256-{HASH} value for CSP
}

func (h *mdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	if h.withSearch && r.URL.Path == "/" && strings.HasPrefix(r.URL.RawQuery, "q=") {
		q := r.URL.Query().Get("q")
		if len(q) < 3 {
			http.Error(w, "Search term is too short", http.StatusBadRequest)
			return
		}
		pat := search.New(language.English, search.Loose).CompileString(q)
		h.renderIndex(w, fmt.Sprintf("Search results for %q", q), dirIndex(h.dir, pat))
		return
	}
	if r.URL.Path == "/" && (h.rootIndex || r.URL.RawQuery == "index") {
		h.renderIndex(w, "Index", dirIndex(h.dir, nil))
		return
	}
	if !strings.HasSuffix(r.URL.Path, mdSuffix) {
		h.fileServer.ServeHTTP(w, r)
		return
	}
	// only markdown files are handled below
	p := path.Clean(r.URL.Path)
	if containsDotDot(p) {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	name := filepath.Join(h.dir, filepath.FromSlash(p))
	rc, mtime, err := h.readerForFile(name)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		log.Printf("read %q: %v", name, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Security-Policy", h.csp(h.hljs))
	http.ServeContent(w, r, "page.html", mtime, rc)
}

func (h *mdHandler) renderIndex(w io.Writer, title string, index []indexRecord) error {
	page := struct {
		Title      string
		StyleHref  string
		Style      template.CSS
		Index      []indexRecord
		WithSearch bool
	}{
		Title:      title,
		Index:      index,
		WithSearch: h.withSearch,
	}
	switch {
	case h.linkStyle:
		page.StyleHref = h.style
	default:
		page.Style = template.CSS(h.style)
	}
	return indexTemplate.Execute(w, page)
}

func (h *mdHandler) csp(withHL bool) string {
	csp := []string{"default-src 'self';img-src http: https: data:;media-src https:"}
	switch {
	case withHL:
		csp = append(csp, "script-src https://cdnjs.cloudflare.com "+
			"'sha256-HGKuhVF4dzwg9Kt9XWXRYCoBYgGWsgnBiY1ynyCokzQ=' "+
			"'sha256-qeFup2+SGOg8HaUXLE/qospaz+lv/lxjtZZVNa2AqTk='", // https://play.golang.org/p/0SUWatm_LGr
		)
		switch {
		case h.linkStyle:
			csp = append(csp, "style-src 'self' https://cdnjs.cloudflare.com")
		default:
			csp = append(csp, "style-src https://cdnjs.cloudflare.com '"+h.styleHash+"'")
		}
	default:
		csp = append(csp, "script-src 'sha256-HGKuhVF4dzwg9Kt9XWXRYCoBYgGWsgnBiY1ynyCokzQ='")
		switch {
		case h.linkStyle:
			csp = append(csp, "style-src 'self'")
		default:
			csp = append(csp, "style-src '"+h.styleHash+"'")
		}
	}
	return strings.Join(csp, ";")
}

// readerForFile returns lazy io.ReadSeeker and mtime to be used as arguments of
// http.ServeContent. It does not use ReadSeeker at all if http client already
// has fresh content as signaled by "If-Modified-Since" request header;
// lazyReadSeeker takes advantage of this by defering any file reading and
// rendering until one of its method is called.
func (h *mdHandler) readerForFile(name string) (*lazyReadSeeker, time.Time, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, time.Time{}, err
	}
	return &lazyReadSeeker{name: name, h: h}, fi.ModTime(), nil
}

type lazyReadSeeker struct {
	name string
	h    *mdHandler
	r    *bytes.Reader // initially nil, initialized with init()
}

func (l *lazyReadSeeker) init() error {
	if l.r != nil {
		return nil
	}
	if testRun {
		log.Print("lazyReadSeeker init()")
	}
	b, err := ioutil.ReadFile(l.name)
	if err != nil {
		return err
	}
	opts := rendererOpts
	if l.h.githubWiki {
		opts.RenderNodeHook = rewriteGithubWikiLinks
	}
	doc := parser.NewWithExtensions(extensions).Parse(b)
	body := markdown.Render(doc, html.NewRenderer(opts))
	body = policy.SanitizeBytes(body)
	title := firstHeaderText(doc)
	if title == "" {
		title = nameToTitle(filepath.Base(l.name))
	}
	withHL := l.h.hljs && bytes.Contains(body, []byte(`<pre><code class=`))
	page := struct {
		Title     string
		StyleHref string
		Style     template.CSS
		Body      template.HTML
		WithHL    bool
	}{
		Title:  title,
		Body:   template.HTML(body),
		WithHL: withHL,
	}
	switch {
	case l.h.linkStyle:
		page.StyleHref = l.h.style
	default:
		page.Style = template.CSS(l.h.style)
	}
	buf := bytes.NewBuffer(b[:0]) // reuse b to reduce allocations
	if err := pageTemplate.Execute(buf, page); err != nil {
		return err
	}
	l.r = bytes.NewReader(buf.Bytes())
	return nil
}

func (l *lazyReadSeeker) Read(p []byte) (n int, err error) {
	if l.r == nil {
		if err := l.init(); err != nil {
			return 0, err
		}
	}
	return l.r.Read(p)
}

func (l *lazyReadSeeker) Seek(offset int64, whence int) (int64, error) {
	if l.r == nil {
		if err := l.init(); err != nil {
			return 0, err
		}
	}
	return l.r.Seek(offset, whence)
}

func dirIndex(dir string, pat *search.Pattern) []indexRecord {
	var matches []string
	fn := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && p != "." && strings.HasPrefix(filepath.Base(p), ".") {
			return filepath.SkipDir
		}
		if info.IsDir() || !strings.HasSuffix(p, mdSuffix) {
			return nil
		}
		matches = append(matches, p)
		return nil
	}
	if err := filepath.Walk(dir, fn); err != nil {
		log.Printf("walk %q: %v", dir, err)
	}
	var index []indexRecord
	if pat == nil {
		index = make([]indexRecord, 0, len(matches))
	}
	for _, s := range matches {
		if pat != nil && !matchPattern(pat, s) {
			continue
		}
		title := documentTitle(s)
		if title == "" {
			title = nameToTitle(filepath.Base(s))
		}
		file, err := filepath.Rel(dir, s)
		if err != nil {
			continue
		}
		index = append(index, indexRecord{
			Title:  title,
			File:   filepath.ToSlash(file),
			Subdir: filepath.ToSlash(filepath.Dir(file)),
			// precalculate sort key to speed up comparisons on sort
			sortKey: strings.ToLower(strings.TrimSuffix(filepath.Base(file), mdSuffix)),
		})
	}
	sort.Slice(index, func(i, j int) bool {
		si, sj := index[i].Subdir, index[j].Subdir
		if si == sj {
			return index[i].sortKey < index[j].sortKey
		}
		return si < sj
	})
	return index
}

type indexRecord struct {
	Title, File string
	Subdir      string // groups index records when rendering template
	sortKey     string // if File is "dir/FileName.md", then sortKey is "filename"
}

// documentTitle extracts h1 header from markdown document
func documentTitle(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	b, err := ioutil.ReadAll(io.LimitReader(f, 1<<17))
	if err != nil {
		return ""
	}
	return firstHeaderText(parser.New().Parse(b))
}

func firstHeaderText(doc ast.Node) string {
	var title string
	walkFn := func(node ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.GoToNext
		}
		switch n := node.(type) {
		case *ast.Heading:
			if n.Level != 1 {
				return ast.GoToNext
			}
			title = string(childLiterals(n))
			return ast.Terminate
		case *ast.Code, *ast.CodeBlock, *ast.BlockQuote:
			return ast.SkipChildren
		}
		return ast.GoToNext
	}
	_ = ast.Walk(doc, ast.NodeVisitorFunc(walkFn))
	return title
}

func childLiterals(node ast.Node) []byte {
	if l := node.AsLeaf(); l != nil {
		return l.Literal
	}
	var out [][]byte
	for _, n := range node.GetChildren() {
		if lit := childLiterals(n); lit != nil {
			out = append(out, lit)
		}
	}
	if out == nil {
		return nil
	}
	return bytes.Join(out, nil)
}

// matchPattern reports whether any line in file matches given pattern. On any
// errors function return false.
func matchPattern(pat *search.Pattern, file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()
	sc := bufio.NewScanner(io.LimitReader(f, 1<<20))
	for sc.Scan() {
		if _, end := pat.Index(sc.Bytes()); end > 0 {
			return true
		}
	}
	return false
}

// rewriteGithubWikiLinks is a html.RenderNodeFunc which renders links
// with github wiki destinations as local ones.
//
// Link with "https://github.com/user/project/wiki/Page" destination would be
// rendered as a link to "Page.md"
func rewriteGithubWikiLinks(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	link, ok := node.(*ast.Link)
	if !ok || !entering {
		return ast.GoToNext, false
	}
	if u, err := url.Parse(string(link.Destination)); err == nil &&
		u.Host == "github.com" && strings.HasSuffix(path.Dir(u.Path), "/wiki") {
		dst := path.Base(u.Path) + mdSuffix
		switch u.Fragment {
		case "":
			fmt.Fprintf(w, "<a href=\"%s\">", url.QueryEscape(dst))
		default:
			fmt.Fprintf(w, "<a href=\"%s#%s\">", url.QueryEscape(dst), url.QueryEscape(u.Fragment))
		}
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

// reportIfMissing tests whether file exists and logs if not
func reportIfMissing(name string) {
	if st, err := os.Stat(name); os.IsNotExist(err) || (st != nil && !st.Mode().IsRegular()) {
		log.Printf("called with -csslink, but path %q does not exist or not a regular file", name)
	}
}

func nameToTitle(name string) string {
	if strings.ContainsAny(name, " ") {
		return strings.TrimSuffix(name, mdSuffix)
	}
	return repl.Replace(strings.TrimSuffix(name, mdSuffix))
}

var repl = strings.NewReplacer("-", " ")

const mdSuffix = ".md"

var indexTemplate = template.Must(template.New("index").Parse(indexTpl))
var pageTemplate = template.Must(template.New("page").Parse(pageTpl))

const indexTpl = `<!doctype html><head><meta charset="utf-8"><title>{{.Title}}</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
{{if .StyleHref}}<link rel="stylesheet" href="{{.StyleHref}}">{{end -}}
{{if .Style}}<style>{{.Style}}</style>{{end}}</head><body id="mdserver-autoindex">{{if .WithSearch}}<form method="get">
<input type="search" name="q" minlength="3" placeholder="Substring search" autofocus required>
<input type="submit"></form>{{end}}
<h1>{{.Title}}</h1><ul>{{$prev := "."}}
{{range .Index}}{{if ne .Subdir $prev}}{{$prev = .Subdir}}</ul><h2>{{.Subdir}}</h2><ul>{{end}}<li><a href="{{.File}}">{{.Title}}</a></li>
{{end}}</ul></body>
`

const pageTpl = `<!doctype html><head><meta charset="utf-8"><title>{{.Title}}</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
{{if .StyleHref}}<link rel="stylesheet" href="{{.StyleHref}}">{{end -}}
{{if .Style}}<style>{{.Style}}</style>{{end}}<script>
document.addEventListener('DOMContentLoaded', function() {
	htmlTableOfContents();
} );
function htmlTableOfContents( documentRef ) {
	var documentRef = documentRef || document;
	var headings = [].slice.call(documentRef.body.querySelectorAll('article h1, article h2, article h3, article h4, article h5, article h6'));
	if (headings.length < 2) { return };
	var toc = documentRef.querySelector("nav#toc details");
	var ul = documentRef.createElement( "ul" );
	headings.forEach(function (heading, index) {
		var ref = heading.getAttribute( "id" );
		var link = documentRef.createElement( "a" );
		link.setAttribute( "href", "#"+ ref );
		link.textContent = heading.textContent;
		var li = documentRef.createElement( "li" );
		li.setAttribute( "class", heading.tagName.toLowerCase() );
		li.appendChild( link );
		ul.appendChild( li );
	});
	toc.appendChild( ul );
}
</script>{{if .WithHL}}
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.6/styles/default.min.css" integrity="sha256-zcunqSn1llgADaIPFyzrQ8USIjX2VpuxHzUwYisOwo8=" crossorigin="anonymous" referrerpolicy="no-referrer">
<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.6/highlight.min.js" integrity="sha256-aYTdUrn6Ow1DDgh5JTc3aDGnnju48y/1c8s1dgkYPQ8=" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<script>
document.addEventListener('DOMContentLoaded', (event) => {
	document.querySelectorAll('pre code[class^="language-"]').forEach((block) => {
		hljs.highlightBlock(block);
	});
});
</script>{{end}}
</head><body><nav id="site"><a href="/?index">index</a></nav>
<nav id="toc"><details open><summary>Contents</summary></details></nav>
<ul id="toc"></ul>
<article>
{{.Body}}
</article></body>
`

const extensions = parser.CommonExtensions | parser.AutoHeadingIDs ^ parser.MathJax

var rendererOpts = html.RendererOptions{Flags: html.CommonFlags}
var policy = bluemonday.UGCPolicy().AllowAttrs("class").OnElements("code")

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, func(r rune) bool { return r == '/' || r == '\\' }) {
		if ent == ".." {
			return true
		}
	}
	return false
}

const style = `body {
	font-family: Charter, Constantia, serif;
	font-size: 1rem;
	line-height: 170%;
	max-width: 45em;
	margin: auto;
	padding-right: 1em;
	padding-left: 1em;
	color: #333;
	background: white;
	text-rendering: optimizeLegibility;
}

@media only screen and (max-width: 480px) {
	body {
		font-size: 125%;
		text-rendering: auto;
	}
}

a {color: #a08941; text-decoration: none;}
a:hover {color: #c6b754; text-decoration: underline;}

h1 a, h2 a, h3 a, h4 a, h5 a {
	text-decoration: none;
	color: gray;
	break-after: avoid;
}
h1 a:hover, h2 a:hover, h3 a:hover, h4 a:hover, h5 a:hover {
	text-decoration: none;
	color: gray;
}
h1, h2, h3, h4, h5 {
	font-weight: bold;
	color: gray;
}

h1 {
	font-size: 150%;
}

h2 {
	font-size: 130%;
}

h3 {
	font-size: 110%;
}

h4, h5 {
	font-size: 100%;
	font-style: italic;
}

pre {
	background-color: rgb(240,240,240);
	color: #111111;
	padding: 0.5em;
	overflow: auto;
}
code, pre {
	font-family: Consolas, "PT Mono", monospace;
}
pre { font-size: 90%; }

hr { border:none; text-align:center; color:gray; }
hr:after {
	content:"\2766";
	display:inline-block;
	font-size:1.5em;
}

dt code {
	font-weight: bold;
}
dd p {
	margin-top: 0;
}

blockquote {
	border-left:thick solid lightgrey;
	color: #111111;
	padding: 0 0.5em;
}

img {display:block;margin:auto;max-width:100%}

table, td, th {
	border:thin solid lightgrey;
	border-collapse:collapse;
	vertical-align:middle;
}
td, th {padding:0.2em 0.5em}
tr:nth-child(even) {background-color: rgba(200,200,200,0.2)}

nav#toc {margin:1em 0 1em 0}
nav#toc summary {font-weight:bold; color:gray}
nav#toc ul:after {
	content:"\2042";
	text-align:center;
	display:block;
	color:gray;
}
nav#toc ul {margin:0; list-style:none; padding-left:0}
nav#toc ul li.h2 {padding-left:1em}
nav#toc ul li.h3 {padding-left:2em}
nav#toc ul li.h4 {padding-left:3em}
nav#toc ul li.h5 {padding-left:4em}
nav#toc ul li.h6 {padding-left:5em}

nav#site {
	font-size:90%;
	text-align:right;
	padding:.5em;
	border-bottom: 1px solid gray;
}
nav#site a:before {content:"\2767\0020"}

footer summary {font-weight:bold; color:gray}

summary {cursor:pointer; outline:none}
summary:only-child {display:none}

@media print {
	nav {display: none}
	pre {overflow-wrap:break-word; white-space:pre-wrap}
}`

var testRun bool // used in tests

//go:generate sh -c "go doc >README"