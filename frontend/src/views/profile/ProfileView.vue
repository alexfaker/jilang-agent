<template>
  <div class="flex-1 overflow-y-auto bg-gray-50">
    <!-- 页面头部 -->
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-6 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">个人中心</h1>
            <p class="text-gray-600 mt-1">管理您的账户信息和资源</p>
          </div>
          
          <!-- 用户积分显示 -->
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-2">
              <i class="fas fa-coins text-yellow-500 text-xl"></i>
              <span class="text-2xl font-bold text-gray-900">{{ pointsData.currentBalance || userProfile.points || 0 }}</span>
              <span class="text-sm text-gray-500">积分</span>
            </div>
            <button class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors">
              <i class="fas fa-bell mr-2"></i>
              通知
            </button>
            <button class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors">
              <i class="fas fa-cog mr-2"></i>
              设置
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 标签页导航 -->
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-6">
        <nav class="flex space-x-8">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'py-4 px-1 border-b-2 font-medium text-sm transition-colors',
              activeTab === tab.id
                ? 'border-indigo-500 text-indigo-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            <i :class="tab.icon + ' mr-2'"></i>
            {{ tab.name }}
          </button>
        </nav>
      </div>
    </div>

    <!-- 标签页内容 -->
    <div class="max-w-7xl mx-auto px-6 py-8">
      <!-- 资源点数管理 -->
      <div v-if="activeTab === 'points'" class="space-y-6">
        <!-- 积分概览卡片 -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- 当前余额 -->
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div class="flex items-center justify-between">
              <div>
                <h3 class="text-lg font-semibold text-gray-900 mb-2">当前余额</h3>
                <div class="flex items-center">
                  <i class="fas fa-coins text-yellow-500 text-2xl mr-3"></i>
                  <span v-if="loading.statistics" class="text-2xl font-bold text-gray-400">
                    <i class="fas fa-spinner fa-spin"></i>
                  </span>
                  <span v-else class="text-3xl font-bold text-gray-900">{{ pointsData.currentBalance.toLocaleString() }}</span>
                </div>
              </div>
              <button 
                @click="$router.push('/recharge')"
                class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
              >
                充值
              </button>
            </div>
          </div>

          <!-- 本月消耗 -->
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-2">本月消耗</h3>
            <div class="flex items-center justify-between">
              <div class="flex items-center">
                <i class="fas fa-chart-line text-green-500 text-2xl mr-3"></i>
                <span v-if="loading.statistics" class="text-2xl font-bold text-gray-400">
                  <i class="fas fa-spinner fa-spin"></i>
                </span>
                <span v-else class="text-3xl font-bold text-gray-900">{{ pointsData.monthlyConsumption.toLocaleString() }}</span>
              </div>
              <span v-if="!loading.statistics" class="text-sm text-green-600 bg-green-100 px-2 py-1 rounded-full">
                近期活跃
              </span>
            </div>
          </div>

          <!-- 累计充值 -->
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-2">累计充值</h3>
            <div class="flex items-center justify-between">
              <div class="flex items-center">
                <i class="fas fa-credit-card text-purple-500 text-2xl mr-3"></i>
                <span v-if="loading.statistics" class="text-2xl font-bold text-gray-400">
                  <i class="fas fa-spinner fa-spin"></i>
                </span>
                <span v-else class="text-3xl font-bold text-gray-900">{{ pointsData.totalRecharge.toLocaleString() }}</span>
              </div>
              <span v-if="!loading.statistics" class="text-sm text-gray-500">
                共充值 {{ pointsData.rechargeCount }} 次
              </span>
            </div>
          </div>
        </div>

        <!-- 资源点数消费明细 -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200">
          <div class="p-6 border-b border-gray-200">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900">资源点数消费明细</h3>
              <div class="flex items-center space-x-4">
                <select class="border border-gray-300 rounded-lg px-3 py-2 text-sm">
                  <option>最近7天</option>
                  <option>最近30天</option>
                  <option>最近90天</option>
                </select>
                <button class="px-3 py-2 text-sm text-gray-600 hover:text-gray-900 border border-gray-300 rounded-lg">
                  <i class="fas fa-download mr-1"></i>
                  导出
                </button>
              </div>
            </div>
          </div>

          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">日期</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">描述</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">工作流</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">点数变化</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">余额</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <!-- 加载状态 -->
                <tr v-if="loading.transactions">
                  <td colspan="5" class="px-6 py-8 text-center">
                    <i class="fas fa-spinner fa-spin text-gray-400 text-xl mr-2"></i>
                    <span class="text-gray-500">正在加载交易记录...</span>
                  </td>
                </tr>
                
                <!-- 空状态 -->
                <tr v-else-if="pointsHistory.length === 0">
                  <td colspan="5" class="px-6 py-8 text-center">
                    <i class="fas fa-history text-gray-300 text-3xl mb-2"></i>
                    <p class="text-gray-500">暂无交易记录</p>
                  </td>
                </tr>
                
                <!-- 交易记录 -->
                <tr v-else v-for="record in pointsHistory" :key="record.id">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ record.date }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ record.description }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ record.workflow }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    <span :class="record.change > 0 ? 'text-green-600' : 'text-red-600'">
                      {{ record.change > 0 ? '+' : '' }}{{ record.change.toLocaleString() }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ record.balance.toLocaleString() }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- 分页 -->
          <div class="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
            <div class="text-sm text-gray-500">
              <span v-if="loading.transactions">加载中...</span>
              <span v-else>
                显示 {{ Math.min((pagination.page - 1) * pagination.pageSize + 1, pagination.total) }} 
                至 {{ Math.min(pagination.page * pagination.pageSize, pagination.total) }} 条，
                共 {{ pagination.total }} 条记录
              </span>
            </div>
            <div v-if="!loading.transactions && pagination.pages > 1" class="flex items-center space-x-2">
              <button 
                @click="handlePageChange(pagination.page - 1)"
                :disabled="pagination.page <= 1"
                :class="[
                  'px-3 py-1 text-sm border border-gray-300 rounded',
                  pagination.page <= 1 
                    ? 'text-gray-400 cursor-not-allowed' 
                    : 'text-gray-600 hover:text-gray-900'
                ]"
              >
                上一页
              </button>
              
              <template v-for="page in Math.min(pagination.pages, 5)" :key="page">
                <button 
                  @click="handlePageChange(page)"
                  :class="[
                    'px-3 py-1 text-sm rounded',
                    pagination.page === page 
                      ? 'bg-indigo-600 text-white' 
                      : 'text-gray-600 hover:text-gray-900 border border-gray-300'
                  ]"
                >
                  {{ page }}
                </button>
              </template>
              
              <button 
                @click="handlePageChange(pagination.page + 1)"
                :disabled="pagination.page >= pagination.pages"
                :class="[
                  'px-3 py-1 text-sm border border-gray-300 rounded',
                  pagination.page >= pagination.pages 
                    ? 'text-gray-400 cursor-not-allowed' 
                    : 'text-gray-600 hover:text-gray-900'
                ]"
              >
                下一页
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 已购工作流 -->
      <div v-else-if="activeTab === 'workflows'" class="space-y-6">
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">已购买的工作流</h3>
          
          <!-- 加载状态 -->
          <div v-if="loading.workflows" class="text-center py-8">
            <i class="fas fa-spinner fa-spin text-gray-400 text-xl mr-2"></i>
            <span class="text-gray-500">正在加载工作流...</span>
          </div>
          
          <!-- 空状态 -->
          <div v-else-if="purchasedWorkflows.length === 0" class="text-center py-8">
            <i class="fas fa-project-diagram text-gray-300 text-3xl mb-4"></i>
            <p class="text-gray-500 mb-4">暂无已购买的工作流</p>
            <router-link 
              to="/agents" 
              class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
            >
              <i class="fas fa-store mr-2"></i>
              浏览工作流市场
            </router-link>
          </div>
          
          <!-- 工作流列表 -->
          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div v-for="workflow in purchasedWorkflows" :key="workflow.id" 
                 class="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 cursor-pointer transition-colors">
              <h4 class="font-medium text-gray-900">{{ workflow.name }}</h4>
              <p class="text-sm text-gray-500 mt-1">{{ workflow.description }}</p>
              <div class="flex items-center justify-between mt-3">
                <span class="text-sm text-gray-500">购买于: {{ workflow.purchaseDate }}</span>
                <button 
                  @click="executeWorkflow(workflow.id)"
                  class="text-indigo-600 hover:text-indigo-500 text-sm font-medium"
                >
                  使用
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 账户安全 -->
      <div v-else-if="activeTab === 'security'" class="space-y-6">
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">账户安全设置</h3>
          <div class="space-y-4">
            <div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
              <div>
                <h4 class="font-medium text-gray-900">登录密码</h4>
                <p class="text-sm text-gray-500">最后修改: 2023-06-01</p>
              </div>
              <button class="text-indigo-600 hover:text-indigo-500 text-sm">修改</button>
            </div>
            <div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
              <div>
                <h4 class="font-medium text-gray-900">两步验证</h4>
                <p class="text-sm text-gray-500">增强账户安全性</p>
              </div>
              <button class="text-indigo-600 hover:text-indigo-500 text-sm">启用</button>
            </div>
          </div>
        </div>
      </div>

      <!-- 通知设置 -->
      <div v-else-if="activeTab === 'notifications'" class="space-y-6">
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">通知设置</h3>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="font-medium text-gray-900">邮件通知</h4>
                <p class="text-sm text-gray-500">接收工作流执行结果通知</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" class="sr-only peer" checked>
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
              </label>
            </div>
            <div class="flex items-center justify-between">
              <div>
                <h4 class="font-medium text-gray-900">系统通知</h4>
                <p class="text-sm text-gray-500">接收系统更新和维护通知</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { authApi, pointsApi, workflowApi } from '../../api/index.js'

