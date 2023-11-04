package http

import (
	"encoding/gob"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/google/uuid"
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

func GetSessionData(session sessions.Session) *SessionData {
	data, ok := session.Get(sessionDataKey).(*SessionData)
	if !ok {
		data = newSessionData()
		session.Set(sessionDataKey, data)
		err := session.Save()
		if err != nil {
			panic(fmt.Errorf("saving session data: %w", err))
		}
	}
	return data
}

func init() {
	gob.Register(&SessionData{})
}
