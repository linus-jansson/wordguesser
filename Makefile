BIN_DIR := bin
TARGET := $(BIN_DIR)/wordguesser
SRC_DIR := wordguesser-go

.PHONY: all build clean

all: build

build:
	mkdir -p $(BIN_DIR)
	cd $(SRC_DIR) && go build -o ../$(TARGET)

clean:
	rm -rf $(BIN_DIR)