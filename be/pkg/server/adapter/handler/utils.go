package handler

import (
	"github.com/google/uuid"
	"github.com/jostrzol/mess/pkg/server/core/id"
	"github.com/jostrzol/mess/pkg/server/core/usrerr"
)

func parseUUID[T id.ID](str string) (T, error) {
	result, err := uuid.Parse(str)
	if err != nil {
		return T{}, usrerr.Wrap(err, "invalid uuid format")
	}
	return T{BaseID: id.BaseID{UUID: result}}, nil
}
