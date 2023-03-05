all: build

clean:
	-rm ./svcbot ./svcbot.linux
	-docker rmi svcbot:latest

local:
	go build -o svcbot

run_local: local
	./svcbot

build_linux:
	GOOS=linux GOARCH=amd64 go build -o svcbot.linux

build: build_linux
	docker build -t svcbot:latest .

deploy: clean build
	helm upgrade --install svcbot ./helm

undeploy:
	helm uninstall svcbot

reset: undeploy clean
