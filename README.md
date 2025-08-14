# www2rss

migrate create -ext sql -dir db/migrations -seq create_users_table

go tool sqlc generate
