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
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

// cspell:disable-next-line
var token = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6NDc5MTA4MjEyMywiaWQiOiI3MTI0MTMxNTQxOTcxMTA3ODYiLCJvcmlnX2lhdCI6MTYzNzQ4MjEyMywidXNlcm5hbWUiOiJ0ZXN0X3VzZXIifQ.pznOSok8X7qv6FSIihJnma_zEy70TerzOs0QDZOq_4RPYOKSEOOYTZ9-VLm2P9XRldS17-7QrLFwjjfXyCodtA"
var problem1ID, submission1ID, submission2ID, submission3ID int

func init() {
	gin.SetMode(gin.TestMode)
	models.Setup()
	judgeservice.Setup()
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
	req, _ := http.NewRequest("POST", "/api/v1/problem", bytes.NewBuffer(data))
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
	req, _ := http.NewRequest("GET", "/api/v1/problem/"+strconv.Itoa(problem1ID), nil)
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
	req, _ := http.NewRequest("POST", "/api/v1/problem", bytes.NewBuffer(data))
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
	req, _ = http.NewRequest("POST", "/api/v1/problem", bytes.NewBuffer(data))
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
	req, _ := http.NewRequest("PATCH", "/api/v1/problem/"+strconv.Itoa(problem1ID), bytes.NewBuffer(data))
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
	req, _ = http.NewRequest("PATCH", "/api/v1/problem/"+strconv.Itoa(problem1ID), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUploadProblemFile(t *testing.T) {
	r := gofight.New()
	r.POST("/api/v1/problem/"+strconv.Itoa(problem1ID)+"/testcase").
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
	r.POST("/api/v1/problem/"+strconv.Itoa(problem1ID)+"/testcase").
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
	r.POST("/api/v1/problem/"+strconv.Itoa(problem1ID)+"/submission").
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
	r.POST("/api/v1/problem/"+strconv.Itoa(problem1ID)+"/submission").
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
	r.POST("/api/v1/problem/"+strconv.Itoa(problem1ID)+"/submission").
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
}

func TestUpdateSubmissionJudgeResult(t *testing.T) {
	r := gofight.New()
	var results []gofight.D

	r.PATCH("/api/v1/submission/"+strconv.Itoa(submission1ID)+"/judge").
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
		"test_case": 1,
	})
	results = append(results, gofight.D{
		"cpu_time":  9,
		"real_time": 21,
		"memory":    8835072,
		"result":    0,
		"test_case": 2,
	})

	r.PATCH("/api/v1/submission/"+strconv.Itoa(submission2ID)+"/judge").
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
		"test_case": 2,
	})

	r.PATCH("/api/v1/submission/"+strconv.Itoa(submission3ID)+"/judge").
		SetJSON(gofight.D{
			"compile_error": 0,
			"results":       results,
		}).
		Run(router.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestCleanup(t *testing.T) {
	e := os.Remove("test.db")
	if e != nil {
		t.Fail()
	}
}
