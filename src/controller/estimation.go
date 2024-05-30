package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kajiLabTeam/mr-platform-relay-server/common"
)

func GetEstimation(c *gin.Context) {
	// マルチパートフォームを取得
	multipartForm, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to get multipart form",
		})
		return
	}

	locationServerUrl := os.Getenv("LOCATION_ESTIMATION_SERVER_URL")
	endPoint := locationServerUrl + "/api/estimation/absolute"

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// マルチパートフォームをコピー
	for _, headers := range multipartForm.File {
		for _, header := range headers {
			file, err := header.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to open file",
				})
				return
			}
			defer file.Close()

			part, err := writer.CreateFormFile("rawDataFile", header.Filename)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to create form file",
				})
				return
			}

			_, err = io.Copy(part, file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to copy file",
				})
				return
			}
		}
	}

	latitude, err := strconv.ParseFloat(c.PostForm("latitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid latitude value",
		})
		return
	}

	longitude, err := strconv.ParseFloat(c.PostForm("longitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid longitude value",
		})
		return
	}

	_ = writer.WriteField("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	_ = writer.WriteField("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))

	err = writer.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to close writer",
		})
		return
	}

	// リクエストを作成
	req, err := http.NewRequest("POST", endPoint, &b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create request",
		})
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to send request",
		})
		return
	}
	defer response.Body.Close()

	// レスポンスを取得
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	newStr := buf.String()

	var absoluteAddress common.AbsoluteAddress
	err = json.Unmarshal([]byte(newStr), &absoluteAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to unmarshal json",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"x": absoluteAddress.X,
		"y": absoluteAddress.Y,
		"z": absoluteAddress.Z,
	})
}
