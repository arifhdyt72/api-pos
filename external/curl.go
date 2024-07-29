package external

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type AuthJWT struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expiredAt"`
}

var Bearer AuthJWT

func GenerateBearerToken() {
	m3sapi := os.Getenv("M3S_URL_API")
	if Bearer.Token == "" || time.Now().After(Bearer.ExpiredAt) {
		body := make(map[string]interface{})
		body["teamCode"] = os.Getenv("TEAM_CODE")
		body["teamPassword"] = os.Getenv("TEAM_PASS")
		jsonString, _ := json.Marshal(body)
		resp, err := CreateHttpReq(m3sapi+"api/v3/auth-token", "POST", "", string(jsonString), "application/json")
		fmt.Println(string(resp))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var rsJson map[string]interface{}
		json.Unmarshal(resp, &rsJson)
		var bodyJson = rsJson["body"].(map[string]interface{})
		Bearer.Token = bodyJson["token"].(string)
		Bearer.ExpiredAt, _ = time.Parse(time.RFC3339, bodyJson["expiredAt"].(string))
	}
}

func DoGetHttpRequest(url string) ([]byte, error) {

	httpRequest, err := http.NewRequest(http.MethodGet, url, bytes.NewReader([]byte("")))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	token, _ := GenerateToken("")
	httpRequest.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	reqBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return reqBody, nil
}

func CreateHttpReq(url string, method string, token string, body string, contentType string) ([]byte, error) {

	httpRequest, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", contentType)
	if token != "" {
		httpRequest.Header.Set("Authorization", "Bearer "+token)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 && res.StatusCode > 299 {
		reqBody, _ := io.ReadAll(res.Body)
		fmt.Println(string(reqBody))
		return nil, errors.New(strconv.Itoa(res.StatusCode) + " : " + res.Status)
	}

	reqBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return reqBody, nil
}
