Stub Server is the easiest way to create a server stubs

How to get

```
#!shell

go get bitbucket.org/filipenos/stub-server
```

How to run

```
#!shell

stub-server -port=8090 -conf=config.json
```

Example of config.json

```
#!json

[{
  "method": "GET",
  "path": "/api/stub",
  "json": "dir_of_json"
}]
```
