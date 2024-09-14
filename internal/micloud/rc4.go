package micloud

import (
	"bytes"
	"encoding/base64"
)

type Rc4 struct {
	idx, jdx int
	ksa      []byte
}

func NewRc4(key string) *Rc4 {
	pwd, _ := base64.StdEncoding.DecodeString(key)
	cnt := len(pwd)
	tempKsa := make([]byte, 256)
	for i := 0; i < 256; i++ {
		tempKsa[i] = byte(i)
	}

	j := 0
	for i := 0; i < 256; i++ {
		j = (j + int(tempKsa[i]) + int(pwd[i%cnt])) & 255
		temp := tempKsa[i]
		tempKsa[i] = tempKsa[j]
		tempKsa[j] = temp
	}

	return &Rc4{
		idx: 0,
		jdx: 0,
		ksa: tempKsa,
	}
}

func (r *Rc4) Crypt(data []byte) []byte {
	ksa := r.ksa
	i := r.idx
	j := r.jdx
	outList := bytes.Buffer{}
	for _, b := range data {
		i = (i + 1) & 255
		j = (j + int(ksa[i])) & 255
		temp := ksa[i]
		ksa[i] = ksa[j]
		ksa[j] = temp
		r := byte(int(b) ^ int(ksa[(int(ksa[i])+int(ksa[j]))&255]))
		outList.WriteByte(r)
	}

	r.idx = i
	r.jdx = j
	r.ksa = ksa
	return outList.Bytes()
}

func (r *Rc4) Init1024() *Rc4 {
	r.Crypt(make([]byte, 1024))
	return r
}
