all: build

clean:
	-rm ./svctop ./svctop.linux
	-docker rmi svctop:latest

local:
	go build -o svctop

run_local: local
	./svctop

build_linux:
	GOOS=linux GOARCH=amd64 go build -o svctop.linux

build: build_linux
	docker build -t svctop:latest .

deploy: clean build
	helm upgrade --install svctop ./helm

undeploy:
	helm uninstall svctop

reset: undeploy clean
