<template>
  <div class="container mx-auto px-4 py-6">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-2">服务购买</h1>
      <p class="text-gray-600">使用积分购买各种AI服务和功能</p>
    </div>

    <!-- 当前积分显示 -->
    <div class="bg-gradient-to-r from-green-500 to-green-600 rounded-lg p-6 text-white mb-8">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-green-100 text-sm">可用积分余额</p>
          <p class="text-3xl font-bold">{{ formatPoints(currentBalance) }}</p>
        </div>
        <div class="text-green-200">
          <i class="fas fa-wallet text-3xl"></i>
        </div>
      </div>
      <div class="mt-4">
        <router-link 
          to="/recharge" 
          class="inline-flex items-center text-sm text-green-100 hover:text-white"
        >
          <i class="fas fa-plus-circle mr-1"></i> 充值积分
        </router-link>
      </div>
    </div>

    <!-- 服务分类 -->
    <div class="mb-8">
      <div class="flex flex-wrap gap-4">
        <button 
          v-for="category in serviceCategories" 
          :key="category.id"
          @click="selectedCategory = category.id"
          :class="[
            'px-4 py-2 rounded-lg font-medium transition-colors',
            selectedCategory === category.id 
              ? 'bg-blue-600 text-white' 
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          ]"
        >
          <i :class="[category.icon, 'mr-2']"></i>
          {{ category.name }}
        </button>
      </div>
    </div>

    <!-- 服务列表 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
      <div 
        v-for="service in filteredServices" 
        :key="service.id"
        class="bg-white rounded-lg shadow-sm border hover:shadow-md transition-shadow"
      >
        <div class="p-6">
          <!-- 服务头部 -->
          <div class="flex items-start justify-between mb-4">
            <div class="flex items-center">
              <div :class="[
                'w-12 h-12 rounded-lg flex items-center justify-center mr-3',
                service.color
              ]">
                <i :class="[service.icon, 'text-white text-xl']"></i>
              </div>
              <div>
                <h3 class="font-medium text-gray-900">{{ service.name }}</h3>
                <p class="text-sm text-gray-500">{{ service.category }}</p>
              </div>
            </div>
            <div v-if="service.popular" class="bg-orange-100 text-orange-800 text-xs px-2 py-1 rounded-full">
              热门
            </div>
          </div>

          <!-- 服务描述 -->
          <p class="text-gray-600 text-sm mb-4">{{ service.description }}</p>

          <!-- 服务特性 -->
          <div class="mb-4">
            <ul class="space-y-1">
              <li 
                v-for="feature in service.features" 
                :key="feature"
                class="text-sm text-gray-600 flex items-center"
              >
                <i class="fas fa-check text-green-500 text-xs mr-2"></i>
                {{ feature }}
              </li>
            </ul>
          </div>

          <!-- 价格和购买 -->
          <div class="flex items-center justify-between">
            <div>
              <span class="text-lg font-semibold text-gray-900">{{ formatPoints(service.price) }}</span>
              <span class="text-sm text-gray-500 ml-1">积分</span>
              <div v-if="service.originalPrice && service.originalPrice > service.price" class="text-xs text-gray-400 line-through">
                原价 {{ formatPoints(service.originalPrice) }} 积分
              </div>
            </div>
            <button 
              @click="purchaseService(service)"
              :disabled="currentBalance < service.price"
              :class="[
                'px-4 py-2 rounded-lg font-medium transition-colors',
                currentBalance >= service.price 
                  ? 'bg-blue-600 hover:bg-blue-700 text-white' 
                  : 'bg-gray-300 text-gray-500 cursor-not-allowed'
              ]"
            >
              {{ currentBalance >= service.price ? '购买' : '积分不足' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 购买历史 -->
    <div class="bg-white rounded-lg shadow-sm border">
      <div class="px-6 py-4 border-b border-gray-200">
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">购买历史</h3>
          <div class="flex items-center space-x-4">
            <!-- 状态筛选 -->
            <select 
              v-model="filters.status" 
              class="border border-gray-300 rounded-md px-3 py-1 text-sm"
              @change="loadPurchaseHistory"
            >
              <option value="">全部状态</option>
              <option value="completed">已完成</option>
              <option value="pending">处理中</option>
              <option value="failed">失败</option>
            </select>
            
            <button 
              @click="loadPurchaseHistory" 
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
        <div v-else-if="purchaseHistory.length > 0" class="space-y-4">
          <div 
            v-for="record in purchaseHistory" 
            :key="record.id"
            class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50"
          >
            <div class="flex items-center">
              <div :class="[
                'w-10 h-10 rounded-lg flex items-center justify-center mr-3',
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
                <p class="font-medium text-gray-900">{{ record.serviceName }}</p>
                <p class="text-sm text-gray-500">
                  {{ formatDate(record.createdAt) }} · {{ getStatusText(record.status) }}
                </p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-semibold text-gray-900">{{ formatPoints(record.price) }} 积分</p>
              <p class="text-sm text-gray-500">{{ record.category }}</p>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="text-center py-8">
          <i class="fas fa-shopping-bag text-4xl text-gray-300 mb-4"></i>
          <p class="text-gray-500">暂无购买记录</p>
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

    <!-- 购买确认模态框 -->
    <div v-if="confirmModal.show" class="fixed inset-0 bg-gray-800 bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden">
        <!-- 模态框头部 -->
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">确认购买</h3>
        </div>
        
        <!-- 模态框内容 -->
        <div class="p-6">
          <div v-if="confirmModal.service" class="space-y-4">
            <div class="flex items-center">
              <div :class="[
                'w-12 h-12 rounded-lg flex items-center justify-center mr-3',
                confirmModal.service.color
              ]">
                <i :class="[confirmModal.service.icon, 'text-white text-xl']"></i>
              </div>
              <div>
                <h4 class="font-medium text-gray-900">{{ confirmModal.service.name }}</h4>
                <p class="text-sm text-gray-500">{{ confirmModal.service.category }}</p>
              </div>
            </div>
            
            <div class="bg-gray-50 rounded-lg p-4">
              <div class="flex justify-between items-center mb-2">
                <span class="text-gray-600">服务价格:</span>
                <span class="font-medium">{{ formatPoints(confirmModal.service.price) }} 积分</span>
              </div>
              <div class="flex justify-between items-center mb-2">
                <span class="text-gray-600">当前余额:</span>
                <span class="font-medium">{{ formatPoints(currentBalance) }} 积分</span>
              </div>
              <hr class="my-2">
              <div class="flex justify-between items-center">
                <span class="text-gray-600">购买后余额:</span>
                <span class="font-medium text-blue-600">{{ formatPoints(currentBalance - confirmModal.service.price) }} 积分</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 模态框底部 -->
        <div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end">
          <button 
            @click="confirmModal.show = false" 
            class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-800 rounded-lg mr-2"
            :disabled="purchaseLoading"
          >
            取消
          </button>
          <button 
            @click="confirmPurchase"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg flex items-center"
            :disabled="purchaseLoading"
          >
            <span v-if="purchaseLoading" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></span>
            {{ purchaseLoading ? '购买中...' : '确认购买' }}
          </button>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { purchaseApi, pointsApi } from '../../api/index'

export default {
  name: 'PurchaseView',
  setup() {
    // 状态管理
    const historyLoading = ref(false)
    const purchaseLoading = ref(false)
    const currentBalance = ref(0)
    const selectedCategory = ref('all')

    // 服务分类
    const serviceCategories = ref([
      { id: 'all', name: '全部服务', icon: 'fas fa-th-large' },
      { id: 'ai_model', name: 'AI模型', icon: 'fas fa-brain' },
      { id: 'storage', name: '存储服务', icon: 'fas fa-database' },
      { id: 'compute', name: '计算资源', icon: 'fas fa-server' },
      { id: 'api', name: 'API调用', icon: 'fas fa-plug' },
      { id: 'premium', name: '高级功能', icon: 'fas fa-crown' }
    ])

    // 可购买服务
    const services = ref([
      {
        id: 1,
        name: 'GPT-4 调用包',
        category: 'AI模型',
        categoryId: 'ai_model',
        description: '1000次GPT-4 API调用额度，适合高质量文本生成需求',
        features: ['1000次API调用', '支持最新GPT-4模型', '30天有效期'],
        price: 5000,
        originalPrice: 6000,
        icon: 'fas fa-robot',
        color: 'bg-blue-500',
        popular: true
      },
      {
        id: 2,
        name: '云存储空间',
        category: '存储服务',
        categoryId: 'storage',
        description: '100GB云端存储空间，用于保存工作流和数据文件',
        features: ['100GB存储空间', '高速访问', '自动备份'],
        price: 2000,
        icon: 'fas fa-cloud',
        color: 'bg-green-500',
        popular: false
      },
      {
        id: 3,
        name: '高级计算资源',
        category: '计算资源',
        categoryId: 'compute',
        description: '专用计算节点，提供更强的工作流执行性能',
        features: ['专用CPU资源', '优先队列', '高内存配置'],
        price: 8000,
        icon: 'fas fa-microchip',
        color: 'bg-purple-500',
        popular: false
      },
      {
        id: 4,
        name: 'API访问包',
        category: 'API调用',
        categoryId: 'api',
        description: '第三方API调用额度，支持各种外部服务集成',
        features: ['10000次API调用', '支持多种第三方服务', '无限制频率'],
        price: 3000,
        icon: 'fas fa-exchange-alt',
        color: 'bg-orange-500',
        popular: true
      },
      {
        id: 5,
        name: '高级分析功能',
        category: '高级功能',
        categoryId: 'premium',
        description: '解锁高级数据分析和可视化功能',
        features: ['高级图表', '数据导出', '自定义报告'],
        price: 4500,
        icon: 'fas fa-chart-line',
        color: 'bg-indigo-500',
        popular: false
      },
      {
        id: 6,
        name: '批量处理套餐',
        category: 'AI模型',
        categoryId: 'ai_model',
        description: '大批量数据处理，适合企业级应用',
        features: ['无限批量处理', '并行执行', '优先处理'],
        price: 12000,
        originalPrice: 15000,
        icon: 'fas fa-layer-group',
        color: 'bg-red-500',
        popular: false
      }
    ])

    // 购买历史
    const purchaseHistory = ref([])

    // 筛选条件
    const filters = reactive({
      status: ''
    })

    // 分页
    const pagination = reactive({
      page: 1,
      limit: 20,
      total: 0
    })

    // 确认购买模态框
    const confirmModal = reactive({
      show: false,
      service: null
    })

    // 通知
    const notification = reactive({
      show: false,
      type: 'info',
      title: '',
      message: ''
    })

    // 计算筛选后的服务
    const filteredServices = computed(() => {
      if (selectedCategory.value === 'all') {
        return services.value
      }
      return services.value.filter(service => service.categoryId === selectedCategory.value)
    })

    // 格式化积分
    const formatPoints = (points) => {
      return points?.toLocaleString() || '0'
    }

    // 格式化日期
    const formatDate = (date) => {
      return new Date(date).toLocaleString('zh-CN')
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

    // 加载当前余额
    const loadCurrentBalance = async () => {
      try {
        const response = await pointsApi.getPoints()
        currentBalance.value = response.data?.balance || 0
      } catch (error) {
        console.error('加载余额失败:', error)
      }
    }

    // 购买服务
    const purchaseService = (service) => {
      if (currentBalance.value < service.price) {
        showNotification('error', '积分不足', '请先充值积分')
        return
      }
      
      confirmModal.service = service
      confirmModal.show = true
    }

    // 确认购买
    const confirmPurchase = async () => {
      if (!confirmModal.service) return

      purchaseLoading.value = true
      
      try {
        const purchaseData = {
          serviceId: confirmModal.service.id,
          serviceName: confirmModal.service.name,
          category: confirmModal.service.category,
          price: confirmModal.service.price
        }

        await purchaseApi.createPurchase(purchaseData)
        
        showNotification('success', '购买成功', `成功购买 ${confirmModal.service.name}`)
        
        // 重新加载余额和历史
        loadCurrentBalance()
        loadPurchaseHistory()
        
        // 关闭模态框
        confirmModal.show = false
        confirmModal.service = null
        
      } catch (error) {
        console.error('购买失败:', error)
        
        let errorMessage = '购买失败，请稍后重试'
        if (error.response && error.response.data && error.response.data.message) {
          errorMessage = error.response.data.message
        }
        
        showNotification('error', '购买失败', errorMessage)
      } finally {
        purchaseLoading.value = false
      }
    }

    // 加载购买历史
    const loadPurchaseHistory = async () => {
      historyLoading.value = true
      try {
        const params = {
          page: pagination.page,
          limit: pagination.limit
        }
        
        if (filters.status) {
          params.status = filters.status
        }
        
        const response = await purchaseApi.getPurchaseHistory(params)
        purchaseHistory.value = response.data?.records || []
        pagination.total = response.data?.total || 0
      } catch (error) {
        console.error('加载购买历史失败:', error)
        showNotification('error', '加载失败', '无法加载购买历史')
      } finally {
        historyLoading.value = false
      }
    }

    // 分页处理
    const changePage = (page) => {
      if (page >= 1 && page <= Math.ceil(pagination.total / pagination.limit)) {
        pagination.page = page
        loadPurchaseHistory()
      }
    }

    // 生命周期
    onMounted(() => {
      loadCurrentBalance()
      loadPurchaseHistory()
    })

    return {
      historyLoading,
      purchaseLoading,
      currentBalance,
      selectedCategory,
      serviceCategories,
      filteredServices,
      purchaseHistory,
      filters,
      pagination,
      confirmModal,
      notification,
      formatPoints,
      formatDate,
      getStatusText,
      purchaseService,
      confirmPurchase,
      loadPurchaseHistory,
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