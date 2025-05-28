<template>
  <div class="container mx-auto px-4 py-6">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-2">积分充值</h1>
      <p class="text-gray-600">选择合适的充值套餐，为您的账户充值积分</p>
    </div>

    <!-- 当前余额显示 -->
    <div class="bg-gradient-to-r from-blue-500 to-blue-600 rounded-lg p-6 text-white mb-8">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-blue-100 text-sm">当前积分余额</p>
          <p class="text-3xl font-bold">{{ formatPoints(currentBalance) }}</p>
        </div>
        <div class="text-blue-200">
          <i class="fas fa-coins text-3xl"></i>
        </div>
      </div>
    </div>

    <!-- 充值套餐 -->
    <div class="bg-white rounded-lg shadow-sm border p-6 mb-8">
      <h3 class="text-lg font-medium text-gray-900 mb-6">选择充值套餐</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="pkg in rechargePackages" 
          :key="pkg.id"
          @click="selectPackage(pkg)"
          :class="[
            'relative cursor-pointer border-2 rounded-lg p-6 transition-all',
            selectedPackage?.id === pkg.id 
              ? 'border-blue-500 bg-blue-50' 
              : 'border-gray-200 hover:border-blue-300'
          ]"
        >
          <!-- 推荐标签 -->
          <div v-if="pkg.recommended" class="absolute -top-3 left-1/2 transform -translate-x-1/2">
            <span class="bg-orange-500 text-white text-xs px-3 py-1 rounded-full">推荐</span>
          </div>

          <!-- 套餐内容 -->
          <div class="text-center">
            <div class="text-2xl font-bold text-gray-900 mb-2">
              {{ formatPoints(pkg.points) }}
            </div>
            <div class="text-sm text-gray-500 mb-4">积分</div>
            <div class="text-xl font-semibold text-blue-600 mb-2">
              ¥{{ pkg.price }}
            </div>
            <div v-if="pkg.bonus > 0" class="text-sm text-green-600 mb-4">
              额外赠送 {{ formatPoints(pkg.bonus) }} 积分
            </div>
            <div class="text-xs text-gray-400">
              单价: ¥{{ (pkg.price / pkg.points).toFixed(4) }}/积分
            </div>
          </div>

          <!-- 选中标识 -->
          <div v-if="selectedPackage?.id === pkg.id" class="absolute top-2 right-2">
            <i class="fas fa-check-circle text-blue-500 text-xl"></i>
          </div>
        </div>
      </div>
    </div>

    <!-- 支付方式选择 -->
    <div v-if="selectedPackage" class="bg-white rounded-lg shadow-sm border p-6 mb-8">
      <h3 class="text-lg font-medium text-gray-900 mb-6">选择支付方式</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div 
          v-for="method in paymentMethods" 
          :key="method.id"
          @click="selectedPaymentMethod = method.id"
          :class="[
            'relative cursor-pointer border-2 rounded-lg p-4 transition-all',
            selectedPaymentMethod === method.id 
              ? 'border-blue-500 bg-blue-50' 
              : 'border-gray-200 hover:border-blue-300'
          ]"
        >
          <div class="flex items-center">
            <i :class="[method.icon, 'text-2xl mr-3', method.color]"></i>
            <div>
              <div class="font-medium text-gray-900">{{ method.name }}</div>
              <div class="text-sm text-gray-500">{{ method.description }}</div>
            </div>
          </div>
          
          <!-- 选中标识 -->
          <div v-if="selectedPaymentMethod === method.id" class="absolute top-2 right-2">
            <i class="fas fa-check-circle text-blue-500"></i>
          </div>
        </div>
      </div>
    </div>

    <!-- 订单确认 -->
    <div v-if="selectedPackage && selectedPaymentMethod" class="bg-white rounded-lg shadow-sm border p-6 mb-8">
      <h3 class="text-lg font-medium text-gray-900 mb-6">订单确认</h3>
      <div class="space-y-4">
        <div class="flex justify-between">
          <span class="text-gray-600">充值套餐:</span>
          <span class="font-medium">{{ formatPoints(selectedPackage.points) }} 积分</span>
        </div>
        <div v-if="selectedPackage.bonus > 0" class="flex justify-between">
          <span class="text-gray-600">赠送积分:</span>
          <span class="font-medium text-green-600">{{ formatPoints(selectedPackage.bonus) }} 积分</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-600">支付方式:</span>
          <span class="font-medium">{{ getPaymentMethodName(selectedPaymentMethod) }}</span>
        </div>
        <hr>
        <div class="flex justify-between text-lg font-semibold">
          <span>支付金额:</span>
          <span class="text-blue-600">¥{{ selectedPackage.price }}</span>
        </div>
      </div>

      <!-- 支付按钮 -->
      <div class="mt-6">
        <button 
          @click="handlePay"
          :disabled="paymentLoading"
          class="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-4 rounded-lg flex items-center justify-center disabled:opacity-50"
        >
          <span v-if="paymentLoading" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></span>
          {{ paymentLoading ? '正在处理...' : '立即支付' }}
        </button>
      </div>
    </div>

    <!-- 充值历史 -->
    <div class="bg-white rounded-lg shadow-sm border">
      <div class="px-6 py-4 border-b border-gray-200">
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">充值历史</h3>
          <button 
            @click="loadRechargeHistory" 
            class="text-blue-600 hover:text-blue-800 text-sm"
            :disabled="historyLoading"
          >
            <i class="fas fa-sync-alt" :class="{ 'animate-spin': historyLoading }"></i> 刷新
          </button>
        </div>
      </div>

      <div class="p-6">
        <!-- 加载状态 -->
        <div v-if="historyLoading" class="text-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
          <p class="text-gray-500 mt-2">加载中...</p>
        </div>

        <!-- 历史记录列表 -->
        <div v-else-if="rechargeHistory.length > 0" class="space-y-4">
          <div 
            v-for="record in rechargeHistory" 
            :key="record.id"
            class="flex items-center justify-between p-4 border border-gray-200 rounded-lg"
          >
            <div class="flex items-center">
              <div :class="[
                'w-10 h-10 rounded-full flex items-center justify-center mr-3',
                record.status === 'completed' ? 'bg-green-100 text-green-600' :
                record.status === 'pending' ? 'bg-yellow-100 text-yellow-600' :
                'bg-red-100 text-red-600'
              ]">
                <i :class="[
                  'fas',
                  record.status === 'completed' ? 'fa-check' :
                  record.status === 'pending' ? 'fa-clock' :
                  'fa-times'
                ]"></i>
              </div>
              <div>
                <p class="font-medium text-gray-900">
                  充值 {{ formatPoints(record.points) }} 积分
                </p>
                <p class="text-sm text-gray-500">
                  {{ formatDate(record.createdAt) }} · {{ getStatusText(record.status) }}
                </p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-semibold text-gray-900">¥{{ record.amount }}</p>
              <p class="text-sm text-gray-500">{{ record.paymentMethod }}</p>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="text-center py-8">
          <i class="fas fa-receipt text-4xl text-gray-300 mb-4"></i>
          <p class="text-gray-500">暂无充值记录</p>
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
import { rechargeApi, pointsApi } from '../../api/index'

