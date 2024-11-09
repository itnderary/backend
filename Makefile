BINARY_NAME=itnderary
KO_DOCKER_REPO=europe-west3-docker.pkg.dev/itnderary/backend
KO_LOCAL_DOCKER_REPO=ko.local/jajcoszek/itnderary-backend
DOCKER_RUN_IMAGE_SUFFIX=v2:latest

build: ${BINARY_NAME}

test:
	go test

run: package
	docker run -p 8080:8080 ${KO_LOCAL_DOCKER_REPO}/${DOCKER_RUN_IMAGE_SUFFIX}

${BINARY_NAME}: *.go
	go build -o ${BINARY_NAME}

package: build
	KO_DOCKER_REPO=${KO_LOCAL_DOCKER_REPO} ko build -B

publish: build
	KO_DOCKER_REPO=${KO_DOCKER_REPO} ko build -B

clean:
	go clean
	rm ${BINARY_NAME}

.PHONY: run build package publish clean
