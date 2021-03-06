package utils

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"grpc-demo/logger"
	"time"
)

const (
	KEY                    string = "JWT-ARY-STARK"
	DEFAULT_EXPIRE_SECONDS int    = 600 //默认过期时间（s）
)

type User struct {
	UserName string `json:"user"`
}

// JWT -- json web token
// HEADER PAYLOAD SIGNATURE
// This struct is the PAYLOAD
type MyCustomClaims struct {
	User
	jwt.StandardClaims
}

//验证jtw token
func ValidateToken(tokenString string) (info User, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(KEY), nil
		})

	if err != nil {
		return User{}, errors.New("ParseWithClaims Failed")
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		logger.Info("%v %v", claims.User, claims.StandardClaims.ExpiresAt)
		logger.Info("token will be expired at ", time.Unix(claims.StandardClaims.ExpiresAt, 0))
		info = claims.User
	} else {
		//fmt.Println("validate tokenString failed !!!", err)
		logger.Info("%v %v", claims.User, claims.StandardClaims.ExpiresAt)
		logger.Info("token will be expired at ", time.Unix(claims.StandardClaims.ExpiresAt, 0))
		return User{}, errors.New("validation failed")
	}
	return info, nil
}

//获取jwt token
func GenerateToken(info *User, expiredSeconds int) (tokenString string) {
	if expiredSeconds == 0 {
		expiredSeconds = DEFAULT_EXPIRE_SECONDS
	}
	// Create the Claims
	mySigningKey := []byte(KEY)

	expireAt := time.Now().Add(time.Second * time.Duration(expiredSeconds)).Unix()

	//fmt.Println("token will be expired at ", time.Unix(expireAt, 0))
	// pass parameter to this func or not
	user := *info

	claims := MyCustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			Issuer:    user.UserName,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//对token进行加密
	tokenStr, err := token.SignedString(mySigningKey)
	if err != nil {
		//fmt.Println("generate json web token failed !! error :", err)
	} else {
		tokenString = tokenStr
		//logs.Info("JWT 设置成功")
		//logs.Info("tokenStr is ",tokenStr)
	}
	return tokenString
}
