package client

var (
	InternalServerStatus   = "Ошибка сервера, повторите попытку позже"
	ForbiddenStatus        = "Не достаточно прав"
	ErrorUnknownStatusCode = "Unknown status code %v"
)

type ResponseError struct {
	StatusCode int
	Message    string
	Err        error
}

func (re ResponseError) Error() string {
	return re.Message
}

type httpError struct {
	Error string `json:"error"`
}
