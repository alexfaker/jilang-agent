# JiLang Agent 前端

JiLang Agent 前端是一个基于 Vue 3 的现代化 Web 应用，提供 AI 代理工作流管理系统的用户界面。

## 技术栈

- Vue 3 (组合式 API)
- Vite (构建工具)
- Vue Router (路由管理)
- Pinia (状态管理)
- Axios (HTTP 客户端)
- Tailwind CSS (样式框架)
- Vue Toastification (通知系统)
- Vuelidate (表单验证)

## 开发环境设置

### 前提条件

- Node.js (v16+)
- npm (v7+) 或 yarn (v1.22+)

### 安装

1. 安装依赖

```bash
npm install
# 或
yarn install
```

2. 启动开发服务器

```bash
npm run dev
# 或
yarn dev
```

应用将在 http://localhost:5173 上运行。

## 项目结构

```
frontend/
├── public/                # 静态资源
├── src/                   # 源代码
│   ├── api/               # API 调用
│   ├── assets/            # 静态资源
│   ├── components/        # Vue 组件
│   │   ├── common/        # 通用组件
│   │   ├── layout/        # 布局组件
│   │   └── ...            # 其他组件
│   ├── router/            # 路由配置
│   ├── stores/            # Pinia 状态存储
│   ├── utils/             # 工具函数
│   ├── views/             # 页面视图
│   │   ├── dashboard/     # 仪表盘相关视图
│   │   ├── settings/      # 设置相关视图
│   │   └── ...            # 其他视图
│   ├── App.vue            # 根组件
│   ├── main.js            # 应用入口
│   └── style.css          # 全局样式
├── .env.development       # 开发环境变量
├── .env.production        # 生产环境变量
├── index.html             # HTML 模板
├── package.json           # 项目依赖
├── tailwind.config.js     # Tailwind 配置
└── vite.config.js         # Vite 配置
```

## 开发指南

### 命名约定

- 组件: PascalCase (例如 `UserProfile.vue`)
- 文件和目录: kebab-case (例如 `user-profile.js`)
- 变量和函数: camelCase (例如 `getUserData`)
- 常量: UPPER_SNAKE_CASE (例如 `API_BASE_URL`)

### 代码风格

- 使用组合式 API (`<script setup>`)
- 使用 ES6+ 特性
- 遵循 Vue 3 风格指南

### 提交规范

提交信息应遵循以下格式:

```
<类型>(<范围>): <描述>

[可选的详细描述]

[可选的脚注]
```

类型:
- feat: 新功能
- fix: 修复 bug
- docs: 文档更改
- style: 不影响代码含义的更改 (空格、格式等)
- refactor: 既不修复 bug 也不添加功能的代码更改
- perf: 提高性能的代码更改
- test: 添加或修改测试
- chore: 对构建过程或辅助工具的更改

## 构建生产版本

```bash
npm run build
# 或
yarn build
```

构建后的文件将位于 `dist` 目录中。

## 许可证

MIT License 