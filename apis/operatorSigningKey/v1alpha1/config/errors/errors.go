package errors

var (
// ErrJetStreamNotEnabled is an error returned when both push and pull consumer are defined in the same consumer configuration.
// OperatorSigningKeyError OperatorSigningKeyError = &operatorSigningKeyError{message: "invalid configuration: cannot define both 'push' and 'pull' consumer"}
)

// OperatorSigningKeyError is an error result that happens when using consumer.
type OperatorSigningKeyError interface {
	error
}

type operatorSigningKeyError struct {
	message string
}

func (err *operatorSigningKeyError) Error() string {
	return err.message
}
