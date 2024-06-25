package common

type AbsoluteAddress struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type RecommendServerResponse struct {
	ContentIds []string `json:"contentsIds"`
}

type DigitalTwinServerRequest struct {
	ContentIds []string `json:"contentIds"`
}

type Rotation struct {
	Row   float64 `json:"row"`
	Pitch float64 `json:"pitch"`
	Yaw   float64 `json:"yaw"`
}

type ResponseHtml2d struct {
	ContentId string   `json:"contentId"`
	Location  AbsoluteAddress `json:"location"`
	Rotation  Rotation `json:"rotation"`
	TextType  string   `json:"textType"`
	TextURL   string   `json:"textUrl"`
	StyleURL  string   `json:"styleUrl"`
}

type ResponseModel3d struct {
	ContentId    string   `json:"contentId"`
	Location     AbsoluteAddress `json:"location"`
	Rotation     Rotation `json:"rotation"`
	PresignedURL string   `json:"presignedUrl"`
}

type DigitalTwinServerResponse struct {
	ResponseHtml2ds  []ResponseHtml2d  `json:"html2d"`
	ResponseModel3ds []ResponseModel3d `json:"model3d"`
}

// クライアントへのレスポンス
type Response struct {
	AbsoluteAddress           AbsoluteAddress           `json:"absoluteAddress"`
	DigitalTwinServerResponse DigitalTwinServerResponse `json:"contents"`
}
