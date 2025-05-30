<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 页面头部 -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-6 py-4">
        <div class="flex items-center space-x-3">
          <router-link 
            to="/profile" 
            class="text-gray-500 hover:text-gray-700 transition-colors"
          >
            <i class="fas fa-arrow-left text-lg"></i>
          </router-link>
          <h1 class="text-2xl font-bold text-gray-900">资源点充值</h1>
        </div>
      </div>
    </div>

    <div class="max-w-4xl mx-auto px-6 py-8">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- 充值选择区域 -->
        <div class="lg:col-span-2">
          <!-- 当前余额 -->
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">当前余额</h2>
            <div class="flex items-center">
              <i class="fas fa-coins text-yellow-500 text-3xl mr-4"></i>
              <div>
                <span class="text-3xl font-bold text-gray-900">{{ currentBalance.toLocaleString() }}</span>
                <span class="text-sm text-gray-500 ml-2">点数</span>
              </div>
            </div>
          </div>

          <!-- 充值套餐 -->
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">选择充值套餐</h2>
            <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
              <div 
                v-for="package_ in rechargePackages" 
                :key="package_.id"
                @click="selectPackage(package_)"
                :class="[
                  'relative p-4 border-2 rounded-lg cursor-pointer transition-all',
                  selectedPackage?.id === package_.id 
                    ? 'border-indigo-500 bg-indigo-50' 
                    : 'border-gray-200 hover:border-gray-300'
                ]"
              >
                <!-- 推荐标签 -->
                <div v-if="package_.recommended" class="absolute -top-2 -right-2">
                  <span class="bg-red-500 text-white text-xs px-2 py-1 rounded-full">推荐</span>
                </div>
                
                <div class="text-center">
                  <div class="text-lg font-bold text-gray-900">¥{{ package_.price }}</div>
                  <div class="text-sm text-gray-600 mt-1">{{ package_.points.toLocaleString() }} 点数</div>
                  <div v-if="package_.bonus > 0" class="text-xs text-green-600 mt-1">
                    赠送 {{ package_.bonus.toLocaleString() }} 点数
                  </div>
                </div>

                <!-- 选中状态 -->
                <div v-if="selectedPackage?.id === package_.id" class="absolute top-2 right-2">
                  <i class="fas fa-check-circle text-indigo-500"></i>
                </div>
              </div>
            </div>

            <!-- 自定义金额 -->
            <div class="mt-6 p-4 border border-gray-200 rounded-lg">
              <div class="flex items-center mb-3">
                <input 
                  type="radio" 
                  id="custom" 
                  name="package" 
                  v-model="isCustomAmount" 
                  :value="true"
                  class="text-indigo-600 focus:ring-indigo-500"
                >
                <label for="custom" class="ml-2 text-sm font-medium text-gray-700">
                  自定义金额
                </label>
              </div>
              <div v-if="isCustomAmount" class="flex items-center space-x-4">
                <div class="flex-1">
                  <div class="relative">
                    <span class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">¥</span>
                    <input 
                      type="number" 
                      v-model="customAmount"
                      placeholder="请输入充值金额"
                      min="1"
                      max="10000"
                      class="w-full pl-8 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
                    >
                  </div>
                </div>
                <div class="text-sm text-gray-600">
                  = {{ Math.floor(customAmount * 100).toLocaleString() }} 点数
                </div>
              </div>
            </div>
          </div>

          <!-- 支付方式 -->
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">选择支付方式</h2>
            <div class="space-y-3">
              <div 
                v-for="method in paymentMethods" 
                :key="method.id"
                @click="selectedPaymentMethod = method"
                :class="[
                  'flex items-center p-4 border rounded-lg cursor-pointer transition-all',
                  selectedPaymentMethod?.id === method.id 
                    ? 'border-indigo-500 bg-indigo-50' 
                    : 'border-gray-200 hover:border-gray-300'
                ]"
              >
                <div class="flex-1 flex items-center">
                  <i :class="method.icon" class="text-2xl mr-3"></i>
                  <div>
                    <div class="font-medium text-gray-900">{{ method.name }}</div>
                    <div class="text-sm text-gray-500">{{ method.description }}</div>
                  </div>
                </div>
                <div v-if="selectedPaymentMethod?.id === method.id">
                  <i class="fas fa-check-circle text-indigo-500"></i>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 订单摘要 -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 sticky top-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">订单摘要</h2>
            
            <div class="space-y-3">
              <div class="flex justify-between text-sm">
                <span class="text-gray-600">充值金额</span>
                <span class="font-medium">¥{{ finalAmount }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-gray-600">获得点数</span>
                <span class="font-medium">{{ finalPoints.toLocaleString() }} 点数</span>
              </div>
              <div v-if="bonusPoints > 0" class="flex justify-between text-sm">
                <span class="text-green-600">赠送点数</span>
                <span class="font-medium text-green-600">+{{ bonusPoints.toLocaleString() }} 点数</span>
              </div>
              <div class="border-t border-gray-200 pt-3">
                <div class="flex justify-between">
                  <span class="text-gray-900 font-medium">总计点数</span>
                  <span class="text-lg font-bold text-indigo-600">{{ totalPoints.toLocaleString() }} 点数</span>
                </div>
              </div>
            </div>

            <div class="mt-6">
              <button 
                @click="handleRecharge"
                :disabled="!canProceed || isProcessing"
                :class="[
                  'w-full py-3 px-4 rounded-lg font-medium transition-colors',
                  canProceed && !isProcessing
                    ? 'bg-indigo-600 text-white hover:bg-indigo-700' 
                    : 'bg-gray-300 text-gray-500 cursor-not-allowed'
                ]"
              >
                <i v-if="isProcessing" class="fas fa-spinner fa-spin mr-2"></i>
                {{ isProcessing ? '处理中...' : `立即支付 ¥${finalAmount}` }}
              </button>
            </div>

            <!-- 安全提示 -->
            <div class="mt-4 p-3 bg-blue-50 rounded-lg">
              <div class="flex items-start">
                <i class="fas fa-shield-alt text-blue-500 mt-1 mr-2"></i>
                <div class="text-sm text-blue-700">
                  <div class="font-medium mb-1">安全保障</div>
                  <ul class="text-xs space-y-1">
                    <li>• 支付全程SSL加密</li>
                    <li>• 7×24小时服务监控</li>
                    <li>• 充值即时到账</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 支付确认对话框 -->
    <div v-if="showPaymentDialog" class="fixed inset-0 z-50 overflow-y-auto">
      <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="showPaymentDialog = false"></div>
        
        <div class="inline-block align-bottom bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full sm:p-6">
          <div class="sm:flex sm:items-start">
            <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-blue-100 sm:mx-0 sm:h-10 sm:w-10">
              <i class="fas fa-credit-card text-blue-600"></i>
            </div>
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
              <h3 class="text-lg leading-6 font-medium text-gray-900">
                确认支付
              </h3>
              <div class="mt-2">
                <p class="text-sm text-gray-500">
                  您即将支付 <strong>¥{{ finalAmount }}</strong>，获得 <strong>{{ totalPoints.toLocaleString() }} 点数</strong>
                </p>
                <p class="text-sm text-gray-500 mt-1">
                  支付方式：{{ selectedPaymentMethod?.name }}
                </p>
              </div>
            </div>
          </div>
          <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
            <button 
              type="button" 
              @click="confirmPayment"
              :disabled="isProcessing"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <i v-if="isProcessing" class="fas fa-spinner fa-spin mr-2"></i>
              {{ isProcessing ? '处理中...' : '确认支付' }}
            </button>
            <button 
              type="button" 
              @click="showPaymentDialog = false"
              :disabled="isProcessing"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              取消
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import { rechargeApi, pointsApi } from '@/api'

