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

type AccInfo struct {
	Coins     int `json:"coins" db:"balance"`
	Inventory []struct {
		Type     string `json:"type" db:"item"`
		Quantity int    `json:"quantity" db:"cnt"`
	} `json:"inventory"`
	CoinHistory struct {
		Received []struct {
			FromUser string `json:"fromUser" db:"sender"`
			Amount   int    `json:"amount" db:"amount"`
		} `json:"received"`
		Sent []struct {
			ToUser string `json:"toUser" db:"recipient"`
			Amount int    `json:"amount" db:"amount"`
		} `json:"sent"`
	} `json:"coinHistory"`
}
