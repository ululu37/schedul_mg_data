package dto

type EvaluationResponse struct {
	Evaluation []EvaluationItem `json:"evaluation"`
}

type EvaluationItem struct {
	ID       int `json:"id"`
	Adtritud int `json:"adtritud"`
}
