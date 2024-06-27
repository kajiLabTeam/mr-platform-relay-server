package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kajiLabTeam/mr-platform-relay-server/common"
	"github.com/kajiLabTeam/mr-platform-relay-server/service"
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

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	lat := c.PostForm("latitude")
	lon := c.PostForm("longitude")

	err = copyMultipartFromData(multipartForm, writer, b ,lat,lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to copy multipart form",
		})
		return
	}

	locationServerUrl := os.Getenv("LOCATION_ESTIMATION_SERVER_URL")
	endPoint := locationServerUrl + "/api/estimation/absolute"

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

	userResponse, err := service.DecideContents(absoluteAddress)
	log.Println(userResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//userResponse.DigitalTwinServerResponse.ResponseHtml2dsとかuserResponse.DigitalTwinServerResponse.ResponseHtml2dsが空の場合、nullにするのではなく、空の配列を返す
	if len(userResponse.DigitalTwinServerResponse.ResponseHtml2ds) == 0 {
		userResponse.DigitalTwinServerResponse.ResponseHtml2ds = []common.ResponseHtml2d{}
	}
	if len(userResponse.DigitalTwinServerResponse.ResponseModel3ds) == 0 {
		userResponse.DigitalTwinServerResponse.ResponseModel3ds = []common.ResponseModel3d{}
	}
	c.JSON(http.StatusOK, userResponse)
}

func copyMultipartFromData(multipartForm *multipart.Form, writer *multipart.Writer, b bytes.Buffer,lat string, lon string) error {
	// マルチパートフォームをコピー
	for _, headers := range multipartForm.File {
		for _, header := range headers {
			file, err := header.Open()
			if err != nil {
				log.Println("failed to open file")
				return err
			}
			defer file.Close()

			part, err := writer.CreateFormFile("rawDataFile", header.Filename)
			if err != nil {
				log.Println("failed to create form file")
				return err
			}

			_, err = io.Copy(part, file)
			if err != nil {
				log.Println("failed to copy file")
				return err
			}
		}
	}

	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Println("failed to parse latitude")
		return err
	}

	longitude, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		log.Println("failed to parse longitude")
		return err
	}

	_ = writer.WriteField("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	_ = writer.WriteField("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))

	err = writer.Close()
	if err != nil {
		log.Println("failed to close writer")
		return err
	}
	return nil
}