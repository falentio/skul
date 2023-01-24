package app

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gofiber/storage"
	"github.com/gofiber/storage/badger"
	"github.com/gofiber/storage/memory"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
	"github.com/falentio/skul/internal/pkg/seeder"
	"github.com/falentio/skul/internal/service/admin"
	"github.com/falentio/skul/internal/service/enterance_token"
	"github.com/falentio/skul/internal/service/examination"
	"github.com/falentio/skul/internal/service/examine_answer"
	"github.com/falentio/skul/internal/service/examine_attatchment"
	"github.com/falentio/skul/internal/service/examine_question"
	"github.com/falentio/skul/internal/service/examine_student"
	"github.com/falentio/skul/internal/service/file"
	"github.com/falentio/skul/internal/service/student"
	"github.com/falentio/skul/internal/service/student_answer"
	"github.com/falentio/skul/web"
)

type Repository struct {
	AdminRepository              domain.AdminRepository
	EnteranceTokenRepository     domain.EnteranceTokenRepository
	ExaminationRepository        domain.ExaminationRepository
	ExamineAnswerRepository      domain.ExamineAnswerRepository
	ExamineAttatchmentRepository domain.ExamineAttatchmentRepository
	ExamineStudentRepository     domain.ExamineStudentRepositoryRead
	ExamineQuestionRepository    domain.ExamineQuestionRepository
	StudentRepository            domain.StudentRepository
	StudentAnswerRepository      domain.StudentAnswerRepository
}

type Application struct {
	Options AppOptions
	Logger  zerolog.Logger

	router     chi.Router
	storage    storage.Storage
	repository Repository
}

func (app *Application) repositoryGuard() {
	v := reflect.ValueOf(app.repository)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.IsNil() {
			err := fmt.Errorf("Application: Application.repository.%s is nil, call (*Application).InitRepository() before (*Application).ListenAndServe()", f.Type().Name())
			log.Fatal().Err(err).Msg("invalid repository")
		}
	}
}

func (app *Application) Handler() http.Handler {
	if app.router == nil {
		app.router = chi.NewRouter()
	}
	return app.router
}

func (app *Application) ListenAndServe() error {
	app.repositoryGuard()
	return http.ListenAndServe(app.Options.Addr, app.Handler())
}

func (app *Application) InitStorage() error {
	switch app.Options.Storage.Driver {
	case "badger":
		if app.Options.Storage.Badger.Database == "" {
			app.Options.Storage.Badger = badger.ConfigDefault
		}
		app.storage = badger.New(app.Options.Storage.Badger)
	case "memory":
		if app.Options.Storage.Memory.GCInterval == 0 {
			app.Options.Storage.Memory = memory.ConfigDefault
		}
		app.storage = memory.New(app.Options.Storage.Memory)
	}
	return nil
}

func (app *Application) InitRepository() error {
	switch app.Options.Database.Driver {
	case "sqlite3", "sqlite":
		dialect := sqlite.Open(app.Options.Database.Dsn)
		return app.initGormRepository(dialect)
	}
	return nil
}

func (app *Application) initGormRepository(dialect gorm.Dialector) error {
	db, err := gorm.Open(dialect)
	if err != nil {
		return err
	}
	app.repository.AdminRepository = &admin.AdminRepositoryGorm{
		DB: db,
	}
	app.repository.StudentRepository = &student.StudentRepositoryGorm{
		DB: db,
	}
	app.repository.ExaminationRepository = &examination.ExaminationRepositoryGorm{
		DB: db,
	}
	app.repository.ExamineAnswerRepository = &examineanswer.ExamineAnswerRepositoryGorm{
		DB: db,
	}
	app.repository.ExamineAttatchmentRepository = &examineattatchment.ExamineAttatchmentRepositoryGorm{
		DB: db,
	}
	app.repository.ExamineQuestionRepository = &examinequestion.ExamineQuestionRepositoryGorm{
		DB: db,
	}
	app.repository.ExamineStudentRepository = &examinestudent.ExamineStudetnRepositoryGorm{
		DB: db,
	}
	app.repository.EnteranceTokenRepository = &enterancetoken.EnteranceTokenRepositoryGorm{
		DB: db,
	}
	app.repository.StudentAnswerRepository = &studentanswer.StudentAnswerRepositoryGorm{
		DB: db,
	}
	if err := db.AutoMigrate(
		&domain.Admin{},
		&domain.EnteranceToken{},
		&domain.Examination{},
		&domain.ExamineAnswer{},
		&domain.ExamineAttatchment{},
		&domain.ExamineStudent{},
		&domain.ExamineQuestion{},
		&domain.Student{},
	); err != nil {
		return err
	}
	return nil
}

