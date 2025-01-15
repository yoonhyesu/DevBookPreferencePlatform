package service

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
)

/*
req에 기본 대상 유형을 준수하는 JSON 인코딩 된 값을 포함하는 본문이 있는 JSON 콘텐츠 유형이 있을 것으로 예상.
대상을 채우거나 오류 반환
*/
func readRequestJSON(req *http.Request, target any) error {
	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	if mediaType != "application/json" {
		log.Printf("expect application/json Content-Type, got %s", mediaType)
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(target)
}

// v를 JSON으로 렌더링하고 이를 w에 응답으로 사용
func renderJSON(w http.ResponseWriter, v any) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
