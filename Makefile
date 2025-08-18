# Makefile for SoJSON

BINARY_NAME=sojson

# 默认目标
.PHONY: all
all: build

# 构建二进制文件
.PHONY: build
build:
	@echo "构建 $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) main.go
	@echo "构建完成: $(BINARY_NAME)"

# 清理构建产物
.PHONY: clean
clean:
	@echo "清理构建文件..."
	rm -f $(BINARY_NAME)
	@echo "清理完成"

# 部署到 supervisor
.PHONY: deploy
deploy:
	@echo "开始部署..."
	./deploy.sh
