package dto

type Passport struct {
	Token   string
	Payload PayLoad
}

type PayLoad struct {
	ID        uint
	Role      int
	HumanType string
}
