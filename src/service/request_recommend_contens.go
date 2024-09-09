package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/kajiLabTeam/mr-platform-relay-server/common"
)

func RequestRecommendContents(userId string, userLocation common.UserLocation) error {
	recommendContentsServerUrl := os.Getenv("RECOMMEND_CONTENTS_SERVER_URL")
	endPoint := recommendContentsServerUrl + "/api/content/recomend"

	// userIdをリクエストボディに設定
	requestBody := common.RequestRecommendContentsServer{
		UserId:       userId,
		UserLocation: userLocation,
	}
	requestBodyJsonStr, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// リクエストを作成
	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(requestBodyJsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 201 以外はエラー
	if response.StatusCode != http.StatusCreated {
		var responseError common.ResponseError
		if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
			return err
		}
		return fmt.Errorf("status code: %d, error message: %s", response.StatusCode, responseError.ErrorMessage)
	}
	return nil
}
