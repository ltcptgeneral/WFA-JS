.PHONY: build clean test

build: clean
	@echo "======================== Building Binary ======================="
	GOOS=js GOARCH=wasm CGO_ENABLED=0 tinygo build -no-debug -opt=2 -target=wasm -o dist/wfa.wasm .

clean:
	@echo "======================== Cleaning Project ======================"
	go clean
	rm -f dist/wfa.wasm cover.prof cpu.prof mem.prof test.test

test:
	@echo "======================== Running Tests ========================="
	go test -v -cover -coverpkg=./pkg/ -coverprofile cover.prof -cpuprofile cpu.prof -memprofile mem.prof ./test/
	@echo "======================= Coverage Report ========================"
	go tool cover -func=cover.prof
	@rm -f cover.prof
	@echo "==================== CPU Performance Report ===================="
	go tool pprof -top cpu.prof
	@rm -f cpu.prof
	@echo "=================== Memory Performance Report =================="
	go tool pprof -top mem.prof
	@rm -f mem.prof

	@rm -f test.test