setup_db:
	docker-compose exec -u root db bash -c "echo \"host all all 192.168.17.0/24 trust\" >> /var/lib/postgresql/data/pg_hba.conf"

create_db:
	docker-compose exec db createdb --username=postgres --owner=postgres trackers

#make PASS=example migrate
migrate:
	docker-compose run api migrate -path db/migrations -database "postgres://postgres:$(PASS)@db:5432/trackers?sslmode=disable" -verbose up

#make TABLE=create_users_table migrate
migration-create:
	docker-compose run api migrate create -ext sql -dir db/migrations $(TABLE)

#make PASS=example migrate-down
migrate-down:
	docker-compose run api migrate -path db/migrations -database "postgres://postgres:$(PASS)@db:5432/trackers?sslmode=disable" -verbose down

# build:
# 	docker-compose exec tracker go build /var/www/html/src/index.go

# run: build
# 	docker-compose exec tracker go run /var/www/html/src/index.go