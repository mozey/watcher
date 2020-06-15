# watcher

CLI wrapper for [fsnotify](https://github.com/fsnotify/fsnotify)

## Quick start

Install

    go get github.com/mozey/watcher

Update

    go get -u github.com/mozey/watcher
    go install github.com/mozey/watcher

Watch files, and print debug logs

    APP_DEBUG=true APP_DIR=$(pwd) $GOPATH/bin/watcher -r -dir testdata
    
Watch files, only output changes

    APP_DIR=$(pwd) $GOPATH/bin/watcher -r -dir testdata

Run inside module, and print debug logs

    APP_DEBUG=true APP_DIR=$(pwd) go run ./main.go -r \
    -dir testdata \
    -include ".*.txt$" \
    -include ".*.json$" \
    -excludeDir ".*exclude.*"
    
Print version

    APP_DEBUG=true APP_DIR=$(pwd) go run ./main.go -version
