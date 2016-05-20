default: test

ci: get test bench

get:
	go get -t -v ./...

test:
	go test -check.vv -cover ./...

test-32:
	GOARCH=386 go test -check.vv -cover ./...

bench:
	go test -check.vv -check.b

bench-32:
	GOARCH=386 go test -check.vv -check.b

