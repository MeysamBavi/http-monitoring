package model

import "strconv"

type ID uint64

func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

func ParseId(str string) (ID, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	return ID(id), err
}
