package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/kajiLabTeam/mr-platform-relay-server/common"
)

func UserLocationEstimation(multipartForm *multipart.Form, latStr string, lonStr string) (common.UserLocation, error) {
	locationServerUrl := os.Getenv("LOCATION_ESTIMATION_SERVER_URL")
	endPoint := locationServerUrl + "/api/estimation/absolute"

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// マルチパートフォームをコピー
	for _, headers := range multipartForm.File {
		for _, header := range headers {
			file, err := header.Open()
			if err != nil {
				return common.UserLocation{}, err
			}
			defer file.Close()

			part, err := writer.CreateFormFile("rawDataFile", header.Filename)
			if err != nil {
				return common.UserLocation{}, err
			}

			_, err = io.Copy(part, file)
			if err != nil {
				return common.UserLocation{}, err
			}
		}
	}

	err := writer.WriteField("lat", latStr)
	if err != nil {
		return common.UserLocation{}, err
	}
	err = writer.WriteField("lon", lonStr)
	if err != nil {
		return common.UserLocation{}, err
	}

	// writerを閉じて、マルチパートメッセージを終了
	err = writer.Close()
	if err != nil {
		return common.UserLocation{}, err
	}

	// リクエストを作成
	req, err := http.NewRequest("POST", endPoint, &b)
	if err != nil {
		return common.UserLocation{}, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return common.UserLocation{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var bodyBytes []byte
		response.Body.Read(bodyBytes)
		return common.UserLocation{}, fmt.Errorf("response status code is not 200. status code: %d\nresponse body: %s", response.StatusCode, string(bodyBytes))
	}

	// レスポンスを取得
	var userLocation common.UserLocation
	err = json.NewDecoder(response.Body).Decode(&userLocation)
	if err != nil {
		return common.UserLocation{}, err
	}

	return userLocation, nil
}
