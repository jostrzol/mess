package handler

import (
	"encoding/gob"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/jostrzol/mess/pkg/server/adapter/schema"
	"github.com/jostrzol/mess/pkg/server/core/id"
)

const SessionKey = "session"

const sessionDataKey = "data"

func newSessionData() *schema.SessionData {
	return &schema.SessionData{
		ID: id.New[id.Session](),
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