const router = useRouter()
const toast = useToast()

// 状态管理
const currentBalance = ref(0)
const selectedPackage = ref(null)
const isCustomAmount = ref(false)
const customAmount = ref('')
const selectedPaymentMethod = ref(null)
const showPaymentDialog = ref(false)
const isProcessing = ref(false)

// 充值套餐
const rechargePackages = ref([
  {
    id: 1,
    price: 10,
    points: 1000,
    bonus: 0,
    recommended: false
  },
  {
    id: 2,
    price: 50,
    points: 5000,
    bonus: 500,
    recommended: true
  },
  {
    id: 3,
    price: 100,
    points: 10000,
    bonus: 1500,
    recommended: false
  },
  {
    id: 4,
    price: 200,
    points: 20000,
    bonus: 4000,
    recommended: false
  },
  {
    id: 5,
    price: 500,
    points: 50000,
    bonus: 12500,
    recommended: false
  },
  {
    id: 6,
    price: 1000,
    points: 100000,
    bonus: 30000,
    recommended: false
  }
])

// 支付方式
const paymentMethods = ref([
  {
    id: 'alipay',
    name: '支付宝',
    description: '支持支付宝扫码支付',
    icon: 'fab fa-alipay text-blue-500'
  },
  {
    id: 'wechat',
    name: '微信支付',
    description: '支持微信扫码支付',
    icon: 'fab fa-weixin text-green-500'
  },
  {
    id: 'credit',
    name: '银行卡',
    description: '支持信用卡、储蓄卡',
    icon: 'fas fa-credit-card text-gray-600'
  }
])

