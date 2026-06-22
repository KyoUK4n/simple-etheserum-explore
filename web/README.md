# etherscan-frontend

一个简易的区块链浏览器前端部分

开发阶段可前后端分离，部署或运行可以通过后端服务反向代理访问即可

# Quick Start

安装依赖
```shell
npm install
```

## 1. 前端单独执行

配置api地址
```shell
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
```

启动
```shell
npm run dev
```

## 2. 由后端服务代理

打包
```shell
npm run build
```

启动后端服务
```shell
go run etherscan.go
```

# 技术栈

+ next.js
+ tailwindcss

