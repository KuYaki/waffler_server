package controller

type UserSave struct {
	Parser struct {
		Type  string `json:"type"`
		Token string `json:"token"`
	} `json:"parser"`
	Locale string `json:"locale"`
}
