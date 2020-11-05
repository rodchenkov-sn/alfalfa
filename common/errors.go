package common

type UserAlreadyExistError struct {
	Login string
}

func (e UserAlreadyExistError) Error() string {
	return e.Login + " is already exist"
}

type UserNotfoundError struct {
	Login string
}

func (e UserNotfoundError) Error() string {
	return e.Login + " was not found"
}

type InvalidPasswordError struct {
	Login string
}

func (e InvalidPasswordError) Error() string {
	return "Password for " + e.Login + " was invalid"
}

type TokenExpiredError struct {}

func (e TokenExpiredError) Error() string {
	return "Token expired"
}

type InvalidTokenError struct {}

func (e InvalidTokenError) Error() string {
	return "Invalid token"
}
