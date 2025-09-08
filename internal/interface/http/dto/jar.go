package dto

type JWTRequestObject struct {
	Request    string `form:"request" json:"request"`
	RequestURI string `form:"request_uri" json:"request_uri"`
}
