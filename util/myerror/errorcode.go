package myerror

/***************例外**********/
// 1.   错误码：-1  通用报错，与前台无关
// 2.   错误码：1 ，表示成功

/**********************数据库*****************************/

const ERROR_INSERT_TABLE_CODE int = 50001
const ERROR_INSERT_TABLE_MESSAGE string = "数据库insert错误"

const ERROR_SELECT_TABLE_CODE int = 50002
const ERROR_SELECT_TABLE_MESSAGE  string = "数据库查询表错误"

const ERROR_INSERT_REDIS_KEY_CODE int = 50003
const ERROR_INSERT_REDIS_KEY_MESSAGE string = "Redis插入key值失败"

const ERROR_DELETE_REDIS_KEY_CODE int = 50004
const ERROR_DELETE_REDIS_KEY_MESSAGE string = "Redis删除key值失败"

const ERROR_UPDATE_TABLE_CODE int = 50005
const ERROR_UPDATE_TABLE_MESSAGE string = "数据库表更新失败"

const ERROR_SELECT_NOTFOUND_CODE int = 50006
const ERROR_SELECT_NOTFOUND_MESSAGE string = "查询结果不存在"

/***********************用户信息管理************************************/
const ERROR_USER_NOTFOUND_CODE int = 10001
const ERROR_USER_NOTFOUND_MESSAGE string = "用户未注册"

const ERROR_USER_WRONGPASSWORD_CODE int = 10002
const ERROR_USER_WRONGPASSWORD_MESSAGE string = "密码错误"

const ERROR_USER_HAVEEXISTED_CODE int = 10003
const ERROR_USER_HAVEEXISTED_MESSAGE string = "用户已存在"

const ERROR_LOGINID_NOTMATCH_CODE int = 10004
const ERROR_LOGINID_NOTMATCH_MESSAGE string = "账号不符合要求"

const ERROR_PASSWORD_NOTMATCH_CODE int = 10005
const ERROR_PASSWORD_NOTMATCH_MESSAGE string = "密码不符合要求"

const ERROR_USERNAME_NOTMATCH_CODE int = 10006
const ERROR_USERNAME_NOTMATCH_MESSAGE string = "用户名不符合要求"

const ERROR_USERSEX_NOTMATCH_CODE int = 10007
const ERROR_USERSEX_NOTMATCH_MESSAGE string = "用户性别不符合要求"

/*********************中间件***********************************/
const ERROR_TOKEN_NOTFOUND_CODE int = 60001
const ERROR_TOKEN_NOTFOUND_MESSAGE string = "未携带token"

const ERROR_TOKEN_TIMEOUT_CODE int = 60002
const ERROR_TOKEN_TIMEOUT_MESSAGE string = "token已过期"

const ERROR_TOKEN_CREATE_CODE int = 60003
const ERROR_TOKEN_CREATE_MESSAGE string = "token生成失败"

const ERROR_TOKEN_ILLEGEL_CODE int = 60004
const ERROR_TOKEN_ILLEGEL_MESSAGE string = "token非法"

/********************gin错误*****************************/

const ERROR_JSONBIND_FAILED_CODE int = 70001
const ERROR_JSONBIND_FAILED_MESSAGE string = "JSON绑定失败"

/****************序列化****************/

const ERROR_STRUCT_TO_STRING_CODE int = 80001
const ERROR_STRUCT_TO_STRING_MESSAGE string = "JSON转字符串失败"

const ERROR_STRING_TO_STRUCT_CODE int = 80002
const ERROR_STRING_TO_STRUCT_MESSAGE string = "字符串转结构体失败"

/****************API调用****************/

const ERROR_PARAMS_MISSING_CODE int = 90001
const ERROR_PARAMS_MISSING_MESSAGE string = "缺少参数"