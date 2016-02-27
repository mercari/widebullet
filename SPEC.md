# Specification for Widebullet

Widebullet is the proxy server between JSON-RPC and RESTful API server. It accepts a HTTP request based JSON-RPC.

## API

Widebullet has the APIs below.

 * [POST /wbt](#post-wbt)
 * [GET /stat/go](#get-statgo)

### POST /wbt

Accepts a HTTP request based JSON-RPC.

The JSON below is a request-body example.

```json
[
    {"jsonrpc": "2.0", "ep": "ep-1",                        "method": "/user/get",    "params": { "user_id": 1 },                   "id": "1"},
    {"jsonrpc": "2.0", "ep": "ep-1", "http_method": "GET",  "method": "/item/get",    "params": { "item_id": 2 },                   "id": "2"},
    {"jsonrpc": "2.0", "ep": "ep-2", "http_method": "POST", "method": "/item/update", "params": { "item_id": 2, "desc": "update" }, "id": "3"}
]
```

The definitions of parameters are below.

|name            |type  |description                              |required|note                              |
|----------------|------|-----------------------------------------|--------|----------------------------------|
|jsonrpc         |string|version number of JSON-RPC               |o       |fixed as 2.0                      |
|ep              |string|endpoint name                            |o       |selected in Endpoints Section     |
|http_method     |string|method string for HTTP                   |o       |HTTP method string. GET by default|
|method          |string|method string                            |o       |URI                               |
|params          |object|parameters for method                    |o       |                                  |
|id              |string|ID string                                |o       |                                  |


When Widebullet receives an invalid request(for example, malformed body is included), a status of response it returns is 400(Bad Request).

### GET /stat/go

Returns a statictics for golang-runtime. See [golang-stats-api-handler](https://github.com/fukata/golang-stats-api-handler) about details.
