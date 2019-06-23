package database
import (

    _ "github.com/go-sql-driver/mysql"
	"log"
	"github.com/go-xorm/xorm"
	_"github.com/go-xorm/core"
)
//因为我们需要在其他地方使用Engine这个变量，所以需要大写代表public
var Engine *xorm.Engine
//初始化方法，每个package 中的init()方法会被自动加载
func init() {
    var err error
	Engine, err = xorm.NewEngine("mysql", "root:liang.0804@tcp(127.0.0.1:3306)/video?charset=utf8")
	if err != nil{
		log.Fatal("数据库连接失败！！",err)
	}
    //连接检测
    err = Engine.Ping()
    if err != nil {
        log.Fatal(err.Error())
    }
}