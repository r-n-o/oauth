# POC for Google OIDC

Basic integration with Google OIDC to understand things in a bit more details.

Seeded from [this POC](https://github.com/BNPrashanth/poc-go-oauth2/tree/google), updated to work with [go-oidc](https://github.com/coreos/go-oidc)

## Starting the app

Create your own config file:
```
cp config.example.yml config.yml
```

```
go run ./cmd/poc/main.go
```
