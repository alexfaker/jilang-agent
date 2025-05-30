-- 为users表添加设置页面需要的新字段
-- 执行前请备份数据库

-- 添加个人简介字段
ALTER TABLE users ADD COLUMN bio TEXT COMMENT '个人简介';

-- 添加时区字段
ALTER TABLE users ADD COLUMN timezone VARCHAR(50) DEFAULT 'Asia/Shanghai' COMMENT '用户时区';

-- 添加语言字段  
ALTER TABLE users ADD COLUMN language VARCHAR(10) DEFAULT 'zh_CN' COMMENT '用户语言偏好';

-- 添加主题字段
ALTER TABLE users ADD COLUMN theme VARCHAR(20) DEFAULT 'light' COMMENT '用户主题偏好';

-- 验证字段添加成功
SELECT COLUMN_NAME, DATA_TYPE, COLUMN_DEFAULT, COLUMN_COMMENT 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'users' 
AND COLUMN_NAME IN ('bio', 'timezone', 'language', 'theme'); 