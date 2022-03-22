package adding

type User struct {
	Name     string `json:"Name"`
	Surname  string `json:"Surname"`
	Email    string `json:"Email"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Admin    bool   `json:"Admin"`
}
