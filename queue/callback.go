package queue

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
    "fmt"
    "time"
	log "github.com/sirupsen/logrus"
)

var httpClient = &http.Client{
	Timeout: Timeout,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: MaxIdleConnsPerHost,
		MaxIdleConns:        MaxIdleConns,
		IdleConnTimeout:     IdleConnTimeout,
	},
}

type callbackRequest struct {
	ID      string `json:"id"`
	Topic   string `json:"topic"`
	Content string `json:"content"`
}

type callbackResponse struct {
	Code int `json:"code"`
}

const (
	CodeSuccess        = 100
	CodeTooManyRequest = 101
)

func post(task *Task) (int, error) {
	request := callbackRequest{
		ID:      task.ID,
		Topic:   task.Topic,
		Content: task.Content,
	}
	data, err := json.Marshal(request)
	if err != nil {
	    fmt.Println("task fail task-id:", task.ID)
		log.WithError(err).Error("json marshal fail task-id:" + task.ID)
		return 0, err
	}

	content := bytes.NewBuffer(data)
	start := time.Now()
	resp, err := httpClient.Post(task.Callback, "application/json", content)
	end   := time.Now()
	if err != nil {
	    fmt.Println("task fail task-id:", task.ID)
		log.WithError(err).Error("http post fail task-id:" + task.ID)
		return 0, err
	}
	
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	date := time.Now().Format("2006-01-02 15:04:05")
	transferTime := end.Sub(start)
// 	fmt.Println(transferTime)
	fmt.Printf("%s result : %s => %s ,cost time %s \n", date, task.ID, result, transferTime)
	log.Infof("result : %s => %s ,cost time %s",task.ID, result ,transferTime)
	if err != nil {
		log.WithError(err).Error("io read from backend fail task-id:" + task.ID)
		return 0, err
	}
	var response callbackResponse
	err = json.Unmarshal(result, &response)
	if err != nil {
		log.WithError(err).Error("result json unmarshal fail task-id:" + task.ID)
		return 0, err
	}
	return response.Code, nil
}
