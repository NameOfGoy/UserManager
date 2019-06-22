package router
import (
    "github.com/gin-gonic/gin"
    ."UserManager/apis"
    "UserManager/util/jwt"
)
func InitRouter() *gin.Engine {
    router := gin.Default()
    //IndexApi为一个Handler
    router.GET("/", Index)
    router.POST("/Register",UserRegister)
    router.POST("/Login",UserLogin)
    router.GET("/gettoken",GetToken)
    router.GET("/ExitLogin",ExitLogin)
    router.GET("/getuser",GetUser)
   

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