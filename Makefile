BINARY_NAME=dsp

build:
	@echo "Building ${BINARY_NAME}..."
	go build -o $(BINARY_NAME) main.go

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)

.PHONY: build clean 