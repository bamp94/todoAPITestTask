package model

import (
	"fmt"
	"sort"
	"strings"

	"cyberzilla_api_task/config"

	"github.com/GuiaBolso/darwin"
	"github.com/gobuffalo/packr"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrInternal      = errors.New("Внутренняя ошибка сервера, повторите попытку позже или обратитесь к системному администратору")
	ErrModelNotFound = errors.New("Запись не найдена")
)

type TodoTask struct {
	ID                 int64
	Task               string
	AuthorizationToken string
}

// Model is data tier of 3-layer architecture
type Model struct {
	db *gorm.DB
}

// New Model constructor
func NewFromConfig(config config.Database) Model {
	db, err := gorm.Open("postgres", config.ConnURL())
	if err != nil {
		logrus.WithField("connURL", config.ConnURL()).WithError(err).Fatal("can't open connection with a database")
	}
	if err := db.DB().Ping(); err != nil {
		logrus.WithError(err).Fatal("can't ping connection with a database")
	}
	return Model{db: db}
}

// CheckMigrations validates database condition
func (m *Model) CheckMigrations() error {
	driver := darwin.NewGenericDriver(m.db.DB(), darwin.PostgresDialect{})
	d := darwin.New(driver, m.getMigrations(), nil)
	if err := d.Validate(); err != nil {
		return err
	}
	migrationInfo, err := d.Info()
	if err != nil {
		return err
	}
	for _, i := range migrationInfo {
		if i.Status == darwin.Applied {
			continue
		}
		return fmt.Errorf("found not applied migration: %s", i.Migration.Description)
	}
	return nil
}

// Migrate applies all migrations to connected database
func (m *Model) Migrate() {
	driver := darwin.NewGenericDriver(m.db.DB(), darwin.PostgresDialect{})
	d := darwin.New(driver, m.getMigrations(), nil)
	if err := d.Migrate(); err != nil {
		logrus.WithError(err).Error("can't migrate")
	}
}

// getMigrations provides migrations in darwin format
func (m *Model) getMigrations() []darwin.Migration {
	// migrationBox is used for embedding the migrations into the binary
	box := packr.NewBox("../etc/migrations")
	var migrations []darwin.Migration
	arr := box.List()
	sort.Strings(arr)
	for i, fileName := range arr {
		if !(strings.HasSuffix(fileName, ".sql") || strings.HasSuffix(fileName, ".SQL")) {
			logrus.Warnf("found file %s with unexpected type, skipping", fileName)
			continue
		}

		migration, err := box.FindString(fileName)
		if err != nil {
			logrus.WithError(err).Error("internal error of packr library")
		}
		migrations = append(migrations, darwin.Migration{
			Version:     float64(i + 1),
			Description: fileName,
			Script:      migration,
		})
	}
	return migrations
}

// Ping connection with database
func (m *Model) Ping() error {
	return m.db.DB().Ping()
}

// TodoList retrieves list of todo tasks
func (m *Model) TodoList(token string) ([]TodoTask, error) {
	var res []TodoTask
	raw := m.db.Raw(`SELECT * FROM todo_tasks tls;`)
	if token != "" {
		raw = m.db.Raw(`
		SELECT * FROM todo_tasks tls
		WHERE tls.authorization_token = ?;`, token)
	}
	if err := raw.Scan(&res).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []TodoTask{}, ErrModelNotFound
		}
		logrus.WithError(err).Errorf("can't get todo tasks by token: %s", token)
		return []TodoTask{}, ErrInternal
	}
	return res, nil
}

// CreateTodoTask creates todo task
func (m *Model) CreateTodoTask(task TodoTask) error {
	var res []TodoTask
	raw := m.db.Raw(`INSERT INTO todo_tasks(task, authorization_token)  VALUES ($1, $2);`, task.Task, task.AuthorizationToken)
	if err := raw.Scan(&res).Error; err != nil {
		logrus.WithError(err).Error("can't create todo task ")
		return ErrInternal
	}
	return nil
}
