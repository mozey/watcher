# logutil

Utils for use with [zerolog](https://github.com/rs/zerolog)


## Quick start

Setup logger with console writer (human readable)
```go
logutil.SetupLogger(true)
```

Log new error with stack trace
```go
err := errors.Errorf("testing")
// Msg must be empty string, if specified it overrides the error message
log.Error().Stack().Err(err).Msg("")
```

Add stack trace to an existing error
```go
err := errors.WithStack(err)
```

When logging to json a stack trace will also be included    


## Multiple writers

Setup additional writers like this, see examples in `log_test.go`
```go
logutil.SetupLogger(true, w1, w2)
```


## Testing

Run all tests
```bash
cd ${GOPATH}/src/github.com/mozey/logutil
gotest -v ./...
```