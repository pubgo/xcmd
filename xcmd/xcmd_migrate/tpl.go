package xcmd_migrate

var initTpl = `package migrations

import (
	"github.com/pubgo/g/xmigrate"
)

var registry = xmigrate.Registry
`
var migrationTpl = `
package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pubgo/g/xerror"
	"gopkg.in/gormigrate.v1"
)

func init() {
	registry(&gormigrate.Migration{
		ID: "%s",
		Migrate: func(tx *gorm.DB) (err error) {
			defer xerror.RespErr(&err)

			return xerror.Wrap(tx.AutoMigrate(nil).Error, "")
		},
		Rollback: func(tx *gorm.DB) (err error) {
			defer xerror.RespErr(&err)

			return xerror.Wrap(tx.DropTable("").Error, "")
		},
	})
}
`
