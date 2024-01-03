package tools

import "github.com/google/uuid"

func StringToUuid(s string) (id uuid.UUID, err error) {
	id, err = uuid.Parse(s)
	return
}

func UuidToString(id uuid.UUID) (s string) {
	s = id.String()
	return
}

func ListUuidToString(ids []uuid.UUID) (list []string) {
	for i := 0; i < len(ids); i++ {
		list = append(list, ids[i].String())
	}
	return
}
