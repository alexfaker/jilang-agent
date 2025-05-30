-- 添加封面图字段到agents表
USE jilang_agent;

-- 添加cover_image字段
ALTER TABLE agents ADD COLUMN cover_image VARCHAR(500) DEFAULT NULL COMMENT '封面图URL';

-- 验证字段添加成功
DESCRIBE agents; 