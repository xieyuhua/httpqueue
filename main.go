package main

import (
    "fmt"
	"flag"
	"time"
	"httpqueue/queue"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/pprof"
)

var (
	redisURL = flag.String("redis", "redis://root:123456@127.0.0.1:6379/3", "redis address")
	address  = flag.String("address", ":2356", "serve listen address")
)

func Run() {
    addr := ":12800"
	fmt.Printf("run status http server %s \n", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("listen stat addr %s err %v", addr, err)
		return
	}

	srv := http.Server{}
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	srv.Handler = mux
	srv.Serve(l)
}


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
	
	done := make(chan struct{}, 1)
	go func() {
		Run()
		done <- struct{}{}
	}()
	
	queue.RunWorker()
	fmt.Println("server listen on :", *address)
	err = queue.ListenAndServe(*address)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
