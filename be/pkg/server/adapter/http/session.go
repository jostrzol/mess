package http

import (
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/ioc"
	"go.uber.org/zap"
)

const SessionKey = "session"

const sessionDataKey = "data"

type SessionData struct {
	ID uuid.UUID
}

func newSessionData() *SessionData {
	return &SessionData{
		ID: uuid.New(),
	}
}

func GetSessionData(c *gin.Context) *SessionData {
	session := sessions.Default(c)
	data, ok := session.Get(sessionDataKey).(*SessionData)
	if !ok {
		data = newSessionData()
		session.Set(sessionDataKey, data)
		err := session.Save()
		if err != nil {
			log := ioc.MustResolve[*zap.SugaredLogger]()
			log.Errorf("saving session data", zap.Error(err))
			return nil
		}
	}
	return data
}

func init() {
	gob.Register(&SessionData{})
}
