package structs

type SendCoinInfo struct {
	From   string `json:"from"`
	To     string `json:"toUser"`
	Amount int    `json:"amount"`
}

type RegisterUserInfo struct {
	Login string `json:"login"`
	Pswd  string `json:"pswd"`
}

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	PswdHash string `db:"password_hash"`
}

type AuthUserInfo struct {
	Login string `json:"login"`
	Pswd  string `json:"pswd"`
}

type UserSecret struct {
	Login     string `db:"login"`
	Secret    string `db:"secret"`
	SessionID string `db:"session_id"`
}
