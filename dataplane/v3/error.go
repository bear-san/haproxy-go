package v3

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

type InvalidResponseError struct {
	Message string
}

func (e *InvalidResponseError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type UnknownError struct {
	Message string
}

func (e *UnknownError) Error() string {
	return e.Message
}

type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return e.Message
}
