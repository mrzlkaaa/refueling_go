package loggining

type User struct {
	ID       uint   `json:"ID"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessToken struct {
	AccessUuid string
	UserId     uint
}
