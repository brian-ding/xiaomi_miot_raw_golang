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
	cloud.GetLoginInfo()
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

func TestLogin(t *testing.T) {
	// arrange
	username := ""
	password := ""
	cloud := NewCloud(username, password)
	info, _ := cloud.GetLoginInfo()

	// act
	got, _ := cloud.Login(info)

	// assert
	if len(got.Location) <= 0 {
		t.Errorf("got %s, wanted %s", got.Location, "a url")
	}
}
