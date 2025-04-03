Clone project

```git clone https://github.com/ogabekkadirov/oauth-server.git```

Make .env file
```cp deploy/example-env.env deploy/.env``` or ```make .env```

```go mod init github.com/ogabekkadirov/oauth-server```

```go mod tidy```

Run project
```make compose-up``` or ```docker-compose -f ./deploy/docker-compose.yml up -d```

Migrate
```make migrateup```

Make new migration
```make migration ${tablename}```

Run seeder
```make seed```

Get token (grant_type:password)

```
curl --request POST \
  --url http://localhost:3030/api/v1/oauth/token \
  --header 'Content-Type: application/json' \
  --data '{
	"grant_type":"password",
	"client_id":"client2",
	"client_secret":"secret2",
	"username":"admin",
	"password":"admin"
}'
```


Get token (grant_type:client_credentials)
```
curl --request POST \
  --url http://localhost:3030/api/v1/oauth/token \
  --header 'Content-Type: application/json' \
  --data '{
	"grant_type":"client_credentials",
	"client_id":"client1",
	"client_secret":"secret1"
}'
```

Get token (grant_type:refresh_token)
```
curl --request POST \
  --url http://localhost:3030/api/v1/oauth/token \
  --header 'Content-Type: application/json' \
  --data '{
	"grant_type":"refresh_token",
	"client_id":"client1",
	"client_secret":"secret1",
	"refresh_token":"${token}"
}'
```

Get auth info
```
curl --request GET \
  --url http://localhost:3030/api/v1/me \
  --header 'Authorization: Bearer ${token}'
```
