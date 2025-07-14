package errwrap

import "fmt"

func Wrap(s string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", s, err)
}

func Wrapf(err error, format string, args ...string) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args)

	return fmt.Errorf("%s:%w", msg, err)
}
