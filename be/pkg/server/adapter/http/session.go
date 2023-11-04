package http

import (
	"encoding/gob"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/adapter/httpschema"
)

const SessionKey = "session"

const sessionDataKey = "data"

func newSessionData() *httpschema.SessionData {
	return &httpschema.SessionData{
		ID: uuid.New(),
	}
}

func GetSessionData(session sessions.Session) *httpschema.SessionData {
	data, ok := session.Get(sessionDataKey).(*httpschema.SessionData)
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
	gob.Register(&httpschema.SessionData{})
}
