<template>
  <div class="flex-1 overflow-y-auto p-6">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-xl font-semibold text-gray-800">设置</h1>
      <p class="text-gray-600 mt-1">管理您的账户和偏好设置</p>
    </div>

    <!-- 设置导航 -->
    <div class="border-b mb-6">
      <nav class="flex space-x-8">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            activeTab === tab.id
              ? 'py-4 px-1 border-b-2 border-indigo-600 text-indigo-600 font-medium text-sm'
              : 'py-4 px-1 text-gray-500 hover:text-indigo-600 font-medium text-sm'
          ]"
        >
          {{ tab.name }}
        </button>
      </nav>
    </div>

    <!-- 账户设置 -->
    <div v-if="activeTab === 'account'" class="space-y-6">
      <!-- 个人信息 -->
      <div class="bg-white rounded-lg shadow-sm p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">个人信息</h2>
        
        <form @submit.prevent="updateProfile">
          <!-- 头像设置 -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-2">头像</label>
            <div class="flex items-center">
              <img 
                :src="profileData.avatar || '/default-avatar.png'" 
                class="w-16 h-16 rounded-full mr-4 object-cover" 
                alt="当前头像"
                @error="handleAvatarError"
              >
              <div>
                <input 
                  type="file" 
                  ref="avatarInput" 
                  @change="handleAvatarChange" 
                  accept="image/*"
                  class="hidden"
                >
                <button 
                  type="button" 
                  @click="$refs.avatarInput.click()"
                  :disabled="uploading"
                  class="bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-400 text-white px-4 py-2 rounded-lg text-sm mr-2"
                >
                  {{ uploading ? '上传中...' : '上传新头像' }}
                </button>
                <button 
                  type="button" 
                  @click="removeAvatar"
                  class="text-gray-600 hover:text-gray-800 px-4 py-2 border border-gray-300 rounded-lg text-sm"
                >
                  移除
                </button>
                <p class="text-xs text-gray-500 mt-1">允许的文件类型: PNG, JPG, GIF. 最大文件大小: 2MB</p>
              </div>
            </div>
          </div>
          
          <!-- 名字和邮箱 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div>
              <label for="username" class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
              <input 
                type="text" 
                id="username" 
                v-model="profileData.username"
                readonly
                class="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50 cursor-not-allowed"
              >
              <p class="text-xs text-gray-500 mt-1">用户名不可修改</p>
            </div>
            <div>
              <label for="email" class="block text-sm font-medium text-gray-700 mb-1">电子邮箱</label>
              <input 
                type="email" 
                id="email" 
                v-model="profileData.email"
                required
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
              >
            </div>
          </div>

          <!-- 全名 -->
          <div class="mb-6">
            <label for="fullName" class="block text-sm font-medium text-gray-700 mb-1">全名</label>
            <input 
              type="text" 
              id="fullName" 
              v-model="profileData.fullName"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
            >
          </div>
          
          <!-- 个人简介 -->
          <div class="mb-6">
            <label for="bio" class="block text-sm font-medium text-gray-700 mb-1">个人简介</label>
            <textarea 
              id="bio" 
              v-model="profileData.bio"
              rows="4" 
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
              placeholder="介绍一下自己..."
            ></textarea>
          </div>
          
          <!-- 时区和语言 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div>
              <label for="timezone" class="block text-sm font-medium text-gray-700 mb-1">时区</label>
              <select 
                id="timezone" 
                v-model="profileData.timezone"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="Asia/Shanghai">中国标准时间 (GMT+8)</option>
                <option value="America/New_York">美国东部时间 (GMT-5)</option>
                <option value="Europe/London">格林威治标准时间 (GMT)</option>
                <option value="Asia/Tokyo">日本标准时间 (GMT+9)</option>
              </select>
            </div>
            <div>
              <label for="language" class="block text-sm font-medium text-gray-700 mb-1">语言</label>
              <select 
                id="language" 
                v-model="profileData.language"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="zh_CN">简体中文</option>
                <option value="en_US">English (US)</option>
                <option value="ja_JP">日本語</option>
                <option value="ko_KR">한국어</option>
              </select>
            </div>
          </div>
          
          <!-- 保存按钮 -->
          <div class="flex justify-end">
            <button 
              type="button" 
              @click="resetProfile"
              class="mr-4 px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
            >
              重置
            </button>
            <button 
              type="submit" 
              :disabled="saving"
              class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-400 text-white rounded-lg"
            >
              {{ saving ? '保存中...' : '保存更改' }}
            </button>
          </div>
        </form>
      </div>

      <!-- 界面设置 -->
      <div class="bg-white rounded-lg shadow-sm p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">界面设置</h2>
        
        <!-- 主题设置 -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">主题</label>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div 
              v-for="theme in themes"
              :key="theme.value"
              @click="profileData.theme = theme.value"
              :class="[
                profileData.theme === theme.value 
                  ? 'border-2 border-indigo-600 bg-indigo-50' 
                  : 'border border-gray-200 hover:bg-gray-50',
                'rounded-lg p-4 flex items-center cursor-pointer transition-all'
              ]"
            >
              <div :class="theme.colorClass" class="w-6 h-6 rounded-full mr-3"></div>
              <span class="text-sm font-medium">{{ theme.name }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 安全设置 -->
    <div v-if="activeTab === 'security'" class="space-y-6">
      <div class="bg-white rounded-lg shadow-sm p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">修改密码</h2>
        
        <form @submit.prevent="changePassword">
          <div class="space-y-4">
            <div>
              <label for="currentPassword" class="block text-sm font-medium text-gray-700 mb-1">当前密码</label>
              <input 
                type="password" 
                id="currentPassword" 
                v-model="passwordData.currentPassword"
                required
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
              >
            </div>
            <div>
              <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">新密码</label>
              <input 
                type="password" 
                id="newPassword" 
                v-model="passwordData.newPassword"
                required
                minlength="8"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
              >
            </div>
            <div>
              <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">确认新密码</label>
              <input 
                type="password" 
                id="confirmPassword" 
                v-model="passwordData.confirmPassword"
                required
                minlength="8"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
              >
              <p v-if="passwordMismatch" class="text-sm text-red-600 mt-1">密码不匹配</p>
            </div>
          </div>
          
          <div class="flex justify-end mt-6">
            <button 
              type="submit" 
              :disabled="changingPassword || passwordMismatch"
              class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-400 text-white rounded-lg"
            >
              {{ changingPassword ? '修改中...' : '修改密码' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 其他标签页暂时显示占位内容 -->
    <div v-if="activeTab === 'workflow'" class="bg-white rounded-lg shadow-sm p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">工作流设置</h2>
      <p class="text-gray-600">工作流相关设置功能正在开发中...</p>
    </div>

    <div v-if="activeTab === 'api'" class="bg-white rounded-lg shadow-sm p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">API密钥</h2>
      <p class="text-gray-600">API密钥管理功能正在开发中...</p>
    </div>

    <div v-if="activeTab === 'notification'" class="bg-white rounded-lg shadow-sm p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">通知设置</h2>
      <p class="text-gray-600">通知设置功能正在开发中...</p>
    </div>

    <!-- 错误提示 -->
    <div v-if="error" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
      <p class="text-red-800">{{ error }}</p>
    </div>

    <!-- 成功提示 -->
    <div v-if="success" class="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
      <p class="text-green-800">{{ success }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { authApi } from '../../api/index.js'

const activeTab = ref('account')
const loading = ref(false)
const saving = ref(false)
const uploading = ref(false)
const changingPassword = ref(false)
const error = ref('')
const success = ref('')

const tabs = [
  { id: 'account', name: '账户设置' },
  { id: 'workflow', name: '工作流设置' },
  { id: 'security', name: '安全设置' },
  { id: 'api', name: 'API密钥' },
  { id: 'notification', name: '通知设置' }
]

const themes = [
  { value: 'light', name: '浅色主题', colorClass: 'bg-white border border-gray-300' },
  { value: 'dark', name: '深色主题', colorClass: 'bg-gray-800 border border-gray-700' },
  { value: 'auto', name: '跟随系统', colorClass: 'bg-gradient-to-r from-white to-gray-800 border border-gray-300' }
]

const profileData = reactive({
  id: 0,
  userID: '',
  username: '',
  email: '',
  fullName: '',
  avatar: '',
  bio: '',
  timezone: 'Asia/Shanghai',
  language: 'zh_CN',
  theme: 'light',
  role: '',
  points: 0
})

const originalProfileData = reactive({})

const passwordData = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordMismatch = computed(() => {
  return passwordData.newPassword && passwordData.confirmPassword && 
         passwordData.newPassword !== passwordData.confirmPassword
})

// 获取用户资料
const fetchUserProfile = async () => {
  try {
    loading.value = true
    error.value = ''
    
    const response = await authApi.getProfile()
    console.log('API Response:', response) // 调试信息
    
    if (response.status === 'success') {
      console.log('User data:', response.data) // 调试信息
      Object.assign(profileData, response.data)
      Object.assign(originalProfileData, response.data)
    } else {
      error.value = response.message || '获取用户资料失败'
    }
  } catch (err) {
    console.error('Fetch profile error:', err) // 调试信息
    error.value = '获取用户资料失败: ' + (err.response?.data?.message || err.message)
  } finally {
    loading.value = false
  }
}

// 更新用户资料
const updateProfile = async () => {
  try {
    saving.value = true
    error.value = ''
    success.value = ''
    
    const response = await authApi.updateProfile({
      email: profileData.email,
      fullName: profileData.fullName,
      bio: profileData.bio,
      timezone: profileData.timezone,
      language: profileData.language,
      theme: profileData.theme
    })
    
    if (response.status === 'success') {
      success.value = '用户资料更新成功'
      Object.assign(originalProfileData, profileData)
      setTimeout(() => { success.value = '' }, 3000)
    } else {
      error.value = response.message || '更新用户资料失败'
    }
  } catch (err) {
    error.value = '更新用户资料失败: ' + (err.response?.data?.message || err.message)
  } finally {
    saving.value = false
  }
}

// 重置表单
const resetProfile = () => {
  Object.assign(profileData, originalProfileData)
  error.value = ''
  success.value = ''
}

// 处理头像上传
const handleAvatarChange = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // 检查文件大小
  if (file.size > 2 * 1024 * 1024) {
    error.value = '文件大小不能超过2MB'
    return
  }

  // 检查文件类型
  if (!['image/jpeg', 'image/png', 'image/gif'].includes(file.type)) {
    error.value = '只支持 JPG, PNG, GIF 格式的图片'
    return
  }

  try {
    uploading.value = true
    error.value = ''
    
    const formData = new FormData()
    formData.append('avatar', file)
    
    const response = await authApi.updateProfileWithAvatar(formData)
    
    if (response.status === 'success') {
      profileData.avatar = response.data.avatar
      success.value = '头像上传成功'
      setTimeout(() => { success.value = '' }, 3000)
    } else {
      error.value = response.message || '头像上传失败'
    }
  } catch (err) {
    error.value = '头像上传失败: ' + (err.response?.data?.message || err.message)
  } finally {
    uploading.value = false
    event.target.value = '' // 清空input
  }
}

// 移除头像
const removeAvatar = () => {
  profileData.avatar = '/default-avatar.png'
}

// 处理头像加载错误
const handleAvatarError = (event) => {
  event.target.src = '/default-avatar.png'
}

// 修改密码
const changePassword = async () => {
  if (passwordMismatch.value) {
    error.value = '密码不匹配'
    return
  }

  try {
    changingPassword.value = true
    error.value = ''
    success.value = ''
    
    const response = await authApi.changePassword({
      currentPassword: passwordData.currentPassword,
      newPassword: passwordData.newPassword
    })
    
    if (response.status === 'success') {
      success.value = '密码修改成功'
      // 清空密码表单
      passwordData.currentPassword = ''
      passwordData.newPassword = ''
      passwordData.confirmPassword = ''
      setTimeout(() => { success.value = '' }, 3000)
    } else {
      error.value = response.message || '密码修改失败'
    }
  } catch (err) {
    error.value = '密码修改失败: ' + (err.response?.data?.message || err.message)
  } finally {
    changingPassword.value = false
  }
}

// 页面加载时获取用户资料
onMounted(() => {
  fetchUserProfile()
})
</script> 