// 当前活跃的标签
const activeTab = ref('points')

// 标签页配置
const tabs = [
  { id: 'points', name: '资源点数管理', icon: 'fas fa-coins' },
  { id: 'workflows', name: '已购工作流', icon: 'fas fa-project-diagram' },
  { id: 'security', name: '账户安全', icon: 'fas fa-shield-alt' },
  { id: 'notifications', name: '通知设置', icon: 'fas fa-bell' }
]

// 加载状态
const loading = reactive({
  profile: false,
  points: false,
  statistics: false,
  transactions: false,
  workflows: false
})

// 用户资料
const userProfile = reactive({
  points: 0
})

// 积分数据
const pointsData = reactive({
  currentBalance: 0,
  monthlyConsumption: 0,
  totalRecharge: 0,
  rechargeCount: 0,
  totalSpent: 0,
  purchasedWorkflowCount: 0,
  recentTransactionCount: 0
})

// 积分历史记录
const pointsHistory = ref([])

// 分页信息
const pagination = reactive({
  total: 0,
  page: 1,
  pageSize: 10,
  pages: 0
})

// 已购买的工作流
const purchasedWorkflows = ref([])

// 获取用户资料
const fetchUserProfile = async () => {
  try {
    loading.profile = true
    const response = await authApi.getProfile()
    if (response.status === 'success') {
      Object.assign(userProfile, response.data)
    }
  } catch (error) {
    console.error('获取用户资料失败:', error)
  } finally {
    loading.profile = false
  }
}

