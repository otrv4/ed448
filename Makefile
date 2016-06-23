PROFILING_FOLDER = profiling

default: test

ci: get lint test bench

get:
	go get -t -v ./...

lint:
	go get -u github.com/golang/lint/golint
	golint

test:
	go test -check.vv -cover ./...

test-32:
	GOARCH=386 go test -check.vv -cover ./...

bench:
	go test -check.vv -check.b

bench-32:
	mkdir $(PROFILING_FOLDER)
	GOARCH=386 go test -check.vv -check.b -outputdir $(PROFILING_FOLDER) -cpuprofile cpu.pprof -memprofile memory.pprof
	mv ed448.test $(PROFILING_FOLDER)
	go tool pprof -top -output=$(PROFILING_FOLDER)/cpu-top.txt $(PROFILING_FOLDER)/ed448.test $(PROFILING_FOLDER)/cpu.pprof
	go tool pprof -top -output=$(PROFILING_FOLDER)/mem-top.txt $(PROFILING_FOLDER)/ed448.test $(PROFILING_FOLDER)/memory.pprof

clean:
	rm -rf $(PROFILING_FOLDER)
