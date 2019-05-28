package bfss_sdk

type BfssException interface {
	Error() string
}

type bfssException struct {
	why string
}

func (self bfssException) GetWhy() string {
	return self.why
}

func (self bfssException) Error() string {
	return self.why
}

func createException(why string) BfssException {
	return bfssException{why: why}
}

func toBfssException(err interface{}) BfssException {
	switch err.(type) {
	case BfssException:
		return err.(BfssException)
	default:
		return bfssException{why: err.(string)}
	}
}

var FileEofException = createException("end of file")
var InvalidSizeException = createException("invalid size")
var AlreadyWrittenException = createException("already written")
var NotOpenException = createException("not open")
var InvalidWhenceException = createException("invalid whence")
var InvalidSchemeException = createException("invalid scheme")
var AccessReadException = createException("deny access to read interface")
var AccessWriteException = createException("deny access to read interface")
var AccessRenameException = createException("deny access to rename interface")
var AccessResizeException = createException("deny access to resize interface")
var ConnectException = createException("cannot connect to server")
var InvalidDataException = createException("invalid data")
var InvalidOffsetException = createException("invalid offset")
