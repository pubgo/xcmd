package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/pubgo/g/pkg/gitutil"
	"github.com/pubgo/g/xenv"
	"github.com/pubgo/g/xerror"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	xerror.Panic(xenv.LoadFile(".env"))
	systray.RunWithAppWindow("代码监控", 1024, 768, onReady, func() {
		fmt.Println("exit ")
	})
}

func onReady() {
	systray.SetIcon(Data)
	systray.SetTitle("代码监控")

	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()

	monitor := systray.AddMenuItem("代码监控", "Change Me")

	go func() {
		for {
			rTo := xerror.PanicErr(git.PlainOpen("/Users/barry/gopath/src/mydocs")).(*git.Repository)
			wTo := xerror.PanicErr(rTo.Worktree()).(*git.Worktree)
			_status := xerror.PanicErr(wTo.Status()).(git.Status)
			if len(_status) == 0 {
				time.Sleep(time.Second)
				continue
			}

			monitor.SetTitle("代码已经改变")
			<-monitor.ClickedCh
			fmt.Println("Commit")
			for k := range _status {
				xerror.PanicErr(wTo.Add(k))
			}

			// 提交commit
			xerror.PanicErr(wTo.Commit("ok", &git.CommitOptions{
				All: true,
				Author: &object.Signature{
					Name:  "kooksee",
					Email: "kooksee@163.com",
					When:  time.Now(),
				},
			}))

			if err := rTo.Push(&git.PushOptions{
				Auth:     xerror.PanicErr(gitutil.GitBasicAuth(xenv.GetEnv("username"), xenv.GetEnv("password"))).(transport.AuthMethod),
				Progress: os.Stdout,
				RefSpecs: []config.RefSpec{"+" + config.DefaultPushRefSpec},
			}); err != nil && err != git.NoErrAlreadyUpToDate && !strings.Contains(err.Error(), "non-fast-forward update") {
				xerror.PanicM(err, "git 仓库 %s push failed", err)
			}

			monitor.SetTitle("代码监控")
			time.Sleep(time.Second)
		}

	}()

	mChange := systray.AddMenuItem("Change Me", "Change Me")
	mChecked := systray.AddMenuItem("Unchecked", "Check Me")
	mEnabled := systray.AddMenuItem("Enabled", "Enabled")

	mEnabled.SetTooltip("测试tooltip")

	systray.AddMenuItem("Ignored", "Ignored")
	mUrl := systray.AddMenuItem("Open UI", "my home")
	mQuit := systray.AddMenuItem("退出", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac.
	mQuit.SetIcon(Data)

	//mQuit.Hide()

	systray.AddSeparator()
	mToggle := systray.AddMenuItem("Toggle", "Toggle the Quit button")
	shown := true

	// We can manipulate the systray in other goroutines
	go func() {
		for {
			select {
			//case <-monitor.ClickedCh:
			//	fmt.Println(monitor.Checked())
			case <-mChange.ClickedCh:
				mChange.SetTitle("I've Changed")
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					mChecked.Uncheck()
					mChecked.SetTitle("Unchecked")
				} else {
					mChecked.Check()
					mChecked.SetTitle("Checked")
				}
			case <-mEnabled.ClickedCh:
				mEnabled.SetTitle("Disabled")
				mEnabled.Disable()
			case <-mUrl.ClickedCh:
				var indexHTML = `
<!DOCTYPE html>
<html>
<head>
<title>测试 - 幕布</title>
<meta charset="utf-8"/>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="renderer" content="webkit"/>
<meta name="author" content="mubu.com"/>
</head>
<body style="margin: 50px 20px;color: #333;font-family: SourceSansPro,-apple-system,BlinkMacSystemFont,'PingFang SC',Helvetica,Arial,'Microsoft YaHei',微软雅黑,黑体,Heiti,sans-serif,SimSun,宋体,serif">
<div class="export-wrapper"><div style="font-size: 22px; padding: 0 15px 0;"><div style="padding-bottom: 24px">测试</div><div style="background: #e5e6e8; height: 1px; margin-bottom: 20px;"></div></div><ul style="list-style: disc outside;"><li style="line-height: 34px;"><span class="content mubu-node" heading="1" style="line-height: 34px; min-height: 34px; font-size: 24px; padding: 2px 0px; display: inline-block; vertical-align: top;"># sssnj</span></li><li style="line-height: 34px;"><span class="content mubu-node" heading="1" style="line-height: 34px; min-height: 34px; font-size: 24px; padding: 2px 0px; display: inline-block; vertical-align: top;"><span class="bold" style="font-weight: bold;">heade1</span></span><ul class="children" style="list-style: disc outside; padding-bottom: 4px;"><li style="line-height: 30px;"><span class="content mubu-node" heading="2" style="line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;">heade2</span><ul class="children" style="list-style: disc outside; padding-bottom: 4px;"><li style="line-height: 27px;"><span class="content mubu-node" heading="3" style="line-height: 27px; min-height: 27px; font-size: 19px; padding: 2px 0px; display: inline-block; vertical-align: top;">heade3</span><ul class="children" style="list-style: disc outside; padding-bottom: 4px;"><li style="line-height: 24px;"><span class="content mubu-node" color="#dc2d1e" style="color: rgb(220, 45, 30); line-height: 24px; min-height: 24px; font-size: 16px; padding: 2px 0px; display: inline-block; vertical-align: top;">三生三世</span></li></ul></li></ul></li></ul></li><li style="line-height: 30px;"><span class="content mubu-node" heading="2" style="line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;"><span class="bold" style="font-weight: bold;">heade2</span></span><ul class="children" style="list-style: disc outside; padding-bottom: 4px;"><li style="line-height: 27px;"><span class="content mubu-node" heading="3" style="line-height: 27px; min-height: 27px; font-size: 19px; padding: 2px 0px; display: inline-block; vertical-align: top;">heade3</span><ul class="children" style="list-style: disc outside; padding-bottom: 4px;"><li style="line-height: 24px;"><span class="content mubu-node" style="line-height: 24px; min-height: 24px; font-size: 16px; padding: 2px 0px; display: inline-block; vertical-align: top;"><span class="bold" style="font-weight: bold;">三生三世</span></span></li></ul></li></ul></li><li style="line-height: 27px;"><span class="content mubu-node" heading="3" style="line-height: 27px; min-height: 27px; font-size: 19px; padding: 2px 0px; display: inline-block; vertical-align: top;"><span class="bold" style="font-weight: bold;">heade3</span></span><ul class="children" style="list-style: disc outside; padding-bottom: 4px;"><li style="line-height: 24px;"><span class="content mubu-node" style="line-height: 24px; min-height: 24px; font-size: 16px; padding: 2px 0px; display: inline-block; vertical-align: top;"><span class="bold underline" style="font-weight: bold; text-decoration: underline;">三生三世</span></span></li></ul></li><li style="line-height: 24px;"><span class="content mubu-node" color="#dc2d1e" style="color: rgb(220, 45, 30); line-height: 24px; min-height: 24px; font-size: 16px; padding: 2px 0px; display: inline-block; vertical-align: top;">三生三世</span><ul class="children" style="list-style: disc outside; padding-bottom: 4px;"><li style="line-height: 24px;"><span class="content mubu-node" color="#dc2d1e" style="color: rgb(220, 45, 30); line-height: 24px; min-height: 24px; font-size: 16px; padding: 2px 0px; display: inline-block; vertical-align: top;">hello</span></li></ul></li><li style="line-height: 30px;"><span class="content mubu-node" color="#333333" heading="2" style="color: rgb(51, 51, 51); line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;">测试</span><ul class="children" style="list-style: disc outside; padding-bottom: 4px;"><li style="line-height: 30px;"><span class="content mubu-node" color="#3da8f5" heading="2" images="%5B%7B%22id%22%3A%221d916f3267c0b118a-40263%22%2C%22oh%22%3A1004%2C%22ow%22%3A742%2C%22uri%22%3A%22document_image%2F7fabd28a-8c59-4ffe-b9f3-ab2ef4c91549-40263.jpg%22%2C%22w%22%3A87%7D%5D" style="color: rgb(61, 168, 245); line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;"><span class="bold italic underline" style="font-weight: bold; text-decoration: underline; font-style: italic;">测试图片</span></span><div style="padding: 3px 0"><img src="https://img.mubu.com/document_image/7fabd28a-8c59-4ffe-b9f3-ab2ef4c91549-40263.jpg" style="max-width: 720px; width: 87px;" class="attach-img"></div></li><li style="line-height: 30px;"><span class="content mubu-node" color="#3da8f5" heading="2" style="color: rgb(61, 168, 245); line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;">是是是</span><ul class="children" style="list-style: disc outside; padding-bottom: 4px;"><li style="line-height: 30px;"><span class="content mubu-node" color="#3da8f5" heading="2" style="color: rgb(61, 168, 245); line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;">ssss</span></li><li style="line-height: 30px;"><span class="content mubu-node" color="#dc2d1e" heading="2" style="color: rgb(220, 45, 30); line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;"><a class="content-link" target="_blank" href="https://mubu.com/doclcoXBPA2D" style="text-decoration: underline; opacity: 0.6; color: inherit;"><span class="bold italic" style="font-weight: bold; font-style: italic;">https://mubu.com/doclcoXBPA2D</span></a></span></li><li style="line-height: 24px;"><span class="content mubu-node" color="#dc2d1e" images="%5B%7B%22id%22%3A%2237d16f3289cb16101%22%2C%22oh%22%3A1004%2C%22ow%22%3A742%2C%22uri%22%3A%22document_image%2F7fabd28a-8c59-4ffe-b9f3-ab2ef4c91549-40263.jpg%22%2C%22w%22%3A87%7D%5D" style="color: rgb(220, 45, 30); line-height: 24px; min-height: 24px; font-size: 16px; padding: 2px 0px; display: inline-block; vertical-align: top;"><a class="content-link" target="_blank" href="https://mubu.com/doclcoXBPA2D" style="text-decoration: underline; opacity: 0.6; color: inherit;"><span class="bold italic" style="font-weight: bold; font-style: italic;">https://mubu.com/doclcoXBPA2D</span></a></span><div style="padding: 3px 0"><img src="https://img.mubu.com/document_image/7fabd28a-8c59-4ffe-b9f3-ab2ef4c91549-40263.jpg" style="max-width: 720px; width: 87px;" class="attach-img"></div></li><li style="line-height: 30px;"><span class="content mubu-node" color="#dc2d1e" heading="2" style="color: rgb(220, 45, 30); line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;"><span class="italic" style="font-style: italic;">sss</span></span></li></ul></li></ul></li><li style="line-height: 30px;"><span class="content mubu-node" color="#dc2d1e" heading="2" style="color: rgb(220, 45, 30); line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;"></span></li><li style="line-height: 30px;"><span class="content mubu-node" color="#dc2d1e" heading="2" style="color: rgb(220, 45, 30); line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;">ok</span></li><li style="line-height: 30px;"><span class="content mubu-node" color="#dc2d1e" heading="2" style="color: rgb(220, 45, 30); line-height: 30px; min-height: 30px; font-size: 21px; padding: 2px 0px; display: inline-block; vertical-align: top;"><a class="content-link" target="_blank" href="https://github.com/alash3al/redix" style="text-decoration: underline; opacity: 0.6; color: inherit;">https://github.com/alash3al/redix</a></span></li><li style="line-height: 27px;"><span class="content mubu-node" color="#dc2d1e" heading="3" style="color: rgb(220, 45, 30); line-height: 27px; min-height: 27px; font-size: 19px; padding: 2px 0px; display: inline-block; vertical-align: top;">标签</span></li></ul></div>

</body>
</html>
`
				//systray.ShowAppWindow("https://www.github.com/getlantern/lantern")
				systray.ShowAppWindow("data:text/html," + url.PathEscape(indexHTML), )
			case <-mToggle.ClickedCh:
				if shown {
					mQuitOrig.Hide()
					mEnabled.Hide()
					shown = false
				} else {
					mQuitOrig.Show()
					mEnabled.Show()
					shown = true
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}

var Data = []byte{0, 0, 1, 0, 1, 0, 16, 16, 0, 0, 1, 0, 32, 0, 104, 4, 0, 0, 22, 0, 0, 0, 40, 0, 0, 0, 16, 0, 0, 0, 32, 0, 0, 0, 1, 0, 32, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 101, 185, 102, 255, 94, 178, 106, 255, 62, 141, 144, 255, 107, 171, 195, 255, 158, 206, 233, 255, 107, 100, 79, 255, 172, 163, 123, 255, 148, 129, 86, 255, 151, 136, 92, 255, 125, 102, 59, 255, 148, 133, 90, 255, 125, 118, 98, 255, 160, 205, 227, 255, 83, 162, 186, 255, 62, 152, 125, 255, 72, 168, 88, 255, 103, 187, 102, 255, 89, 167, 99, 255, 68, 154, 179, 255, 86, 182, 209, 255, 147, 199, 221, 255, 111, 135, 139, 255, 126, 100, 57, 255, 132, 103, 60, 255, 118, 86, 42, 255, 146, 123, 78, 255, 128, 100, 56, 255, 133, 158, 168, 255, 151, 198, 219, 255, 79, 175, 206, 255, 66, 161, 165, 255, 74, 189, 109, 255, 107, 186, 104, 255, 89, 161, 97, 255, 74, 174, 200, 255, 78, 187, 217, 255, 120, 191, 213, 255, 133, 167, 177, 255, 115, 81, 37, 255, 123, 85, 39, 255, 120, 84, 39, 255, 119, 82, 38, 255, 116, 81, 36, 255, 150, 181, 190, 255, 125, 189, 212, 255, 76, 184, 215, 255, 68, 164, 180, 255, 93, 200, 137, 255, 82, 142, 90, 255, 81, 141, 88, 255, 75, 169, 188, 255, 77, 182, 214, 255, 84, 188, 219, 255, 104, 142, 143, 255, 110, 80, 38, 255, 124, 122, 105, 255, 138, 157, 152, 255, 120, 130, 117, 255, 108, 80, 39, 255, 120, 147, 142, 255, 86, 188, 217, 255, 79, 190, 219, 255, 72, 170, 178, 255, 102, 205, 139, 255, 128, 180, 142, 255, 111, 175, 111, 255, 84, 167, 161, 255, 78, 183, 215, 255, 83, 192, 231, 255, 74, 178, 206, 255, 73, 143, 157, 255, 120, 163, 183, 255, 125, 162, 186, 255, 121, 163, 181, 255, 73, 142, 157, 255, 81, 190, 216, 255, 78, 191, 221, 255, 76, 182, 211, 255, 74, 163, 156, 255, 97, 188, 130, 255, 94, 177, 104, 255, 76, 143, 84, 255, 93, 163, 102, 255, 81, 188, 212, 255, 84, 191, 229, 255, 85, 189, 218, 255, 116, 188, 218, 255, 137, 182, 205, 255, 141, 179, 200, 255, 139, 183, 206, 255, 118, 193, 215, 255, 83, 192, 219, 255, 79, 190, 224, 255, 76, 180, 206, 255, 91, 184, 131, 255, 107, 202, 127, 255, 103, 174, 91, 255, 95, 162, 90, 255, 105, 172, 94, 255, 91, 166, 117, 255, 100, 183, 204, 255, 150, 197, 215, 255, 197, 218, 220, 255, 222, 235, 233, 255, 227, 240, 234, 255, 216, 233, 225, 255, 192, 215, 216, 255, 135, 192, 207, 255, 88, 182, 208, 255, 96, 178, 150, 255, 117, 191, 119, 255, 100, 182, 104, 255, 114, 179, 94, 255, 103, 169, 90, 255, 97, 165, 90, 255, 120, 176, 130, 255, 130, 173, 192, 255, 180, 205, 209, 255, 218, 231, 228, 255, 224, 236, 233, 255, 225, 238, 231, 255, 224, 238, 231, 255, 221, 234, 229, 255, 176, 203, 206, 255, 116, 160, 177, 255, 102, 159, 109, 255, 119, 186, 111, 255, 96, 165, 98, 255, 109, 178, 94, 255, 100, 167, 91, 255, 87, 147, 82, 255, 101, 141, 119, 255, 139, 176, 194, 255, 176, 210, 216, 255, 199, 222, 223, 255, 210, 226, 225, 255, 211, 227, 225, 255, 209, 226, 225, 255, 204, 223, 224, 255, 182, 209, 218, 255, 126, 166, 185, 255, 76, 114, 92, 255, 106, 170, 110, 255, 115, 191, 124, 255, 120, 194, 109, 255, 101, 175, 93, 255, 81, 149, 82, 255, 13, 21, 15, 255, 92, 116, 125, 255, 161, 206, 225, 255, 163, 203, 230, 255, 162, 203, 226, 255, 156, 196, 223, 255, 160, 202, 228, 255, 150, 193, 216, 255, 159, 201, 222, 255, 70, 89, 95, 255, 59, 85, 62, 255, 98, 159, 113, 255, 119, 200, 124, 255, 105, 194, 111, 255, 95, 170, 93, 255, 54, 99, 57, 255, 13, 20, 13, 255, 91, 115, 124, 255, 139, 175, 192, 255, 130, 159, 172, 255, 144, 185, 206, 255, 142, 183, 208, 255, 140, 178, 202, 255, 115, 145, 157, 255, 148, 187, 203, 255, 42, 56, 63, 255, 22, 34, 24, 255, 124, 190, 132, 255, 117, 200, 120, 255, 95, 204, 124, 255, 84, 159, 97, 255, 77, 145, 88, 255, 19, 36, 23, 255, 60, 73, 77, 255, 124, 164, 182, 255, 153, 196, 216, 255, 159, 202, 222, 255, 163, 207, 226, 255, 150, 194, 214, 255, 149, 190, 211, 255, 136, 173, 194, 255, 30, 41, 46, 255, 6, 12, 7, 255, 80, 129, 88, 255, 107, 196, 118, 255, 90, 208, 134, 255, 76, 162, 106, 255, 48, 93, 62, 255, 25, 39, 27, 255, 13, 16, 14, 255, 78, 113, 122, 255, 104, 151, 165, 255, 148, 192, 212, 255, 153, 195, 216, 255, 143, 181, 202, 255, 94, 121, 138, 255, 75, 98, 115, 255, 8, 11, 11, 255, 34, 55, 41, 255, 96, 161, 113, 255, 108, 196, 127, 255, 70, 158, 105, 255, 71, 155, 107, 255, 75, 155, 112, 255, 27, 57, 40, 255, 5, 8, 6, 255, 31, 42, 45, 255, 113, 139, 153, 255, 162, 204, 228, 255, 94, 122, 140, 255, 99, 128, 147, 255, 33, 44, 49, 255, 22, 31, 36, 255, 0, 3, 2, 255, 21, 37, 28, 255, 97, 172, 124, 255, 111, 199, 132, 255, 79, 173, 114, 255, 96, 223, 146, 255, 94, 225, 153, 255, 57, 130, 94, 255, 25, 46, 37, 255, 26, 39, 31, 255, 51, 77, 67, 255, 126, 162, 178, 255, 29, 39, 41, 255, 24, 36, 39, 255, 5, 12, 5, 255, 1, 8, 2, 255, 15, 27, 20, 255, 49, 88, 74, 255, 101, 207, 140, 255, 105, 202, 125, 255, 70, 153, 93, 255, 79, 186, 121, 255, 81, 193, 122, 255, 88, 208, 134, 255, 81, 176, 121, 255, 82, 167, 107, 255, 51, 112, 74, 255, 23, 33, 31, 255, 1, 5, 2, 255, 7, 13, 7, 255, 64, 124, 83, 255, 88, 170, 128, 255, 107, 222, 165, 255, 90, 197, 148, 255, 109, 240, 167, 255, 111, 228, 147, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
