#!/bin/bash

# =============================================================================
# 脚本名称: app/gen.sh (精准定位版)
# 适用结构: app/<service>/cmd/api/desc/<module>/<file>.api
# =============================================================================

GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

ACTION=$1
API_FILE=$2

# 1. 参数预处理
if [ -f "$1" ] && [ -z "$2" ]; then
    ACTION="code"
    API_FILE=$1
fi

if [ -z "$API_FILE" ] || [ ! -f "$API_FILE" ]; then
    echo -e "${RED}❌ Usage: ./gen.sh <code> <path_to_api_file>${NC}"
    exit 1
fi

# 2. 确定性路径截取
# 逻辑：获取绝对路径后，截取到 "cmd/api" 为止
ABS_API_FILE=$(realpath "$API_FILE")
# 使用模式匹配截掉 cmd/api/ 之后的所有内容，再补回 cmd/api
TARGET_DIR="${ABS_API_FILE%/cmd/api/*}/cmd/api"

if [ ! -d "$TARGET_DIR" ]; then
    echo -e "${RED}❌ Error: 路径中未发现 cmd/api 目录: $ABS_API_FILE${NC}"
    exit 1
fi

# 3. 执行代码生成
case $ACTION in
    "code")
        echo -e "${BLUE}==> 服务根目录 (cmd/api): $TARGET_DIR${NC}"

        # 1. 格式化 API 定义
        goctl api format -dir "$ABS_API_FILE"

        # 2. 生成 Go 代码 [cite: 2026-02-03]
        # -api 指定具体的子文件，-dir 指定 cmd/api 根目录
        goctl api go -api "$ABS_API_FILE" -dir "$TARGET_DIR" -style goZero

        # 3. 依赖清理
        # 自动向根目录寻找 go.mod
        MOD_DIR="$TARGET_DIR"
        while [[ "$MOD_DIR" != "/" ]]; do
            if [ -f "$MOD_DIR/go.mod" ]; then
                echo -e "${BLUE}==> 整理依赖 (go mod tidy)...${NC}"
                (cd "$MOD_DIR" && go mod tidy)
                break
            fi
            MOD_DIR=$(dirname "$MOD_DIR")
        done
        ;;

    "fmt")
        goctl api format -dir "$ABS_API_FILE"
        ;;

    *)
        echo -e "${RED}❌ 未知操作: $ACTION${NC}"
        exit 1
        ;;
esac

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✨ 代码已在 $TARGET_DIR 成功更新!${NC}"
fi