package micloud

import (
	"bytes"
	"crypto/md5"
	rand "crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	randV2 "math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Cloud struct {
	Username     string
	Password     string
	DeviceId     string
	ServiceToken string
	UserId       int
	Security     string
	httpClient   *http.Client
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
		DeviceId:   getDeviceId(),
		httpClient: client,
	}
}

type AuthInfo struct {
	SignToken string `json:"_sign"`
	CallBack  string `json:"callback"`
	QueryStr  string `json:"qs"`
	Sid       string `json:"sid"`
}

func (cloud *Cloud) GetAuthInfo() (AuthInfo, error) {
	var result AuthInfo
	deviceId := getDeviceId()
	domainUrl := "https://account.xiaomi.com/"
	domainUrlObj, _ := url.Parse(domainUrl)

	cookies := make([]*http.Cookie, 2)
	cookies[0] = &http.Cookie{Name: "sdkVersion", Value: "3.9"}
	cookies[1] = &http.Cookie{Name: "deviceId", Value: deviceId}

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

type AuthResult struct {
	Location string `json:"location"`
	Security string `json:"ssecurity"`
	UserId   int    `json:"userId"`
}

func (cloud *Cloud) Auth(info AuthInfo) (AuthResult, error) {
	var result AuthResult
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

	cloud.UserId = result.UserId
	cloud.Security = result.Security

	return result, nil
}

func (cloud *Cloud) Login(loginUrl string) error {
	response, err := cloud.httpClient.Get(loginUrl)
	if err != nil {
		return errors.New(err.Error())
	}

	for key, values := range response.Header {
		if key != "Set-Cookie" {
			continue
		}

		for _, value := range values {

			if strings.HasPrefix(value, "serviceToken") {
				cloud.ServiceToken = strings.Replace(strings.Split(value, ";")[0], "serviceToken=", "", 1)

				return nil
			}
		}
	}

	return errors.New("not found")
}

func hash(password string) string {
	hash := md5.Sum([]byte(password))
	hashStr := fmt.Sprintf("%x", hash)

	return strings.ToUpper(hashStr)
}

func getDeviceId() string {
	letters := "ABCDEF"
	count := 13
	resultList := make([]string, count)

	for i := 0; i < count; i++ {
		index := randV2.IntN(len(letters))
		letter := string(letters[index])
		resultList = append(resultList, letter)
	}

	result := strings.Join(resultList, "")

	return result
}

func getNonce() (string, error) {
	var result string
	firstByteArray := make([]byte, 8)
	_, err := rand.Read(firstByteArray)
	if err != nil {
		return result, errors.New(err.Error())
	}

	seconds := int(time.Now().Unix())
	secondByteArray := make([]byte, 4)
	for i := range secondByteArray {
		secondByteArray[i] = byte((seconds >> uint(8*i)) & 0xFF)
	}
	finalByteArray := append(firstByteArray[:], secondByteArray...)
	result = base64.StdEncoding.EncodeToString(finalByteArray)

	return result, nil
}

func signNonce(security string, nonce string) string {
	secretBytes, _ := base64.StdEncoding.DecodeString(security)
	nonceBytes, _ := base64.StdEncoding.DecodeString(nonce)
	finalBytes := append(secretBytes, nonceBytes...)

	h := sha256.New()
	h.Write(finalBytes)
	bytes2 := h.Sum(nil)
	finalResult := base64.StdEncoding.EncodeToString(bytes2)

	return finalResult
}

func sha1Sign(urlStr string, dat map[string]string, nonce string, method string) string {
	uri, _ := url.Parse(urlStr)
	path := uri.Path
	if strings.HasPrefix(path, "/app/") {
		path = path[4:]
	}

	arr := []string{strings.ToUpper(method), path}
	for key, value := range dat {
		arr = append(arr, key+"="+value)
	}

	arr = append(arr, nonce)
	sign := strings.Join(arr, "&")

	h := sha1.New()
	h.Write([]byte(sign))
	bytes := h.Sum(nil)
	result := base64.StdEncoding.EncodeToString(bytes)

	return result
}

func encryptData(key string, data string) string {

	bytes := NewRc4(key).Init1024().Crypt([]byte(data))
	result := base64.StdEncoding.EncodeToString(bytes)

	return result
}

type GetDeivceDto struct {
	GetVirtualModel  bool `json:"getVirtualModel"`
	GetHuamiDevices  int  `json:"getHuamiDevices"`
	GetSplitDevice   bool `json:"get_split_device"`
	SupportSmartHome bool `json:"support_smart_home"`
}

func getRC4Params(method string, urlStr string, dto GetDeivceDto, security string) map[string]string {
	result := make(map[string]string, 1)
	jsonBytes, _ := json.Marshal(dto)
	result["data"] = string(jsonBytes)

	nonce, _ := getNonce()
	signedNonce := signNonce(security, nonce)
	result["rc4_hash__"] = sha1Sign(urlStr, result, signedNonce, method)
	for k, v := range result {
		result[k] = encryptData(signedNonce, v)
	}
	result["signature"] = sha1Sign(urlStr, result, signedNonce, method)
	result["_nonce"] = nonce
	result["ssecurity"] = security
	result["signedNonce"] = signedNonce

	return result
}

func (cloud *Cloud) GetDeivces(method string, dto GetDeivceDto) {
	baseUrlStr := "https://api.io.mi.com/app/"
	urlStr := baseUrlStr + "home/device_list"
	param := getRC4Params(method, urlStr, dto, cloud.Security)

	form := url.Values{}
	for k, v := range param {
		form.Set(k, v)
	}
	encoded := form.Encode()
	payload := bytes.NewReader([]byte(encoded))

	request, _ := http.NewRequest(http.MethodPost, urlStr, payload)
	// request.Header.Add("Content-Length", strconv.Itoa(len(encoded)))
	request.Header.Add("User-Agnet", "APP/com.xiaomi.mihome APPV/6.0.103 iosPassportSDK/3.9.0 iOS/14.4 miHSTS")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("X-XIAOMI-PROTOCAL-FLAG-CLI", "PROTOCAL-HTTP2")
	request.Header.Add("Accept-Encoding", "identity")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("MIOT-ENCRYPT-ALGORITHM", "ENCRYPT-RC4")
	request.Header.Add("Connection", "keep-alive")

	cookies := make([]*http.Cookie, 10)
	cookies[0] = &http.Cookie{Name: "userId", Value: strconv.Itoa(cloud.UserId)}
	cookies[1] = &http.Cookie{Name: "serviceToken", Value: cloud.ServiceToken}
	cookies[2] = &http.Cookie{Name: "yetAnotherServiceToken", Value: cloud.ServiceToken}
	cookies[3] = &http.Cookie{Name: "is_daylight", Value: "0"}
	cookies[4] = &http.Cookie{Name: "channel", Value: "MI_APP_STORE"}
	cookies[5] = &http.Cookie{Name: "dst_offset", Value: "0"}
	cookies[6] = &http.Cookie{Name: "locale", Value: "zh_CN"}
	cookies[7] = &http.Cookie{Name: "timezone", Value: "GMT+8:00"}
	cookies[8] = &http.Cookie{Name: "sdkVersion", Value: "3.9"}
	cookies[9] = &http.Cookie{Name: "deviceId", Value: cloud.DeviceId}

	urlObj, _ := url.Parse(baseUrlStr)
	cloud.httpClient.Jar.SetCookies(urlObj, cookies)
	response, _ := cloud.httpClient.Do(request)
	body, _ := io.ReadAll(response.Body)
	bodyStr := string(body)

	baseData, _ := base64.StdEncoding.DecodeString(bodyStr)
	pBytes := NewRc4(param["signedNonce"]).Init1024().Crypt(baseData)
	pStr := string(pBytes)
	fmt.Println(pStr)

	// prefix := "&&&START&&&"
	// if !strings.HasPrefix(bodyStr, prefix) {
	// 	return result, errors.New("response prefix does not match")
	// }

	// bodyStr = strings.Replace(bodyStr, prefix, "", 1)
	// err = json.Unmarshal([]byte(bodyStr), &result)
	// if err != nil {
	// 	return result, errors.New(err.Error())
	// }
}
