package structs

type SendCoinReqBody struct {
	To     string `json:"toUser"`
	Amount int    `json:"amount"`
}

type RegisterReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
