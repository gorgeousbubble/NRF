package util

import "github.com/google/uuid"

func checkUUID(id uuid.UUID) (b bool, err error) {
	_, err = uuid.Parse(id.String())
	if err != nil {
		return b, err
	}
	return b, err
}
