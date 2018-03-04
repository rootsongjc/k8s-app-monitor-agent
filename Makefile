.PHONY: build run clean all

export GOOS= linux

export GOARCH=amd64

repo=jimmysong
imageName=k8s-app-monitor-agent
tag=`git rev-parse --short HEAD`
imageWholeName=${repo}/${imageName}:${tag}
port=8888

build:
	go build
	docker build -t jimmysong/k8s-app-monitor-agent:${tag} .

run:
	docker run -d --name ${imageName} -p ${port}:${port} ${imageWholeName}
clean:
	docker rm -f ${imageName}

all: build run
