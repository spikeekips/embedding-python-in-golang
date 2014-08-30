## Lock the python-related code with `sync.Mutex`

```
$ cd src/github.com/spikeekips/embedding-python-in-golang
$ make sync_mutex.go
2014/08/29 20:51:21 python embedding test in golang.
2014/08/29 20:51:21 run inside goroutine
2014/08/29 20:51:21 got json string:
[
  [
    "oldboy",
    "나는 오대수"
  ],
  {
    "I want to eat something alive.": "살아 있는 것을 먹고 싶다.",
    "who are you?": "누구냐, 넌?"
  }
]

2014/08/29 20:51:21 run outside goroutine
2014/08/29 20:51:21 got json string:
[
  [
    "oldboy",
    "나는 오대수"
  ],
  {
    "I want to eat something alive.": "살아 있는 것을 먹고 싶다.",
    "who are you?": "누구냐, 넌?"
  }
]
```

`sync_mutex/main.go` just encoding golang data in python and decoding python data in
golang. It locks the `embed_function` using `sync.Mutex` and `embed_function`
`import`s `json_dump.py` using `go-python`, so simple.

This is `json_dump.py`,

```
import json

def run (*a) :
    return json.dumps(a, )


```

Everything was fine whether goroutine or not.



