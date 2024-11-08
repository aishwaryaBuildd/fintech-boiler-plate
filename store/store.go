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

	CreateFolder(context context.Context, folder models.Folder) error
	UpdateFolder(context context.Context, folder models.Folder) error
	ListFolder(context context.Context) ([]models.Folder, error)
	GetFolder(context context.Context, id string) (models.Folder, error)
	DeleteFolder(context context.Context, id string) error

	AddMessage(context context.Context, message models.Message) error
}
