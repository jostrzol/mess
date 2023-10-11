package eventtest

import (
	"testing"

	"github.com/jostrzol/mess/pkg/event"
	"github.com/stretchr/testify/assert"
)

type MockObserver struct {
	observed []event.Event
	t        *testing.T
}

func NewMockObserver(t *testing.T) *MockObserver {
	return &MockObserver{t: t}
}

func (m *MockObserver) Handle(event event.Event) {
	m.observed = append(m.observed, event)
}

func (m *MockObserver) Observed(event event.Event) bool {
	ok := false
	for _, e := range m.observed {
		if e == event {
			ok = true
			break
		}
	}
	assert.Truef(m.t, ok, "Observation not made\nhave (%v)\nwant (%v)", m.observed, event)
	return ok
}

func (m *MockObserver) ObservedMatch(events ...event.Event) bool {
	return assert.ElementsMatchf(
		m.t,
		events,
		m.observed,
		"Observtions don't match\nhave (%v)\nwant (%v)", m.observed, events)
}

func (m *MockObserver) Reset() {
	m.observed = make([]event.Event, 0)
}
