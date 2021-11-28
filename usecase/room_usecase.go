package usecase

import (
	"chat/domain"
	"chat/utils/errors"
	"context"
	_ "reflect"
	"time"
)
type roomUseCase struct {
	rr      domain.RoomRepository
	timeout time.Duration
}

func (u *roomUseCase) SaveRoom(ctx context.Context, room *domain.ChatRoom) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	err := u.rr.SaveRoom(ctx , room )
	if err!=nil {
		return err
	}

	anyExistingRoom, err := u.rr.GetRoom(ctx , room.RoomID )
	if anyExistingRoom != nil {
		return errors.NewBadRequestError("the room already exist")
	}

   return err
}
