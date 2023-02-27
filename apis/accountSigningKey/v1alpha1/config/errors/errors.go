package errors

var (
// ErrJetStreamNotEnabled is an error returned when both push and pull consumer are defined in the same consumer configuration.
// AccountSigningKeyError AccountSigningKeyError = &accountSigningKeyError{message: "invalid configuration: cannot define both 'push' and 'pull' consumer"}
)

// AccountSigningKeyError is an error result that happens when using consumer.
type AccountSigningKeyError interface {
	error
}

type accountSigningKeyError struct {
	message string
}

func (err *accountSigningKeyError) Error() string {
	return err.message
}
