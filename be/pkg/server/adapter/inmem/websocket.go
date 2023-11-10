package inmem

import (
	"fmt"
	"sync"

	"github.com/golobby/container/v3"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"go.uber.org/zap"
)

type WsRepository struct {
	channels map[id.Session]chan<- (schema.Event)
	logger   *zap.Logger `container:"type"`
	mutex    sync.Mutex
}

func NewWsRepository() *WsRepository {
	repo := WsRepository{channels: make(map[id.Session]chan<- (schema.Event))}
	container.MustFill(container.Global, &repo)
	return &repo
}

func init() {
	container.MustSingletonLazy(container.Global, func() *WsRepository {
		return NewWsRepository()
	})
}

func (r *WsRepository) New(sessionID id.Session) <-chan (schema.Event) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	old, ok := r.channels[sessionID]
	if ok {
		r.logger.Warn("closing old websocket channel", zap.Stringer("session", sessionID))
		close(old)
	}
	channel := make(chan (schema.Event))
	r.channels[sessionID] = channel
	return channel
}

func (r *WsRepository) Send(sessionID id.Session, event schema.Event) error {
	var c chan<- (schema.Event)
	func() {
		r.mutex.Lock()
		defer r.mutex.Unlock()
		c = r.channels[sessionID]
	}()
	if c == nil {
		return fmt.Errorf("sending to a nonexistant socket")
	}
	c <- event
	return nil
}

func (r *WsRepository) Close(sessionID id.Session) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	c, ok := r.channels[sessionID]
	if ok {
		delete(r.channels, sessionID)
		close(c)
	}
}
