package errors

type ConflictError struct {
	Err string
}

type NotFoundError struct {
	Err string
}

type InternalError struct {
	Err string
}

type BadRequest struct {
	Err string
}

func (c ConflictError) Error() string {
	return c.Err
}

func (n NotFoundError) Error() string {
	return n.Err
}

func (i InternalError) Error() string {
	return i.Err
}

func (b BadRequest) Error() string {
	return b.Err
}
