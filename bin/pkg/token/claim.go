package token

type ResToken struct {
	TokenType string  `json:"token_type"`
	Duration  float64 `json:"duration"`
	Token     string  `json:"access_token"`
}

type Claim struct {
	UserID        string `json:"userId"`
	RoleName	  string `json:"roleName"`
	Authorization string `json:"authorization"`
}
