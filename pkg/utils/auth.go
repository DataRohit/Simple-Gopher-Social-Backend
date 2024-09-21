package utils

import (
	"errors"
)

func VerifyOwnership(userID, authUserID string) error {
	if userID != authUserID {
		return errors.New("user does not own the resource")
	}
	return nil
}
