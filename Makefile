BIN_DIR := bin
TARGET := $(BIN_DIR)/wordguesser
SRC_DIR := wordguesser-go

.PHONY: all build clean

all: build

build:
	mkdir -p $(BIN_DIR)
	gzip -c words/valid-words.json > wordguesser-go/valid-words.json.gz
	cd $(SRC_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../$(TARGET)
	rm wordguesser-go/valid-words.json.gz

clean:
	rm -rf $(BIN_DIR)