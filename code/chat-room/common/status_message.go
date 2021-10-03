package common

import "errors"

const (
	// status code for login
	ServerError   = 500
	PasswordError = 400
	UserNotExist  = 401
	LoginSucceed  = 200

	// status code for register
	UserHasExited    = 400
	PasswordNotMatch = 401
	RegisterSucceed  = 200
)

//根据业务逻辑需要，自定义一些错误
//这个nb 这样我就不用自己定义红了
var (
	ERROR_USER_DOES_NOT_EXIST = errors.New("User does not exist!")
	ERROR_USER_PWD            = errors.New("Password is invalid!")

	// status code for register
	ERROR_USER_ALREADY_EXISTS     = errors.New("Username already exists!")
	ERROR_PASSWORD_DOES_NOT_MATCH = errors.New("Password does not match!")

	Register_Success = "Register Success"
	Login_Success    = "Login Success"
	Login_Failed     = "Login Failed"
)
