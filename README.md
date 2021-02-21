# watcher

CLI wrapper for [fsnotify](https://github.com/fsnotify/fsnotify)

## Quick start

Install

    go get github.com/mozey/watcher

Update

    go get -u github.com/mozey/watcher
    go install github.com/mozey/watcher
    
Print version
    
    $GOPATH/bin/watcher -version

Watch files, only output changes

    $GOPATH/bin/watcher -dir testdata
    
Watch files recursively, pass in `APP_DIR`

    APP_DIR=$(pwd) $GOPATH/bin/watcher -r -dir testdata

Watch files recursively, and print debug logs

    APP_DEBUG=true APP_DIR=$(pwd) $GOPATH/bin/watcher -r -dir testdata
    
Use `go run` inside module
    
    APP_DEBUG=true APP_DIR=$(pwd) go run ./main.go -r -dir testdata

...another example with filters

    APP_DEBUG=true APP_DIR=$(pwd) go run ./main.go -r -dir testdata \
    -include ".*.txt$" \
    -include ".*.json$" \
    -excludeDir ".*exclude.*" \
    -exclude ".*\/d.txt$"
    
