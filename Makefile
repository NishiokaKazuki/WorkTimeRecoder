docker-up:
	cd docker &&\
	docker-compose up -d

docker-init-db:
	cd docker &&\
	docker-compose exec db /bin/bash -psecret -c "chmod 0775 docker-entrypoint-initdb.d/init-db.sh" &&\
	docker-compose exec db /bin/bash -psecret -c "sh ./docker-entrypoint-initdb.d/init-db.sh"

docker-stop:
	cd docker &&\
	docker-compose stop

docker-rm:
	cd docker &&\
	docker-compose rm -f
	rm -rf mysql/data/*

docker-build:
	cd docker &&\
	docker-compose build

docker-ps:
	cd docker &&\
	docker-compose ps

docker-exec-db:
	docker exec -it WorkTimeRecoder-db /bin/bash

docker-exec-server:
	docker exec -it WorkTimeRecoder-server /bin/bash

run-server:
	cd server &&\
	go run main.go