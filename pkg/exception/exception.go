package exception

type AppException struct {
	Code    int
	Error   error
	Context map[string]any
}

func NewAppException(code int, err error, context map[string]any) *AppException {
	return &AppException{code, err, context}
}
