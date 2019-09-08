arm:
	CGO_ENABLED=0 GOARM=7 GOOS=linux GOARCH=arm  go build -v -o ipproxypool -a -ldflags "-s -w" main.go
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -v -o ipproxypool -a -ldflags "-s -w" main.go
windows32:
	CGO_ENABLED=0 GOOS=windows GOARCH=386  go build -v -o ipproxypool.exe -a -ldflags "-s -w" main.go  
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build -v -o ipproxypool.exe -a -ldflags "-s -w" main.go  
dev:
	go build main.go
dockerpush:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -v -o docker/ipproxypool -a -ldflags "-s -w" main.go
	cd docker && docker build -t="suconghou/ipproxypool" . && docker images && docker push suconghou/ipproxypool