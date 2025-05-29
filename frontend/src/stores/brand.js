import { defineStore } from 'pinia'

export const useBrandStore = defineStore('brand', {
  state: () => ({
    // 品牌基础信息
    brandName: 'JiLang Agent',
    brandDescription: '智能工作流平台',
    brandSlogan: '让您轻松创建和管理工作流程',
    
    // 版权信息
    copyrightYear: new Date().getFullYear(),
    copyrightText: '保留所有权利',
    
    // 联系信息
    website: 'https://jilang-agent.com',
    supportEmail: 'support@jilang-agent.com',
    
    // 社交媒体
    socialLinks: {
      github: '#',
      twitter: '#',
      linkedin: '#'
    }
  }),
  
  getters: {
    // 完整的版权信息
    fullCopyright: (state) => `© ${state.copyrightYear} ${state.brandName}. ${state.copyrightText}`,
    
    // 品牌完整描述
    fullDescription: (state) => `${state.brandName} - ${state.brandDescription}，${state.brandSlogan}`,
    
    // 欢迎消息
    welcomeMessage: (state) => `欢迎使用${state.brandName}系统`
  },
  
  actions: {
    // 更新品牌名称
    updateBrandName(newName) {
      this.brandName = newName
    },
    
    // 更新品牌描述
    updateBrandDescription(newDescription) {
      this.brandDescription = newDescription
    },
    
    // 更新品牌标语
    updateBrandSlogan(newSlogan) {
      this.brandSlogan = newSlogan
    },
    
    // 批量更新品牌信息
    updateBrandInfo(brandInfo) {
      Object.keys(brandInfo).forEach(key => {
        if (this.hasOwnProperty(key)) {
          this[key] = brandInfo[key]
        }
      })
    }
  }
}) 