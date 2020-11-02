package service

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
