package views

import (
	"net/http"
	"strconv"

	"github.com/NCNUCodeOJ/BackendQuestionDatabase/models"
	"github.com/gin-gonic/gin"
	"github.com/vincentinttsh/replace"
	"github.com/vincentinttsh/zero"
)

// Pong test server is operating
func Pong(c *gin.Context) {
	if models.Ping() != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

//CreateProblem 創建題目
func CreateProblem(c *gin.Context) {
	var problem models.Problem
	userID := c.MustGet("userID").(uint)
	data := struct {
		ProblemName       *string   `json:"problem_name"`
		Description       *string   `json:"description"`
		InputDescription  *string   `json:"input_description"`
		OutputDescription *string   `json:"output_description"`
		MemoryLimit       *uint     `json:"memory_limit"`
		CPUTime           *uint     `json:"cpu_time"`
		Layer             *uint8    `json:"layer"`
		SampleInput       []*string `json:"sample_input"`
		SampleOutput      []*string `json:"sample_output"`
		TagsList          []*string `json:"tags_list"`
	}{}

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
	if len(data.SampleInput) != len(data.SampleOutput) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "sample input length != output length",
		})
		return
	}

	replace.Replace(&problem, &data)
	problem.Author = userID

	if err := models.AddProblem(&problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "題目創建失敗",
		})
		return
	}

	for _, tag := range data.TagsList {
		if err := models.AddTag2Problem(*tag, problem.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "題目創建失敗",
			})
			return
		}
	}

	for pos := range data.SampleInput {
		var sample models.Sample

		sample.Input = *(data.SampleInput[pos])
		sample.Output = *(data.SampleOutput[pos])
		sample.ProblemID = problem.ID

		if err := models.AddSample(&sample); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "題目創建失敗",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
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

	if problem, err := models.GetProblemByID(problemID); err != nil {
		var SampleIn, SampleOut, tags []string

		if samples, err := models.GetAllProblemSamples(problemID); err == nil {
			for _, sample := range samples {
				SampleIn = append(SampleIn, sample.Input)
				SampleOut = append(SampleOut, sample.Output)
			}
		}

		if TagList, err := models.GetProblemAllTags(problemID); err == nil {
			for _, tag := range TagList {
				tags = append(tags, tag.Name)
			}
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
			"sample_input":       SampleIn,
			"sample_output":      SampleOut,
			"tags_list":          tags,
		})

		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "無此題目",
	})

	return
}
