package jsonresponder

type ErrorContainer map[string][]string

func NewErrorContainer() ErrorContainer {
	return ErrorContainer{}
}

func (e ErrorContainer) Add(field, message string) {
	if _, ok := e[field]; !ok {
		e[field] = []string{}
	}

	e[field] = append(e[field], message)
}

func (e ErrorContainer) NotEmpty() bool {
	return len(e) > 0
}
