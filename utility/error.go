package utility

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

func LogError(errString string, err error) error {
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Print(errString, ": ", err)
	}
	return err
}
