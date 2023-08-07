package errs

import "errors"

var (
	ConfiguredBadEnvError    = errors.New("configured bad key env variable")
	NotMatchCredentialsError = errors.New("not match credentials")
	NilError                 = errors.New("nil error")
)
