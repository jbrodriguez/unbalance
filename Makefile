default: build

build: vet
	go build -v -o ./dist/unbalance ./server/boot.go

clean:
	rm -rf ./dist

run: clean build
	cp -r client/ dist
	cd dist && ./unbalance

vet:
	go vet ./server/...