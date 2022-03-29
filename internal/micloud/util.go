package micloud

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"strings"
	"time"
)

// Time based nonce
func getNonce() string {
	b1 := make([]byte, 8)
	rand.Read(b1)

	b2 := make([]byte, 4)
	println(time.Now().Unix() / 60)
	binary.BigEndian.PutUint32(b2, uint32(time.Now().Unix()/60))

	b3 := append(b1, b2...)
	nonce := base64.StdEncoding.EncodeToString(b3)
	return nonce
}

// nonce signed with secret
func signNonce(secret string, nonce string) string {
	h := sha256.New()

	b1, _ := base64.StdEncoding.DecodeString(secret)
	h.Write(b1)

	b2, _ := base64.StdEncoding.DecodeString(nonce)
	h.Write(b2)

	result := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return result
}

// request signature based on url, signedNonce, nonce and data
func getSignature(url string, signedNonce string, nonce string, data string) string {
	sign := strings.Join([]string{url, signedNonce, nonce, "data=" + data}, "&")

	b1, _ := base64.StdEncoding.DecodeString(signedNonce)
	h := hmac.New(sha256.New, b1)

	b2 := []byte(sign)
	h.Write(b2)
	signature := h.Sum(nil)

	result := base64.StdEncoding.EncodeToString(signature)
	return result
}
