package inmem

import (
	"github.com/golobby/container/v3"
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/adapter/ws"
	"go.uber.org/zap"
)

type WsRepository struct {
	channels map[uuid.UUID]chan<- (ws.Event)
	logger   *zap.SugaredLogger `container:"type"`
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
	old, ok := r.channels[sessionID]
	if ok {
		r.logger.Warn("closing old websocket channel", zap.String("session", sessionID.String()))
		close(old)
	}
	channel := make(chan (ws.Event))
	r.channels[sessionID] = channel
	return channel
}

func (r *WsRepository) Get(sessionID uuid.UUID) chan<- (ws.Event) {
	return r.channels[sessionID]
}
