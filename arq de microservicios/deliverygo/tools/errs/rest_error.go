package errs

// son var globales que representan errores comunes en la app
// Permiten reutilizar estos errores en toda la aplicación sin repetir código.
// Unauthorized el usuario no esta autorizado al recurso
var Unauthorized = NewRestError(401, "Unauthorized")

// NotFound cuando un registro no se encuentra en la db
var NotFound = NewRestError(404, "Document not found")

// AlreadyExist cuando no se puede ingresar un registro a la db
var AlreadyExist = NewRestError(400, "Already exist")

// Internal esta aplicación no sabe como manejar el error
var Internal = NewRestError(500, "Internal server error")

// RestError es una interfaz para definir errores custom, o personalizados
type RestError interface {
	Status() int
	Error() string
}

// NewRestError creates a new errCustom
func NewRestError(status int, message string) RestError {
	return &restError{
		status:  status,
		Message: message,
	}
}

// restError es un error personalizado para http
type restError struct {
	status  int //codigo http asociado al error
	Message string `json:"error"`
}


//metodos de resterror
//devuelve el memsaje de error como un string
func (e *restError) Error() string {
	return e.Message
}

// Status http status code
func (e *restError) Status() int {
	return e.status
}