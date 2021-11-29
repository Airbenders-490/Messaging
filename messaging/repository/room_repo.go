// Package repository todo: this must be removed
package repository

import (
	"chat/domain"
	"context"
)

type roomRepository struct {
}

func (r roomRepository) SaveRoom(ctx context.Context, room *domain.ChatRoom) error {
	return nil
}

func (r roomRepository) GetRoom(ctx context.Context, roomID string) (*domain.ChatRoom, error) {
	panic("implement me")
}

func (r roomRepository) GetRoomsFor(ctx context.Context, studentID string) (*domain.StudentChatRooms, error) {
	return nil, nil
}

func (r roomRepository) EditChatRoomParticipants(ctx context.Context, roomID string, student []domain.Student) error {
	panic("implement me")
}

func (r roomRepository) AddRoomForParticipants(ctx context.Context, roomID string, student []domain.Student) error {
	panic("implement me")
}

func (r roomRepository) RemoveRoomForParticipants(ctx context.Context, roomID string, student []domain.Student) error {
	panic("implement me")
}

func (r roomRepository) DeleteRoom(ctx context.Context, roomID string) error {
	panic("implement me")
}

// NewRoomRepo temp!
func NewRoomRepo() domain.RoomRepository {
	return &roomRepository{}
}
