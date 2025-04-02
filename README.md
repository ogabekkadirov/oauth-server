Clone project

```git clone https://github.com/ogabekkadirov/oauth-server.git```

Make .env file
```cp deploy/example-env.env deploy/.env``` or ```make .env```

```go mod tidy```

Run project
```make compose-up``` or ```docker-compose -f ./deploy/docker-compose.yml up -d```

Migrate
```make migrateup```

Make new migration
```make migration ${tablename}```

Run seeder
```make seed```

