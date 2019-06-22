package models

import(
	db "UserManager/database"
	_"errors"
	"log"
	e "UserManager/util/myerror"
)

type Userinfo struct {
	Userid   string    `xorm:"not null pk VARCHAR(100)"`
	Username string `xorm:"VARCHAR(255)"`
	Birthday string `xorm:"VARCHAR(255)"`
	Sex      string `xorm:"VARCHAR(255)"`
	Address  string `xorm:"VARCHAR(255)"`
	Email    string `xorm:"VARCHAR(255)"`
	Phonenum string `xorm:"VARCHAR(100)"`
}

func (u *Userinfo)AddUserinfo() (err error) {
	_,err = db.Engine.Insert(u)
	if err != nil {
		err = &e.Myerror{
			Code : e.ERROR_INSERT_TABLE_CODE,
			Message : "Userinfo表插入错误",
		}
		return
	 }
	return
}


type User struct {
	Loginid  string `xorm:"not null pk VARCHAR(20)"`
	Password string `xorm:"VARCHAR(255)"`
	Userid   string    `xorm:"VARCHAR(100)"`
}

func (u *User)AddUser() (err error) {
	_,err = db.Engine.Insert(u)
	if err != nil {
		err = &e.Myerror{
			Code : e.ERROR_INSERT_TABLE_CODE,
			Message : "User表插入错误",
		}
		return
	 }
	return
}

func (u *User)GetUser() (user User,err error) {
	log.Println(u)
	ret,err := db.Engine.ID(u.Loginid).Get(&user)
	if err != nil{
		err = &e.Myerror{
			Code : e.ERROR_SELECT_TABLE_CODE,
			Message : "user表查询错误", 
		}
		return 
	}
	if !ret {
		err = &e.Myerror{
			Code : e.ERROR_USER_NOTFOUND_CODE,
			Message : e.ERROR_USER_NOTFOUND_MESSAGE,
		}
	}
	return
}

type LoginInfo struct {
	Loginid string  		`json:"loginid"`
    Password string 		`json:"password"`
}


func RegistProcedure (user *User,userinfo *Userinfo) (err error) {
	
	// 先查询账号是否存在, nil则表示GetUser方法执行没出错，找得到用户
	_,err = user.GetUser()
	if err == nil {
		err = &e.Myerror{
			Code : e.ERROR_USER_HAVEEXISTED_CODE,
			Message : e.ERROR_USER_HAVEEXISTED_MESSAGE,
		}
		return
	}

	session := db.Engine.NewSession()
	defer session.Close()

	// add Begin() before any action
	err = session.Begin()
	if err != nil {
		err = &e.Myerror{
			Code : -1,
			Message : "注册事务session开启失败",
		}
		return
	}

	_,err = session.Insert(userinfo)
	if err != nil {
		session.Rollback()
		log.Println("err : ",err,"    rolled back !")
		err = &e.Myerror{
			Code : e.ERROR_INSERT_TABLE_CODE,
			Message : "Userinfo表插入错误",
		}
		return
	}

	_,err = session.Insert(user)
	if err != nil {
		session.Rollback()
		log.Println("err : ",err,"    rolled back !")
		err = &e.Myerror{
			Code : e.ERROR_INSERT_TABLE_CODE,
			Message : "User表插入错误",
		}
		return
	}

	// add Commit() after all actions
	err = session.Commit()
	if err != nil {
		err = &e.Myerror{
			Code : -1,
			Message : "注册事务commit错误",
		}
		return
	}

	return
}