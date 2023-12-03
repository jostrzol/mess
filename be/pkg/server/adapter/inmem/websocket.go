package inmem

import (
	"fmt"
	"sync"
	"time"

	"github.com/golobby/container/v3"
	"github.com/jostrzol/mess/configs/serverconfig"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type WsRepository struct {
	logger   *zap.Logger          `container:"type"`
	config   *serverconfig.Config `container:"type"`
	channels map[id.Session]*websocket
	mutex    sync.Mutex
}

type websocket struct {
	channel    chan (schema.Event)
	errorCount int
}

func NewWsRepository() *WsRepository {
	repo := WsRepository{channels: make(map[id.Session]*websocket)}
	container.MustFill(container.Global, &repo)
	go repo.heartbeatTask()
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
		close(old.channel)
	}
	ws := &websocket{
		channel:    make(chan (schema.Event)),
		errorCount: 0,
	}
	r.channels[sessionID] = ws
	return ws.channel
}

func (r *WsRepository) Send(sessionID id.Session, event schema.Event) error {
	var ws *websocket
	func() {
		r.mutex.Lock()
		defer r.mutex.Unlock()
		ws = r.channels[sessionID]
	}()
	if ws == nil {
		return fmt.Errorf("sending to a nonexistant socket")
	}
	ws.channel <- event
	return nil
}

func (r *WsRepository) heartbeatTask() {
	for {
		time.Sleep(r.config.HeartbeatPeriod)
		var websockets map[id.Session]*websocket
		func() {
			r.mutex.Lock()
			defer r.mutex.Unlock()
			websockets = maps.Clone(r.channels)
		}()
		for sessionID, ws := range websockets {
			if ws.errorCount >= r.config.MaxWebsocketErrors {
				r.logger.Info(
					"websocket error count exceeded maximum; removing",
					zap.Stringer("session", sessionID),
					zap.Int("maximum", r.config.MaxWebsocketErrors))
				close(ws.channel)
				delete(r.channels, sessionID)
				continue
			}
			ws.channel <- &schema.Heartbeat{}
		}
	}
}

func (r *WsRepository) IndicateError(sessionID id.Session) {
	ws := r.channels[sessionID]
	if ws != nil {
		ws.errorCount++
	} else {
		r.logger.Error("indicate error on non-existant channel", zap.Stringer("session", sessionID))
	}
}

func (r *WsRepository) Close(sessionID id.Session) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	c, ok := r.channels[sessionID]
	if ok {
		delete(r.channels, sessionID)
		close(c.channel)
	}
}
