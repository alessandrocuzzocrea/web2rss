# www2rss

migrate create -ext sql -dir db/migrations -seq create_users_table

migrate -database "sqlite3://data/www2rss.sqlite3" -path db/migrations up

go tool sqlc generate
