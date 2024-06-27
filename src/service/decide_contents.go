package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/kajiLabTeam/mr-platform-relay-server/common"
)

func DecideContents(absoluteAddress common.AbsoluteAddress) (common.Response, error) {
	recommendServerResponse, err := getRecommendContents(absoluteAddress)
	if err != nil {
		return common.Response{}, err
	}

	// コンテンツが存在しない場合
	if len(recommendServerResponse.ContentIds) == 0 {
		userResponse := common.Response{
			AbsoluteAddress:           absoluteAddress,
			DigitalTwinServerResponse: common.DigitalTwinServerResponse{},
		}
		return userResponse, nil
	}

	// コンテンツが存在しない場合
	digitalTwinServerResponse, err := getContentFromDigitalTwinServer(recommendServerResponse)
	if err != nil {
		return common.Response{}, err
	}

	// コンテンツが存在しない場合
	if len(digitalTwinServerResponse.ResponseHtml2ds) == 0 && len(digitalTwinServerResponse.ResponseModel3ds) == 0 {
		userResponse := common.Response{
			AbsoluteAddress:           absoluteAddress,
			DigitalTwinServerResponse: common.DigitalTwinServerResponse{},
		}
		return userResponse, nil
	}

	// Response の構造体を作成
	userResponse := common.Response{
		AbsoluteAddress:           absoluteAddress,
		DigitalTwinServerResponse: digitalTwinServerResponse,
	}
	// デジタルツインポータルサーバーから返ってきたコンテンツをクライアントに渡す
	return userResponse, nil
}

func getRecommendContents(absoluteAddress common.AbsoluteAddress) (common.RecommendServerResponse, error) {
	// ユーザーIDと絶対座標を元にリクエストを送る
	RECOMMEND_CONTENTS_SERVER_URL := os.Getenv("RECOMMEND_CONTENTS_SERVER_URL")
	endPoint := RECOMMEND_CONTENTS_SERVER_URL + "/api/recommend/contents"

	jsonStr, err := json.Marshal(absoluteAddress)
	if err != nil {
		return common.RecommendServerResponse{}, err
	}

	// POSTリクエストを送る
	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return common.RecommendServerResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	// レスポンスを受け取る
	client := &http.Client{
		Timeout: time.Second * 30, // タイムアウトを設定
	}
	resp, err := client.Do(req)
	if err != nil {
		return common.RecommendServerResponse{}, err
	}
	defer resp.Body.Close()

	// レスポンスを取得
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	var contentIds common.RecommendServerResponse
	err = json.Unmarshal([]byte(newStr), &contentIds)
	if err != nil {
		return common.RecommendServerResponse{}, err
	}

	return contentIds, nil
}

func getContentFromDigitalTwinServer(recommendServerResponse common.RecommendServerResponse) (common.DigitalTwinServerResponse, error) {
	// レスポンスで返ってきたコンテンツIDをデジタルツインポータルサーバーに送る
	DIGITAL_TWIN_SERVER_URL := os.Getenv("DIGITAL_TWIN_SERVER_URL")
	endPoint := DIGITAL_TWIN_SERVER_URL + "/api/space/user/content/get"

	var digitalTwinServerRequest common.DigitalTwinServerRequest
	digitalTwinServerRequest.ContentIds = recommendServerResponse.ContentIds

	jsonStr, err := json.Marshal(digitalTwinServerRequest)
	if err != nil {
		return common.DigitalTwinServerResponse{}, err
	}

	// POSTリクエストを送る
	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return common.DigitalTwinServerResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	// レスポンスを受け取る
	client := &http.Client{
		Timeout: time.Second * 30, // タイムアウトを設定
	}
	resp, err := client.Do(req)
	if err != nil {
		return common.DigitalTwinServerResponse{}, err
	}

	defer resp.Body.Close()

	// 204の場合は空のレスポンスを返す
	if resp.StatusCode == http.StatusNoContent {
		return common.DigitalTwinServerResponse{}, nil
	}

	// レスポンスを取得
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	var digitalTwinServerResponse common.DigitalTwinServerResponse
	err = json.Unmarshal([]byte(newStr), &digitalTwinServerResponse)
	if err != nil {
		return common.DigitalTwinServerResponse{}, err
	}

	return digitalTwinServerResponse, nil
}