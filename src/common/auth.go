package common

import "net/http"

func AuthWithGetID(header http.Header) (string, error) {
	// TODO
	// OAuth 2.0を使ってIDを取得する処理を記述する
	userId := header.Get("Authorization")
	return userId, nil
}
