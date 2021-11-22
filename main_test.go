package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NCNUCodeOJ/BackendQuestionDatabase/models"
	router "github.com/NCNUCodeOJ/BackendQuestionDatabase/routers"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

// cspell:disable-next-line
var token = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6NDc5MTA4MjEyMywiaWQiOiI3MTI0MTMxNTQxOTcxMTA3ODYiLCJvcmlnX2lhdCI6MTYzNzQ4MjEyMywidXNlcm5hbWUiOiJ0ZXN0X3VzZXIifQ.pznOSok8X7qv6FSIihJnma_zEy70TerzOs0QDZOq_4RPYOKSEOOYTZ9-VLm2P9XRldS17-7QrLFwjjfXyCodtA"

func init() {
	gin.SetMode(gin.TestMode)
	models.Setup()
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
		"memory_limit":       123,
		"cpu_time":           123,
		"layer":              1,
		"sample_input":       ["123"],
		"sample_output":      ["456"],
		"tags_list":          ["簡單"]
	}`)
	r := router.SetupRouter()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/v1/problem", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
