

httpqueue is a redis base delay queue

todo 
```
多节点 服务负载均衡
redis，etcd 集群  主主复制同步
日志对接es,ck,mogodb
```

```
#!/bin/bash
binName="bin/httpqueue"
GOOS=linux GOARCH=amd64 go build -o "$binName"_linux
GOOS=darwin GOARCH=amd64 go build -o "$binName"_macos
GOOS=freebsd GOARCH=amd64 go build -o "$binName"_freebsd
GOOS=linux GOARCH=arm go build -o "$binName"_arm
GOOS=windows GOARCH=amd64 go build -o "$binName"_win64.exe
```

### Usege
golang version: 1.9+
```
$: ./httpqueue -h
Usage of ./httpqueue:
  -address string
    	serve listen address (default ":2356")
  -redis string
    	redis address (default "redis://:895623@127.0.0.1:6379/3")
    	
    	
:12800/debug/pprof/

```



### Frontend API

Response http code: **200** success, **400** request invalid, **404** task not found, **500** internal error

- Create Task

  ```
  Request:
  POST /create
  {
  	"topic":"order",
  	"delay":15, // second
  	"retry":3,  // max retry 3 times, interval 10,20,40... seconds
  	"callback":"http://127.0.0.1:8888/", // http post to target url
  	"content":"hello" // content to post
  }
  Response:
  {
      "id": "35adbde5-77c4-4d65-adac-0082d91f2554"
  }
  ```

- Delete Task

  ```
  Request:
  POST /delete
  {
  	"id":"35adbde5-77c4-4d65-adac-0082d91f2554"
  }
  ```

- Query Task

  ```
  Request:
  POST /query
  {
  	"id":"35adbde5-77c4-4d65-adac-0082d91f2554"
  }
  Response:
  {
      "id": "cb9aefdd-5bd1-4bf3-8c94-1ed5c2ea638e",
      "topic": "order",
      "execute_time": 1504934230,
      "max_retry": 1,
      "has_retry": 0,
      "callback": "http://127.0.0.1:8888/success",
      "content": "hello",
      "creat_time": 1504934220
  }
  ```

## Backend API

- Callback

  ```
  Request:
  POST /?
  {
    "id": "57e177ff-454c-42d6-93ab-65895b950dbf",
    "topic": "order",
    "content": "hello"
  }
  Response:
  {
      "code":100 // 100: success,101: too many request,other: fail
  }
  ```
  
```
[root@iZjf2a9cZ httpqueue]# ./httpqueue 
server listen on : :8080

createTask :{"ID":"queue:164c9e82-d3bd-44c1-b1f5-de578c4b2fd8","Topic":"order","ExecuteTime":1684488198,"MaxRetry":3,"HasRetry":0,"Callback":"http://###/reg/index","Content":"hello","CreatTime":1684488198} 
createTask :{"ID":"queue:034a4675-33ea-42ca-8269-3b898bb4de7e","Topic":"order","ExecuteTime":1684488198,"MaxRetry":3,"HasRetry":0,"Callback":"http:://###/reg/index","Content":"hello","CreatTime":1684488198} 
createTask :{"ID":"queue:71819883-3c57-4360-a7a8-625ae3906ad8","Topic":"order","ExecuteTime":1684488198,"MaxRetry":3,"HasRetry":0,"Callback":"http://###/reg/index","Content":"hello","CreatTime":1684488198} 
2023-05-22 09:43:24 result : orer:91ad86a4-f86f-4ba9-ab7f-4d6881f4590e => {"code":100} ,cost time 59.882415ms 
2023-05-22 09:43:24 result : orer:b4bd4648-3644-4033-8c3e-ed939455c42a => {"code":100} ,cost time 61.810333ms 
2023-05-22 09:43:24 result : orer:1a5839c0-880c-494e-bf65-3d549b6de7a1 => {"code":100} ,cost time 168.571621ms 


```
