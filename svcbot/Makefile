NAMESPACE ?= default
SVC_NAME = svcbot

all: build

clean_local:
	-rm ./$(SVC_NAME)

clean: clean_local
	-rm ./$(SVC_NAME).linux
	-docker rmi $(SVC_NAME):latest

local: clean_local
	go build -o $(SVC_NAME)

run_local: local
	./$(SVC_NAME)

build_linux:
	GOOS=linux GOARCH=amd64 go build -o $(SVC_NAME).linux

build: build_linux
	docker build -t $(SVC_NAME):latest .

deploy: clean build
	helm -n $(NAMESPACE) upgrade --install $(SVC_NAME) ./helm

undeploy:
	helm -n $(NAMESPACE) uninstall $(SVC_NAME)

reset: undeploy clean

list:
	telepresence -n $(NAMESPACE) list

intercept:
	telepresence -n $(NAMESPACE) intercept $(SVC_NAME) --port 8080:http

leave:
	telepresence leave $(SVC_NAME)-$(NAMESPACE)
