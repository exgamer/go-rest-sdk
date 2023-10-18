package httpHelperStruct

type HttpResponse struct {
	Status     string
	Body       []byte
	StatusCode int
	Url        string
	Method     string
	Headers    map[string]string
}
