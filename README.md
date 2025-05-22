# JiLang Agent

JiLang Agent 是一个强大的AI代理工作流管理系统，它可以帮助您创建、管理和执行基于AI的工作流，实现自动化处理各种任务。

## 项目结构

```
jilang-agent/
├── backend/             # Go后端
│   ├── api/             # API相关代码
│   │   ├── handlers/    # 请求处理器
│   │   ├── middleware/  # 中间件
│   │   └── routes/      # 路由定义
│   ├── config/          # 配置管理
│   ├── models/          # 数据模型
│   └── pkg/             # 工具包
│       ├── database/    # 数据库连接
│       └── logger/      # 日志工具
├── frontend/            # Vue.js前端
│   ├── public/          # 静态资源
│   └── src/             # 源代码
│       ├── api/         # API调用
│       ├── assets/      # 静态资源
│       ├── components/  # Vue组件
│       ├── router/      # 路由管理
│       ├── store/       # 状态管理
│       └── views/       # 页面视图
└── main.go              # 应用入口
```

## 技术栈

### 后端
- Go 1.20+
- Chi (路由)
- MySQL / PostgreSQL (数据库)
- JWT (认证)
- Zap (日志)

### 前端
- Vue 3
- Vue Router
- Axios
- Vite

## 功能特性

- 用户认证与授权
- 工作流管理 (创建、编辑、删除、执行)
- 执行历史记录
- 代理管理
- 统计数据分析
- 响应式UI设计

## 开发环境配置

### 后端

1. 安装Go (1.20+): [https://go.dev/doc/install](https://go.dev/doc/install)
2. 设置数据库 (MySQL/PostgreSQL)
3. 克隆仓库并进入项目目录

```bash
git clone https://github.com/alexfaker/jilang-agent.git
cd jilang-agent
```

4. 安装依赖

```bash
go mod download
```

5. 创建并配置 `.env` 文件 (参考 `.env.development`)
6. 启动后端服务

```bash
go run main.go
```

### 前端

1. 安装Node.js (16+): [https://nodejs.org/](https://nodejs.org/)
2. 进入前端目录

```bash
cd frontend
```

3. 安装依赖

```bash
npm install
```

4. 启动开发服务器

```bash
npm run dev
```

## 生产环境部署

### 后端

1. 编译Go程序

```bash
go build -o jilang-agent main.go
```

2. 配置生产环境变量 (参考 `.env.production`)
3. 运行编译后的程序

```bash
./jilang-agent
```

### 前端

1. 编译前端资源

```bash
cd frontend
npm run build
```

2. 将 `frontend/dist` 目录中的文件部署到Web服务器

## API文档

API文档可以通过访问 `/api/docs` 路径获取 (在开发模式下)。

## 贡献指南

欢迎提交Pull Request或Issue来改进项目!

## 许可证

MIT License