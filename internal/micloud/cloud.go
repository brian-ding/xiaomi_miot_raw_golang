package micloud

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

type Cloud struct {
	Username   string
	Password   string
	httpClient *http.Client
}

func NewCloud(username string, password string) *Cloud {
	options := &cookiejar.Options{}
	jar, err := cookiejar.New(options)
	if err != nil {
		log.Fatal(err.Error())
	}

	client := &http.Client{
		Jar: jar,
	}

	return &Cloud{
		Username:   username,
		Password:   password,
		httpClient: client,
	}
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
	domainUrlObj, _ := url.Parse(domainUrl)

	cookies := make([]*http.Cookie, 2)
	cookies[0] = &http.Cookie{Name: "sdkVersion", Value: "3.9"}
	cookies[1] = &http.Cookie{Name: "deviceId", Value: clientId}

	cloud.httpClient.Jar.SetCookies(domainUrlObj, cookies)

	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", domainUrl, "pass/serviceLogin?sid=xiaomiio&_json=true"), nil)
	request.Header.Add("User-Agnet", "APP/com.xiaomi.mihome APPV/6.0.103 iosPassportSDK/3.9.0 iOS/14.4 miHSTS")

	response, err := http.DefaultClient.Do(request)
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
		return result, errors.New("response prefix does not match")
	}

	bodyStr = strings.Replace(bodyStr, prefix, "", 1)
	err = json.Unmarshal([]byte(bodyStr), &result)
	if err != nil {
		return result, errors.New(err.Error())
	}

	return result, nil
}

type LoginResult struct {
	Location string `json:"location"`
}

func (cloud *Cloud) Login(info LoginInfo) (LoginResult, error) {
	var result LoginResult
	form := url.Values{}
	form.Set("_json", "true")
	form.Set("_sign", info.SignToken)
	form.Set("callback", info.CallBack)
	form.Set("qs", info.QueryStr)
	form.Set("sid", info.Sid)
	form.Set("user", cloud.Username)
	form.Set("hash", hash(cloud.Password))
	vEncoded := form.Encode()
	payload := bytes.NewReader([]byte(vEncoded))

	domainUrl := "https://account.xiaomi.com/"
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", domainUrl, "pass/serviceLoginAuth2"), payload)
	request.Header.Add("Content-Length", strconv.Itoa(len(vEncoded)))
	request.Header.Add("User-Agnet", "APP/com.xiaomi.mihome APPV/6.0.103 iosPassportSDK/3.9.0 iOS/14.4 miHSTS")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := cloud.httpClient.Do(request)
	if err != nil {
		println(err.Error())
		return result, errors.New(err.Error())
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, errors.New(err.Error())
	}

	bodyStr := string(body)

	prefix := "&&&START&&&"
	if !strings.HasPrefix(bodyStr, prefix) {
		return result, errors.New("response prefix does not match")
	}

	bodyStr = strings.Replace(bodyStr, prefix, "", 1)
	err = json.Unmarshal([]byte(bodyStr), &result)
	if err != nil {
		return result, errors.New(err.Error())
	}

	return result, nil
}

func hash(password string) string {
	hash := md5.Sum([]byte(password))
	hashStr := fmt.Sprintf("%x", hash)

	return strings.ToUpper(hashStr)
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
