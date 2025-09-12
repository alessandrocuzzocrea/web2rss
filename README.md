# web2rss

migrate create -ext sql -dir db/migrations -seq create_xxxx_table

migrate -database "sqlite3://data/web2rss.sqlite3" -path db/migrations up

go tool sqlc generate
