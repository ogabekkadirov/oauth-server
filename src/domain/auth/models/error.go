package models

type Error struct {
	Code      int
	Err       error
	Validator *ValidatorError
}

type ValidatorError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
	// Param     string
	// Condition string
}

type ErrBody struct {
	Error []ValidatorError `json:"error"`
}
