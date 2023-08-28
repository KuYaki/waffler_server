package controller

type UserSave struct {
	Parser `json:"parser,omitempty"`
	Locale string `json:"locale,omitempty"`
}

type Parser struct {
	Type  string `json:"type,omitempty"`
	Token string `json:"token,omitempty"`
}
