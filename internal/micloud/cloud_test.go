package micloud

import "testing"

func TestLogin(t *testing.T) {
	cloud := Cloud{}
	cloud.GetLoginInfo()
}
