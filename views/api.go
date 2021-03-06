package views

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/NCNUCodeOJ/BackendQuestionDatabase/judgeservice"
	"github.com/NCNUCodeOJ/BackendQuestionDatabase/models"
	"github.com/NCNUCodeOJ/BackendQuestionDatabase/styleservice"
	"github.com/gin-gonic/gin"
	"github.com/vincentinttsh/replace"
	"github.com/vincentinttsh/zero"
)

//CreateProblem 創建題目
func CreateProblem(c *gin.Context) {
	var problem models.Problem
	userID := c.MustGet("userID").(uint)
	data := problemAPIRequest{}
	if err := c.BindJSON(&data); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫或未使用json",
		})
		return
	}
	if zero.IsZero(data) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未填寫完成",
		})
		return
	}

	replace.Replace(&problem, &data)
	problem.Author = userID
	// check program name
	isValidProgramName := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	if !isValidProgramName(problem.ProgramName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ProgramName can only contain letters and numbers",
		})
		return
	}
	// check memorylimit
	if problem.MemoryLimit > 512*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Out of Max memorylimit",
		})
		return
	}
	if err := models.AddProblem(&problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "題目創建失敗",
		})
		return
	}

	for _, tag := range data.TagsList {
		if err := models.AddTag2Problem(problem.ID, tag); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "題目創建失敗",
			})
			return
		}
	}

	for i, sampleData := range data.Sample {
		var sample models.Sample

		sample.Input = sampleData.Input
		sample.Output = sampleData.Output
		sample.Sort = uint(i + 1)
		sample.ProblemID = problem.ID

		if err := models.AddSample(&sample); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "題目創建失敗",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "題目創建成功",
		"problem_id": problem.ID,
	})
}

// GetProblemByID 讀取題目
func GetProblemByID(c *gin.Context) {
	var problemID uint

	if ID, err := strconv.Atoi(c.Params.ByName("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	} else {
		problemID = uint(ID)
	}

	if problem, err := models.GetProblemByID(problemID); err == nil {
		var tags []string
		var err error
		var samples []models.SampleData

		if samples, err = models.GetProblemAllSamples(problemID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "題目讀取失敗",
			})
			return
		}

		if tags, err = models.GetProblemAllTags(problemID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "題目讀取失敗",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"problem_id":         problem.ID,
			"problem_name":       problem.ProblemName,
			"description":        problem.Description,
			"input_description":  problem.InputDescription,
			"output_description": problem.OutputDescription,
			"author":             problem.Author,
			"memory_limit":       problem.MemoryLimit,
			"cpu_time":           problem.CPUTime,
			"layer":              problem.Layer,
			"has_test_case":      problem.HasTestCase,
			"samples":            samples,
			"tags_list":          tags,
		})

		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "無此題目",
	})

	return
}

// GetSubmissionByID 讀取提交
func GetSubmissionByID(c *gin.Context) {
	var err error
	var submission models.SubmissionStatus
	var submissionID uint
	var wrong = make([]gin.H, 0)
	var testcase = make([]gin.H, 0)

	if ID, err := strconv.Atoi(c.Params.ByName("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	} else {
		submissionID = uint(ID)
	}

	if submission, err = models.GetSubmissionByID(submissionID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "無此提交",
		})
		return
	}

	for _, w := range submission.Wrong {
		wrong = append(wrong, gin.H{
			"line":        w.Line,
			"col":         w.Col,
			"rule":        w.Rule,
			"description": w.Description,
		})
	}

	for _, t := range submission.TestCase {
		testcase = append(testcase, gin.H{
			"cpu_time":  t.CPUTime,
			"memory":    t.Memory,
			"status":    t.Status,
			"test_case": t.TestCase,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"submission_id": submission.SubmissionID,
		"problem_id":    submission.ProblemID,
		"author":        submission.Author,
		"language":      submission.Language,
		"code":          submission.SourceCode,
		"status":        submission.Status,
		"cpu_time":      submission.CPUTime,
		"memory":        submission.Memory,
		"score":         submission.Score,
		"wrong":         wrong,
		"testcase":      testcase,
	})
}

