default: build

build: vet
	go build -v -o ./dist/unbalance ./server/boot.go

buildx: vet
	GOOS=linux GOARCH=amd64 go build -v -o ./dist/unbalance ./server/boot.go

clean:
	rm -rf ./dist

install: clean buildx
	cp -r client/* dist
	scp -pr ./dist/* wopr:/boot/custom/unbalance

run: clean build
	cp -r client/ dist
	cd dist && ./unbalance

dev: clean build
	cp -r client/* dist
	cd dist && http-server

vet:
	go vet ./server/...