package errors

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/getsentry/raven-go"
)

type wrappedError struct {
	message    string
	next       error
	frame      xerrors.Frame
	stackTrace *raven.Stacktrace
	extra      map[string]interface{}
}

func New(message string) error {
	return makeExtendedError(message, nil, nil)
}

func NewWithExtra(message string, extra map[string]interface{}) error {
	return makeExtendedError(message, nil, extra)
}

func Wrap(err error, message string) error {
	return makeExtendedError(message, err, nil)
}

func WrapWithExtra(err error, message string, extra map[string]interface{}) error {
	return makeExtendedError(message, err, extra)
}

func Extras(err error) map[string]interface{} {
	var wrappedErrorInstance *wrappedError
	if xerrors.As(err, &wrappedErrorInstance) {
		return wrappedErrorInstance.extra
	}
	return nil
}

func makeExtendedError(message string, next error, extra map[string]interface{}) error {
	return &wrappedError{
		message:    message,
		next:       next,
		frame:      xerrors.Caller(2),
		stackTrace: raven.GetOrNewStacktrace(next, 2, 0, nil),
		extra:      mergeMaps(Extras(next), extra),
	}
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, mapInstance := range maps {
		for key, value := range mapInstance {
			result[key] = value
		}
	}

	return result
}

func (err *wrappedError) Error() string {
	return fmt.Sprint(err)
}

// GetStacktrace implements "github.com/evalphobia/logrus_sentry.Stacktracer"
func (err *wrappedError) GetStacktrace() *raven.Stacktrace {
	return err.stackTrace
}

// Unwrap implements "golang.org/x/xerrors.Wrapper"
func (err *wrappedError) Unwrap() error {
	return err.next
}

// Format implements fmt.Formatter
func (err *wrappedError) Format(f fmt.State, c rune) {
	xerrors.FormatError(err, f, c)
}

// FormatError implements xerrors.Formatter
func (err *wrappedError) FormatError(p xerrors.Printer) (next error) {
	p.Print(err.message)
	err.frame.Format(p)
	return err.next
}