// GetProblemsByTag 讀取屬於該 tag 的題目
func GetProblemsByTag(c *gin.Context) {
	var tagName string
	var problems []models.Problem
	var err error
	var returnData []map[string]interface{}

	if tagName = c.Params.ByName("tagName"); tagName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "tag name required",
		})
		return
	}

	if problems, err = models.GetProblemsByTag(tagName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查無此 tag",
		})
		return
	}

	for _, problem := range problems {
		returnData = append(returnData, map[string]interface{}{
			"problem_id":   problem.ID,
			"problem_name": problem.ProblemName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"problems": returnData,
	})
}

// EditProblem 編輯題目
func EditProblem(c *gin.Context) {
	var problemID uint
	var err error
	var problem models.Problem
	var oldTags []string
	var samplesLen int64

	if id, err := strconv.Atoi(c.Params.ByName("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	} else {
		problemID = uint(id)
	}

	data := problemAPIRequest{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫或未使用json",
			"err":     err.Error(),
		})
		return
	}

	if problem, err = models.GetProblemByID(problemID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "無此題目",
		})
		return
	}

	replace.Replace(&problem, &data)
	// check program name
	isValidProgramName := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	if !isValidProgramName(problem.ProgramName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ProgramName can only contain letters and numbers",
		})
		return
	}
	// check memorylimit
	if problem.MemoryLimit > 512*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Out of Max memorylimit",
		})
		return
	}
	models.UpdateProblem(&problem)

	if data.TagsList != nil {
		if oldTags, err = models.GetProblemAllTags(problemID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "題目編輯失敗-伺服器錯誤-get tag",
			})
			return
		}
		// 舊的沒有在最新的中，就刪除
		for _, tag := range oldTags {
			if !contains(data.TagsList, tag) {
				if err = models.DeleteProblemTag(problemID, tag); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "題目編輯失敗-伺服器錯誤-delete old tag",
					})
					return
				}
			}
		}
		// 新的沒有在舊的中，就新增
		for _, tag := range data.TagsList {
			if !contains(oldTags, tag) {
				if err = models.AddTag2Problem(problemID, tag); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "題目編輯失敗-伺服器錯誤-add new tag",
					})
					return
				}
			}
		}
	}

	if data.Sample != nil {
		if samplesLen, err = models.GetProblemSampleCount(problemID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "題目編輯失敗-伺服器錯誤-get problem sample count",
			})
			return
		}

		for i, sampleData := range data.Sample {
			if int64(i) < samplesLen {
				if err = models.UpdateSample(problemID, uint(i+1), sampleData.Input, sampleData.Output); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "題目創建失敗-伺服器錯誤-update sample",
						"err":     err.Error(),
					})
					return
				}
			} else {
				var sample models.Sample

				sample.Input = sampleData.Input
				sample.Output = sampleData.Output
				sample.Sort = uint(i + 1)
				sample.ProblemID = problem.ID

				if err = models.AddSample(&sample); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "題目創建失敗-伺服器錯誤-add new sample",
					})
					return
				}
			}
		}
		// 舊的沒有在最新的中，就刪除
		for i := len(data.Sample); int64(i) < samplesLen; i++ {
			if err = models.DeleteProblemSample(problemID, uint(i+1)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "題目編輯失敗-伺服器錯誤-delete old sample",
					"err":     err.Error(),
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "題目修改成功",
		"problem_id": problemID,
	})
}

