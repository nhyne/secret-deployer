FROM alpine:3.9.2

ADD ./build/linux-secret-deployer /usr/local/bin/secret-deployer

ENTRYPOINT ["secret-deployer"]
