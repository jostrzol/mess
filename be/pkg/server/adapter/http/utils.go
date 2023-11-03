package http

import (
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
)

func parseUUID(str string) (uuid.UUID, error) {
	result, err := uuid.Parse(str)
	if err != nil {
		return uuid.UUID{}, usrerr.Wrap(err, "invalid uuid format")
	}
	return result, nil
}
