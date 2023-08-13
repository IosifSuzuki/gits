package errs

import "errors"

var (
	ConfiguredBadEnvError             = errors.New("configured bad key env variable")
	NotMatchCredentialsError          = errors.New("not match credentials")
	NilError                          = errors.New("nil error")
	UnknownFileNameError              = errors.New("unknown file name error")
	UnsuccessfulUploadAttachmentError = errors.New("unsuccessful upload file")
	UnsuccessfulCreateArticleError    = errors.New("unsuccessful create article")
)