// UploadProblemTestCase upload problem test case
func UploadProblemTestCase(c *gin.Context) {
	var problemID uint
	var problem models.Problem
	var err error
	var id int
	var file *multipart.FileHeader
	var dir, filePath string
	var files map[string]string

	if id, err = strconv.Atoi(c.Params.ByName("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}

	problemID = uint(id)
	if problem, err = models.GetProblemByID(problemID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "無此題目",
		})
		return
	}

	if file, err = c.FormFile("testcase"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "無檔案",
		})
		return
	}

	dir = filepath.Join(os.Getenv("TESTCASEDIR"), "tmp_"+strconv.Itoa(id))
	if err = os.Mkdir(dir, 0700); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
			"error":   err.Error(),
		})
		return
	}

	if needLog {
		log.Println("dir:", dir)
	} else {
		defer os.RemoveAll(dir)
	}

	filePath = filepath.Join(dir, "case.zip")

	if err = c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "上傳失敗",
		})
		return
	}

	if files, err = unZip(filePath, dir); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "解壓縮失敗" + err.Error(),
		})
		return
	}

	var start int = 1

	testCasePath := filepath.Join(os.Getenv("TESTCASEDIR"), strconv.Itoa(id))
	if _, err := os.Stat(testCasePath); err == nil {
		os.RemoveAll(testCasePath)
	}
	if err = os.Mkdir(testCasePath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
			"error":   err.Error(),
		})
		return
	}

	infoData := struct {
		TestCaseNumber int                         `json:"test_case_number"`
		Spj            bool                        `json:"spj"`
		TestCases      map[string]testCaseTemplate `json:"test_cases"`
	}{}
	infoData.Spj = false
	infoData.TestCaseNumber = 0
	infoData.TestCases = make(map[string]testCaseTemplate)

	for true {
		inData, inOK := files[strconv.Itoa(start)+".in"]
		outData, outOK := files[strconv.Itoa(start)+".out"]
		if !inOK || !outOK {
			break
		}
		infoData.TestCaseNumber++
		var testcaseInfo testCaseTemplate

		testcaseInfo.InputName = strconv.Itoa(start) + ".in"
		testcaseInfo.OutputName = strconv.Itoa(start) + ".out"
		testcaseInfo.InputSize = len(inData)
		testcaseInfo.OutputSize = len(outData)
		testcaseInfo.OutputMD5 = fmt.Sprintf("%x", md5.Sum([]byte(outData)))
		testcaseInfo.StrippedOutputMD5 = fmt.Sprintf("%x", md5.Sum([]byte(strings.TrimSpace(outData))))

		if needLog {
			log.Printf(
				"%s -> %s",
				filepath.Join(dir, strconv.Itoa(start)+".in"),
				filepath.Join(testCasePath, strconv.Itoa(start)+".in"),
			)
			log.Printf(
				"%s -> %s",
				filepath.Join(dir, strconv.Itoa(start)+".out"),
				filepath.Join(testCasePath, strconv.Itoa(start)+".out"),
			)
		}

		err = os.Rename(
			filepath.Join(dir, strconv.Itoa(start)+".in"),
			filepath.Join(testCasePath, strconv.Itoa(start)+".in"),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "server error",
			})
			if needLog {
				log.Println(err.Error())
			}
			return
		}
		err = os.Rename(
			filepath.Join(dir, strconv.Itoa(start)+".out"),
			filepath.Join(testCasePath, strconv.Itoa(start)+".out"),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "server error",
			})
			if needLog {
				log.Println(err.Error())
			}
			return
		}

		infoData.TestCases[strconv.Itoa(start)] = testcaseInfo

		start++
	}

	infoFile, _ := json.MarshalIndent(infoData, "", " ")
	infoFilePath := filepath.Join(testCasePath, "info")
	err = ioutil.WriteFile(infoFilePath, infoFile, 0644)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
			"error":   err.Error(),
		})
		return
	}

	problem.HasTestCase = true

	if err = models.UpdateProblem(&problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":          "上傳成功",
		"problem_id":       problemID,
		"test_case_number": infoData.TestCaseNumber,
	})
}

