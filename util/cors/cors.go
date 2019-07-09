package cors

import(
	"github.com/gin-gonic/gin"
	//"fmt"
	//"strings"
	"net/http"
)

func Cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        method := c.Request.Method

        // 请求进来的时候，header添加以下内容
        c.Header("Access-Control-Allow-Origin", "*") //允许所有域名
        c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")// 允许所列header
        c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS") //允许所列方法
        c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
        c.Header("Access-Control-Allow-Credentials", "true")
    
        //放行所有OPTIONS方法
        if method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
        }
        // 处理请求
        c.Next()
    }

}