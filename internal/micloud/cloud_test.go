package micloud

import (
	"strings"
	"testing"
)

func TestGetLoginInfo(t *testing.T) {
	// arrange
	username := "username"
	password := "password"
	cloud := NewCloud(username, password)

	// act
	cloud.GetAuthInfo()
}

func TestHash(t *testing.T) {
	// arrange
	raw := "hello, world"

	// act
	got := hash(raw)

	// assert
	want := strings.ToUpper("e4d7f1b4ed2e42d15898f4b27b019da4")
	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestAuth(t *testing.T) {
	// arrange
	username := ""
	password := ""
	cloud := NewCloud(username, password)
	info, _ := cloud.GetAuthInfo()

	// act
	got, _ := cloud.Auth(info)

	// assert
	if len(got.Location) <= 0 {
		t.Errorf("got %s, wanted %s", got.Location, "a url")
	}
}

func TestLogin(t *testing.T) {
	// arrange
	username := ""
	password := ""
	cloud := NewCloud(username, password)
	info, _ := cloud.GetAuthInfo()
	auth, _ := cloud.Auth(info)

	// act
	cloud.Login(auth.Location)

	// assert
}

func TestGetNonce(t *testing.T) {
	// arrange

	// act
	got, _ := getNonce()

	// assert
	if len(got) <= 0 {
		t.Errorf("got %s,", got)
	}
}

func TestGetDevices(t *testing.T) {
	// arrange
	username := ""
	password := ""
	cloud := NewCloud(username, password)
	info, _ := cloud.GetAuthInfo()
	auth, _ := cloud.Auth(info)
	cloud.Login(auth.Location)

	// act
	dto := GetDeivceDto{
		GetVirtualModel:  true,
		GetHuamiDevices:  1,
		GetSplitDevice:   false,
		SupportSmartHome: true,
	}
	cloud.GetDeivces("POST", dto)

	// assert
}
