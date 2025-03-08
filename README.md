# Poker

## Run tests
```
$ go test ./...
```

## Build and run project
```
$ go build
$ ./poker
```

In a new terminal, in the root directory of the repo:
```
$ python3 -m http.server -d . 8082
```

Visit `http://localhost:8082` in a web browser, open multiple tabs and join the same table to simulate multiple players.
