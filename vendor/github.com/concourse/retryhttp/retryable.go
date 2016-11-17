package retryhttp

import (
	"errors"
	"net"
	"strings"
	"syscall"
)

func retryable(err error) bool {
	if neterr, ok := err.(net.Error); ok {
		if neterr.Temporary() {
			return true
		}
	}

	s := err.Error()
	for _, retryableError := range retryableErrors {
		if strings.HasSuffix(s, retryableError.Error()) {
			return true
		}
	}

	return false
}

var retryableErrors = []error{
	syscall.ECONNREFUSED,
	syscall.ECONNRESET,
	syscall.ETIMEDOUT,
	errors.New("i/o timeout"),
	errors.New("no such host"),
	errors.New("remote error: handshake failure"),
}
