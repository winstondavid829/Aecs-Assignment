package data

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type CustomerToken struct {
	Access_uuid string  `json:"access_uuid"`
	Authorized  bool    `json:"authorized"`
	Email       string  `json:"email"`
	Exp         float64 `json:"exp"`
	UserID      string  `json:"userID"`
	SessionID   string  `json:"sessionID"`
}

type AccessDetails struct {
	AccessUuid string
	UserId     string
}

type TokenResponse1 struct {
	AccessToken  string
	RefreshToken string
}
