package utils

import "github.com/google/uuid"

func GenIntUUID() (uu uint32, err error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return
	}
	uu = u.ID()
	return
}

func GenStrUUID() (uu string, err error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return
	}
	uu = u.String()
	return
}
