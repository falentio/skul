package domain

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/pkg/response"
)

var (
	ErrStudentConflict = errors.New("student: student data already exists")
	ErrStudentNotFound = errors.New("student: can not found student")
)

const StudentIDPrefix = "stu"

type Student struct {
	Model

	AdminID raid.Raid `json:"adminID" gorm:"type:varchar(32);not null"`

	Name           string `json:"name"`
	Username       string `json:"username" validate:"max=32" gorm:"type:varchar(32)"`
	Class          string `json:"class"`
	Grade          string `json:"grade"`
	PresenceNumber int    `json:"presenceNumber"`
	PasswordHash   string `json:"-"`
	Password       string `json:"password,omitempty" validate:"min=8" gorm:"-"`

	EnteranceTokens []*EnteranceToken `json:"enteranceTokens" gorm:"many2many:examine_student"`
	ExamineAnswer   []*ExamineAnswer  `json:"examineAnswer" gorm:"many2many:student_examine_answer"`
	Admin           *Admin            `json:"admin"`
}

type ListStudentOptions struct {
	PaginateOptions

	AdminID        raid.Raid `json:"adminID"`
	Name           string    `json:"name"`
	Class          string    `json:"class"`
	Grade          string    `json:"grade"`
	PresenceNumber int       `json:"presenceNumber"`
}

type StudentRepositoryRead interface {
	GetStudent(ctx context.Context, studentID raid.Raid) (*Student, error)
	GetStudentByUsername(ctx context.Context, username string) (*Student, error)
	ListStudent(ctx context.Context, opts *ListStudentOptions) ([]*Student, error)
}

type StudentRepositoryWrite interface {
	CreateStudent(ctx context.Context, student *Student) error
	BatchCreateStudent(ctx context.Context, students []*Student) error
	DeleteStudent(ctx context.Context, studentID raid.Raid) error
	UpdateStudent(ctx context.Context, student *Student) error
}

type StudentRepository interface {
	StudentRepositoryRead
	StudentRepositoryWrite
}

type StudentServiceRead interface {
	GetStudent(ctx context.Context, studentID raid.Raid) (response.Response, error)
	ListStudent(ctx context.Context, opts *ListStudentOptions) (response.Response, error)
}

type StudentServiceWrite interface {
	CreateStudent(ctx context.Context, student *Student) (response.Response, error)
	BatchCreateStudent(ctx context.Context, students []*Student) (response.Response, error)
	DeleteStudent(ctx context.Context, studentID raid.Raid) (response.Response, error)
	UpdateStudent(ctx context.Context, student *Student) (response.Response, error)
	LoginStudent(ctx context.Context, student *Student) (res response.Response, err error)
}

type StudentService interface {
	StudentServiceRead
	StudentServiceWrite
}