// GetSourceCodeAndAuthor get source code and author
func GetSourceCodeAndAuthor(c *gin.Context) {
	var data struct {
		SubmissionIDs []string `json:"submission_ids"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "json format error",
		})
		return
	}

	if zero.IsZero(data) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "data is not complete",
		})
		return
	}

	submissionIDs := make([]uint, len(data.SubmissionIDs))
	for i, id := range data.SubmissionIDs {
		tmp, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "submission id is invalid",
			})
		}

		submissionIDs[i] = uint(tmp)
	}

	submissions, err := models.GetSourceCodeAndAuthor(submissionIDs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	var responseData []gin.H

	url := userHost + privateBaseURL + "/username"
	var reqBody []byte
	reqBodyStruct := struct {
		UserID []string `json:"user_id"`
	}{}

	for _, submission := range submissions {
		reqBodyStruct.UserID = append(reqBodyStruct.UserID, strconv.Itoa(int(submission.Author)))
	}

	reqBody, err = json.Marshal(reqBodyStruct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	req.Header.Set("Content-Type", c.GetHeader("Content-Type"))

	userID2username := make(map[uint]string)

	if gin.Mode() != "test" {
		res, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		resData := struct {
			UserList []struct {
				UserID   uint   `json:"user_id"`
				UserName string `json:"username"`
			} `json:"user_list"`
		}{}

		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &resData)
		for _, user := range resData.UserList {
			userID2username[user.UserID] = user.UserName
		}
		if needLog {
			log.Println(userID2username)
		}
	}

	for _, submission := range submissions {
		username := userID2username[submission.Author]
		if username == "" {
			username = "Unknown_" + strconv.Itoa(int(submission.Author))
		}
		responseData = append(responseData, gin.H{
			"name":    username,
			"content": submission.SourceCode,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"submission_list": responseData,
	})
}

// CreateSubmission create problem submission
func CreateSubmission(c *gin.Context) {
	var problemID uint
	var err error
	var data submissionAPIRequest
	userID := c.MustGet("userID").(uint)

	var judgeTask judgeservice.JudgeTask
	var problem models.Problem
	var submission models.Submission

	if ID, err := strconv.Atoi(c.Params.ByName("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
	} else {
		problemID = uint(ID)
	}

	if problem, err = models.GetProblemByID(problemID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "無此題目",
		})
		return
	}

	if !problem.HasTestCase {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "problem has no test case",
		})
		return
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未按照格式填寫或未使用json",
		})
		return
	}
	if zero.IsZero(data) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未填寫完成",
		})
		return
	}

	replace.Replace(&judgeTask, &data)
	if err = judgeTask.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	submission.Author = userID
	submission.ProblemID = problem.ID
	submission.Language = *data.Language
	submission.SourceCode = *data.SourceCode

	if err = models.CreateSubmission(&submission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}

	judgeTask.ProblemID = problem.ID
	judgeTask.ProgramName = problem.ProgramName
	judgeTask.CPUTime = problem.CPUTime
	judgeTask.MemoryLimit = problem.MemoryLimit
	judgeTask.SubmissionID = submission.ID

	if err = judgeTask.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "提交成功",
		"submission_id": submission.ID,
	})
}

// UpdateSubmissionJudgeResult update submission judge result
func UpdateSubmissionJudgeResult(c *gin.Context) {
	var submissionID uint
	var language, code string
	var status int
	var err error
	var data models.SubmissionResult
	var styleTask styleservice.StyleTask

	if ID, err := strconv.Atoi(c.Params.ByName("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
			"error":   err.Error(),
		})
	} else {
		submissionID = uint(ID)
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "json error",
			"error":   err.Error(),
		})
		return
	}
	// fmt.Printf("%+v\n", data)
	if language, code, status, err = models.UpdateSubmissionJudgeResult(submissionID, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
			"error":   err.Error(),
		})
		return
	}

	styleTask.Language = language
	styleTask.SourceCode = code
	styleTask.SubmissionID = submissionID

	if err = styleTask.Validate(); err == nil && status == 0 {
		styleTask.Run()
	} else {
		var result models.StyleResult

		if status == 0 {
			result.Score = "10.00"
		} else {
			result.Score = "0.00"
		}

		models.UpdateSubmissionStyleResult(submissionID, &result)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update success",
	})
}

// UpdateSubmissionStyleResult update submission style result
func UpdateSubmissionStyleResult(c *gin.Context) {
	var submissionID uint
	var data models.StyleResult

	if ID, err := strconv.Atoi(c.Params.ByName("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
			"error":   err.Error(),
		})
	} else {
		submissionID = uint(ID)
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "json error",
			"error":   err.Error(),
		})
		return
	}
	// fmt.Printf("%+v\n", data)
	if err := models.UpdateSubmissionStyleResult(submissionID, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "system error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update success",
	})
}
