package api

type Credentials struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
