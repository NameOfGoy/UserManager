package database

import(
	"github.com/garyburd/redigo/redis"
)
//全局引用
var Conn redis.Conn

func init()  {
	conn,err := redis.Dial("tcp","127.0.0.1:6379")
	Conn = conn
	if err != nil {
		panic(err)
	}
	return
}