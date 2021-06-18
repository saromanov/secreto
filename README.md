# secreto
Store of secrets

## Environment variables

`REST_SECRET` - 32 symbol of the secret key


## API

### POST /api/secrets

```
{
    "key":"one",
    "value":"two"
}
```

### GET /api/secrets

Params:
key=<key for secret>

### GET /api/keys

#### Response
```
{
    "error":"false",
    "result":"8c6f56de30f34cfd897d813d3731488d"
}
```
