package utility

import (
	"errors"
	"log"
)

func LogError(errString string, err error) error {
	log.Print(errString, ": ", err)
	return errors.New(errString)
}
