package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	ProblemId     uint           `json:"problemId" gorm:"column:problem_id;type:varchar(36);"`      // 问题的id
	CategoryId    uint           `json:"categoryId" gorm:"column:category_id;type:varchar(36);"`    // 分类的id
	CategoryBasic *CategoryBasic `json:"categoryBasic" gorm:"foreignKey:id;references:category_id"` // 关联分类的基础信息表
}

func (ProblemCategory) TableName() string {
	return "problem_category"
}
