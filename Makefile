.PHONY: build clean test

build: clean
	@echo "======================== Building Binary ======================="
	GOOS=js GOARCH=wasm CGO_ENABLED=0 tinygo build -no-debug -opt=2 -target=wasm -o dist/wfa.wasm .

clean:
	@echo "======================== Cleaning Project ======================"
	go clean
	rm -f dist/wfa.wasm

test:
	@echo "======================== Running Tests ========================="
	go test -v -cover -coverpkg=./pkg/ -coverprofile coverage ./test/
	@echo "======================= Coverage Report ========================"
	go tool cover -func=coverage
	@rm -f coverage