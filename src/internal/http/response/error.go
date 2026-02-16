package response

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func (e *ValidationError) Add(field, message string) {
	if e.Fields == nil {
		e.Fields = make(map[string]string)
	}
	e.Fields[field] = message
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Fields) > 0
}

type NotFoundError struct {
	Resource string 
}

func (e NotFoundError) Error() string {
	if e.Resource == "" {
		return "resource not found"
	}
	return e.Resource + " not found"
}