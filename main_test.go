package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/NCNUCodeOJ/BackendQuestionDatabase/judgeservice"
	"github.com/NCNUCodeOJ/BackendQuestionDatabase/models"
	router "github.com/NCNUCodeOJ/BackendQuestionDatabase/routers"
	"github.com/NCNUCodeOJ/BackendQuestionDatabase/styleservice"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

// cspell:disable-next-line
var token = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjo0NzYwNjk2NDkyLCJpZCI6IjEiLCJvcmlnX2lhdCI6MTYzODYzMjQ5MiwidGVhY2hlciI6dHJ1ZSwidXNlcm5hbWUiOiJ2aW5jZW50In0.SUnwDQX_wkWlZdTMyCjhqIX4TIIzYrrY7lTiR_E2K8tvQBU1pyUgja60K0xcF1_x0m-egvRJQmhix5l6wdoR6g"
var problem1ID, submission1ID, submission2ID, submission3ID, submission4ID, submission5ID, submission6ID int

func init() {
	gin.SetMode(gin.TestMode)
	models.Setup()
	judgeservice.Setup()
	styleservice.Setup()
}

type Problem struct {
	ProblemID   int    `json:"problem_id"`
	ProblemName string `json:"problem_name"`
}

func (tp *Problem) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Printf("Error while decoding %v\n", err)
		return err
	}
	tp.ProblemID = int(v["problem_id"].(float64))
	tp.ProblemName = v["problem_name"].(string)
	return nil
}

