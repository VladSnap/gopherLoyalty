package api

type EmptyBody struct{}

type KeyContext string

type HttpStatusResponse struct {
	httpStatus          int
	excpectedHttpStatus []int
}

func (r *HttpStatusResponse) WithHttpStatus(httpStatus int) *HttpStatusResponse {
	r.httpStatus = httpStatus
	return r
}

func (r *HttpStatusResponse) WithExcpectedHttpStatus(excpectedHttpStatus []int) *HttpStatusResponse {
	r.excpectedHttpStatus = excpectedHttpStatus
	return r
}

func (r *HttpStatusResponse) HTTPStatus() int {
	return r.httpStatus
}

func (r *HttpStatusResponse) ExpectedHTTPStatuses() []int {
	return r.excpectedHttpStatus
}
