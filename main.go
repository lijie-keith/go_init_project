package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lijie-keith/go_init_project/common"
	"github.com/lijie-keith/go_init_project/config"
	"github.com/lijie-keith/go_init_project/middleware"
	"github.com/lijie-keith/go_init_project/router"
	"log"
	"net/http"
	"runtime/debug"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.LoggerToFile(), gin.Recovery()) // 日志文件

	router.InitRouter(engine)

	return engine
}

// Recover
/**
 * global exception catch
 */
func Recover(c *gin.Context) {
	// 加载defer异常处理
	defer func() {
		if err := recover(); err != nil {
			// 异常日志
			log.Printf("出现异常: %v\n", err)

			// 打印错误堆栈信息
			debug.PrintStack()

			// 返回统一的Json风格
			c.JSON(http.StatusOK, common.Err)
			//终止后续操作
			c.Abort()
		}
	}()
	//继续操作
	c.Next()
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(config.APP_PORT)
	if err != nil {
		return
	}
}
