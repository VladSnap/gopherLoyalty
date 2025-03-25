package api

type EmptyBody struct{}

type KeyContext string

type HTTPStatusResponse struct {
	httpStatus         int
	expectedHTTPStatus []int
}

func (r *HTTPStatusResponse) WithHTTPStatus(httpStatus int) *HTTPStatusResponse {
	r.httpStatus = httpStatus
	return r
}

func (r *HTTPStatusResponse) WithExcpectedHTTPStatus(expectedHTTPStatus []int) *HTTPStatusResponse {
	r.expectedHTTPStatus = expectedHTTPStatus
	return r
}

func (r *HTTPStatusResponse) HTTPStatus() int {
	return r.httpStatus
}

func (r *HTTPStatusResponse) ExpectedHTTPStatuses() []int {
	return r.expectedHTTPStatus
}
