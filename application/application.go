package application

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"cyberzilla_api_task/config"
	"cyberzilla_api_task/model"
)

// Application tier of 3-layer architecture
type Application struct {
	model  model.Model
	config config.Main
}

// QueryResp contains channel response
type QueryResp struct {
	Addr string
	Time float64
}

// ProxyServersStatus contains proxy server status
type ProxyServersStatus struct {
	ProxyServer string
	Status      string
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

// UpdateTodoTask updates todo task
func (a *Application) UpdateTodoTask(token string, task model.TodoTask) error {
	if _, err := a.model.TodoTask(task.ID, token); err != nil {
		return err
	}
	return a.model.UpdateTodoTask(task)
}

// DeleteTodoTask deletes todo task
func (a *Application) DeleteTodoTask(taskID int64, token string) error {
	if _, err := a.model.TodoTask(taskID, token); err != nil {
		return err
	}
	return a.model.DeleteTodoTask(taskID)
}

// ProxyServersStatus retrieves proxy servers working statuses
func (a *Application) ProxyServersStatuses(proxyServersList []string) []ProxyServersStatus {
	//runtime.GOMAXPROCS(4)
	var statuses []ProxyServersStatus
	resp_chan := make(chan QueryResp, 10)
	var proxyNumber int
	for _, server := range proxyServersList {
		if !validateIPPort(server) {
			statuses = append(statuses, ProxyServersStatus{server, "Invalid"})
			continue
		}
		address := strings.SplitN(server, string(':'), 2)
		ip, port := address[0], address[1]
		go a.checkProxyServer(ip, port, resp_chan)
		proxyNumber++
	}
	for i := 0; i < proxyNumber; i++ {
		r := <-resp_chan
		if r.Time > 0 {
			statuses = append(statuses, ProxyServersStatus{r.Addr, "Active"})
			continue
		}
		statuses = append(statuses, ProxyServersStatus{r.Addr, "Disabled"})
	}
	return statuses
}

// validateIPPort check if text validate on ip:port
func validateIPPort(text string) bool {
	re := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):[0-9]+$`)
	return re.MatchString(text)
}

// checkProxyServer checks if it work
func (a *Application) checkProxyServer(ip string, port string, c chan QueryResp) {
	start_ts := time.Now()
	var timeout = time.Duration(a.config.ProxyCheckTimeout * time.Second)
	host := fmt.Sprintf("%s:%s", ip, port)
	url_proxy := &url.URL{Host: host}
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(url_proxy)},
		Timeout:   timeout}
	resp, err := client.Get("http://err.taobao.com/error1.html")
	if err != nil {
		c <- QueryResp{Addr: host, Time: float64(-1)}
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	time_diff := time.Now().UnixNano() - start_ts.UnixNano()
	if strings.Contains(string(body), "baidu.com") {
		c <- QueryResp{Addr: host, Time: float64(time_diff) / 1e9}
	} else {
		c <- QueryResp{Addr: host, Time: float64(-1)}
	}
}
