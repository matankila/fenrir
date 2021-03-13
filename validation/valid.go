package validation

import "sync"

type Validator interface {
	IsValid(rawObj []byte, ns string) error
}

var (
	o sync.Once
	l = []Validator{}
)

func GetValidators() []Validator {
	o.Do(func() {
		l = append(l, Pod{})
		l = append(l, Service{})
	})

	return l
}
