package api

type EmptyBody struct{}

type KeyContext string

type HTTPStatusResponse struct {
	httpStatus          int
	expectedHTTPStatus []int
}

func (r *HTTPStatusResponse) WithHttpStatus(httpStatus int) *HTTPStatusResponse {
	r.httpStatus = httpStatus
	return r
}

func (r *HTTPStatusResponse) WithExcpectedHttpStatus(expectedHttpStatus []int) *HTTPStatusResponse {
	r.expectedHTTPStatus = expectedHttpStatus
	return r
}

func (r *HTTPStatusResponse) HTTPStatus() int {
	return r.httpStatus
}

func (r *HTTPStatusResponse) ExpectedHTTPStatuses() []int {
	return r.expectedHTTPStatus
}
