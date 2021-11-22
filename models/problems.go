package models

import "gorm.io/gorm"

//Problem Database - database
type Problem struct {
	gorm.Model
	ProblemName       string `gorm:"type:text;"`
	Description       string `gorm:"type:text;"`
	InputDescription  string `gorm:"type:text;"`
	OutputDescription string `gorm:"type:text"`
	Author            uint   `gorm:"NOT NULL;"`
	MemoryLimit       uint   `gorm:"NOT NULL;"`
	CPUTime           uint   `gorm:"NOT NULL;"`
	Layer             uint8  `gorm:"NOT NULL;"`
}

//AddProblem 創建題目
func AddProblem(problem *Problem) (err error) {
	err = DB.Create(&problem).Error
	return
}

//UpdateProblem 更新題目
func UpdateProblem(problem *Problem) (err error) {
	err = DB.Save(&problem).Error
	return
}

//ListProblem 列出所有題目
func ListProblem() (problems []Problem, err error) {
	err = DB.Find(&problems).Error
	return
}

//GetProblemByID 查詢題目用 problem id
func GetProblemByID(id uint) (problem Problem, err error) {
	err = DB.First(&problem, id).Error
	return
}
