test:
	go test ./...
	GOOS=windows go build
	GOOS=linux GOARCH=amd64 go build
	GOOS=linux GOARCH=mips64 go build
	GOOS=linux GOARCH=arm64 go build
	#GOOS=linux GOARCH=sw64 go build
