VERSION=$(shell cat VERSION)
RELEASE_TYPE=patch

install:
	go install ./main.go

build_binary:
	go build -o build/secret-deployer ./main.go

docker:
	GOOS=linux GOARCH=amd64 go build -o build/linux-secret-deployer ./main.go
	docker build -t nhyne/secret-deployer:${VERSION} .
	docker push nhyne/secret-deployer:${VERSION}

release: install
	./release.sh ${RELEASE_TYPE}
