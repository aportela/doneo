package utils

import "github.com/gofrs/uuid"

func UUID() string {
	id := func() string { u, _ := uuid.NewV7(); return u.String() }()
	return (id)
}