export default {
  name: 'RechargeView',
  setup() {
    // 状态管理
    const paymentLoading = ref(false)
    const historyLoading = ref(false)
    const currentBalance = ref(0)

    // 选择状态
    const selectedPackage = ref(null)
    const selectedPaymentMethod = ref('')

    // 充值套餐
    const rechargePackages = ref([
      {
        id: 1,
        points: 1000,
        price: 10,
        bonus: 0,
        recommended: false
      },
      {
        id: 2,
        points: 5000,
        price: 45,
        bonus: 500,
        recommended: true
      },
      {
        id: 3,
        points: 10000,
        price: 80,
        bonus: 1500,
        recommended: false
      },
      {
        id: 4,
        points: 20000,
        price: 150,
        bonus: 4000,
        recommended: false
      },
      {
        id: 5,
        points: 50000,
        price: 350,
        bonus: 12500,
        recommended: false
      },
      {
        id: 6,
        points: 100000,
        price: 650,
        bonus: 30000,
        recommended: false
      }
    ])

    // 支付方式
    const paymentMethods = ref([
      {
        id: 'alipay',
        name: '支付宝',
        description: '安全便捷',
        icon: 'fab fa-alipay',
        color: 'text-blue-500'
      },
      {
        id: 'wechat',
        name: '微信支付',
        description: '快速支付',
        icon: 'fab fa-weixin',
        color: 'text-green-500'
      },
      {
        id: 'bank',
        name: '银行卡',
        description: '网银支付',
        icon: 'fas fa-credit-card',
        color: 'text-gray-500'
      }
    ])

    // 充值历史
    const rechargeHistory = ref([])

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

    // 获取支付方式名称
    const getPaymentMethodName = (id) => {
      const method = paymentMethods.value.find(m => m.id === id)
      return method?.name || ''
    }

    // 获取状态文本
    const getStatusText = (status) => {
      const statusMap = {
        'pending': '处理中',
        'completed': '已完成',
        'failed': '失败',
        'cancelled': '已取消'
      }
      return statusMap[status] || status
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

    // 选择套餐
    const selectPackage = (pkg) => {
      selectedPackage.value = pkg
    }

    // 加载当前余额
    const loadCurrentBalance = async () => {
      try {
        const response = await pointsApi.getPoints()
        currentBalance.value = response.data?.balance || 0
      } catch (error) {
        console.error('加载余额失败:', error)
      }
    }

    // 处理支付
    const handlePay = async () => {
      if (!selectedPackage.value || !selectedPaymentMethod.value) {
        showNotification('error', '支付失败', '请选择充值套餐和支付方式')
        return
      }

      paymentLoading.value = true
      
      try {
        const paymentData = {
          packageId: selectedPackage.value.id,
          paymentMethod: selectedPaymentMethod.value,
          amount: selectedPackage.value.price,
          points: selectedPackage.value.points + selectedPackage.value.bonus
        }

        const response = await rechargeApi.createRecharge(paymentData)
        
        showNotification('success', '支付成功', '积分已充值到您的账户')
        
        // 重新加载余额和历史
        loadCurrentBalance()
        loadRechargeHistory()
        
        // 清除选择
        selectedPackage.value = null
        selectedPaymentMethod.value = ''
        
      } catch (error) {
        console.error('支付失败:', error)
        
        let errorMessage = '支付失败，请稍后重试'
        if (error.response && error.response.data && error.response.data.message) {
          errorMessage = error.response.data.message
        }
        
        showNotification('error', '支付失败', errorMessage)
      } finally {
        paymentLoading.value = false
      }
    }

    // 加载充值历史
    const loadRechargeHistory = async () => {
      historyLoading.value = true
      try {
        const response = await rechargeApi.getRechargeHistory({
          page: 1,
          limit: 20
        })
        rechargeHistory.value = response.data?.records || []
      } catch (error) {
        console.error('加载充值历史失败:', error)
        showNotification('error', '加载失败', '无法加载充值历史')
      } finally {
        historyLoading.value = false
      }
    }

    // 生命周期
    onMounted(() => {
      loadCurrentBalance()
      loadRechargeHistory()
    })

    return {
      paymentLoading,
      historyLoading,
      currentBalance,
      selectedPackage,
      selectedPaymentMethod,
      rechargePackages,
      paymentMethods,
      rechargeHistory,
      notification,
      formatPoints,
      formatDate,
      getPaymentMethodName,
      getStatusText,
      selectPackage,
      handlePay,
      loadRechargeHistory
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