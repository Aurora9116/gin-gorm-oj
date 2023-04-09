package test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`

	jwt.StandardClaims
}

var myKey = "gin-gorm-oj"

// 生成token
func TestGenerateToken(t *testing.T) {
	userClaims := &UserClaims{
		Identity:       "user_1",
		Name:           "Get",
		StandardClaims: jwt.StandardClaims{},
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenStrign, err := claims.SignedString([]byte(myKey))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokenStrign)
}

// 解析token
func TestAnalyseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.kEQctWkYxI4UXVevGaZuquesF5aJeFlRu5BBuPdSfg0"

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(myKey), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if claims, ok := token.Claims.(*UserClaims); token.Valid && ok {
		fmt.Println(claims)
	}

}
