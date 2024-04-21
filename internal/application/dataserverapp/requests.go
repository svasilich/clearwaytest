package dataserverapp

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
