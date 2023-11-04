package handler

import (
	"encoding/gob"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
)

const SessionKey = "session"

const sessionDataKey = "data"

func newSessionData() *schema.SessionData {
	return &schema.SessionData{
		ID: uuid.New(),
	}
}

func GetSessionData(session sessions.Session) *schema.SessionData {
	data, ok := session.Get(sessionDataKey).(*schema.SessionData)
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
	gob.Register(&schema.SessionData{})
}
