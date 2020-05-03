package application

import (
	"cyberzilla_api_task/config"
	"cyberzilla_api_task/model"
)

// Application tier of 3-layer architecture
type Application struct {
	model  model.Model
	config config.Main
}

// New Application constructor
func New(m model.Model, c config.Main) Application {
	return Application{
		model:  m,
		config: c,
	}
}

// PingDatabase ensures db connection is valid
func (a *Application) PingDatabase() error {
	return a.model.Ping()
}

// TodoTasksList retrieves list of todo tasks
func (a *Application) TodoTasksList(token string) ([]model.TodoTask, error) {
	if _, err := a.model.TodoList(token); err != nil {
		return []model.TodoTask{}, err
	}
	return a.model.TodoTasks(token)
}

// CreateTodoTask creates todo task
func (a *Application) CreateTodoTask(token string, task model.TodoTask) error {
	todoList, err := a.model.TodoList(token)
	if err != nil {
		return err
	}
	return a.model.CreateTodoTask(todoList.ID, task)
}

// TodoTask retrieves todo task
func (a *Application) TodoTask(id int64, token string) (model.TodoTask, error) {
	return a.model.TodoTask(id, token)
}

// CreateTodoTask creates todo task
func (a *Application) UpdateTodoTask(token string, task model.TodoTask) error {
	if _, err := a.model.TodoTask(task.ID, token); err != nil {
		return err
	}
	return a.model.UpdateTodoTask(task)
}
