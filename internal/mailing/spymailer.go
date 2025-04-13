package mailing

import "slices"

type Spy struct {
	receivers []string
}

func NewSpy() *Spy {
	return &Spy{}
}

func (m *Spy) Send(_, to, _, _ string) error {
	m.receivers = append(m.receivers, to)
	return nil
}

func (m *Spy) HasReceiver(to string) bool {
	return slices.Contains(m.receivers, to)
}
