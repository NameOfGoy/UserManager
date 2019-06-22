package check

import(
	"encoding/json"
	e "UserManager/util/myerror"
	"regexp"
)

func RegistParamCheck(v interface{})  (err error) {
	registinfo := struct {
		Loginid string  		`json:"loginid"`
		Password string 		`json:"password"`
		Username string 		`json:"username"`
		Birthday string			`json:"birthday"`
		Sex string				`json:"sex"`
		Address string			`json:"address"`
		Email string      		`json:"email"`
		PhoneNum string			`json:"phoneNum"`
	}{}
	b, err := json.Marshal(v)
	if err != nil {
		err = &e.Myerror{
			Code : e.ERROR_STRUCT_TO_STRING_CODE,
			Message : e.ERROR_STRUCT_TO_STRING_MESSAGE,
		}
		return
	}
	err = json.Unmarshal(b, &registinfo)
	if err != nil {
		err = &e.Myerror{
			Code : e.ERROR_STRING_TO_STRUCT_CODE,
			Message : e.ERROR_STRING_TO_STRUCT_MESSAGE,
		}
		return
	}

	// 检测账号合法性
	match, _:= regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_-]{7,19}$",registinfo.Loginid)
	if !match {
		err = &e.Myerror{
			Code : e.ERROR_LOGINID_NOTMATCH_CODE,
			Message : e.ERROR_LOGINID_NOTMATCH_MESSAGE,
		}
		return
	}

	if registinfo.Password == "" {
		err = &e.Myerror{
			Code : e.ERROR_PASSWORD_NOTMATCH_CODE,
			Message : e.ERROR_PASSWORD_NOTMATCH_MESSAGE,
		}
		return
	}

	if registinfo.Username == "" {
		err = &e.Myerror{
			Code : e.ERROR_USERNAME_NOTMATCH_CODE,
			Message : e.ERROR_USERNAME_NOTMATCH_MESSAGE,
		}
		return
	}

	if registinfo.Sex != "" {
		if registinfo.Sex != "male" && registinfo.Sex != "female" {
			err = &e.Myerror{
				Code : e.ERROR_USERSEX_NOTMATCH_CODE,
				Message : e.ERROR_USERSEX_NOTMATCH_MESSAGE,
			}
			return
		}
	}
	return
}