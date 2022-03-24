package micloud

type signResponse struct {
	Sign string `json:"_sign"`
}

type authenticateResponse struct {
	Security string `json:"ssecurity"`
	Token    string `json:"passToken"`
	Result   string `json:"result"`
	UserId   string `json:"userId"`
	CUserId  string `json:"cUserId"`
	Location string `json:"location"`
	Code     int    `json:"code"`
}
