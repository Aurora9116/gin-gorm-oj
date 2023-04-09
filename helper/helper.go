package helper

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = "gin-gorm-oj"

// GenerateToken 生成token
func GenerateToken(name string, identity string, isAdmin int) (string, error) {
	userClaims := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenString, err := claims.SignedString([]byte(myKey))
	if err != nil {
		return "", err
	}
	fmt.Println(tokenString)
	return tokenString, nil
}

// AnalyseToken 解析token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(myKey), nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if claims, ok := token.Claims.(*UserClaims); token.Valid && ok {
		return claims, nil
	}
	return nil, errors.New("parseToken err")
}

func SendEmail(toUsers []string, code string) error {
	e := email.NewEmail()
	e.From = "from <949244762@qq.com>" // 发送者
	e.To = toUsers                     // 接收者
	e.Subject = "验证码，请不要泄露"
	e.HTML = []byte("您的验证码是：<b>" + code + "</b>")
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "949244762@qq.com", "qtssdidxkuytbcah", "smtp.qq.com"))
	//mtkubqxqhwfxbbff
	//err := e.SendWithTLS("smtp.qq.com:587", smtp.PlainAuth("", "Aurora@gmail.com", "qtssdidxkuytbcah", "smtp.qq.com"),
	//	&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		return err
	}
	return nil
}

// GetUUID 生成uuid
func GetUUID() string {
	return uuid.NewV4().String()
}

// GetRand 生成验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 10; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

// CodeSave 代码保存
func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write(code)
	defer f.Close()
	return path, nil
}
