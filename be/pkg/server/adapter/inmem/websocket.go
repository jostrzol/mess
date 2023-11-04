package inmem

import (
	"fmt"
	"sync"

	"github.com/golobby/container/v3"
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/adapter/ws"
	"go.uber.org/zap"
)

type WsRepository struct {
	channels map[uuid.UUID]chan<- (ws.Event)
	logger   *zap.Logger `container:"type"`
	mutex    sync.Mutex
}

func NewWsRepository() *WsRepository {
	repo := WsRepository{channels: make(map[uuid.UUID]chan<- (ws.Event))}
	container.MustFill(container.Global, &repo)
	return &repo
}

func init() {
	container.MustSingletonLazy(container.Global, func() *WsRepository {
		return NewWsRepository()
	})
}

func (r *WsRepository) New(sessionID uuid.UUID) <-chan (ws.Event) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	old, ok := r.channels[sessionID]
	if ok {
		r.logger.Warn("closing old websocket channel", zap.Stringer("session", sessionID))
		close(old)
	}
	channel := make(chan (ws.Event))
	r.channels[sessionID] = channel
	return channel
}

func (r *WsRepository) Send(sessionID uuid.UUID, event ws.Event) error {
	var c chan<- (ws.Event)
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

func (r *WsRepository) Close(sessionID uuid.UUID) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	c, ok := r.channels[sessionID]
	if ok {
		delete(r.channels, sessionID)
		close(c)
	}
}