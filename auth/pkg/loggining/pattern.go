package loggining

type User struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type UserData struct {
	ID        uint
	Name      string
	Surname   string
	Email     string
	Username  string
	PswdHash  []byte
	Moderator bool
	Admin     bool
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
