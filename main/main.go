
package main
import (
    //这里讲db作为go/databases的一个别名，表示数据库连接池
    db "UserManager/database"
    . "UserManager/router"
    _"github.com/icattlecoder/godaemon"
)
func main() {
    //当整个程序完成之后关闭数据库连接
    defer db.Engine.Close()
    router := InitRouter()
    router.Run(":8088")
}   