centralized socket based chatroom

- Both socket and websocket based chatroom are supported
- Multi-users in a chatroom
- Receive and display chat history once login
- Asynchronized between message sending and broadcasting
- Timestamp for each chat content
- OpenID should be unique in a chatroom
- Leave and re-connect with the same openID is allowed

Usage:

- server:
	- cd ./chatserver
	- go get -u golang.org/x/net
	- go build
	- ./chatserver -h localhost -p 1024
- client:
	- cd ./chatclient
	- go build
	- ./chatclient -h localhost -p 1024 -n kewei // Join chatroom with openID "kewei"
