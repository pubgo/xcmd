package xcmd_migrate

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pubgo/g/pkg/fileutil"
	"github.com/pubgo/g/xconfig"
	"github.com/pubgo/g/xconfig/xconfig_rds"
	"github.com/pubgo/g/xdi"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/g/xmigrate"
	"github.com/pubgo/xcmd/xcmd"
	"gopkg.in/gormigrate.v1"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func Init() *xcmd.Command {
	var _migrate = false
	var _rollback = false
	var _list = false
	var _create = "test"
	var db *gorm.DB
	var migrationsPath = "migrations"

	var migrations []*gormigrate.Migration
	var _Migrate = func(db *gorm.DB) error {
		return xerror.Wrap(gormigrate.New(db, gormigrate.DefaultOptions, migrations).Migrate(), "db xcmd_migrate error")
	}

	var _Rollback = func(db *gorm.DB) (err error) {
		defer xerror.RespErr(&err)

		_timeParse := func(u int) time.Time {
			_i, err := time.Parse("2006_01_02_15_04", strings.Join(strings.Split(migrations[u].ID, "_")[:5], "_"))
			xerror.PanicM(err, "time parse error, data: %s", migrations[u].ID)
			return _i
		}

		sort.Slice(migrations, func(i, j int) bool {
			return _timeParse(i).Sub(_timeParse(j)) < 0
		})

		return xerror.Wrap(gormigrate.New(db, gormigrate.DefaultOptions, migrations).RollbackLast(), "db rollback error")
	}

	var args = xcmd.Args(func(cmd *xcmd.Command) {
		cmd.Flags().StringVarP(&migrationsPath, "path", "p", migrationsPath, "迁移脚本的目录路径")
		cmd.Flags().StringVarP(&_create, "create", "c", _create, "数据库迁移生成")
		cmd.Flags().BoolVarP(&_migrate, "migrate", "m", _migrate, "数据库迁移")
		cmd.Flags().BoolVarP(&_rollback, "rollback", "r", _rollback, "数据库回滚")
		cmd.Flags().BoolVar(&_list, "ls", _list, "数据库迁移脚本列表")
	})

	return args(&xcmd.Command{
		Use:   "migrate",
		Short: "migrate or rollback",
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			// 检查迁移脚本目录
			xerror.Panic(fileutil.IsNotExistMkDir(migrationsPath))
			// 检查init文件
			_init := filepath.Join(migrationsPath, "init.go")
			if fileutil.CheckNotExist(_init) {
				xerror.PanicM(ioutil.WriteFile(_init, []byte(initTpl), 0644), "gen init file error")
			}

			xerror.Panic(xdi.Invoke(func(cfg *xconfig.Config, rds *xconfig_rds.Rds) (err error) {
				defer xerror.RespErr(&err)

				db = rds.GetRDS()
				xerror.PanicT(db == nil, "db is nil, please init db")

				gormigrate.DefaultOptions.TableName = cfg.Rds.Prefix + "_migrations"
				gormigrate.DefaultOptions.IDColumnSize = 150
				gormigrate.DefaultOptions.UseTransaction = true
				gormigrate.DefaultOptions.ValidateUnknownMigrations = true

				return
			}))

			// 导入数据库迁移脚本列表
			migrations = xmigrate.Migrations()

			if _migrate {
				xerror.PanicM(_Migrate(db), "Migrate Error")
				return
			}

			if _rollback {
				xerror.PanicM(_Rollback(db), "Rollback Error")
				return
			}

			if _list {
				for _, m := range migrations {
					ms := strings.SplitN(m.ID, "--", 2)
					t, err := time.Parse("2006_01_02_15_04", ms[0])
					xerror.Panic(err)
					fmt.Printf("%s\n  %s\n", t.String(), ms[1])
				}
				return
			}

			if _create != "" && _create != "test" {
				_time := time.Now().Format("2006_01_02_15_04")
				_file := filepath.Join(migrationsPath, "m_"+_time)
				logger.Info().Msgf("migration file: %s", _file+".go")
				xerror.PanicM(ioutil.WriteFile(_file+".go", []byte(fmt.Sprintf(migrationTpl, _time+"--"+_create)), 0644), "gen migration error")
				return
			}
			return
		}})
}
