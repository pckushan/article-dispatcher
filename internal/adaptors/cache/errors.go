package cache

type TypeAssertError struct {
	error
}

type DuplicateError struct {
	error
}

type DataNotFoundError struct {
	error
}

type InvalidDataError struct {
	error
}
