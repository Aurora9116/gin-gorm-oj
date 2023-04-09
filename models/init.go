package models

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB = Init()

var RDB = InitRedis()

func Init() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gin-gorm-oj?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	_ = db.AutoMigrate(&CategoryBasic{}, &ProblemCategory{}, &ProblemBasic{}, &UserBasic{}, &SubmitBasic{})
	return db
}

func InitRedis() *redis.Client {
	var rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}
