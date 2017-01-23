PROFILING_FOLDER = profiling

default: test

ci: get lint test

get:
	go get -u github.com/golang/lint/golint
	go get -t -v ./...

lint:
	golint

test:
	go test -check.vv -cover ./...

bench:
	mkdir -p $(PROFILING_FOLDER)
	go test -check.vv -check.b -outputdir $(PROFILING_FOLDER) -cpuprofile cpu.pprof -memprofile memory.pprof $(RUN)
	mv ed448.test $(PROFILING_FOLDER)
	go tool pprof -top -output=$(PROFILING_FOLDER)/cpu-top.txt $(PROFILING_FOLDER)/ed448.test $(PROFILING_FOLDER)/cpu.pprof
	go tool pprof -top -output=$(PROFILING_FOLDER)/mem-top.txt $(PROFILING_FOLDER)/ed448.test $(PROFILING_FOLDER)/memory.pprof

clean:
	rm -rf $(PROFILING_FOLDER)
