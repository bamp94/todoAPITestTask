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

type TodoList struct {
	ID                 int64   `json:"id" example:"1"`
	Name               string  `json:"name" example:"Dayly tasks"`
	AuthorizationToken *string `json:"authorizationToken" example:"kad3u452iubdifadDWA"`
}

type TodoTask struct {
	ID         int64  `json:"id" example:"1"`
	TodoListID int    `json:"todoListID" example:"1"`
	Task       string `json:"task" example:"Do my homework"`
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

// TodoList retrieves todo list
func (m *Model) TodoList(token string) (TodoList, error) {
	var res TodoList
	if err := m.db.Raw(`
		SELECT * FROM todo_lists tls
		WHERE tls.authorization_token = ?
		LIMIT 1;`, token).Scan(&res).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return TodoList{}, ErrModelNotFound
		}
		logrus.WithError(err).Errorf("can't get todo list by token: %s", token)
		return TodoList{}, ErrInternal
	}
	return res, nil
}

// TodoTasks retrieves todo tasks by token
func (m *Model) TodoTasks(token string) ([]TodoTask, error) {
	var res []TodoTask
	if err := m.db.Raw(`
		SELECT * FROM todo_tasks tts
		JOIN todo_lists tls ON tls.id = tts.todo_list_id AND  tls.authorization_token = ?;`, token).
		Scan(&res).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []TodoTask{}, ErrModelNotFound
		}
		logrus.WithError(err).Errorf("can't get todo tasks by token: %s", token)
		return []TodoTask{}, ErrInternal
	}
	return res, nil
}

// CreateTodoTask creates todo task
func (m *Model) CreateTodoTask(todoListID int64, task TodoTask) error {
	var res []TodoTask
	raw := m.db.Raw(`INSERT INTO todo_tasks(todo_list_id, task)  VALUES ($1, $2);`, todoListID, task.Task)
	if err := raw.Scan(&res).Error; err != nil {
		logrus.WithError(err).Error("can't create todo task ")
		return ErrInternal
	}
	return nil
}

// TodoTask retrieves todo task
func (m *Model) TodoTask(id int64, token string) (TodoTask, error) {
	var res TodoTask
	if err := m.db.Raw(`
		SELECT * FROM todo_tasks tts 
		JOIN todo_lists tls ON tls.id = tts.todo_list_id AND  tls.authorization_token = $1
		WHERE tts.id = $2;`, token, id).Scan(&res).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return TodoTask{}, ErrModelNotFound
		}
		logrus.WithError(err).Errorf("can't get todo task by id: %v", id)
		return TodoTask{}, ErrInternal
	}
	return res, nil
}

// UpdateTodoTask updates todo task
func (m *Model) UpdateTodoTask(task TodoTask) error {
	var res []TodoTask
	if err := m.db.Raw(`
		UPDATE todo_tasks SET task = $1 
		WHERE id = $2;`, task.Task, task.ID).
		Scan(&res).Error; err != nil {
		logrus.WithError(err).Error("can't update todo task ")
		return ErrInternal
	}
	return nil
}

// DeleteTodoTask deletes todo task
func (m *Model) DeleteTodoTask(taskID int64) error {
	var res []TodoTask
	if err := m.db.Raw(`
		DELETE FROM todo_tasks 
		WHERE id = ?;`, taskID).
		Scan(&res).Error; err != nil {
		logrus.WithError(err).Error("can't delete todo task ")
		return ErrInternal
	}
	return nil
}
