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

type deviceListResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Result  deviceList `json:"result"`
}

type deviceList struct {
	List []miDevice `json:"list"`
}

type miDevice struct {
	DID         string        `json:"did"`
	Token       string        `json:"token"`
	Longitude   string        `json:"longitude"`
	Latitude    string        `json:"latitude"`
	Name        string        `json:"name"`
	PID         string        `json:"pid"`
	LocalIP     string        `json:"localip"`
	MAC         string        `json:"mac"`
	SSID        string        `json:"ssid"`
	BSSID       string        `json:"bssid"`
	ParentID    string        `json:"parent_id"`
	ParentModel string        `json:"parent_model"`
	ShowMode    int           `json:"show_mode"`
	Model       string        `json:"model"`
	AdminFlag   int           `json:"adminFlag"`
	ShareFlag   int           `json:"shareFlag"`
	PermitLevel int           `json:"permitLevel"`
	IsOnline    bool          `json:"isOnline"`
	Description string        `json:"desc"`
	Extra       miDeviceExtra `json:"extra"`
	UID         int64         `json:"uid"`
	PDID        int           `json:"pd_id"`
	Password    string        `json:"password"`
	P2PID       string        `json:"p2p_id"`
	RSSI        int           `json:"rssi"`
	FamilyID    int           `json:"family_id"`
	ResetFlag   int           `json:"reset_flag"`
}

type miDeviceExtra struct {
	IsSetPincode      int    `json:"isSetPincode"`
	PincodeType       int    `json:"pincodeType"`
	FWVersion         string `json:"fw_version"`
	NeedVerifyCode    int    `json:"needVerifyCode"`
	IsPasswordEncrypt int    `json:"isPasswordEncrypt"`
}
