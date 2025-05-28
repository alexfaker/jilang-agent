# JiLang Agent API 文档

## 概述

JiLang Agent 是一个AI代理工作流管理系统的后端API，提供用户认证、工作流管理、执行历史、代理管理和统计分析等功能。

**基础URL**: `http://localhost:8080/api`

## 认证

大部分API端点需要JWT认证。在请求头中包含：
```
Authorization: Bearer <your-jwt-token>
```

## API 端点

### 健康检查

#### GET /health
检查服务器状态

**响应**:
```json
{
  "status": "success",
  "data": {
    "message": "服务运行正常"
  }
}
```

### 认证相关

#### POST /api/auth/register
用户注册

**请求体**:
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "role": "user"
    },
    "token": "jwt-token-string"
  }
}
```

#### POST /api/auth/login
用户登录

**请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "role": "user"
    },
    "token": "jwt-token-string"
  }
}
```

#### POST /api/auth/refresh
刷新JWT令牌

**请求体**:
```json
{
  "token": "current-jwt-token"
}
```

### 用户相关 🔒

#### GET /api/user/profile
获取当前用户资料

**响应**:
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "fullName": "Test User",
    "avatar": "/static/avatars/default.png",
    "role": "user",
    "createdAt": "2023-11-15T10:00:00Z",
    "lastLoginAt": "2023-11-15T10:00:00Z"
  }
}
```

#### PUT /api/user/profile
更新用户资料

**请求体**:
```json
{
  "email": "string",
  "fullName": "string",
  "avatar": "string"
}
```

#### POST /api/user/change-password
修改密码

**请求体**:
```json
{
  "currentPassword": "string",
  "newPassword": "string"
}
```

#### GET /api/user/:id
根据ID获取用户信息

### 工作流相关 🔒

#### GET /api/workflows
获取工作流列表

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 10, 最大: 100)
- `name`: 名称筛选
- `tag`: 标签筛选
- `status`: 状态筛选 (active, inactive, archived, draft, all)

**响应**:
```json
{
  "status": "success",
  "data": {
    "workflows": [
      {
        "id": 1,
        "name": "示例工作流",
        "description": "这是一个示例工作流",
        "definition": {},
        "status": "active",
        "userId": 1,
        "createdAt": "2023-11-15T10:00:00Z",
        "updatedAt": "2023-11-15T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 1,
      "page": 1,
      "page_size": 10,
      "pages": 1
    }
  }
}
```

#### POST /api/workflows
创建新工作流

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "definition": {},
  "status": "draft|active",
  "tags": ["string"]
}
```

#### GET /api/workflows/:id
获取工作流详情

#### PUT /api/workflows/:id
更新工作流

#### DELETE /api/workflows/:id
删除工作流

### 执行相关 🔒

#### GET /api/executions
获取执行历史列表

**查询参数**:
- `page`: 页码
- `page_size`: 每页数量
- `workflow_id`: 工作流ID筛选
- `status`: 状态筛选

#### POST /api/workflows/:id/execute
执行工作流

**请求体**:
```json
{
  "inputs": {},
  "config": {}
}
```

#### GET /api/executions/:id
获取执行详情

#### DELETE /api/executions/:id
删除执行记录

### 代理相关

#### GET /api/agent-categories
获取代理分类列表（公开接口）

**响应**:
```json
{
  "status": "success",
  "data": [
    "数据处理",
    "文本分析",
    "图像处理",
    "API集成"
  ]
}
```

#### GET /api/agents 🔒
获取代理列表

**查询参数**:
- `page`: 页码
- `page_size`: 每页数量
- `category`: 分类筛选
- `is_public`: 是否公开
- `search`: 搜索关键词

#### POST /api/agents 🔒
创建新代理

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "type": "string",
  "category": "string",
  "icon": "string",
  "definition": {},
  "isPublic": false
}
```

#### GET /api/agents/:id 🔒
获取代理详情

#### PUT /api/agents/:id 🔒
更新代理

#### DELETE /api/agents/:id 🔒
删除代理

### 统计相关 🔒

#### GET /api/stats/dashboard
获取仪表盘统计数据

**响应**:
```json
{
  "status": "success",
  "data": {
    "total_workflows": 10,
    "total_executions": 50,
    "success_rate": 85.5,
    "recent_executions": [],
    "daily_stats": [
      {
        "date": "2023-11-15",
        "count": 5,
        "succeeded": 4,
        "failed": 1
      }
    ]
  }
}
```

#### GET /api/stats/workflows
获取工作流统计数据

#### GET /api/stats/executions
获取执行统计数据

## 错误响应

所有错误响应都遵循以下格式：

```json
{
  "status": "error",
  "message": "错误描述信息"
}
```

常见HTTP状态码：
- `200`: 成功
- `201`: 创建成功
- `400`: 请求参数错误
- `401`: 未认证
- `403`: 权限不足
- `404`: 资源不存在
- `500`: 服务器内部错误

## 数据模型

### User (用户)
```json
{
  "id": "int64",
  "username": "string",
  "email": "string",
  "fullName": "string",
  "avatar": "string",
  "role": "string",
  "createdAt": "datetime",
  "updatedAt": "datetime",
  "lastLoginAt": "datetime"
}
```

### Workflow (工作流)
```json
{
  "id": "int64",
  "name": "string",
  "description": "string",
  "definition": "json",
  "status": "string",
  "userId": "int64",
  "createdAt": "datetime",
  "updatedAt": "datetime"
}
```

### WorkflowExecution (工作流执行)
```json
{
  "id": "int64",
  "workflowId": "int64",
  "status": "string",
  "inputs": "json",
  "outputs": "json",
  "error": "string",
  "startedAt": "datetime",
  "completedAt": "datetime",
  "duration": "int64"
}
```

### Agent (代理)
```json
{
  "id": "int64",
  "name": "string",
  "description": "string",
  "type": "string",
  "category": "string",
  "icon": "string",
  "definition": "json",
  "isPublic": "boolean",
  "userId": "int64",
  "usageCount": "int",
  "createdAt": "datetime",
  "updatedAt": "datetime"
}
```

## 开发环境设置

1. 确保MySQL数据库运行在 `localhost:3306`
2. 创建数据库 `jilang_agent`
3. 配置文件位于 `config/config.development.json`
4. 启动服务器：`go run main.go`
5. 服务器将在 `http://localhost:8080` 启动

## 注意事项

- 🔒 标记的端点需要JWT认证
- 所有时间戳使用ISO 8601格式
- JSON字段使用驼峰命名法
- 分页从第1页开始
- 默认每页返回10条记录 