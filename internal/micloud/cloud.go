package micloud

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

	urlObj, _ := url.Parse(fmt.Sprintf("%s%s", domainUrl, "pass/serviceLogin?sid=xiaomiio&_json=true"))
	request := http.Request{Method: http.MethodGet, URL: urlObj}

	request.Header = map[string][]string{
		"User-Agent": {"APP/com.xiaomi.mihome APPV/6.0.103 iosPassportSDK/3.9.0 iOS/14.4 miHSTS"},
	}

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

type LoginRequestDto struct {
	LoginInfo
	Username       string `json:"user"`
	HashedPassword string `json:"hash"`
}

func (cloud *Cloud) Login(info LoginInfo) string {
	values := make(map[string][]string, 6)
	values["_sign"] = []string{info.SignToken}
	values["callback"] = []string{info.CallBack}
	values["qs"] = []string{info.QueryStr}
	values["sid"] = []string{info.Sid}
	values["user"] = []string{cloud.Username}
	values["hash"] = []string{hash(cloud.Password)}

	domainUrl := "https://account.xiaomi.com/"
	response, err := cloud.httpClient.PostForm(fmt.Sprintf("%s%s", domainUrl, "pass/serviceLoginAuth2?_json=true"), values)
	if err != nil {
		return err.Error()
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	bodyStr := string(body)
	return bodyStr
}

func hash(password string) string {
	h := md5.New()
	io.WriteString(h, password)
	hash := h.Sum(nil)

	return hex.EncodeToString(hash)
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
