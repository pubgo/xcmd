package abc

import (
	"github.com/pubgo/g/models"
	"github.com/pubgo/xcmd/xcmd/xcmd_res/res"
)

// IResLinkStore
type IResLinkStore interface {
	Cfg(*res.Config) error
	Driver() string
	Create(name string) (models.ResLink, error)
	Walk(func(models.ResLink) error) error
	Update(models.ResLink) error
	Delete(name string) error
	Find(name string) (models.ResLink, error)
	Search(match string, cb func(models.ResLink) error) error
}

type IResLink interface {
	// 平台名称
	Platform() string
	// 遍历数据
	Walk(func(models.ResLink) error) error
	// 获取单篇数据, 根据传递的格式内容，转换成对应的格式
	GetRes(tpy string) (models.ResLink, error)
	// 删除单篇
	DelRes() error
	// 更新单篇
	UpdateRes(tpy string) error
	// 创建
	CreateRes(tpy string) error
	// 搜索
	Search(kw string) ([]models.ResLink, error)
}

type IResLinkCmd interface {
	// 平台名称
	Platform() string
	// 遍历数据
	Walk(func(models.ResLink) error) error
	// 获取单篇数据, 根据传递的格式内容，转换成对应的格式
	GetRes(tpy string) (models.ResLink, error)
	// 删除单篇
	DelRes() error
	// 更新单篇
	UpdateRes(tpy string) error
	// 创建
	CreateRes(tpy string) error
	// 搜索
	Search(kw string) ([]models.ResLink, error)
}
