package ssw

type httpResponse struct {
	Success bool `json:"success"`
	Code    int  `json:"code"`
	Data    any  `json:"data"`
}
