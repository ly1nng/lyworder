#!/bin/bash

# 前端构建和部署脚本

echo "🚀 开始前端构建和部署流程..."

# 检查是否在正确的目录
if [ ! -f "package.json" ]; then
    echo "❌ 错误：请在frontend目录下运行此脚本"
    exit 1
fi

echo "📦 正在安装依赖..."
npm install

if [ $? -ne 0 ]; then
    echo "❌ 依赖安装失败"
    exit 1
fi

echo "🔨 正在构建生产版本..."
npm run build

if [ $? -ne 0 ]; then
    echo "❌ 构建失败"
    exit 1
fi

echo "🚚 正在部署到后端..."
# 复制构建后的文件到后端static目录
rsync -avz dist/* ../static/

echo "✅ 前端构建和部署完成！"

echo "💡 要启动服务，请运行："
echo "   cd .. && go run main.go"
echo "   然后访问 http://localhost:8080"