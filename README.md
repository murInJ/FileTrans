# FiTrans
go file transport toolkits

## contributor
MurInJ

## function
- send or receive file

now only surport tcp

## install
```shell
go get github.com/murInJ/FiTrans
```

## quick start

```go
//client
trans, err := FiTrans.NewFiTrans("tcp", 8002, true)
	if err != nil {
		log.Fatal(err)
	}

err = trans.Send("xxx.xxx.xxx.xxx:8002", "./pic.jpg", 10240)
if err != nil {
    log.Fatal(err)
}	
```

```go
//server
trans, err := FiTrans.NewFiTrans("tcp", 8002, true)
	if err != nil {
		log.Fatal(err)
	}
go trans.StartServer()

fileChannel := trans.GetFileChannel()
file := <-fileChannel

f, err := os.Create("./" + file.FileName)
if err != nil {
    return
}
f.Write(file.Data)

trans.CloseServer()
```