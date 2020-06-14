# watcher

CLI wrapper for [fsnotify](https://github.com/fsnotify/fsnotify)

## Quick start

    APP_DEBUG=true APP_DIR=/Users/mozey/pro/watcher go run ./main.go -r \
    -dir testdata \
    -include ".*.txt$" \
    -include ".*.json$" \
    -excludeDir ".*exclude.*"