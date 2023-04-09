package service

import (
	"bytes"
	"fmt"
	"gin-gorm-oj/define"
	"gin-gorm-oj/helper"
	"gin-gorm-oj/models"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query string false "status"
// @Success 200 {string} json "{"code", "data"}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("status"))
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	list := make([]*models.ProblemBasic, 0)
	var count int64
	err := tx.Count(&count).Offset((page - 1) * size).Limit(size).Find(&list).Error
	if err != nil {
		log.Println(err)
		return
	}
	//c.String(http.StatusOK, "%d", 2)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// Submit
// @Tags 用户私有方法
// @Summary 代码提交
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code", "data"}"
// @Router /user/submit [post]
func Submit(c *gin.Context) {
	problemIdentity := c.Query("problem_identity")
	code, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Read Code Error:" + err.Error(),
		})
		return
	}
	// 代码保存
	path, err := helper.CodeSave(code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Code Save Error:" + err.Error(),
		})
		return
	}
	u, _ := c.Get("user")
	userClaims := u.(*helper.UserClaims)
	sb := &models.SubmitBasic{
		Identity:        helper.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaims.Identity,
		Path:            path,
	}
	// 代码判断
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCase").Find(pb).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Problem Error:" + err.Error(),
		})
		return
	}
	WA := make(chan int)
	OOM := make(chan int)
	CE := make(chan int)
	passCount := 0
	var msg string
	var lock sync.Mutex
	for _, testCase := range pb.TestCase {
		test := testCase
		go func() {
			// go run code-user/main.go
			cmd := exec.Command("go", "run", path)
			var out, stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &out

			stdinPip, err := cmd.StdinPipe()
			if err != nil {
				fmt.Println(err)
			}
			io.WriteString(stdinPip, test.Input)
			var bm runtime.MemStats
			runtime.ReadMemStats(&bm)
			if err := cmd.Run(); err != nil {
				log.Println(err, stderr.String())
				if err.Error() == "exit status 2" {
					CE <- 1
					msg = stderr.String()
					return
				}
			}
			var em runtime.MemStats
			runtime.ReadMemStats(&em)
			if test.Output != out.String() {
				msg = "答案错误"
				WA <- 1
				return
			}
			if em.Alloc/1024-bm.Alloc/1024 > uint64(pb.MaxMem) {
				msg = "运行超内存"
				OOM <- 1
				return
			}
			lock.Lock()
			passCount++
			lock.Unlock()
		}()
	}
	select {
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		if passCount == len(pb.TestCase) {
			sb.Status = 1
		} else {
			sb.Status = 3
		}
	case <-CE:
		sb.Status = 5

	}
	err = models.DB.Create(sb).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Submit Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"status": sb.Status,
			"msg":    msg,
		},
	})
}
