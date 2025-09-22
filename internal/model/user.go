package model

type User struct {
	ID int `json:"user_id"`

	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`

	Balance  int `json:"balance"`  // in kopecks
	Withdraw int `json:"withdraw"` // in kopecks
}

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetUserBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
