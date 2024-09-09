package common

type UserLocation struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Height float64 `json:"height"`
}

type Content struct {
	ContentId   string      `json:"contentId"`
	ContentType string      `json:"contentType"`
	Content     interface{} `json:"content"`
}

type ResponseClient struct {
	UserLocation UserLocation `json:"userLocation"`
}

type RequestLocationEstimationServerCurrentLocation struct {
	RawData string  `json:"raw_data"`
	AppId   string  `json:"appId"`
	UserId  string  `json:"userId"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type ResponseLocationEstimationServerCurrentLocation struct {
	UserLocation UserLocation `json:"userLocation"`
}

type RequestRecommendContentsServer struct {
	UserId       string       `json:"userId"`
	UserLocation UserLocation `json:"userLocation"`
}

type ResponseError struct {
	ErrorMessage string `json:"error"`
}
