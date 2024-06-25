package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/kajiLabTeam/mr-platform-relay-server/common"
)

func DecideContents(absoluteAddress common.AbsoluteAddress) (common.Response, error) {
	// ユーザーIDと絶対座標を元にリクエストを送る
	RECOMMEND_CONTENTS_SERVER_URL := os.Getenv("RECOMMEND_CONTENTS_SERVER_URL")
	endPoint := RECOMMEND_CONTENTS_SERVER_URL + "/api/recommend/contents"

	jsonStr, err := json.Marshal(absoluteAddress)
	if err != nil {
		return common.Response{}, err
	}

	// POSTリクエストを送る
	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return common.Response{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	// レスポンスを受け取る
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return common.Response{}, err
	}
	defer resp.Body.Close()

	// レスポンスを取得
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	var contentsId common.RecommendServerResponse
	err = json.Unmarshal([]byte(newStr), &contentsId)
	if err != nil {
		return common.Response{}, err
	}

	// レスポンスで返ってきたコンテンツIDをデジタルツインポータルサーバーに送る
	DIGITAL_TWIN_SERVER_URL := os.Getenv("DIGITAL_TWIN_SERVER_URL")
	endPoint = DIGITAL_TWIN_SERVER_URL + "/api/space/user/get"

	digitalTwinServerRequest := contentsId.ContentIds

	jsonStr, err = json.Marshal(digitalTwinServerRequest)
	if err != nil {
		return common.Response{}, err
	}

	// POSTリクエストを送る
	req, err = http.NewRequest("POST", endPoint, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return common.Response{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	// レスポンスを受け取る
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return common.Response{}, err
	}

	defer resp.Body.Close()

	// レスポンスを取得
	buf = new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr = buf.String()

	var digitalTwinServerResponse common.DigitalTwinServerResponse
	err = json.Unmarshal([]byte(newStr), &digitalTwinServerResponse)
	if err != nil {
		return common.Response{}, err
	}

	// Response の構造体を作成
	userResponse := common.Response{
		AbsoluteAddress:           absoluteAddress,
		DigitalTwinServerResponse: digitalTwinServerResponse,
	}
	// デジタルツインポータルサーバーから返ってきたコンテンツをクライアントに渡す
	return userResponse, nil
}
