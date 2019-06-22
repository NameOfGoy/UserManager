package apis

import(
	"github.com/gin-gonic/gin"
	"UserManager/util/check"
	"net/http"
	"github.com/Unknown/goconfig"
	e "UserManager/util/myerror"
	"UserManager/models"
	"encoding/json"
	_"strings"
)

func GetVideoInfo(c *gin.Context) {
	// 检查token
	_,err := check.CheckToken(c)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"err" : err,
		})
		return
	} 
	//获取配置文件里的API的url
	cfg, err := goconfig.LoadConfigFile("../util/conf/conf.ini")
	if err != nil{
		c.JSON(http.StatusOK,gin.H{
			"err" : err,
		})
		return
	}
	url,err := cfg.GetValue("vod-video-api","VideoInfoAPI")
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"err" : err,
		})
		return
	}

	// 判断参数
	if videoid := c.Query("VideoId"); videoid == "" {
		c.JSON(http.StatusOK,gin.H{
			"err" : &e.Myerror{
				Code : e.ERROR_PARAMS_MISSING_CODE,
				Message : "缺少参数VideoId",
			},
		})
		return
	} 

	// 构建http请求
	req,err := http.NewRequest("GET",url+"?VideoId="+c.Query("VideoId"),nil)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"err" : &e.Myerror{
				Code : -1,
				Message : "http请求构建失败",
			},
		})
		return
	}
	client := &http.Client{}
	res,err := client.Do(req) //发送请求
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"err" : &e.Myerror{
				Code : -1,
				Message : "http请求连接失败",
			},
		})
	}

	defer res.Body.Close() //函数执行完毕后，关闭连接

	buf := make([]byte,1024)
	n,err := res.Body.Read(buf) // 读取获取到的body
	if err != nil && err.Error() != "EOF" { // 读取完成时，err为EOF，意思时文件尾部，不是err
		c.JSON(http.StatusOK,gin.H{
			"err" : &e.Myerror{
				Code : -1,
				Message : "读取buf失败",
			},
		})
		return
	}
	buf2 := make([]byte,n) // 重新用一个确定真实容量的[]byte
	for i:=0;i<n;i++ {
        buf2[i] = buf[i]
	}
	
	// 构造视频信息结构体
	videoinfo := struct{
		FileId     string `json:"FileId"`
		Name     string `json:"Name"`
		ClassId int64 `json:"ClassId"`
		ClassName      string `json:"ClassName"`
		ClassPath string `json:"ClassPath"`
		CoverUrl string `json:"CoverUrl"`
		Description string `json:"Description"`
		Duration float64 `json:"Duration"`
		CreateTime string `json:"CreateTime"`
	}{}

	err = json.Unmarshal(buf2, &videoinfo)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"err" : &e.Myerror{
				Code : -1,
				Message : "反序列化失败",
			},
		})
		return
	} 

	c.JSON(http.StatusOK,gin.H{
		"response" : &models.ResponseSet{
			Code : 1,
			Message : "success",
			Data : videoinfo,
		},
	})
}