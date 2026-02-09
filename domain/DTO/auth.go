package dto

type Passport struct {
	Token   string  `json:"token"`
	Payload PayLoad `json:"payload"`
}

type PayLoad struct {
	ID        uint   `json:"id"`
	Role      int    `json:"role"`
	HumanType string `json:"human_type"`
}
