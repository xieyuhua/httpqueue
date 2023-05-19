

httpqueue is a redis base delay queue

### Usege
golang version: 1.9+
```
$: ./httpqueue -h
Usage of ./httpqueue:
  -address string
    	serve listen address (default ":2356")
  -redis string
    	redis address (default "redis://root:123456@127.0.0.1:6379/3")
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
result : queue:c81967ca-1211-4f1e-b048-8e3363eef88b => {"code":100} 
result : queue:805e2ac1-eb07-4fad-bb8c-808b257b10f2 => {"code":100} 
result : queue:9163d064-f9b8-471d-9289-30d4b8f1a49c => {"code":100} 

```