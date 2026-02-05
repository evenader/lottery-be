#!/usr/bin/env bash

# =============================================================================
# 脚本位置: app/usercenter/genModel.sh
# 目标目录: app/usercenter/model (与 cmd 平级)
# =============================================================================
# 示例
#./genModel.sh looklook_usercenter user
# 1. 基础配置
host="127.0.0.1"
port="33069"
username="root"
passwd="PXDN93VRKUm8TeE7"
dbname=$1
tables=$2

# 2. 路径处理
SCRIPT_DIR=$(cd "$(dirname "$0")"; pwd)
# 直接生成到 usercenter/model 目录下
modeldir="${SCRIPT_DIR}/model"
# 向上两级寻找项目根目录下的 goctl 模板
#template="${SCRIPT_DIR}/../../goctl/1.6.1"
template=$(cd "${SCRIPT_DIR}/goctl/1.9.2"; pwd)
# 3. 执行生成
echo -e "\033[0;34m==> 正在生成 Model 到: $modeldir\033[0m"

goctl model mysql datasource \
    -url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" \
    -table="${tables}" \
    -dir="${modeldir}" \
    -cache=true \
    --home="${template}" \
    --style=goZero

# 4. 依赖整理
if [ $? -eq 0 ]; then
    (cd "${SCRIPT_DIR}" && go mod tidy)
    echo -e "\033[0;32m✅ 生成成功！\033[0m"
fi