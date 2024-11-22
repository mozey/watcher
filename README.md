# watcher

CLI wrapper for [fsnotify](https://github.com/fsnotify/fsnotify)

You probably want to use one of the following tools instead
- [air](https://github.com/air-verse/air)
- [wgo](https://github.com/bokwoon95/wgo)

The WatcherParams include a `Change` callback that's useful when using [Mage](https://magefile.org/) instead of a Makefile


## Quick start

Install
```bash
# Since Go 1.20.3 
# "'go get' is no longer supported outside a module"
go install github.com/mozey/watcher@latest
```

Print version
```bash  
$GOPATH/bin/watcher -version
```

Watch files, only output changes
```bash
$GOPATH/bin/watcher -dir testdata
```

Watch files recursively, and print debug logs
```bash
APP_DEBUG=true $GOPATH/bin/watcher -r -dir testdata
```


## Testing

Recursively watch for change in both `./testdata` and `./testdata2`
```bash
# APP_DIR env sets base dir
APP_DEBUG=true APP_DIR=$(pwd) go run ./main.go -r -dir testdata -dir testdata2
# Base dir set from flag
APP_DEBUG=true go run ./main.go -r -b $APP_DIR -dir testdata -dir testdata2
```

Using absolute path
```bash
cd ${PRO_PATH}/watcher # Base dir defaults to working dir
APP_DEBUG=true go run ./main.go -r \
    -dir testdata \
    -dir ${PRO_PATH}/watcher/testdata2
```

Example with filters
```bash
APP_DEBUG=true APP_DIR=$(pwd) go run ./main.go -r -dir testdata \
-include ".*.txt$" \
-include ".*.json$" \
-excludeDir ".*exclude.*" \
-exclude ".*\/d.txt$"
```

**TODO** See comments in main_test.go
