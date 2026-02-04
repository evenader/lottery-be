#!/usr/bin/env bash

# =============================================================================
# 脚本位置: app/usercenter/genModel.sh
# 使用方法:
#   生成到 API 层: ./genModel.sh lottery user api
#   生成到 RPC 层: ./genModel.sh lottery user rpc
# =============================================================================

# 1. 基础配置
host="127.0.0.1"
port="33069"
username="root"
passwd="PXDN93VRKUm8TeE7"

dbname=$1
tables=$2
service_type=${3:-"api"} # 默认生成到 api 层，也可以传入 rpc

# 2. 路径处理 (基于你图中的位置)
# SCRIPT_DIR 就是 app/usercenter
SCRIPT_DIR=$(cd "$(dirname "$0")"; pwd)

# 确定生成的 model 放置目录
modeldir="${SCRIPT_DIR}/cmd/${service_type}/internal/model"

# 向上三级寻找项目根目录下的 goctl 模板 (app -> lottery-be)
# 对应你图中的相对路径：../../goctl/1.6.1
template="${SCRIPT_DIR}/../../goctl/1.6.1"

# 3. 检查并创建目录
if [ ! -d "$modeldir" ]; then
    mkdir -p "$modeldir"
fi

echo -e "\033[0;34m==> 数据库: $dbname, 表: $tables\033[0m"
echo -e "\033[0;34m==> 目标路径: $modeldir\033[0m"

# 4. 执行 goctl 命令
# 注意：直接在目标目录下生成，goctl 会自动处理 package model 声明 [cite: 2026-02-03]
goctl model mysql datasource \
    -url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" \
    -table="${tables}" \
    -dir="${modeldir}" \
    -cache=true \
    --home="${template}" \
    --style=goZero

# 5. 自动整理依赖
if [ -f "${SCRIPT_DIR}/go.mod" ]; then
    (cd "${SCRIPT_DIR}" && go mod tidy)
fi

echo -e "\033[0;32m✅ 执行完成！代码已存入 $modeldir\033[0m"