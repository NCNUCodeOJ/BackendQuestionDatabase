package models

import "gorm.io/gorm"

//Sample Database - database
type Sample struct {
	gorm.Model
	Input     string `gorm:"type:text;"`
	Output    string `gorm:"type:text;"`
	Sort      uint   `gorm:"NOT NULL;"`
	ProblemID uint   `gorm:"NOT NULL;"`
}

// SampleData -- json response data structure
type SampleData struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Sort   uint   `json:"sort"`
}

// SampleInternalData -- sample internal data structure
type SampleInternalData struct {
	Input  string
	Output string
	Sort   uint
	ID     uint
}

//AddSample 增加範例
func AddSample(sample *Sample) (err error) {
	err = DB.Create(&sample).Error
	return
}

//DeleteProblemAllSamples 用 problem id 直接有這個 problem id 的範例
func DeleteProblemAllSamples(problemID uint) (err error) {
	err = DB.Where(&Sample{ProblemID: problemID}).Delete(Sample{}).Error
	return
}

//DeleteProblemSample 用 problem id 直接有這個 problem id 的範例
func DeleteProblemSample(problemID uint, sort uint) (err error) {
	err = DB.Where(&Sample{ProblemID: problemID, Sort: sort}).Delete(&Sample{}).Error
	return
}

// GetProblemSampleCount 用 problem id 找 sample 數量
func GetProblemSampleCount(problemID uint) (count int64, err error) {
	err = DB.Model(&Sample{}).Where(&Sample{ProblemID: problemID}).Count(&count).Error
	return
}

// UpdateSample 更新範例
func UpdateSample(problemID uint, sort uint, input string, output string) (err error) {
	var sample SampleData
	if err = DB.Model(&Sample{}).
		Where(&Sample{ProblemID: problemID, Sort: sort}).
		First(&sample).Error; err != nil {

		return
	}
	if sample.Input != input || sample.Output != output {
		err = DB.Model(&Sample{}).
			Where(&Sample{ProblemID: problemID, Sort: sort}).
			Updates(Sample{Input: input, Output: output}).Error
		return
	}
	return
}

//GetProblemAllSamples 用 problem id 找 sample
func GetProblemAllSamples(problemID uint) (sample []SampleData, err error) {
	err = DB.Model(&Sample{}).Where(&Sample{ProblemID: problemID}).Find(&sample).Error
	return
}

//GetProblemAllSamplesHaveSampleID 用 problem id 找 sample
func GetProblemAllSamplesHaveSampleID(ProblemID uint) (sample []SampleInternalData, err error) {
	err = DB.Model(&Sample{}).Where(&Sample{ProblemID: ProblemID}).Find(&sample).Error
	return
}
