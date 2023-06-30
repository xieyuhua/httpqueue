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
	redisURL = flag.String("redis", "redis://:895623@127.0.0.1:6379/3", "redis address")
	address  = flag.String("address", ":2356", "serve listen address")
	port  = flag.String("port", ":12860", "listen stat addr")
)

func Run() {
    
	fmt.Printf("run status http server %s \n", *port)
	l, err := net.Listen("tcp", *port)
	if err != nil {
		fmt.Printf("listen stat addr %s err %v", *port, err)
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
        // path+".%Y%m%d%H%M",
        rotatelogs.WithLinkName(path),
        rotatelogs.WithRotationTime(24*time.Hour),
        // rotatelogs.WithRotationTime(time.Minute),
        rotatelogs.WithRotationCount(30),
        rotatelogs.WithRotationSize(100*1024*1024),
    )
    log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
    log.SetLevel(log.InfoLevel)
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
