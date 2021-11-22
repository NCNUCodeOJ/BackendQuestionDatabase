package models

import "gorm.io/gorm"

//Sample Database - database
type Sample struct {
	gorm.Model
	Input     string `gorm:"type:text;"`
	Output    string `gorm:"type:text;"`
	ProblemID uint   `gorm:"NOT NULL;"`
}

//AddSample 增加範例
func AddSample(sample *Sample) (err error) {
	err = DB.Create(&sample).Error
	return
}

//DeleteSample 用 problem id 直接有這個 problem id 的範例
func DeleteSample(id uint) (err error) {
	err = DB.Where(&Sample{ProblemID: id}).Delete(Sample{}).Error
	return
}

//GetAllProblemSamples 用 problem id 找 sample
func GetAllProblemSamples(ProblemID uint) (sample []Sample, err error) {
	err = DB.Where(&Sample{ProblemID: ProblemID}).Find(&sample).Error
	return
}
