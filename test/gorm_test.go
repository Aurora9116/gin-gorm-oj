package test

import (
	"gin-gorm-oj/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormTest(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gin-gorm-oj?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&models.TestCase{}, &models.UserBasic{}, &models.ProblemBasic{}, &models.CategoryBasic{}, &models.ProblemCategory{}, &models.SubmitBasic{})
	//db.Delete(test)
}
