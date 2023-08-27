db-up:
	docker compose -f deploy/db/pg-composer.yaml up --force-recreate -d
db-down:
	docker compose -f deploy/db/pg-composer.yaml down 
app-up:
	docker compose -f deploy/db/pg-composer.yaml up --force-recreate -d
app-down:
	docker compose -f deploy/db/pg-composer.yaml down
app-destroy:
	docker compose -f deploy/db/pg-composer.yaml down
	docker rmi spiral_app:latest
