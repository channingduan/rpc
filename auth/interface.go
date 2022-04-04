package auth

type Detail struct {
}
type Token struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	AccessUuid    string `json:"access_uuid"`
	RefreshUuid   string `json:"refresh_uuid"`
	AccessExpire  int64  `json:"access_expire"`
	RefreshExpire int64  `json:"refresh_expire"`
}

type IAuth interface {
	CreateToken(id uint) (*Token, error)
	RefreshToken()
	ValidToken()
	RevokedToken()
}
