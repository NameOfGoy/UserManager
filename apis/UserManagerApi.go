package apis

import(
	"github.com/gin-gonic/gin"
	"net/http"
	"UserManager/models"
	"time"
	"math/rand"
	myjwt "UserManager/util/jwt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"fmt"
	"strconv"
	db "UserManager/database"
	"UserManager/util/check"
	e "UserManager/util/myerror"
	"encoding/json"
	"log"
)

func Index(c *gin.Context) {
	c.String(http.StatusOK,"UserManager API Project works ！")
}

func UserRegister(c *gin.Context) {
	registinfo := struct {
		Loginid string  		`json:"loginid"`
		Password string 		`json:"password"`
		Username string 		`json:"username"`
		Birthday string			`json:"birthday"`
		Sex string				`json:"sex"`
		Address string			`json:"address"`
		Email string      		`json:"email"`
		PhoneNum string			`json:"phonenum"`
	}{}
	err := c.ShouldBind(&registinfo)  // 绑定json
	if  err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &e.Myerror{
				Code : e.ERROR_JSONBIND_FAILED_CODE,
				Message : e.ERROR_JSONBIND_FAILED_MESSAGE,
			},
		})
		return
	}
	
	// 校验参数合法性 
	err = check.RegistParamCheck(&registinfo)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}

	// 转换成功，赋值给对应表的model
	userid := strconv.FormatInt(time.Now().Unix(),10) + strconv.Itoa(rand.Intn(1000)) // 时间戳+后三位随机数
	fmt.Println(userid)
	user := models.User{
		Loginid : registinfo.Loginid,
		Password : registinfo.Password,
		Userid : userid,
	}
	userinfo := models.Userinfo{
		Userid   :  userid,
		Username :	registinfo.Username,
		Birthday :	registinfo.Birthday,
		Sex		 :	registinfo.Sex,
		Address	 :	registinfo.Address,
		Email	 :	registinfo.Email,
		Phonenum :	registinfo.PhoneNum,
	}
	// 插2张表，事务操作，使用存储过程
	log.Println(user)
	log.Println(userinfo)
	err = models.RegistProcedure(&user,&userinfo)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"response" : &models.ResponseSet{
			Code : 1,
			Message : "注册成功",
		},
	})
	return
}

func UserLogin(c *gin.Context) {
	logininfo := struct {
		Loginid string  		`json:"loginid"`
    	Password string 		`json:"password"`
	}{}

	// 绑定json
	err := c.ShouldBind(&logininfo)
	if  err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &e.Myerror{
				Code : e.ERROR_JSONBIND_FAILED_CODE,
				Message : e.ERROR_JSONBIND_FAILED_MESSAGE,
			},
		})
		return
	}

	// 转换成功，登录操作
	// 赋值
	user := models.User {
		Loginid  : logininfo.Loginid,
		Password : logininfo.Password,
	}
	fmt.Println(user)
	// 查找数据库
	user2,err := user.GetUser()
	if  err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	// 对比密码
	if user.Password != user2.Password {
		c.JSON(http.StatusOK,gin.H{
			"response" : &e.Myerror{
				Code : e.ERROR_USER_WRONGPASSWORD_CODE,
				Message : e.ERROR_USER_WRONGPASSWORD_MESSAGE,
			},
		})
		return
	}

	//密码正确，生成token
	token,err := generateToken(c,user)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}

	//存token进redis
	_,err = db.Conn.Do("SET",user.Loginid,token,"EX","3600") // key:登录名  value : token
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &e.Myerror{
				Code : e.ERROR_INSERT_REDIS_KEY_CODE,
				Message : e.ERROR_INSERT_REDIS_KEY_MESSAGE,
			},
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"response" : &models.ResponseSet{
			Code : 1,
			Message : "登录成功",
			Data : token,
		},
	})
	return
}

func ExitLogin(c *gin.Context) {
	//检查token是否合法
	claims,err := check.CheckToken(c)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	//删除redis的token缓存
	_,err = db.Conn.Do("DEL", claims.ID)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &e.Myerror{
				Code : e.ERROR_DELETE_REDIS_KEY_CODE,
				Message : e.ERROR_DELETE_REDIS_KEY_MESSAGE,
			},
		})
		return
	}
	//返回
	c.JSON(http.StatusOK,gin.H{
		"response" : &models.ResponseSet{
			Code : 1,
			Message : "已退出登录",
		},
	})
}

//生成令牌
func generateToken(c *gin.Context, user models.User) (token string,err error){
    j := &myjwt.JWT{
        SigningKey : []byte("newtrekWang"),
    }
    StandardClaims := myjwt.CustomClaims{
		ID : user.Loginid,
		StandardClaims : jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 100), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
			Issuer:    "goy",                   //签名的发行者
			Audience: user.Loginid,
		},      
    }
 
	token, err = j.CreateToken(StandardClaims) // 取得token
	if err != nil{
		err = &e.Myerror{
			Code : e.ERROR_TOKEN_CREATE_CODE,
			Message : e.ERROR_TOKEN_CREATE_MESSAGE,
		}
		return
	}
    return
}

