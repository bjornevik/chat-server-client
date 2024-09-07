# Chat Server-Client

## Run in Docker
```sh
docker-compose build
docker-compose up
```

To connect to the docker instances and run the applications:
```sh
docker exec -it server /app/server
docker exec -it client1 /app/client
docker exec -it client2 /app/client
docker exec -it client3 /app/client
```
Then you can chat from there. If the server is taken offline the connected clients
will attempt to reconnect every 5 seconds for a minute before exiting themselves.

If you'd like to inspect the chat-logs:
```sh
docker exec -it server /bin/sh
cat server_log.txt
```

## Run locally
If you'd rather run locally
```sh
cd server # or client
go run main.go
```
OR
```sh
cd server # or client
go build
./server # or client
```
