package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kajiLabTeam/mr-platform-relay-server/common"
	"github.com/kajiLabTeam/mr-platform-relay-server/service"
)

func UserContentsEstimation(c *gin.Context) {
	userId, err := common.AuthWithGetID(c.Request.Header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// マルチパートフォームを取得
	multipartForm, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 文字列の取得
	latStr := c.PostForm("lat")
	lonStr := c.PostForm("lon")

	// 絶対座標推定
	userLocation, err := service.UserLocationEstimation(multipartForm, latStr, lonStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// コンテンツ選定を要求
	if err = service.RequestRecommendContents(userId, userLocation); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, common.ResponseClient{UserLocation: userLocation})
}
