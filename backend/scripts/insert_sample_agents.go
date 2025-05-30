package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alexfaker/jilang-agent/config"
	"github.com/alexfaker/jilang-agent/models"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	// 连接数据库
	db, err := database.ConnectGormDB(cfg.Database)
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	// 插入示例数据
	if err := insertSampleAgents(db); err != nil {
		log.Fatalf("插入示例数据失败: %v", err)
	}

	fmt.Println("成功插入示例工作流市场数据！")
}

func insertSampleAgents(db *gorm.DB) error {
	agents := []models.Agent{
		{
			Name:        "智能文档处理器",
			Description: "自动处理PDF、Word文档，提取关键信息并生成摘要。支持多种文件格式，批量处理效率高。",
			Type:        "automation",
			Category:    "data",
			Icon:        "document-text",
			Definition: json.RawMessage(`{
				"steps": [
					{"type": "file_input", "name": "文件输入", "config": {"formats": ["pdf", "docx"]}},
					{"type": "text_extract", "name": "文本提取", "config": {"method": "ocr"}},
					{"type": "nlp_process", "name": "信息提取", "config": {"extract": ["summary", "keywords"]}},
					{"type": "output", "name": "结果输出", "config": {"format": "json"}}
				],
				"version": "1.0"
			}`),
			Price:         0,
			PurchaseCount: 245,
			Rating:        4.8,
			IsPublic:      true,
			CreatedAt:     time.Now().AddDate(0, -2, -15),
			UpdatedAt:     time.Now().AddDate(0, 0, -2),
		},
		{
			Name:        "社交媒体内容生成器",
			Description: "根据输入的主题和风格，自动生成吸引人的社交媒体内容。支持微博、微信、抖音等多平台。",
			Type:        "generator",
			Category:    "content",
			Icon:        "megaphone",
			Definition: json.RawMessage(`{
				"steps": [
					{"type": "topic_input", "name": "主题输入", "config": {"required": true}},
					{"type": "style_select", "name": "风格选择", "config": {"options": ["专业", "活泼", "幽默"]}},
					{"type": "content_generate", "name": "内容生成", "config": {"platforms": ["weibo", "wechat", "douyin"]}},
					{"type": "format_output", "name": "格式化输出", "config": {"include_hashtags": true}}
				],
				"version": "1.2"
			}`),
			Price:         50,
			PurchaseCount: 189,
			Rating:        4.6,
			IsPublic:      true,
			CreatedAt:     time.Now().AddDate(0, -1, -25),
			UpdatedAt:     time.Now().AddDate(0, 0, -5),
		},
		{
			Name:        "数据可视化大师",
			Description: "将Excel、CSV数据自动转换为精美的图表和报告。支持多种图表类型，一键生成专业报告。",
			Type:        "analytics",
			Category:    "analysis",
			Icon:        "chart-bar",
			Definition: json.RawMessage(`{
				"steps": [
					{"type": "data_import", "name": "数据导入", "config": {"formats": ["csv", "xlsx", "json"]}},
					{"type": "data_clean", "name": "数据清洗", "config": {"auto_detect": true}},
					{"type": "chart_generate", "name": "图表生成", "config": {"types": ["bar", "line", "pie", "scatter"]}},
					{"type": "report_build", "name": "报告构建", "config": {"template": "professional"}}
				],
				"version": "2.0"
			}`),
			Price:         120,
			PurchaseCount: 156,
			Rating:        4.9,
			IsPublic:      true,
			CreatedAt:     time.Now().AddDate(0, -1, -10),
			UpdatedAt:     time.Now().AddDate(0, 0, -1),
		},
		{
			Name:        "邮件营销助手",
			Description: "智能邮件营销工具，个性化内容生成，A/B测试，效果分析。提高邮件开启率和转换率。",
			Type:        "marketing",
			Category:    "automation",
			Icon:        "envelope",
			Definition: json.RawMessage(`{
				"steps": [
					{"type": "audience_segment", "name": "受众分析", "config": {"auto_segment": true}},
					{"type": "content_personalize", "name": "内容个性化", "config": {"use_ai": true}},
					{"type": "ab_test", "name": "A/B测试", "config": {"split_ratio": 0.5}},
					{"type": "send_email", "name": "邮件发送", "config": {"schedule": true}},
					{"type": "analytics", "name": "效果分析", "config": {"metrics": ["open_rate", "click_rate"]}}
				],
				"version": "1.5"
			}`),
			Price:         200,
			PurchaseCount: 98,
			Rating:        4.4,
			IsPublic:      true,
			CreatedAt:     time.Now().AddDate(0, -3, -5),
			UpdatedAt:     time.Now().AddDate(0, 0, -3),
		},
		{
			Name:        "语言翻译专家",
			Description: "支持100+语言的智能翻译工具。保持上下文准确性，专业术语识别，批量翻译文档。",
			Type:        "processor",
			Category:    "nlp",
			Icon:        "language",
			Definition: json.RawMessage(`{
				"steps": [
					{"type": "text_input", "name": "文本输入", "config": {"max_length": 10000}},
					{"type": "language_detect", "name": "语言检测", "config": {"confidence_threshold": 0.9}},
					{"type": "translate", "name": "智能翻译", "config": {"preserve_format": true}},
					{"type": "quality_check", "name": "质量检查", "config": {"grammar_check": true}},
					{"type": "output", "name": "结果输出", "config": {"formats": ["text", "docx", "pdf"]}}
				],
				"version": "3.1"
			}`),
			Price:         0,
			PurchaseCount: 312,
			Rating:        4.7,
			IsPublic:      true,
			CreatedAt:     time.Now().AddDate(0, -2, -20),
			UpdatedAt:     time.Now().AddDate(0, 0, -1),
		},
		{
			Name:        "图像风格转换器",
			Description: "将普通照片转换为艺术风格图像。支持油画、水彩、素描等多种艺术风格，一键美化图片。",
			Type:        "transformer",
			Category:    "content",
			Icon:        "photo",
			Definition: json.RawMessage(`{
				"steps": [
					{"type": "image_upload", "name": "图像上传", "config": {"formats": ["jpg", "png", "webp"]}},
					{"type": "style_select", "name": "风格选择", "config": {"styles": ["oil_painting", "watercolor", "sketch", "cartoon"]}},
					{"type": "ai_process", "name": "AI处理", "config": {"quality": "high"}},
					{"type": "preview", "name": "效果预览", "config": {"allow_adjust": true}},
					{"type": "download", "name": "下载结果", "config": {"resolution": "original"}}
				],
				"version": "2.3"
			}`),
			Price:         80,
			PurchaseCount: 167,
			Rating:        4.5,
			IsPublic:      true,
			CreatedAt:     time.Now().AddDate(0, -1, -18),
			UpdatedAt:     time.Now().AddDate(0, 0, -4),
		},
		{
			Name:        "网站SEO优化器",
			Description: "全面分析网站SEO状况，提供优化建议。关键词分析，竞争对手研究，排名监控。",
			Type:        "analyzer",
			Category:    "analysis",
			Icon:        "search",
			Definition: json.RawMessage(`{
				"steps": [
					{"type": "url_input", "name": "网站输入", "config": {"validate": true}},
					{"type": "crawl_site", "name": "网站爬取", "config": {"depth": 3}},
					{"type": "seo_analyze", "name": "SEO分析", "config": {"check_all": true}},
					{"type": "keyword_research", "name": "关键词研究", "config": {"include_competitors": true}},
					{"type": "report_generate", "name": "报告生成", "config": {"format": "detailed"}}
				],
				"version": "1.8"
			}`),
			Price:         150,
			PurchaseCount: 134,
			Rating:        4.6,
			IsPublic:      true,
			CreatedAt:     time.Now().AddDate(0, -1, -30),
			UpdatedAt:     time.Now().AddDate(0, 0, -6),
		},
		{
			Name:        "代码质量检查器",
			Description: "自动检查代码质量，发现潜在问题。支持多种编程语言，代码规范检查，安全漏洞扫描。",
			Type:        "validator",
			Category:    "automation",
			Icon:        "code",
			Definition: json.RawMessage(`{
				"steps": [
					{"type": "code_input", "name": "代码输入", "config": {"languages": ["python", "javascript", "java", "go"]}},
					{"type": "syntax_check", "name": "语法检查", "config": {"strict_mode": true}},
					{"type": "quality_scan", "name": "质量扫描", "config": {"rules": "comprehensive"}},
					{"type": "security_audit", "name": "安全审计", "config": {"vulnerability_check": true}},
					{"type": "report_output", "name": "报告输出", "config": {"include_fixes": true}}
				],
				"version": "2.1"
			}`),
			Price:         100,
			PurchaseCount: 89,
			Rating:        4.3,
			IsPublic:      true,
			CreatedAt:     time.Now().AddDate(0, -2, -8),
			UpdatedAt:     time.Now().AddDate(0, 0, -2),
		},
		{
			Name:        "会议记录转录器",
			Description: "将语音会议自动转录为文字，生成会议纪要。支持多人识别，智能提取行动项和决策。",
			Type:        "transcriber",
			Category:    "nlp",
			Icon:        "microphone",
			Definition: json.RawMessage(`{
				"steps": [
					{"type": "audio_upload", "name": "音频上传", "config": {"formats": ["mp3", "wav", "m4a"]}},
					{"type": "speech_to_text", "name": "语音转文字", "config": {"multi_speaker": true}},
					{"type": "content_structure", "name": "内容结构化", "config": {"identify_topics": true}},
					{"type": "action_extract", "name": "行动项提取", "config": {"smart_detection": true}},
					{"type": "summary_generate", "name": "摘要生成", "config": {"format": "minutes"}}
				],
				"version": "1.7"
			}`),
			Price:         0,
			PurchaseCount: 278,
			Rating:        4.8,
			IsPublic:      true,
			CreatedAt:     time.Now().AddDate(0, -1, -12),
			UpdatedAt:     time.Now().AddDate(0, 0, -1),
		},
		{
			Name:        "电商数据分析师",
			Description: "专为电商平台设计的数据分析工具。销售趋势分析，用户行为洞察，库存优化建议。",
			Type:        "analyzer",
			Category:    "analysis",
			Icon:        "shopping-cart",
			Definition: json.RawMessage(`{
				"steps": [
					{"type": "data_connect", "name": "数据连接", "config": {"platforms": ["shopify", "woocommerce", "magento"]}},
					{"type": "sales_analyze", "name": "销售分析", "config": {"period": "monthly"}},
					{"type": "user_behavior", "name": "用户行为分析", "config": {"track_journey": true}},
					{"type": "inventory_optimize", "name": "库存优化", "config": {"predict_demand": true}},
					{"type": "dashboard_create", "name": "仪表板创建", "config": {"real_time": true}}
				],
				"version": "2.2"
			}`),
			Price:         300,
			PurchaseCount: 67,
			Rating:        4.7,
			IsPublic:      true,
			CreatedAt:     time.Now().AddDate(0, -3, -12),
			UpdatedAt:     time.Now().AddDate(0, 0, -7),
		},
	}

	// 批量插入数据
	for _, agent := range agents {
		var existingAgent models.Agent
		result := db.Where("name = ?", agent.Name).First(&existingAgent)
		if result.Error == gorm.ErrRecordNotFound {
			// 记录不存在，插入新记录
			if err := db.Create(&agent).Error; err != nil {
				return fmt.Errorf("插入工作流 '%s' 失败: %v", agent.Name, err)
			}
			fmt.Printf("✓ 插入工作流: %s\n", agent.Name)
		} else if result.Error != nil {
			return fmt.Errorf("检查工作流 '%s' 时出错: %v", agent.Name, result.Error)
		} else {
			fmt.Printf("- 工作流 '%s' 已存在，跳过\n", agent.Name)
		}
	}

	return nil
}
