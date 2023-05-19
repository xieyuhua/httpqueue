package queue

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
    "fmt"
	log "github.com/sirupsen/logrus"
)

var httpClient = &http.Client{
	Timeout: CallbackTTR,
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
	resp, err := httpClient.Post(task.Callback, "application/json", content)
	if err != nil {
	    fmt.Println("task fail task-id:", task.ID)
		log.WithError(err).Error("http post fail task-id:" + task.ID)
		return 0, err
	}
	
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("result : %s => %s \n",task.ID, result)
	log.Infof("result : %s => %s",task.ID, result)
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
