package responses

type ErrorResponse struct {
	StatusCode  int    `json:"-"`
	Code        int    `json:"code"`
	Description string `json:"description"`
	Trace       string `json:"trace"`
}
