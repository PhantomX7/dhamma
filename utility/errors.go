package utility

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

func LogError(errString string, err error) error {
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Print(errString, ": ", err)
	}
	return err
}
