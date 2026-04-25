package models

type BogusRequest struct {
	Params    string `json:"params"`
	Body      string `json:"body"`
	UserAgent string `json:"userAgent"`
	CFP       int    `json:"cfp"`
}

type GnarlyRequest struct {
	Params    string `json:"params"`
	Body      string `json:"body"`
	UserAgent string `json:"userAgent"`
	Version   string `json:"version"`
	CFP       int    `json:"cfp"`
}

type DataRequest struct {
	Data string `json:"data"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}