// 计算属性
const finalAmount = computed(() => {
  if (isCustomAmount.value && customAmount.value) {
    return parseFloat(customAmount.value) || 0
  }
  return selectedPackage.value?.price || 0
})

const finalPoints = computed(() => {
  if (isCustomAmount.value && customAmount.value) {
    return Math.floor((parseFloat(customAmount.value) || 0) * 100)
  }
  return selectedPackage.value?.points || 0
})

const bonusPoints = computed(() => {
  if (isCustomAmount.value) {
    return 0
  }
  return selectedPackage.value?.bonus || 0
})

const totalPoints = computed(() => {
  return finalPoints.value + bonusPoints.value
})

const canProceed = computed(() => {
  const hasValidAmount = isCustomAmount.value 
    ? (customAmount.value && parseFloat(customAmount.value) >= 1) 
    : selectedPackage.value
  
  return hasValidAmount && selectedPaymentMethod.value
})

// 方法
const selectPackage = (package_) => {
  selectedPackage.value = package_
  isCustomAmount.value = false
  customAmount.value = ''
}

const handleRecharge = () => {
  if (!canProceed.value) return
  showPaymentDialog.value = true
}

const confirmPayment = async () => {
  if (isProcessing.value) return
  
  isProcessing.value = true
  
  try {
    const rechargeData = {
      amount: Math.round(finalAmount.value * 100), // 转换为分
      points: totalPoints.value,
      paymentMethod: selectedPaymentMethod.value.id,
      packageId: selectedPackage.value?.id || null
    }
    
    const response = await rechargeApi.createRecharge(rechargeData)
    
    if (response.status === 'success') {
      toast.success('充值订单创建成功，正在跳转到支付页面...')
      
      // 模拟支付流程（实际应该跳转到支付页面）
      setTimeout(() => {
        toast.success('充值成功！点数已到账')
        router.push('/profile')
      }, 2000)
    } else {
      throw new Error(response.message || '充值失败')
    }
  } catch (error) {
    console.error('充值失败:', error)
    toast.error(error.message || '充值失败，请稍后重试')
  } finally {
    isProcessing.value = false
    showPaymentDialog.value = false
  }
}

const fetchCurrentBalance = async () => {
  try {
    const response = await pointsApi.getPointsBalance()
    if (response.status === 'success') {
      currentBalance.value = response.data.points
    }
  } catch (error) {
    console.error('获取余额失败:', error)
  }
}

// 生命周期
onMounted(() => {
  fetchCurrentBalance()
  // 默认选择第一个支付方式
  selectedPaymentMethod.value = paymentMethods.value[0]
})
</script>

<style scoped>
/* 自定义样式 */
.sticky {
  position: -webkit-sticky;
  position: sticky;
}
</style> 