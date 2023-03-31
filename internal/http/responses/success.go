package responses

type SuccessResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}
