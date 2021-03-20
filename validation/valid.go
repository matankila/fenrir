package validation

import "sync"

type Validator interface {
	IsValid(rawObj []byte, ns string) error
}

var (
	o sync.Once
	l = map[string]Validator{}
)

func Init() {
	o.Do(func() {
		l["Pod"] = Pod{}
		l["Service"] = Service{}
	})
}

func Get(validator string) (Validator, bool) {
	v, ok := l[validator]
	return v, ok
}
