package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wuxiaobai24/gin-blog/config"
	"github.com/wuxiaobai24/gin-blog/models"
	"github.com/wuxiaobai24/gin-blog/router"
)

func main() {

	if err := config.InitConfig(); err != nil {
		panic(err)
		os.Exit(1)
	}

	if err := models.InitDB(); err != nil {
		panic(err)
		os.Exit(2)
	}
	defer models.CloseDB()

	mode := strings.ToUpper(config.Conf.Base.Mode)
	switch mode {
	case "DEBUG":
		gin.SetMode(gin.DebugMode)
	case "RELEASE":
		gin.SetMode(gin.ReleaseMode)
	default:
		fmt.Println("Mode is %v, which is not support. Please use `debug` or `release`.", mode)
	}

	url := config.Conf.Base.Url
	port := config.Conf.Base.Port
	addr := fmt.Sprintf("%s:%d", url, port)
	if err := router.Run(addr); err != nil {
		panic(err)
		os.Exit(3)
	}

}
