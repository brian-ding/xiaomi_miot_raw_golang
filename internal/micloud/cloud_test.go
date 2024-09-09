package micloud

import "testing"

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
	want := "e4d7f1b4ed2e42d15898f4b27b019da4"
	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestLogin(t *testing.T) {
	// arrange
	username := "18652962260"
	password := "Brian1028@XM"
	cloud := NewCloud(username, password)
	info, _ := cloud.GetLoginInfo()

	// act
	got := cloud.Login(info)

	// assert
	want := ""
	if got != want {

	}
}
