all: build

clean_local:
	-rm ./svcmid

clean: clean_local
	-rm ./svcmid.linux
	-docker rmi svcmid:latest

local: clean_local
	go build -o svcmid

run_local: local
	./svcmid

build_linux:
	GOOS=linux GOARCH=amd64 go build -o svcmid.linux

build: build_linux
	docker build -t svcmid:latest .

deploy: clean build
	helm upgrade --install svcmid ./helm

undeploy:
	helm uninstall svcmid

reset: undeploy clean
