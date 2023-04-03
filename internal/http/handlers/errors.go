package handlers

type InvalidPayload struct {
	error
}

type ValidationError struct {
	error
}

type ResponseMarshalError struct {
	error
}
