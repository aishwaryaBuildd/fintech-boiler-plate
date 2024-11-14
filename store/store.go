package store

import (
	"context"
	"fintech/store/models"
	"time"
)

type Store interface {
	GetUserByPhoneNumber(context context.Context, phoneNumber string) (models.User, error)
	CreateUser(context context.Context, phoneNumber, otp string, otpExpiry time.Time, role string) error
	UpdateOTP(context context.Context, phoneNumber string, code string, expiry time.Time) error

	CreateCourse(context context.Context, course models.Course) error
	UpdateCourse(context context.Context, course models.Course) error
	ListCourse(context context.Context) ([]models.Course, error)
	GetCourse(context context.Context, id string) (models.Course, error)
	DeleteCourse(context context.Context, id string) error

	CreateSection(context context.Context, section models.Section) error
	UpdateSection(context context.Context, section models.Section) error
	ListSection(context context.Context) ([]models.Section, error)
	GetSection(context context.Context, id string) (models.Section, error)
	DeleteSection(context context.Context, id string) error

	CreateLecture(context context.Context, lecture models.Lecture) error
	UpdateLecture(context context.Context, lecture models.Lecture) error
	ListLecture(context context.Context) ([]models.Lecture, error)
	GetLecture(context context.Context, id string) (models.Lecture, error)
	DeleteLecture(context context.Context, id string) error

	GetOrCreateSession(context context.Context, message models.Message) (int, error)
	AddMessage(context context.Context, message models.Message) error
	GetChatSessions(context context.Context, userID int) ([]models.ChatSession, error)
	GetChatSessionsMessages(context context.Context, sessionID int) ([]models.Message, error)
	MarkChatSessionsAsRead(context context.Context, ChatSessionID int) error
}
