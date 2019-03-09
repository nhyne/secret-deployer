VERSION=$(shell cat VERSION)
RELEASE_TYPE:=patch
build:
	go build -o build/secret-deployer ./main.go

docker:
	GOOS=linux GOARCH=amd64 go build -o build/linux-secret-deployer ./main.go
	docker build -t nhyne/secret-deployer:${VERSION} .

release: docker
	./release.sh ${RELEASE_TYPE}
	docker push nhyne/secret-deployer:${VERSION}
