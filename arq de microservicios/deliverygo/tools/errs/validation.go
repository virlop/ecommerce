package errs

import (
	"encoding/json"
)

// Validation es una interfaz para definir errores custom
// Validation es un error de validaciones de parameteros o de campos
type Validation interface {
	Add(path string, message string) Validation
	Error() string //devuelve una repres del error como un json
}

func NewValidation() Validation {
	return &ValidationErr{
		Messages: []errField{}, //inicializa lista vacia
	}
}

type ValidationErr struct {
	Messages []errField `json:"messages"`
}

func (e *ValidationErr) Error() string {
	body, err := json.Marshal(e)
	if err != nil {
		return "ErrValidation invalid."
	}
	return string(body)
}

// Add agrega errores a un validation error
func (e *ValidationErr) Add(path string, message string) Validation {
	err := errField{
		Path:    path,
		Message: message,
	}
	e.Messages = append(e.Messages, err)
	return e
}

// errField define un campo inválido. path y mensaje de error
type errField struct {
	Path    string `json:"path"` //campo o parametro
	Message string `json:"message"`
}
