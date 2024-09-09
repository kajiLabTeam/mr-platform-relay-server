package router

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kajiLabTeam/mr-platform-relay-server/controller"
)

func Init() {
	f, _ := os.Create("../log/server.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!!")
	})
	r.POST("/api/contents", controller.UserContentsEstimation)

	// サーバーの起動状態を表示
	if err := r.Run("0.0.0.0:8000"); err != nil {
		fmt.Println("サーバーの起動に失敗しました:", err)
	} else {
		fmt.Println("サーバーが正常に起動しました")
	}
}
