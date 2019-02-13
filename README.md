Centralized websocket/socket based chatroom

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
	- go get -u golang.org/x/net/websocket
	- go build
	- ./chatserver -h localhost -p 1024
	- Another choice is just "make run host=[host] port=[port]" under ./chatserver
- client:
	- cd ./chatclient
	- go build
	- ./chatclient -h localhost -p 1024 -n kewei
	- Another choice is just "make run name=[name] host=[host] port=[port]" under ./chatclient
