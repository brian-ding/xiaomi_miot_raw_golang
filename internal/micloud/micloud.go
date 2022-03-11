package micloud

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	u "net/url"
	"strings"
)

const (
	letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
)

type MiCloud struct {
	username string
	password string
	deviceId string
	agentId  string
	sign     string
}

func NewMiCloud(username string, password string) *MiCloud {
	cloud := MiCloud{username: username, password: password}
	cloud.generateDeviceId()
	cloud.generateAgentId()

	return &cloud
}

func (cloud *MiCloud) LogIn() {
	cloud.sign = cloud.getPayload()
	cloud.authenticate()
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

func (cloud *MiCloud) getPayload() string {
	url := "https://account.xiaomi.com/pass/serviceLogin?sid=xiaomiio&_json=true"
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header = make(http.Header)
	request.Header["User-Agent"] = []string{fmt.Sprintf("Android-7.1.1-1.0.0-ONEPLUS A3010-136-%s APP/xiaomi.smarthome APPV/62830", cloud.agentId)}
	request.AddCookie(&http.Cookie{Name: "sdkVersion", Value: "3.8.6"})
	request.AddCookie(&http.Cookie{Name: "deviceId", Value: cloud.deviceId})

	response, _ := http.DefaultClient.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	return string(body)
}

func (cloud *MiCloud) authenticate() string {
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
	fmt.Println(string(body))
	return string(body)
}

func hash(password string) string {
	d := []byte(password)
	m := md5.New()
	m.Write(d)
	hash := hex.EncodeToString(m.Sum(nil))

	return strings.ToUpper(hash)
}