func (app *Application) SeedRepository() error {
	if err := seeder.AdminRepository(app.repository.AdminRepository); err != nil {
		return err
	}
	if err := seeder.StudentRepository(app.repository.StudentRepository); err != nil {
		return err
	}
	return nil
}

func (app *Application) InitHandler() {
	secret, err := base64.RawURLEncoding.DecodeString(app.Options.JWTSecret)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("failed to decode jwt secret")
	}
	auth := &auth.Auth{
		Name:          "session",
		Secure:        false,
		SigningMethod: jwt.SigningMethodHS512,
		Secret:        secret,
		Logger:        app.Logger,
	}
	adminRouter := &admin.AdminRouter{
		Auth: auth,
		AdminService: &admin.AdminService{
			AdminRepository: app.repository.AdminRepository,
			Auth:            auth,
			Logger:          app.Logger,
		},
	}
	studentRouter := &student.StudentRouter{
		Auth: auth,
		StudentService: &student.StudentService{
			StudentRepository: app.repository.StudentRepository,
			Auth:              auth,
			Logger:            app.Logger,
		},
	}
	fileRouter := &file.FileRouter{
		Auth: auth,
		FileService: &file.FileService{
			Auth:    auth,
			Logger:  app.Logger,
			Storage: app.storage,
		},
	}
	examinationRouter := &examination.ExaminationRouter{
		Auth: auth,
		ExaminationService: &examination.ExaminationService{
			ExaminationRepository: app.repository.ExaminationRepository,
		},
	}
	examineQuestionRouter := &examinequestion.ExamineQuestionRouter{
		Auth: auth,
		ExamineQuestionService: &examinequestion.ExamineQuestionService{
			Auth:                         auth,
			Logger:                       app.Logger,
			ExamineQuestionRepository:    app.repository.ExamineQuestionRepository,
			ExamineAnswerRepository:      app.repository.ExamineAnswerRepository,
			ExamineAttatchmentRepository: app.repository.ExamineAttatchmentRepository,
		},
	}
	examineAnswerRouter := &examineanswer.ExamineAnswerRouter{
		Auth: auth,
		ExamineAnswerService: &examineanswer.ExamineAnswerService{
			Auth:                    auth,
			ExamineAnswerRepository: app.repository.ExamineAnswerRepository,
		},
	}
	examineAttatchmentRouter := &examineattatchment.ExamineAttatchmentRouter{
		Auth: auth,
		ExamineAttatchmentService: &examineattatchment.ExamineAttatchmentService{
			Auth:                         auth,
			Logger:                       app.Logger,
			ExamineAttatchmentRepository: app.repository.ExamineAttatchmentRepository,
		},
	}
	enteranceTokenRouter := &enterancetoken.EnteranceTokenRouter{
		Auth: auth,
		EnteranceTokenService: &enterancetoken.EnteranceTokenService{
			Auth:                     auth,
			EnteranceTokenRepository: app.repository.EnteranceTokenRepository,
		},
	}
	studentAnswerRouter := &studentanswer.StudentAnswerRouter{
		Auth: auth,
		StudentAnswerService: &studentanswer.StudentAnswerService{
			Auth:                    auth,
			StudentAnswerRepository: app.repository.StudentAnswerRepository,
		},
	}

	r := chi.NewRouter()

	// register middewares
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 5))
	r.Use(middleware.AllowContentType("application/json", "multipart/form-data"))
	r.Use(cors.Handler(cors.Options{
		AllowCredentials: true,
		MaxAge:           7200,
	}))

	r.With(middleware.NoCache).Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.NewOK("ok").ServeHTTP(w, r)
	})
	r.With(middleware.NoCache).Get("/logout", auth.Logout)

	// register service handler
	r.Route("/api", func (r chi.Router) {
		r.Route("/admin", adminRouter.Route)
		r.Route("/file", fileRouter.Route)
		r.Route("/examination", examinationRouter.Route)
		r.Route("/examine-question", examineQuestionRouter.Route)
		r.Route("/examine-answer", examineAnswerRouter.Route)
		r.Route("/examine-attatchment", examineAttatchmentRouter.Route)
		r.Route("/enterance-token", enteranceTokenRouter.Route)
		r.Route("/student", studentRouter.Route)
		r.Route("/student-answer", studentAnswerRouter.Route)
	})

	// register website handler
	r.Handle("/*", web.FileServer)

	app.router = r
}
