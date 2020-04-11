## Start server

```bash
go run cmd/web/!(*_test).go
```

## Run tests

```bash
go test ./...

# Run only specific test with regex
go test -v -run="^TestPing$"
```
