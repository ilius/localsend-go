PROJECT_NAME := localsend-go

SRC_DIR := ./cmd/localsend-go/

OUT_DIR := ./bin

GO := go

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: all
all: clean build

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)

.PHONY: build
build:
	mkdir -p $(OUT_DIR)
	$(GO) build -o $(OUT_DIR)/ $(SRC_DIR)

#$(PLATFORMS):
#	GOOS=$(word 1, $(subst /, ,$@)) GOARCH=$(word 2, $(subst /, ,$@)) \
#	$(GO) build -o $(OUT_DIR)/$(PROJECT_NAME)-$(word 1, $(subst /, ,$@))-$(word 2, $(subst /, ,$@))$(if $(findstring windows,$@),.exe) $(SRC_DIR)

.PHONY: test
test:
	$(GO) test ./...

.PHONY: deps
deps:
	$(GO) mod tidy

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make            - 编译所有平台的可执行文件"
	@echo "  make clean      - 清理输出目录"
	@echo "  make build      - 编译所有平台的可执行文件"
	@echo "  make test       - 运行测试"
	@echo "  make deps       - 安装依赖"
	@echo "  make help       - 显示此帮助信息"
