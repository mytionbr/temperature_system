package model

type StatusError struct {
	Code    int
	Err     error
	Message string
	Caller  string
}
