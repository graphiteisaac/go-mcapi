# go-mcapi
Super simple HTTP API for basic Minecraft (Java) server info

`v1.1.1`

## Usage

The optional `?{response}` query string param can be `full`, `players`, or `status`. Defaults to `full`

`GET /v1/ping/{ip}?{response}` Ping a server for basic details

`GET /v1/icon/{ip}` Get the PNG server-icon
