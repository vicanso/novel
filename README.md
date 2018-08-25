# novel
web framework by iris

## test

```bash
GO_ENV=test VIPER_INIT_TEST=true go test -race -coverprofile=test.out ./... && go tool cover --html=test.out
```