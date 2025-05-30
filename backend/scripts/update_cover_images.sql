-- 为现有工作流添加封面图
USE jilang_agent;

-- 更新各个工作流的封面图
UPDATE agents SET cover_image = 'https://images.unsplash.com/photo-1586953208448-b95a79798f07?w=500&h=300&fit=crop&crop=center' WHERE name = '智能文档处理器';

UPDATE agents SET cover_image = 'https://images.unsplash.com/photo-1611224923853-80b023f02d71?w=500&h=300&fit=crop&crop=center' WHERE name = '社交媒体内容生成器';

UPDATE agents SET cover_image = 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=500&h=300&fit=crop&crop=center' WHERE name = '数据可视化大师';

UPDATE agents SET cover_image = 'https://images.unsplash.com/photo-1596526131083-e8c633c948d2?w=500&h=300&fit=crop&crop=center' WHERE name = '邮件营销助手';

UPDATE agents SET cover_image = 'https://images.unsplash.com/photo-1520637836862-4d197d17c93a?w=500&h=300&fit=crop&crop=center' WHERE name = '语言翻译专家';

UPDATE agents SET cover_image = 'https://images.unsplash.com/photo-1541961017774-22349e4a1262?w=500&h=300&fit=crop&crop=center' WHERE name = '图像风格转换器';

UPDATE agents SET cover_image = 'https://images.unsplash.com/photo-1432888622747-4eb9a8efeb07?w=500&h=300&fit=crop&crop=center' WHERE name = '网站SEO优化器';

UPDATE agents SET cover_image = 'https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=500&h=300&fit=crop&crop=center' WHERE name = '代码质量检查器';

UPDATE agents SET cover_image = 'https://images.unsplash.com/photo-1573164574572-cb89e39749b4?w=500&h=300&fit=crop&crop=center' WHERE name = '会议记录转录器';

UPDATE agents SET cover_image = 'https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?w=500&h=300&fit=crop&crop=center' WHERE name = '电商数据分析师';

-- 验证更新结果
SELECT id, name, cover_image FROM agents WHERE cover_image IS NOT NULL; 