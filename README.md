# ChatSystem
This is a simple chat system which build by Golang. 

# How to use it
On your Linux(or Mac) system, create a server by these two command:

```
$ go build -o server main.go server.go user.go
$ ./server
```

Create a client by:
```
$ go build -o client client.go`
$ ./client
```

You can open serveral client on your terminals

# Function
Once you open a client program, you can do these things:
* Public Chat
* Private Chat
* Rename
* Quit

Press digit to determine what you want to do.
