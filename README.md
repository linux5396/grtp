# grtp
A go routine pool .Although the cost of go routine create is very less,using the go routine pool can solve some special scenes.

## quick start

```shell
go get -u github.com/linux5396/grtp
```

import

```go
import "github.com/linux5396/grtp"
```



- async

```go
	p, _ := grtp.NewAsync(8, 40)//create async queue.
	go p.Run()
	task := grtp.NewTask(func() {
		//action
	})
	p.Commit(task)
```

- sync

```go
	p, _ := grtp.New(8)
	go p.Run()
	task := grtp.NewTask(func() {
		//action
	})
	p.Commit(task)
```

----

## FAQ