// 获取积分统计数据
const fetchPointsStatistics = async () => {
  try {
    loading.statistics = true
    const response = await pointsApi.getPointsStatistics()
    if (response.status === 'success') {
      const data = response.data
      pointsData.currentBalance = data.currentBalance || 0
      pointsData.totalRecharge = data.totalRecharge || 0
      pointsData.totalSpent = data.totalSpent || 0
      pointsData.purchasedWorkflowCount = data.purchasedWorkflowCount || 0
      pointsData.recentTransactionCount = data.recentTransactionCount || 0
      
      // 计算本月消耗（这里简化为最近30天的消耗）
      pointsData.monthlyConsumption = Math.min(data.totalSpent || 0, 1000)
      
      // 估算充值次数（简化计算）
      pointsData.rechargeCount = Math.ceil((data.totalRecharge || 0) / 1000)
    }
  } catch (error) {
    console.error('获取积分统计失败:', error)
  } finally {
    loading.statistics = false
  }
}

// 获取积分交易历史
const fetchPointsTransactions = async (page = 1) => {
  try {
    loading.transactions = true
    const response = await pointsApi.getPointsHistory({
      limit: pagination.pageSize,
      offset: (page - 1) * pagination.pageSize
    })
    
    if (response.status === 'success') {
      const data = response.data
      pointsHistory.value = (data.transactions || []).map(transaction => ({
        id: transaction.id,
        date: formatDate(transaction.createdAt),
        description: getTransactionDescription(transaction.type, transaction.description),
        workflow: getWorkflowName(transaction),
        change: transaction.amount,
        balance: transaction.balance
      }))
      
      // 更新分页信息
      if (data.pagination) {
        pagination.total = data.pagination.total
        pagination.page = data.pagination.page
        pagination.pages = data.pagination.pages
      }
    }
  } catch (error) {
    console.error('获取积分交易历史失败:', error)
  } finally {
    loading.transactions = false
  }
}

