include .env

app-deploy:
	set -a
	go mod vendor
	docker compose -f deploy/app/app-composer.yaml up --force-recreate -d
down:
	set -a
	docker compose -f deploy/app/app-composer.yaml down
destroy:
	set -a
	docker compose -f deploy/app/app-composer.yaml down
	docker rmi spiral_app:latest
vendor:
	go mod vendor
test:
	go clean -testcache
	go test -v ./...

