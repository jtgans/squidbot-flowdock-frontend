SRCS := $(wildcard *.go brain/*.go)

all: squidbot-flowdock-frontend-rpi squidbot-flowdock-frontend

squidbot-flowdock-frontend-rpi: $(SRCS)
	GOOS=linux GOARCH=arm go build -v
	mv squidbot-flowdock-frontend squidbot-flowdock-frontend-rpi

squidbot-flowdock-frontend: $(SRCS)
	go build

docker: squidbot-flowdock-frontend
	docker build . --force-rm --tag jtgans/squidbot-flowdock-frontend:latest

docker-rpi: squidbot-flowdock-frontend-rpi
	docker build . -f Dockerfile.rpi --force-rm --tag jtgans/squidbot-flowdock-frontend-rpi:latest

docker-all: docker docker-arm

clean:
	rm -f squidbot-flowdock-frontend squidbot-flowdock-frontend-rpi

mrclean: clean
	find -name \*~ -exec rm -f '{}' ';'

.PHONY: all clean mrclean docker docker-arm docker-all
