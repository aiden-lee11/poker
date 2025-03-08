# Poker

## Run tests
```
$ cd backend
$ go test ./...
```

## Build and run project
```
$ cd backend
$ go build
$ ./poker
```

In a new terminal:
```
$ cd backend
$ python3 -m http.server -d . 8082
```

Visit `http://localhost:8082` in a web browser, open multiple tabs and join the same table to simulate multiple players.
