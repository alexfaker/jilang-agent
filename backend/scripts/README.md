# 数据库脚本使用说明

## 插入示例工作流市场数据

本目录包含用于插入示例工作流市场数据的脚本，用于前端页面调试。

### 方法一：使用MySQL命令行（推荐）

```bash
# 1. 连接到MySQL数据库
mysql -u your_username -p

# 2. 执行SQL脚本
mysql> source /path/to/Backend/scripts/insert_sample_agents.sql;

# 或者直接在命令行执行
mysql -u your_username -p jilang_agent < Backend/scripts/insert_sample_agents.sql
```

### 方法二：使用MySQL工具

1. 打开 MySQL Workbench、phpMyAdmin 或其他MySQL管理工具
2. 连接到数据库
3. 打开 `insert_sample_agents.sql` 文件
4. 执行SQL脚本

### 方法三：使用Go脚本（需Go 1.24.3+）

```bash
cd Backend
go run scripts/insert_sample_agents.go
```

## 插入的示例数据

脚本将插入10个不同类型的工作流模板：

1. **智能文档处理器** (免费) - 数据处理类
2. **社交媒体内容生成器** (50积分) - 内容生成类
3. **数据可视化大师** (120积分) - 数据分析类
4. **邮件营销助手** (200积分) - 自动化类
5. **语言翻译专家** (免费) - 自然语言处理类
6. **图像风格转换器** (80积分) - 内容生成类
7. **网站SEO优化器** (150积分) - 数据分析类
8. **代码质量检查器** (100积分) - 自动化类
9. **会议记录转录器** (免费) - 自然语言处理类
10. **电商数据分析师** (300积分) - 数据分析类

### 特点

- 包含不同分类（data、content、analysis、automation、nlp）
- 既有免费也有付费工作流
- 不同的购买次数和评分
- 真实的工作流定义JSON
- 设置为公开可见（is_public = true）

## 验证数据插入

执行脚本后，可以运行以下查询验证：

```sql
SELECT id, name, category, price, purchase_count, rating 
FROM agents 
WHERE is_public = true 
ORDER BY purchase_count DESC;
```

## 清理测试数据

如需删除测试数据：

```sql
DELETE FROM agents WHERE name IN (
    '智能文档处理器', '社交媒体内容生成器', '数据可视化大师',
    '邮件营销助手', '语言翻译专家', '图像风格转换器',
    '网站SEO优化器', '代码质量检查器', '会议记录转录器',
    '电商数据分析师'
);
``` 