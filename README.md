# watcher

CLI wrapper for [fsnotify](https://github.com/fsnotify/fsnotify)

## Quick start

Install
```sh
go get github.com/mozey/watcher
```

Update
```sh
go get -u github.com/mozey/watcher
go install github.com/mozey/watcher
```

Print version
```sh    
$GOPATH/bin/watcher -version
```

Watch files, only output changes
```sh
$GOPATH/bin/watcher -dir testdata
```

Watch files recursively, pass in `APP_DIR`
```sh
APP_DIR=$(pwd) $GOPATH/bin/watcher -r -dir testdata
```

Watch files recursively, and print debug logs
```sh
APP_DEBUG=true APP_DIR=$(pwd) $GOPATH/bin/watcher -r -dir testdata
```

Use `go run` inside module
```sh
APP_DEBUG=true APP_DIR=$(pwd) go run ./main.go -r -b $APP_DIR -dir testdata -dir testdata2
```

...another example with filters
```sh
APP_DEBUG=true APP_DIR=$(pwd) go run ./main.go -r -dir testdata \
-include ".*.txt$" \
-include ".*.json$" \
-excludeDir ".*exclude.*" \
-exclude ".*\/d.txt$"
```