func TestPing(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProblemCreate(t *testing.T) {
	var data = []byte(`{
		"problem_name":       "接龍遊戲",
		"description":        "開始接龍",
		"input_description":  "567",
		"output_description": "789",
		"memory_limit":       134217728,
		"cpu_time":           1000,
		"program_name":	      "Main",
		"layer":              1,
		"sample":             [
			{"input": "123", "output": "456"},
			{"input": "456", "output": "789"}
		],
		"tags_list":          ["簡單"]
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/private/v1/problem", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	s := struct {
		ProblemID int    `json:"problem_id"`
		Message   string `json:"message"`
	}{}
	json.Unmarshal(body, &s)
	problem1ID = s.ProblemID
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetProblemByID(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("GET", "/api/private/v1/problem/"+strconv.Itoa(problem1ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateMultipleProblem(t *testing.T) {
	var data = []byte(`{
		"problem_name":       "接龍遊戲",
		"description":        "開始接龍",
		"input_description":  "567",
		"output_description": "789",
		"memory_limit":       134217728,
		"cpu_time":           1000,
		"layer":              1,
		"program_name":	      "Main",
		"sample":             [
			{"input": "123", "output": "456"},
			{"input": "123", "output": "456"}
		],
		"tags_list":          ["簡單"]
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/private/v1/problem", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	data = []byte(`{
		"problem_name":       "接龍遊戲",
		"description":        "開始接龍",
		"input_description":  "567",
		"output_description": "789",
		"memory_limit":       134217728,
		"cpu_time":           1000,
		"layer":              1,
		"program_name":	      "Main",
		"sample":             [
			{"input": "123", "output": "456"},
			{"input": "123", "output": "456"}
		],
		"tags_list":          ["難"]
	}`)
	w = httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ = http.NewRequest("POST", "/api/private/v1/problem", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUpdateProblem(t *testing.T) {
	var data = []byte(`{
		"problem_name":       "龍遊戲",
		"sample":             [
			{"input": "456", "output": "789"},
			{"input": "123", "output": "456"},
			{"input": "789", "output": "123"}
		],
		"tags_list":          ["難"]
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("PATCH", "/api/private/v1/problem/"+strconv.Itoa(problem1ID), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	data = []byte(`{
		"problem_name":       "龍遊戲",
		"sample":             [
			{"input": "789", "output": "123"}
		]
	}`)
	w = httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ = http.NewRequest("PATCH", "/api/private/v1/problem/"+strconv.Itoa(problem1ID), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUploadProblemFile(t *testing.T) {
	if os.Getenv("gitlab") == "1" {
		assert.Equal(t, os.Getenv("gitlab"), "1")
		var problem models.Problem
		var err error
		if problem, err = models.GetProblemByID(uint(problem1ID)); err != nil {
			assert.Equal(t, err, nil)
		}
		problem.HasTestCase = true
		if err = models.UpdateProblem(&problem); err != nil {
			assert.Equal(t, err, nil)
		}
		return
	}
	r := gofight.New()
	r.POST("/api/private/v1/problem/"+strconv.Itoa(problem1ID)+"/testcase").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetFileFromPath([]gofight.UploadFile{
			{
				Path: "./test/testcase2.zip",
				Name: "testcase",
			}}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r = gofight.New()
	r.POST("/api/private/v1/problem/"+strconv.Itoa(problem1ID)+"/testcase").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetFileFromPath([]gofight.UploadFile{
			{
				Path: "./test/testcase.zip",
				Name: "testcase",
			}}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestCreateSubmission(t *testing.T) {
	r := gofight.New()
	r.POST("/api/private/v1/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission1ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r.POST("/api/private/v1/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission2ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r.POST("/api/private/v1/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission3ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r.POST("/api/private/v1/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission4ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r.POST("/api/private/v1/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission5ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	r.POST("/api/private/v1/problem/"+strconv.Itoa(problem1ID)+"/submission").
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		SetJSON(gofight.D{
			"source_code": "a, b = map(int,input().split())\nprint(a+b)",
			"language":    "python3",
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())

			id, _ := (jsonparser.GetInt(data, "submission_id"))
			submission6ID = int(id)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestUpdateSubmissionJudgeResult(t *testing.T) {
	r := gofight.New()
	var results []gofight.D

	r.PATCH("/api/private/v1/submission/"+strconv.Itoa(submission1ID)+"/judge").
		SetJSON(gofight.D{
			"compile_error": 1,
			"results":       results,
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	results = append(results, gofight.D{
		"real_time": 19,
		"memory":    8826880,
		"result":    0,
		"test_case": "1",
	})
	results = append(results, gofight.D{
		"cpu_time":  9,
		"real_time": 21,
		"memory":    8835072,
		"result":    0,
		"test_case": "2",
	})

	r.PATCH("/api/private/v1/submission/"+strconv.Itoa(submission2ID)+"/judge").
		SetJSON(gofight.D{
			"compile_error": 0,
			"results":       results,
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	r.PATCH("/api/private/v1/submission/"+strconv.Itoa(submission5ID)+"/judge").
		SetJSON(gofight.D{
			"compile_error": 0,
			"results":       results,
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	r.PATCH("/api/private/v1/submission/"+strconv.Itoa(submission3ID)+"/judge").
		SetJSON(gofight.D{
			"compile_error": 0,
			"results":       results,
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	results = append(results, gofight.D{
		"cpu_time":  9,
		"real_time": 21,
		"memory":    8835072,
		"result":    -1,
		"test_case": "2",
	})

	r.PATCH("/api/private/v1/submission/"+strconv.Itoa(submission4ID)+"/judge").
		SetJSON(gofight.D{
			"compile_error": 0,
			"results":       results,
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestUpdateSubmissionStyleResult(t *testing.T) {
	r := gofight.New()
	var results []gofight.D

	r.PATCH("/api/private/v1/submission/"+strconv.Itoa(submission2ID)+"/style").
		SetJSON(gofight.D{
			"score": "10.00",
			"wrong": results,
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	results = append(results, gofight.D{
		"line":        "2",
		"col":         "0",
		"rule":        "C0304",
		"description": "https://vald-phoenix.github.io/pylint-errors/plerr/errors/format/C0304",
	})
	results = append(results, gofight.D{
		"line":        "1",
		"col":         "0",
		"rule":        "C0304",
		"description": "https://vald-phoenix.github.io/pylint-errors/plerr/errors/format/C0304",
	})

	r.PATCH("/api/private/v1/submission/"+strconv.Itoa(submission3ID)+"/style").
		SetJSON(gofight.D{
			"score": "5.12",
			"wrong": results,
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGetSubmission1(t *testing.T) {
	r := gofight.New()

	r.GET("/api/private/v1/submission/"+strconv.Itoa(submission1ID)).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			data := []byte(r.Body.String())

			status, _ := jsonparser.GetInt(data, "status")
			assert.Equal(t, -2, int(status))
			score, _ := jsonparser.GetString(data, "score")
			assert.Equal(t, "0.00", score)

			length := 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "testcase")
			assert.Equal(t, 0, length)

			length = 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "wrong")
			assert.Equal(t, 0, length)
		})
}

func TestGetSubmission2(t *testing.T) {
	r := gofight.New()

	r.GET("/api/private/v1/submission/"+strconv.Itoa(submission2ID)).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			data := []byte(r.Body.String())

			status, _ := jsonparser.GetInt(data, "status")
			assert.Equal(t, 0, int(status))
			score, _ := jsonparser.GetString(data, "score")
			assert.Equal(t, "10.00", score)

			length := 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "testcase")
			assert.Equal(t, 2, length)

			length = 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "wrong")
			assert.Equal(t, 0, length)
		})
}

func TestGetSubmission3(t *testing.T) {
	r := gofight.New()

	r.GET("/api/private/v1/submission/"+strconv.Itoa(submission3ID)).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			data := []byte(r.Body.String())

			status, _ := jsonparser.GetInt(data, "status")
			assert.Equal(t, 0, int(status))
			score, _ := jsonparser.GetString(data, "score")
			assert.Equal(t, "5.12", score)

			length := 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "testcase")
			assert.Equal(t, 2, length)

			length = 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "wrong")
			assert.Equal(t, 2, length)
		})
}

func TestGetSubmission4(t *testing.T) {
	r := gofight.New()

	r.GET("/api/private/v1/submission/"+strconv.Itoa(submission4ID)).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			data := []byte(r.Body.String())

			status, _ := jsonparser.GetInt(data, "status")
			assert.Equal(t, -1, int(status))
			score, _ := jsonparser.GetString(data, "score")
			assert.Equal(t, "0.00", score)

			length := 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "testcase")
			assert.Equal(t, 3, length)

			length = 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "wrong")
			assert.Equal(t, 0, length)
		})
}
func TestGetSubmission5(t *testing.T) {
	r := gofight.New()

	r.GET("/api/private/v1/submission/"+strconv.Itoa(submission5ID)).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			data := []byte(r.Body.String())

			status, _ := jsonparser.GetInt(data, "status")
			assert.Equal(t, 0, int(status))
			score, _ := jsonparser.GetString(data, "score")
			assert.Equal(t, "", score)

			length := 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "testcase")
			assert.Equal(t, 2, length)

			length = 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "wrong")
			assert.Equal(t, 0, length)
		})
}
func TestGetSubmission6(t *testing.T) {
	r := gofight.New()

	r.GET("/api/private/v1/submission/"+strconv.Itoa(submission6ID)).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			data := []byte(r.Body.String())

			status, _ := jsonparser.GetInt(data, "status")
			assert.Equal(t, 0, int(status))
			score, _ := jsonparser.GetString(data, "score")
			assert.Equal(t, "", score)

			length := 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "testcase")
			assert.Equal(t, 0, length)

			length = 0
			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "wrong")
			assert.Equal(t, 0, length)
		})
}

func TestGetSubmissionNotFound(t *testing.T) {
	r := gofight.New()

	r.GET("/api/private/v1/submission/"+strconv.Itoa(100)).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
		})
}

func TestGetSubmissionCode(t *testing.T) {
	r := gofight.New()

	r.POST("/api/private/v1/submission/code").
		SetJSON(gofight.D{
			"submission_ids": []string{
				strconv.Itoa(submission1ID),
				strconv.Itoa(submission2ID),
				strconv.Itoa(submission3ID),
			},
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			length := 0

			jsonparser.ArrayEach(data,
				func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					length++
				}, "submission_list")

			assert.Equal(t, 3, length)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
func TestCleanup(t *testing.T) {
	e := os.Remove("test.db")
	if e != nil {
		t.Fail()
	}
}
