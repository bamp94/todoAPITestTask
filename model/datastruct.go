package model

type TodoList struct {
	ID          int
	Name        string
	AccessToken string
}

type TodoTask struct {
	ID         int
	TodoListID int
	Task       string
}
