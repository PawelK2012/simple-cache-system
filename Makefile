build:
	go get && go build -ldflags "-s -w" -trimpath -buildvcs=false -o prci
run:
	go run .
test:
	go test -v ./...