package openrouter

type APIError struct {
	Code    string
	Message string
}

func (err APIError) Error() string {
	return err.Code + ": " + err.Message
}
