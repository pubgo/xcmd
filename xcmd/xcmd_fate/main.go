package main

import (
	"context"
	"github.com/godcong/fate"
)

func main() {
	//使用前请导入pre的数据（测试字库尚未完善，生成姓名后可以去一些测名网站验证下）
	//连接mysql数据库
	fate.InitAll()
	eng := fate.InitMysql("192.168.1.161:3306", "root", "111111")
	//生日
	c := chronos.New("2020/01/23 11:31")
	//姓名的最少笔画数（可不设）
	fate.DefaultStrokeMin = 3
	//姓名的最大笔画数（可不设）
	fate.DefaultStrokeMax = 15

	//设定数据库：fate.Database(eng)
	//开启八卦过滤：fate.BaGuaFilter()
	//开启生肖过滤：fate.ZodiacFilter()
	//开启喜用神过滤：fate.SupplyFilter()
	//第一参数：姓
	//第二参数：生日
	f := fate.NewFate("王", c.Solar().Time(), fate.Database(eng), fate.BaGuaFilter(), fate.ZodiacFilter(), fate.SupplyFilter())

	e := f.MakeName(context.Background())
	if e != nil {
		t.Fatal(e)
	}
}
