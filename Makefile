# Builds executables
build:
	rm -rf ./dist/bserver ./dist/bsclient
	mkdir -p ./dist/bserver ./dist/bsclient
	go build -o ./dist/bserver ./cmd/bserver
	go build -o ./dist/bsclient ./cmd/bsclient

# Builds for distribution
build-dist: build
	rm -rf ./dist/build/client ./dist/build/server
	mkdir -p ./dist/build/client/ ./dist/build/server/
	cp -r ./dist/bsclient/bsclient ./dist/build/client
	cp -r ./dist/bserver/bserver ./dist/build/server
	cp -r ./scripts/install/install-server.sh ./dist/build/server
	cp -r ./systemd/bserver.service ./dist/build/server

docker-build:
	docker build -t broadcast .
docker-run: 
	docker run --rm -it broadcast:latest