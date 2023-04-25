GO ?= go
# testing is by subpackages
TESTFOLDER = $(shell find * -maxdepth 10 -type d | grep /tests | xargs -I {} echo "feed-service/{}")
# testing is by files
E2ETESTFOLDER = $(shell find * -maxdepth 10 -type f | grep /e2e)
SOURCE_FILES = $(shell find * -maxdepth 10 | grep .go$)


test:
	bash test.sh

e2e-test:
	bash e2e.sh

transform-struct:
	bash transform-struct.sh

format:
	find services -maxdepth 5 -type f | grep "\.go$$" | xargs -I {} -P8 golines -m 64 -w {} 

dev:
	go run cmd/main/main.go

dev-m1:
	go run --tags dynamic cmd/main/main.go

seed:
	go run cmd/seed/main.go

build.exe:
	go build -o build.exe cmd/main/main.go

transform-typescript:
	bash transform-typescript.sh