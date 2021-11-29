package models

import (
	"encoding/json"
	"fmt"

	"github.com/NCNUCodeOJ/BackendQuestionDatabase/pkg"
	"gorm.io/gorm"
)

// Submission Database - database
type Submission struct {
	gorm.Model
	ProblemID  uint   `gorm:"NOT NULL;"`
	Author     uint   `gorm:"NOT NULL;"`
	Language   string `gorm:"type:text;NOT NULL"`
	SourceCode string `gorm:"type:text;NOT NULL"`
	Status     int    `gorm:"NOT NULL"`
	CPUTime    uint   `gorm:"NOT NULL"`
	Memory     uint   `gorm:"NOT NULL"`
	Score      uint   `gorm:"NOT NULL"`
}

// SubTask 子任務
type SubTask struct {
	gorm.Model
	SubmissionID   uint `gorm:"NOT NULL;"`
	TestCaseNumber uint `gorm:"NOT NULL;"`
	CPUTime        uint `gorm:"NOT NULL"`
	Memory         uint `gorm:"NOT NULL"`
	Status         int  `gorm:"NOT NULL"`
}

// SubmissionResult 提交結果 來自 judge service
type SubmissionResult struct {
	CompileError int             `json:"compile_error"`
	SubResults   []SubTaskResult `json:"results"`
}

// SubTaskResult 子任務結果 來自 judge service
type SubTaskResult struct {
	CPUTime        uint `json:"real_time"`
	Memory         uint `json:"memory"`
	Status         int  `json:"result"`
	TestCaseNumber uint `json:"test_case"`
}

// UnmarshalJSON 將 json 轉換成 struct
func (tp *SubTaskResult) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Printf("Error while decoding %v\n", err)
		return err
	}
	tp.CPUTime = uint(v["real_time"].(float64))
	tp.Memory = uint(v["memory"].(float64))
	tp.Status = int(v["result"].(float64))
	tp.TestCaseNumber = uint(v["test_case"].(float64))
	return nil
}

//CreateSubmission 創建提交
func CreateSubmission(submission *Submission) (err error) {
	err = DB.Create(&submission).Error
	return
}

// UpdateSubmissionJudgeResult 更新提交
func UpdateSubmissionJudgeResult(id uint, result *SubmissionResult) (err error) {
	var submission Submission

	if err = DB.First(&submission, id).Error; err != nil {
		return
	}
	if result.CompileError == 1 {
		submission.Status = -2
	} else {
		for _, v := range result.SubResults {
			var subTask SubTask

			subTask.CPUTime = v.CPUTime
			subTask.Memory = v.Memory
			subTask.Status = v.Status
			subTask.TestCaseNumber = v.TestCaseNumber
			subTask.SubmissionID = id
			if err = DB.Create(&subTask).Error; err != nil {
				return
			}

			if v.Status != 0 && submission.Status != -1 {
				submission.Status = v.Status
			}
			submission.CPUTime = pkg.Max(submission.CPUTime, v.CPUTime)
			submission.Memory = pkg.Max(submission.Memory, v.Memory)
		}
	}
	err = DB.Save(&submission).Error
	return
}
