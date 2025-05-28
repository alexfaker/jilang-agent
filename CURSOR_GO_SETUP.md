# Cursor IDE Go语言跳转修复指南

## 问题描述
Cursor IDE中Go语言无法跳转到变量和方法定义的位置。

## 解决方案

### 1. 重启Cursor IDE
首先完全关闭Cursor IDE，然后重新打开项目。

### 2. 使用工作区文件
打开项目根目录下的 `jilang-agent.code-workspace` 文件，而不是直接打开文件夹。

### 3. 安装Go扩展
确保安装了官方的Go扩展：
- 在Cursor IDE中按 `Cmd+Shift+P` (Mac) 或 `Ctrl+Shift+P` (Windows/Linux)
- 输入 "Extensions: Install Extensions"
- 搜索并安装 "Go" 扩展 (由Google发布)

### 4. 重新加载窗口
- 按 `Cmd+Shift+P` (Mac) 或 `Ctrl+Shift+P` (Windows/Linux)
- 输入 "Developer: Reload Window"
- 选择并执行

### 5. 检查Go语言服务器状态
- 按 `Cmd+Shift+P` (Mac) 或 `Ctrl+Shift+P` (Windows/Linux)
- 输入 "Go: Show Language Server Status"
- 确保gopls正在运行

### 6. 重启语言服务器
如果跳转仍然不工作：
- 按 `Cmd+Shift+P` (Mac) 或 `Ctrl+Shift+P` (Windows/Linux)
- 输入 "Go: Restart Language Server"
- 选择并执行

### 7. 验证配置
检查以下文件是否存在并配置正确：
- `backend/.vscode/settings.json` - Go语言配置
- `backend/.vscode/launch.json` - 调试配置
- `backend/go.work` - Go工作区文件
- `jilang-agent.code-workspace` - Cursor工作区配置

### 8. 环境变量检查
确保以下环境变量正确设置：
```bash
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org
```

### 9. 清理和重建
在backend目录下运行：
```bash
go clean -modcache
go mod download
go mod tidy
```

## 测试跳转功能

### 测试步骤：
1. 打开 `backend/main.go` 文件
2. 将鼠标悬停在 `config.LoadConfig()` 上
3. 按住 `Cmd` (Mac) 或 `Ctrl` (Windows/Linux) 并点击
4. 应该能跳转到 `config.go` 文件中的 `LoadConfig` 函数定义

### 其他测试：
- 在变量上右键选择 "Go to Definition"
- 使用 `F12` 快捷键跳转到定义
- 使用 `Shift+F12` 查看所有引用

## 常见问题

### Q: 仍然无法跳转怎么办？
A: 尝试以下步骤：
1. 完全关闭Cursor IDE
2. 删除 `backend/.vscode` 目录
3. 重新打开项目
4. 重新创建配置文件

### Q: 提示找不到包怎么办？
A: 运行以下命令：
```bash
cd backend
go mod tidy
go mod download
```

### Q: gopls进程占用过多资源？
A: 在设置中调整gopls配置：
```json
{
  "gopls": {
    "build.directoryFilters": ["-node_modules", "-vendor"],
    "ui.diagnostic.staticcheck": false
  }
}
```

## 配置文件说明

### `.vscode/settings.json`
包含Go语言服务器的详细配置，包括gopls设置、代码检查、格式化等。

### `go.work`
Go工作区文件，帮助gopls理解项目结构。

### `jilang-agent.code-workspace`
Cursor IDE工作区配置，包含多文件夹项目设置和Go语言特定配置。

## 验证成功
如果配置成功，你应该能够：
- 跳转到函数定义
- 跳转到变量声明
- 查看函数/变量的所有引用
- 获得智能代码补全
- 看到实时的语法错误提示
- 使用代码重构功能

## 注意事项
- 确保使用工作区文件打开项目，而不是直接打开文件夹
- 第一次打开项目时，gopls可能需要一些时间来索引代码
- 如果项目很大，建议增加gopls的内存限制 