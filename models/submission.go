package models

import (
	"strconv"

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
	Score      string `gorm:"type:char(5) NOT NULL"`
}

// SubTask 子任務
type SubTask struct {
	gorm.Model
	SubmissionID uint   `gorm:"NOT NULL;"`
	TestCase     string `gorm:"type:text;NOT NULL;"`
	CPUTime      uint   `gorm:"NOT NULL"`
	Memory       uint   `gorm:"NOT NULL"`
	Status       int    `gorm:"NOT NULL"`
}

// Wrong style wrong
type Wrong struct {
	gorm.Model
	SubmissionID uint   `gorm:"NOT NULL;"`
	Line         uint   `gorm:"NOT NULL;"`
	Col          uint   `gorm:"NOT NULL;"`
	Description  string `gorm:"type:text;NOT NULL"`
	Rule         string `gorm:"type:text;NOT NULL"`
}

// SubmissionResult 提交結果 來自 judge service
type SubmissionResult struct {
	CompileError int             `json:"compile_error"`
	SubResults   []subTaskResult `json:"results"`
}

// StyleResult 樣式結果 來自 style service
type StyleResult struct {
	Score        string                 `json:"score"`
	WrongResults []wrongResultsTemplate `json:"wrong"`
}

// wrongResultsTemplate 錯誤結果 來自 style service
type wrongResultsTemplate struct {
	Line        string `json:"line"`
	Col         string `json:"col"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

type subTaskResult struct {
	CPUTime  uint   `json:"real_time"`
	Memory   uint   `json:"memory"`
	Status   int    `json:"result"`
	TestCase string `json:"test_case"`
}

// SubmissionStatus 提交狀態
type SubmissionStatus struct {
	SubmissionID uint   `json:"submission_id"`
	ProblemID    uint   `json:"problem_id"`
	Author       uint   `json:"author"`
	Language     string `json:"language"`
	SourceCode   string `json:"source_code"`
	Status       int    `json:"status"`
	CPUTime      uint   `json:"cpu_time"`
	Memory       uint   `json:"memory"`
	Score        string `json:"score"`
	Wrong        []wrongResultsTemplate
	TestCase     []subTaskResult
}

// GetSubmissionByID 獲取提交狀態
func GetSubmissionByID(id uint) (status SubmissionStatus, err error) {
	var submission Submission
	var subTasks []SubTask
	var wrongs []Wrong

	if err = DB.First(&submission, id).Error; err != nil {
		return
	}
	if err = DB.Where(&SubTask{SubmissionID: id}).Find(&subTasks).Error; err != nil {
		return
	}
	if err = DB.Where(&Wrong{SubmissionID: id}).Find(&wrongs).Error; err != nil {
		return
	}

	status.SubmissionID = id
	status.ProblemID = submission.ProblemID
	status.Author = submission.Author
	status.Language = submission.Language
	status.SourceCode = submission.SourceCode
	status.Status = submission.Status
	status.CPUTime = submission.CPUTime
	status.Memory = submission.Memory
	status.Score = submission.Score
	for _, w := range wrongs {
		var wrong wrongResultsTemplate
		wrong.Line = strconv.Itoa(int(w.Line))
		wrong.Col = strconv.Itoa(int(w.Col))
		wrong.Rule = w.Rule
		wrong.Description = w.Description
		status.Wrong = append(status.Wrong, wrong)
	}
	for _, s := range subTasks {
		var subTask subTaskResult
		subTask.CPUTime = s.CPUTime
		subTask.Memory = s.Memory
		subTask.Status = s.Status
		subTask.TestCase = s.TestCase
		status.TestCase = append(status.TestCase, subTask)
	}

	return
}

//CreateSubmission 創建提交
func CreateSubmission(submission *Submission) (err error) {
	err = DB.Create(&submission).Error
	return
}

// UpdateSubmissionJudgeResult 更新提交 - judge service
func UpdateSubmissionJudgeResult(id uint, result *SubmissionResult) (lang, code string, status int, err error) {
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
			subTask.TestCase = v.TestCase
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

	lang = submission.Language
	code = submission.SourceCode
	status = submission.Status

	return
}

// UpdateSubmissionStyleResult 更新提交 - style service
func UpdateSubmissionStyleResult(id uint, result *StyleResult) (err error) {
	var submission Submission

	if err = DB.First(&submission, id).Error; err != nil {
		return
	}

	submission.Score = result.Score

	for _, v := range result.WrongResults {
		var wrong Wrong
		var tmp int

		if tmp, err = strconv.Atoi(v.Col); err != nil {
			return
		}
		wrong.Col = uint(tmp)

		if tmp, err = strconv.Atoi(v.Line); err != nil {
			return
		}
		wrong.Line = uint(tmp)

		wrong.Description = v.Description
		wrong.Rule = v.Rule
		wrong.SubmissionID = id

		if err = DB.Create(&wrong).Error; err != nil {
			return
		}
	}
	err = DB.Save(&submission).Error

	return
}
