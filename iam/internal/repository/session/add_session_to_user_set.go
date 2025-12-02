package session

import (
	"context"
)

func (r *repository) AddToUserSet(ctx context.Context, userUUID, sessionUUID string) error {
	err := r.cache.SAdd(ctx, userUUID, sessionUUID)
	if err != nil {
		return err
	}
	return nil
}
