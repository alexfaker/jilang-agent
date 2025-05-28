<template>
  <div class="container mx-auto px-4 py-6">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-2">积分管理</h1>
      <p class="text-gray-600">查看您的积分余额和使用历史</p>
    </div>

    <!-- 积分概览卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <!-- 当前余额 -->
      <div class="bg-gradient-to-r from-blue-500 to-blue-600 rounded-lg p-6 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-blue-100 text-sm">当前余额</p>
            <p class="text-3xl font-bold">{{ formatPoints(pointsData.balance) }}</p>
          </div>
          <div class="text-blue-200">
            <i class="fas fa-coins text-3xl"></i>
          </div>
        </div>
        <div class="mt-4">
          <router-link 
            to="/recharge" 
            class="inline-flex items-center text-sm text-blue-100 hover:text-white"
          >
            <i class="fas fa-plus-circle mr-1"></i> 充值积分
          </router-link>
        </div>
      </div>

      <!-- 今日消耗 -->
      <div class="bg-white rounded-lg p-6 shadow-sm border">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">今日消耗</p>
            <p class="text-2xl font-semibold text-gray-900">{{ formatPoints(pointsData.todayUsed) }}</p>
          </div>
          <div class="text-orange-500">
            <i class="fas fa-arrow-down text-2xl"></i>
          </div>
        </div>
      </div>

      <!-- 累计获得 -->
      <div class="bg-white rounded-lg p-6 shadow-sm border">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">累计获得</p>
            <p class="text-2xl font-semibold text-gray-900">{{ formatPoints(pointsData.totalEarned) }}</p>
          </div>
          <div class="text-green-500">
            <i class="fas fa-arrow-up text-2xl"></i>
          </div>
        </div>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="bg-white rounded-lg shadow-sm border p-6 mb-8">
      <h3 class="text-lg font-medium text-gray-900 mb-4">快捷操作</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <router-link 
          to="/recharge"
          class="flex items-center justify-center p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors"
        >
          <div class="text-center">
            <i class="fas fa-credit-card text-2xl text-blue-500 mb-2"></i>
            <p class="text-sm font-medium text-gray-700">充值积分</p>
          </div>
        </router-link>
        
        <router-link 
          to="/purchase"
          class="flex items-center justify-center p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-green-500 hover:bg-green-50 transition-colors"
        >
          <div class="text-center">
            <i class="fas fa-shopping-cart text-2xl text-green-500 mb-2"></i>
            <p class="text-sm font-medium text-gray-700">购买服务</p>
          </div>
        </router-link>
        
        <button 
          @click="showGiftModal = true"
          class="flex items-center justify-center p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-purple-500 hover:bg-purple-50 transition-colors"
        >
          <div class="text-center">
            <i class="fas fa-gift text-2xl text-purple-500 mb-2"></i>
            <p class="text-sm font-medium text-gray-700">每日签到</p>
          </div>
        </button>
        
        <button 
          @click="showInviteModal = true"
          class="flex items-center justify-center p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-yellow-500 hover:bg-yellow-50 transition-colors"
        >
          <div class="text-center">
            <i class="fas fa-user-plus text-2xl text-yellow-500 mb-2"></i>
            <p class="text-sm font-medium text-gray-700">邀请好友</p>
          </div>
        </button>
      </div>
    </div>

    <!-- 积分历史 -->
    <div class="bg-white rounded-lg shadow-sm border">
      <div class="px-6 py-4 border-b border-gray-200">
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">积分历史</h3>
          <div class="flex items-center space-x-4">
            <!-- 筛选器 -->
            <select 
              v-model="filters.type" 
              class="border border-gray-300 rounded-md px-3 py-1 text-sm"
              @change="loadPointsHistory"
            >
              <option value="">全部类型</option>
              <option value="earn">获得</option>
              <option value="spend">消耗</option>
              <option value="refund">退还</option>
            </select>
            
            <button 
              @click="loadPointsHistory" 
              class="text-blue-600 hover:text-blue-800 text-sm"
              :disabled="historyLoading"
            >
              <i class="fas fa-sync-alt" :class="{ 'animate-spin': historyLoading }"></i> 刷新
            </button>
          </div>
        </div>
      </div>

      <div class="p-6">
        <!-- 加载状态 -->
        <div v-if="historyLoading" class="text-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
          <p class="text-gray-500 mt-2">加载中...</p>
        </div>

        <!-- 历史记录列表 -->
        <div v-else-if="pointsHistory.length > 0" class="space-y-4">
          <div 
            v-for="record in pointsHistory" 
            :key="record.id"
            class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50"
          >
            <div class="flex items-center">
              <div :class="[
                'w-10 h-10 rounded-full flex items-center justify-center mr-3',
                record.type === 'earn' ? 'bg-green-100 text-green-600' :
                record.type === 'spend' ? 'bg-red-100 text-red-600' :
                'bg-blue-100 text-blue-600'
              ]">
                <i :class="[
                  'fas',
                  record.type === 'earn' ? 'fa-plus' :
                  record.type === 'spend' ? 'fa-minus' :
                  'fa-undo'
                ]"></i>
              </div>
              <div>
                <p class="font-medium text-gray-900">{{ record.description }}</p>
                <p class="text-sm text-gray-500">{{ formatDate(record.createdAt) }}</p>
              </div>
            </div>
            <div class="text-right">
              <p :class="[
                'font-semibold',
                record.type === 'earn' ? 'text-green-600' :
                record.type === 'spend' ? 'text-red-600' :
                'text-blue-600'
              ]">
                {{ record.type === 'spend' ? '-' : '+' }}{{ formatPoints(record.amount) }}
              </p>
              <p class="text-sm text-gray-500">
                余额: {{ formatPoints(record.balanceAfter) }}
              </p>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="text-center py-8">
          <i class="fas fa-history text-4xl text-gray-300 mb-4"></i>
          <p class="text-gray-500">暂无积分历史记录</p>
        </div>

        <!-- 分页 -->
        <div v-if="pagination.total > pagination.limit" class="mt-6 flex justify-center">
          <nav class="flex items-center space-x-2">
            <button 
              @click="changePage(pagination.page - 1)"
              :disabled="pagination.page === 1"
              class="px-3 py-1 border border-gray-300 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50"
            >
              上一页
            </button>
            <span class="text-sm text-gray-600">
              第 {{ pagination.page }} 页，共 {{ Math.ceil(pagination.total / pagination.limit) }} 页
            </span>
            <button 
              @click="changePage(pagination.page + 1)"
              :disabled="pagination.page >= Math.ceil(pagination.total / pagination.limit)"
              class="px-3 py-1 border border-gray-300 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50"
            >
              下一页
            </button>
          </nav>
        </div>
      </div>
    </div>

    <!-- 通知提示 -->
    <div v-if="notification.show" class="fixed top-4 right-4 max-w-md z-50">
      <div :class="[
        'p-4 rounded-lg shadow-lg',
        notification.type === 'success' ? 'bg-green-100 text-green-800 border-l-4 border-green-500' : 
        'bg-red-100 text-red-800 border-l-4 border-red-500'
      ]">
        <div class="flex items-start">
          <i :class="[
            'fas mt-0.5 mr-3',
            notification.type === 'success' ? 'fa-check-circle text-green-500' : 'fa-exclamation-circle text-red-500'
          ]"></i>
          <div class="flex-1">
            <p class="font-medium">{{ notification.title }}</p>
            <p class="text-sm">{{ notification.message }}</p>
          </div>
          <button @click="notification.show = false" class="ml-4 text-gray-500 hover:text-gray-700">
            <i class="fas fa-times"></i>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { pointsApi } from '../../api/index'

