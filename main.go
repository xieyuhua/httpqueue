package main

import (
    "fmt"
	"flag"
	"time"
	"httpqueue/queue"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

var (
	redisURL = flag.String("redis", "redis://root:123456@127.0.0.1:6379/3", "redis address")
	address  = flag.String("address", ":2356", "serve listen address")
)


func init() {
    path := "log/mylog.log"
    writer, _ := rotatelogs.New(
        path+".%Y%m%d",
        rotatelogs.WithLinkName(path),
        rotatelogs.WithMaxAge(time.Duration(180)*time.Second),
        rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
    )
    log.SetOutput(writer)
    // log.SetReportCaller(true) 
	flag.Parse()
}

func main() {
	err := queue.InitRedis(*redisURL)
	if err != nil {
		log.Fatal(err)
	}
	
	
	queue.RunWorker()
	fmt.Println("server listen on :", *address)
	err = queue.ListenAndServe(*address)
	if err != nil {
		log.Fatal(err)
	}
}
