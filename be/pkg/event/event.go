package event

type Event interface{}

type Subject interface {
	IsObserving(observer Observer) bool
	Observe(observer Observer)
	Unobserve(observer Observer)
	Notify(event Event)
}

type subjectImpl struct {
	observers map[Observer]struct{}
}

func NewSubject() Subject {
	return &subjectImpl{
		observers: make(map[Observer]struct{}),
	}
}

func (s *subjectImpl) IsObserving(observer Observer) bool {
	_, present := s.observers[observer]
	return present
}

func (s *subjectImpl) Observe(observer Observer) {
	s.observers[observer] = struct{}{}
}

func (s *subjectImpl) Unobserve(observer Observer) {
	delete(s.observers, observer)
}

func (s *subjectImpl) Notify(event Event) {
	for observer := range s.observers {
		observer.Handle(event)
	}
}

type Observer interface {
	Handle(event Event)
}
