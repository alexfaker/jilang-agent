import api from './index.js'

// 设置相关API
export const settingsApi = {
  // 获取用户资料
  getUserProfile() {
    return api.get('/api/settings/profile')
  },

  // 更新用户资料
  updateUserProfile(data) {
    return api.put('/api/settings/profile', data)
  },

  // 上传头像
  uploadAvatar(formData) {
    return api.post('/api/settings/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 修改密码
  changePassword(data) {
    return api.post('/api/settings/password', data)
  },

  // 获取系统设置
  getSettings() {
    return api.get('/api/settings')
  },

  // 更新系统设置
  updateSettings(data) {
    return api.put('/api/settings', data)
  }
}

export default settingsApi 