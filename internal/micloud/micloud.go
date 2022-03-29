package micloud

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	u "net/url"
	"strings"
)

const (
	letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
)

type MiCloud struct {
	username             string
	password             string
	deviceId             string
	agentId              string
	sign                 string
	authenticateResponse authenticateResponse
	serviceToken         string
}

type Step1Model struct {
}

func NewMiCloud(username string, password string) *MiCloud {
	cloud := MiCloud{username: username, password: password}
	cloud.generateDeviceId()
	cloud.generateAgentId()

	return &cloud
}

func (cloud *MiCloud) LogIn() {
	cloud.getSign()
	cloud.authenticate()
	cloud.getServiceToken()
}

func (cloud *MiCloud) generateDeviceId() {
	candidates := strings.ToLower(letters)
	for i := 0; i < 6; i++ {
		index := rand.Intn(len(candidates))
		cloud.deviceId += string(candidates[index])
	}
}

func (cloud *MiCloud) generateAgentId() {
	candidates := letters
	for i := 0; i < 13; i++ {
		index := rand.Intn(len(candidates))
		cloud.agentId += string(candidates[index])
	}
}

// step 1
func (cloud *MiCloud) getSign() {
	url := "https://account.xiaomi.com/pass/serviceLogin?sid=xiaomiio&_json=true"
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header = make(http.Header)
	request.Header["User-Agent"] = []string{fmt.Sprintf("Android-7.1.1-1.0.0-ONEPLUS A3010-136-%s APP/xiaomi.smarthome APPV/62830", cloud.agentId)}
	request.AddCookie(&http.Cookie{Name: "sdkVersion", Value: "3.8.6"})
	request.AddCookie(&http.Cookie{Name: "deviceId", Value: cloud.deviceId})

	response, _ := http.DefaultClient.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body[11:]))
	responseModel := signResponse{}
	json.Unmarshal(body[11:], &responseModel)

	cloud.sign = responseModel.Sign
}

// step 2
func (cloud *MiCloud) authenticate() {
	url := "https://account.xiaomi.com/pass/serviceLoginAuth2?_json=true"

	form := u.Values{}
	form.Add("sid", "xiaomiio")
	form.Add("hash", hash(cloud.password))
	form.Add("callback", "https://sts.api.io.mi.com/sts")
	form.Add("qs", "%3Fsid%3Dxiaomiio%26_json%3Dtrue")
	form.Add("user", cloud.username)
	form.Add("_json", "true")
	form.Add("_sign", "0psXfr43eNI0IX6q9Suk3qWbRqU=")

	request, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(form.Encode()))
	// request.Header = make(http.Header)
	request.Header.Add("User-Agent", fmt.Sprintf("Android-7.1.1-1.0.0-ONEPLUS A3010-136-%s APP/xiaomi.smarthome APPV/62830", cloud.agentId))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.AddCookie(&http.Cookie{Name: "deviceId", Value: cloud.deviceId})
	request.AddCookie(&http.Cookie{Name: "pass_ua", Value: "web"})
	request.AddCookie(&http.Cookie{Name: "sdkVersion", Value: "3.8.6"})
	request.AddCookie(&http.Cookie{Name: "uLocale", Value: "zh_CN"})
	request.AddCookie(&http.Cookie{Name: "userId", Value: cloud.username})

	client := &http.Client{}
	response, _ := client.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body[11:]))
	responseModel := authenticateResponse{}
	json.Unmarshal(body[11:], &responseModel)

	if responseModel.Result == "ok" {
		cloud.authenticateResponse = responseModel
	}
}

func hash(password string) string {
	d := []byte(password)
	m := md5.New()
	m.Write(d)
	hash := hex.EncodeToString(m.Sum(nil))

	return strings.ToUpper(hash)
}

// step 3
func (cloud *MiCloud) getServiceToken() {
	url := cloud.authenticateResponse.Location

	request, _ := http.NewRequest(http.MethodGet, url, nil)

	client := &http.Client{}
	response, _ := client.Do(request)
	for _, cookie := range response.Cookies() {
		if cookie.Name == "serviceToken" {
			cloud.serviceToken = cookie.Value
			break
		}
	}
}

func (cloud *MiCloud) GetDevices() []miDevice {
	//url := "https://{server}.api.io.mi.com/app/home/device_list"
	proxyStr := "http://127.0.0.1:8888"
	proxyURL, _ := url.Parse(proxyStr)
	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	baseurl := "https://api.io.mi.com/app"
	url := "/home/device_list"
	data := "{\"getVirtualModel\":false,\"getHuamiDevices\":0}"

	nonce := getNonce()
	//nonce := "xvHcbWEVIOMBozql"

	signedNonce := signNonce(cloud.authenticateResponse.Security, nonce)
	//signedNonce := "ZQKAg0kcuEhJj8HAffWIzs1SUx2bCUSPqOG/PUngj5k="

	signature := getSignature(url, signedNonce, nonce, data)
	//signature := "VTBVTehbFi4RetxrRjtS5we3pQYYW5vqnEIZlAtskDE="
	// params = {
	// 	'data': '{"getVirtualModel":true,"getHuamiDevices":1,"get_split_device":false,"support_smart_home":true}'
	// }

	form := u.Values{}
	form.Add("signature", signature)
	form.Add("_nonce", nonce)
	form.Add("data", data)

	request, _ := http.NewRequest(http.MethodPost, baseurl+url, strings.NewReader(form.Encode()))
	request.Header = make(http.Header)
	request.Header["User-Agent"] = []string{fmt.Sprintf("Android-7.1.1-1.0.0-ONEPLUS A3010-136-%s APP/xiaomi.smarthome APPV/62830", cloud.agentId)}
	//request.Header["Accept-Encoding"] = []string{"identity"}
	request.Header["x-xiaomi-protocal-flag-cli"] = []string{"PROTOCAL-HTTP2"}
	request.Header["content-type"] = []string{"application/x-www-form-urlencoded"}
	//request.Header["MIOT-ENCRYPT-ALGORITHM"] = []string{"ENCRYPT-RC4"}
	request.AddCookie(&http.Cookie{Name: "userId", Value: cloud.username})
	//request.AddCookie(&http.Cookie{Name: "yetAnotherServiceToken", Value: cloud.serviceToken})
	request.AddCookie(&http.Cookie{Name: "serviceToken", Value: cloud.serviceToken})
	request.AddCookie(&http.Cookie{Name: "locale", Value: "zh_CN"})
	//request.AddCookie(&http.Cookie{Name: "timezone", Value: "GMT +08:00"})
	//request.AddCookie(&http.Cookie{Name: "is_daylight", Value: "0"})
	//request.AddCookie(&http.Cookie{Name: "dst_offset", Value: "0"})
	//request.AddCookie(&http.Cookie{Name: "channel", Value: "MI_APP_STORE"})

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}
	fmt.Println(client)
	response, _ := http.DefaultClient.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	responseModel := deviceListResponse{}
	json.Unmarshal(body, &responseModel)
	if responseModel.Code == 0 && responseModel.Message == "ok" {
		return responseModel.Result.List
	}

	return nil
}
