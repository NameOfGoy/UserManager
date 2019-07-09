package router
import (
    "github.com/gin-gonic/gin"
    ."UserManager/apis"
    "UserManager/util/jwt"
    "UserManager/util/cors"
)
func InitRouter() *gin.Engine {
    router := gin.Default()
    
    // 使用跨域中间件，解决复杂跨域问题。浏览器会发两次请求，第一次是option
    router.Use(cors.Cors())

    router.GET("/", Index)
    router.POST("/Register",UserRegister)
    router.POST("/Login",UserLogin)
    router.GET("/gettoken",GetToken)
    router.GET("/ExitLogin",ExitLogin)
    router.GET("/getuser",GetUser)
    router.POST("/ModifyPassword",ModifyPassword)
    router.POST("/ModifyUserinfo",ModifyUserinfo)
    router.GET("/Userinfo",GetUserinfo)
   
    taR := router.Group("/data")
    // 路径是/data时，会把jwt.JWTAuth()调用
    taR.Use(jwt.JWTAuth()) // router.Use是使用中间件的意思

    {
        taR.GET("/dataByTime",GetDataByTime)
    }
    
    video := router.Group("/video")
    {
        video.GET("/videoinfo",GetVideoInfo)
    }

    return router
}