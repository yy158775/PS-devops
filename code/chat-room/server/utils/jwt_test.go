package utils

import (
	"log"
	"testing"
)

func TestJwt(t *testing.T) {
	eles := []User{
		{"yy"},
		{"jj"},
		{"kk"},
	}
	for _, user := range eles {
		token := GenerateToken(&user, 600)
		infouser, err := ValidateToken(token)
		if err != nil {
			log.Println("validate failed")
		} else {
			log.Println(infouser.UserName)
		}
	}
}