// 获取已购买的工作流
const fetchPurchasedWorkflows = async () => {
  try {
    loading.workflows = true
    const response = await workflowApi.getWorkflows({ 
      limit: 20,
      // 只获取从代理购买的工作流
      purchased: true
    })
    
    if (response.status === 'success') {
      const workflows = response.data.workflows || []
      purchasedWorkflows.value = workflows
        .filter(workflow => workflow.agentId) // 只显示从代理购买的工作流
        .map(workflow => ({
          id: workflow.id,
          name: workflow.name,
          description: workflow.description || '暂无描述',
          purchaseDate: formatDate(workflow.createdAt, 'YYYY-MM-DD')
        }))
    }
  } catch (error) {
    console.error('获取已购买工作流失败:', error)
  } finally {
    loading.workflows = false
  }
}

// 工具函数
const formatDate = (dateString, format = 'YYYY-MM-DD HH:mm') => {
  if (!dateString) return '未知'
  const date = new Date(dateString)
  
  if (format === 'YYYY-MM-DD') {
    return date.toLocaleDateString('zh-CN')
  }
  
  return date.toLocaleString('zh-CN')
}

const getTransactionDescription = (type, description) => {
  const typeMap = {
    'recharge': '充值',
    'purchase': '购买工作流',
    'execution': '执行工作流',
    'refund': '退款'
  }
  
  return typeMap[type] || description || '未知交易'
}

const getWorkflowName = (transaction) => {
  // 根据交易类型和描述推断工作流名称
  if (transaction.type === 'execution' || transaction.type === 'purchase') {
    return transaction.description || '工作流'
  }
  return '-'
}

// 分页处理
const handlePageChange = (page) => {
  fetchPointsTransactions(page)
}

// 执行工作流
const executeWorkflow = async (workflowId) => {
  try {
    const response = await workflowApi.executeWorkflow(workflowId, {})
    if (response.status === 'success') {
      // 可以添加成功提示或跳转到执行详情页
      console.log('工作流执行成功:', response)
    }
  } catch (error) {
    console.error('执行工作流失败:', error)
  }
}

// 页面加载时获取数据
onMounted(async () => {
  await Promise.all([
    fetchUserProfile(),
    fetchPointsStatistics(),
    fetchPointsTransactions(),
    fetchPurchasedWorkflows()
  ])
})
</script>

<style scoped>
/* 组件特定样式 */
.transition-colors {
  transition-property: color, background-color, border-color, text-decoration-color, fill, stroke;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
}
</style> 