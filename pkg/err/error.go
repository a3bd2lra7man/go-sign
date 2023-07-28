package errs

type UnExpectedError struct {
}

func (err UnExpectedError) Error() string {
	return "UnExpected"
}
