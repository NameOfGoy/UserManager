package myerror

import(
	"encoding/json"
)

type Myerror struct {
	Code int
	Message string
}

func (err *Myerror) Error() string {
	e,_:= json.Marshal(err)
	return string(e)
}