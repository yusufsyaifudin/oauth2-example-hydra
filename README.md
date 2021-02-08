# Sample implementation of OAUTH2 using Hydra

This repository will explain you about:

* How to create an Authentication Server and then use ory/hydra as Authorization server
* How to implement "Sign With {put-your-app-name}" for single identity management for your app. This exactly the same as building "Sign Up using Google" or similar services for your own app.

## How to run, step by step

Clone this repository. Then run `go mod download`.

### Run ory/hydra server 

This will run using Postgres as database and Jaeger as a tracing vendor.

```shell script
$ docker-compose up
```

### Create an OAUTH2 client

This will create an OAUTH2 client using `hydra` cli. Otherwise, you can also create using REST API.

```shell script
$ docker-compose exec hydra hydra clients create \
--endpoint http://127.0.0.1:4445/ \
--name MyApp \
--id myclient \
--secret mysecret \
--grant-types authorization_code,refresh_token \
--response-types code,id_token \
--callbacks http://localhost:1234/callbacks \
--token-endpoint-auth-method client_secret_post \
--scope offline,users.write,users.read,users.edit,users.delete
```

or

```shell script
curl -X POST 'http://localhost:4445/clients' \
-H 'Content-Type: application/json' \
--data-raw '{
    "client_id": "myclient",
    "client_name": "MyApp",
    "client_secret": "mysecret",
    "grant_types": ["authorization_code", "refresh_token"],
    "redirect_uris": ["http://localhost:1234/callbacks"],
    "response_types": ["code", "id_token"],
    "scope": "offline users.write users.read users.edit users.delete",
    "token_endpoint_auth_method": "client_secret_post"
}'
```

Please note `scope` is separated by space.

> If you made wrong request params, you cannot create with the same client id,
> You must delete it first by doing this:
>
>
> ```docker-compose exec hydra hydra clients delete --endpoint http://127.0.0.1:4445/ <client-id>```
> 
> or
>
> ```curl -X DELETE 'http://localhost:4445/clients/<client-id>'```


### Run Identity Provider (Resource Server)

```shell script
$ go run cmd/authc/main.go
```

### Run OAuth 2.0 Client App

```shell script
$ REDIRECT_URL=http://localhost:1234/callbacks CLIENT_ID=myclient CLIENT_SECRET=mysecret go run cmd/frontend/main.go
```

### Access the client App

Access http://localhost:1234 then try to authorize the application.

## Include HTML in Binary

```
$ go get -u github.com/shuLhan/go-bindata/...
$ go-bindata -o views/view.go -ignore=view.go views/...
```