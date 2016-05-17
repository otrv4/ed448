default: test

ci: get test

get:
	go get -t -v ./...

test:
	go test -check.vv -cover ./... -check.b

test-32:
	GOARCH=386 go test -check.vv -cover ./... -check.b

