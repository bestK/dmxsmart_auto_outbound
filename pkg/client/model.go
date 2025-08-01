package client

type Response struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
	Data         any    `json:"data"`
	Total        int    `json:"total"`
}

// CaptchaResponse represents the response from the captcha API
type CaptchaResponse struct {
	Response
	Data struct {
		UUID string `json:"uuid"`
		Img  string `json:"img"`
	} `json:"data"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	Captcha     string  `json:"captcha"`
	UUID        string  `json:"uuid"`
	LoginType   string  `json:"loginType"`
	DeviceToken *string `json:"deviceToken"`
	Lang        string  `json:"lang"`
}

// LoginResponse 登录响应数据
type LoginResponse struct {
	Response
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}
