package model

type ID string

func (id ID) String() string {
	return string(id)
}

func ParseId(str string) (ID, error) {
	return ID(str), nil
}
