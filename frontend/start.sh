#!/bin/bash

# Vue前端项目初始化和启动脚本

echo "🚀 正在初始化Vue前端项目..."

# 检查是否在正确的目录
if [ ! -f "package.json" ]; then
    echo "❌ 错误：请在frontend目录下运行此脚本"
    exit 1
fi

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "❌ 错误：请先安装Node.js"
    exit 1
fi

# 检查npm是否安装
if ! command -v npm &> /dev/null; then
    echo "❌ 错误：请先安装npm"
    exit 1
fi

echo "📦 正在安装依赖..."
npm install

if [ $? -ne 0 ]; then
    echo "❌ 依赖安装失败"
    exit 1
fi

echo "✅ 依赖安装完成"

echo "🔧 正在检查项目配置..."

# 检查必要的配置文件
if [ ! -f "vite.config.js" ]; then
    echo "❌ 错误：缺少vite.config.js配置文件"
    exit 1
fi

if [ ! -f "src/main.js" ]; then
    echo "❌ 错误：缺少src/main.js入口文件"
    exit 1
fi

echo "✅ 项目配置检查完成"

echo "🎨 正在启动开发服务器..."

# 启动开发服务器
npm run dev

echo "🎉 Vue前端项目已启动！"
echo "📱 请在浏览器中访问：http://localhost:5173"
echo "🔧 开发服务器支持热重载，修改代码后会自动刷新"