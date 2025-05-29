/**
 * 品牌配置文件
 * 集中管理品牌相关的配置信息
 */

export const BRAND_CONFIG = {
  // 基础品牌信息
  name: 'JiLang Agent',
  description: '智能工作流平台',
  slogan: '让您轻松创建和管理工作流程',
  
  // 版权信息
  copyright: {
    year: new Date().getFullYear(),
    text: '保留所有权利',
    owner: 'JiLang Agent'
  },
  
  // 联系信息
  contact: {
    website: 'https://jilang-agent.com',
    email: 'support@jilang-agent.com',
    phone: '+86 400-123-4567'
  },
  
  // 社交媒体链接
  social: {
    github: 'https://github.com/jilang-agent',
    twitter: 'https://twitter.com/jilang_agent',
    linkedin: 'https://linkedin.com/company/jilang-agent',
    wechat: '#',
    weibo: '#'
  },
  
  // 品牌颜色
  colors: {
    primary: '#4F46E5',    // indigo-600
    secondary: '#6B7280',  // gray-500
    accent: '#10B981',     // emerald-500
    warning: '#F59E0B',    // amber-500
    error: '#EF4444',      // red-500
    success: '#10B981'     // emerald-500
  },
  
  // Logo配置
  logo: {
    icon: 'CpuChipIcon',  // Heroicons图标名称
    sizes: {
      sm: { icon: 'h-6 w-6', text: 'text-base' },
      base: { icon: 'h-8 w-8', text: 'text-xl' },
      lg: { icon: 'h-10 w-10', text: 'text-2xl' },
      xl: { icon: 'h-12 w-12', text: 'text-3xl' }
    }
  },
  
  // SEO信息
  seo: {
    title: 'JiLang Agent - 智能工作流平台',
    description: 'JiLang Agent是一个强大的智能工作流平台，帮助您轻松创建、管理和执行自动化工作流程。',
    keywords: ['工作流', '自动化', 'AI', '智能平台', '流程管理'],
    author: 'JiLang Agent Team'
  },
  
  // 功能模块配置
  features: {
    workflows: {
      name: '工作流管理',
      description: '创建和管理您的自定义工作流',
      icon: 'RectangleGroupIcon'
    },
    marketplace: {
      name: '工作流市场',
      description: '浏览和使用预构建的工作流',
      icon: 'BuildingStorefrontIcon'
    },
    executions: {
      name: '执行历史',
      description: '查看工作流执行历史记录',
      icon: 'ClockIcon'
    },
    analytics: {
      name: '统计分析',
      description: '获取工作流性能指标和统计数据',
      icon: 'ChartBarIcon'
    },
    profile: {
      name: '个人中心',
      description: '管理账户设置和个人偏好',
      icon: 'UserIcon'
    }
  }
}

// 导出便捷访问函数
export const getBrandName = () => BRAND_CONFIG.name
export const getBrandDescription = () => BRAND_CONFIG.description
export const getBrandSlogan = () => BRAND_CONFIG.slogan
export const getCopyright = () => `© ${BRAND_CONFIG.copyright.year} ${BRAND_CONFIG.copyright.owner}. ${BRAND_CONFIG.copyright.text}`
export const getFullDescription = () => `${BRAND_CONFIG.name} - ${BRAND_CONFIG.description}，${BRAND_CONFIG.slogan}`

// 默认导出
export default BRAND_CONFIG 