func ModifyPassword(c *gin.Context){
	//检查token是否合法
	claims,err := check.CheckToken(c)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}

	oldnewpassword := struct {
		Oldpassword string  `json:"oldpassword"`
    	Newpassword string 	`json:"newpassword"`
	}{}
	err = c.ShouldBind(&oldnewpassword)
	if  err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &e.Myerror{
				Code : e.ERROR_JSONBIND_FAILED_CODE,
				Message : e.ERROR_JSONBIND_FAILED_MESSAGE,
			},
		})
		return
	}
	// 从数据库中取出用户密码，并且对比，正确就更新旧密码
	user := models.User{
		Loginid : claims.ID,
	}
	user,err = user.GetUser()
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	log.Println(user)
	if user.Password != oldnewpassword.Oldpassword {
		c.JSON(http.StatusOK,gin.H{
			"response" : &e.Myerror{
				Code : e.ERROR_USER_WRONGPASSWORD_CODE,
				Message : e.ERROR_USER_WRONGPASSWORD_MESSAGE,
			},
		})
		return
	}
	user.Password = oldnewpassword.Newpassword
	err = user.UpdateUser()  //更新
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	// 成功
	c.JSON(http.StatusOK,gin.H{
		"response" : &models.ResponseSet{
			Code : 1,
			Message : "密码修改成功",
		},
	})
}

func ModifyUserinfo(c *gin.Context){

	//检查token是否合法
	claims,err := check.CheckToken(c)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	userinfo := models.Userinfo{}
	err = c.ShouldBind(&userinfo)
	if  err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &e.Myerror{
				Code : e.ERROR_JSONBIND_FAILED_CODE,
				Message : e.ERROR_JSONBIND_FAILED_MESSAGE,
			},
		})
		return
	}
	log.Println(userinfo)
	//根据token存储的登陆账号，获取该账号对应的userid，然后根据userid更新表
	user := models.User{
		Loginid : claims.ID,
	}
	user,err = user.GetUser()
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	userinfo.Userid = user.Userid
	err = userinfo.ModifyUserinfo()
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	// 成功
	c.JSON(http.StatusOK,gin.H{
		"response" : &models.ResponseSet{
			Code : 1,
			Message : "修改成功",
		},
	})
}

func GetUserinfo(c *gin.Context) {

	//检查token是否合法
	claims,err := check.CheckToken(c)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	// 根据token的loginid找到userid
	user := models.User{
		Loginid : claims.ID,
	}
	user,err = user.GetUser()
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	userinfo := models.Userinfo{
		Userid : user.Userid,
	}
	userinfo,err = userinfo.GetUserinfo()
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	// 成功
	remix := struct{
		Loginid string
		Userid string
		Username string
		Birthday string
		Sex string
		Address string
		Email string
		Phonenum string
	}{
		Loginid : user.Loginid,
		Userid : userinfo.Userid,
		Username : userinfo.Username,
		Birthday : userinfo.Birthday,
		Sex : userinfo.Sex,
		Address : userinfo.Address,
		Email : userinfo.Email,
		Phonenum : userinfo.Phonenum,
	}
	c.JSON(http.StatusOK,gin.H{
		"response" : &models.ResponseSet{
			Code : 1,
			Message : "获取用户信息成功",
			Data : &remix,
		},
	})
}

//
func GetUser(c *gin.Context) {
	user := models.User{
		Loginid : c.Query("loginid"),  // 取得url后面附带的参数loginid
	}
	u,err := user.GetUser()    // 调用查询方法
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &err,
		})
		return
	}
	userstr,err := json.Marshal(u) // 结构体序列化
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : &e.Myerror{
				Code : e.ERROR_STRUCT_TO_STRING_CODE,
				Message : e.ERROR_STRUCT_TO_STRING_MESSAGE,
			},
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"response" : &models.ResponseSet{
			Code : 1,
			Message : "已退出登录",
			Data : string(userstr),
		},
	})
}

// GetDataByTime 一个需要token认证的测试接口
func GetDataByTime(c *gin.Context) {
	_,err := check.CheckToken(c)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{
            "status": -1,
            "msg":    "token验证失败",
            "response":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "token验证通过",
	})
	return
	
}

func GetToken(c *gin.Context) {
	j := &myjwt.JWT{
        SigningKey : []byte("newtrekWang"),
    }
    StandardClaims := myjwt.CustomClaims{
		ID : "12313",
		StandardClaims : jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix()), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 60), // 过期时间 一小时
			Issuer:    "goy",                   //签名的发行者
			Audience: "fff",
		},      
    }
	token, err := j.CreateToken(StandardClaims) // 取得token
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"response" : err.Error(),
			"message" : "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message" : "获取token成功",
		"token" : token,
	})
}

