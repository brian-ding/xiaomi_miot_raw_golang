package micloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
)

type Cloud struct {
}

type LoginInfo struct {
	SignToken string `json:"_sign"`
	CallBack  string `json:"callback"`
	QueryStr  string `json:"qs"`
	Sid       string `json:"sid"`
}

func (cloud *Cloud) GetLoginInfo() (LoginInfo, error) {
	var result LoginInfo
	clientId := getClientId()
	domainUrl := "https://account.xiaomi.com/"
	cookie1 := http.Cookie{Name: "sdkVersion", Value: "3.9"}
	cookie2 := http.Cookie{Name: "deviceId", Value: clientId}

	requestUrl, _ := url.Parse(fmt.Sprintf("%s%s", domainUrl, "pass/serviceLogin?sid=xiaomiio&_json=true"))
	request := http.Request{Method: http.MethodGet, URL: requestUrl}

	request.Header = map[string][]string{
		"User-Agent": {"APP/com.xiaomi.mihome APPV/6.0.103 iosPassportSDK/3.9.0 iOS/14.4 miHSTS"},
	}
	request.AddCookie(&cookie1)
	request.AddCookie(&cookie2)

	response, err := http.DefaultClient.Do(&request)
	if err != nil {
		return result, errors.New(err.Error())
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, errors.New(err.Error())
	}
	bodyStr := string(body)

	prefix := "&&&START&&&"
	if !strings.HasPrefix(bodyStr, prefix) {
		return result, fmt.Errorf("response prefix does not match")
	}

	bodyStr = strings.Replace(bodyStr, prefix, "", 1)
	err = json.Unmarshal([]byte(bodyStr), &result)
	if err != nil {
		return result, errors.New(err.Error())
	}

	return result, nil
}

func getClientId() string {
	letters := "ABCDEF"
	count := 13
	resultList := make([]string, count)

	for i := 0; i < count; i++ {
		index := rand.IntN(len(letters))
		letter := string(letters[index])
		resultList = append(resultList, letter)
	}

	result := strings.Join(resultList, "")

	return result
}