export default {
  name: 'PointsView',
  setup() {
    // 状态管理
    const historyLoading = ref(false)
    const showGiftModal = ref(false)
    const showInviteModal = ref(false)

    // 积分数据
    const pointsData = reactive({
      balance: 0,
      todayUsed: 0,
      totalEarned: 0
    })

    // 积分历史
    const pointsHistory = ref([])

    // 筛选条件
    const filters = reactive({
      type: ''
    })

    // 分页
    const pagination = reactive({
      page: 1,
      limit: 20,
      total: 0
    })

    // 通知
    const notification = reactive({
      show: false,
      type: 'info',
      title: '',
      message: ''
    })

    // 格式化积分
    const formatPoints = (points) => {
      return points?.toLocaleString() || '0'
    }

    // 格式化日期
    const formatDate = (date) => {
      return new Date(date).toLocaleString('zh-CN')
    }

    // 显示通知
    const showNotification = (type, title, message) => {
      notification.show = true
      notification.type = type
      notification.title = title
      notification.message = message
      
      setTimeout(() => {
        notification.show = false
      }, 5000)
    }

    // 加载积分信息
    const loadPointsData = async () => {
      try {
        const response = await pointsApi.getPoints()
        Object.assign(pointsData, response.data)
      } catch (error) {
        console.error('加载积分信息失败:', error)
        showNotification('error', '加载失败', '无法加载积分信息')
      }
    }

    // 加载积分历史
    const loadPointsHistory = async () => {
      historyLoading.value = true
      try {
        const params = {
          page: pagination.page,
          limit: pagination.limit
        }
        
        if (filters.type) {
          params.type = filters.type
        }
        
        const response = await pointsApi.getPointsHistory(params)
        pointsHistory.value = response.data?.records || []
        pagination.total = response.data?.total || 0
      } catch (error) {
        console.error('加载积分历史失败:', error)
        showNotification('error', '加载失败', '无法加载积分历史')
      } finally {
        historyLoading.value = false
      }
    }

    // 分页处理
    const changePage = (page) => {
      if (page >= 1 && page <= Math.ceil(pagination.total / pagination.limit)) {
        pagination.page = page
        loadPointsHistory()
      }
    }

    // 生命周期
    onMounted(() => {
      loadPointsData()
      loadPointsHistory()
    })

    return {
      historyLoading,
      showGiftModal,
      showInviteModal,
      pointsData,
      pointsHistory,
      filters,
      pagination,
      notification,
      formatPoints,
      formatDate,
      loadPointsHistory,
      changePage
    }
  }
}
</script>

<style scoped>
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style> 