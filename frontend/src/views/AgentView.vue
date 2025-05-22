<template>
  <div class="container mx-auto px-4 py-6">
    <!-- 页面标题和操作按钮 -->
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">代理管理</h1>
      <button 
        @click="openCreateModal" 
        class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg flex items-center"
      >
        <i class="fas fa-plus-circle mr-2"></i> 创建代理
      </button>
    </div>

    <!-- 筛选和搜索区域 -->
    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- 分类筛选 -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">分类</label>
          <select 
            v-model="filters.category" 
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
          >
            <option value="">全部分类</option>
            <option v-for="category in categories" :key="category" :value="category">
              {{ category }}
            </option>
          </select>
        </div>
        
        <!-- 公开/私有筛选 -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">可见性</label>
          <select 
            v-model="filters.isPublic" 
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
          >
            <option :value="null">全部</option>
            <option :value="true">公开</option>
            <option :value="false">私有</option>
          </select>
        </div>
        
        <!-- 搜索框 -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">搜索</label>
          <div class="relative">
            <input 
              v-model="filters.query" 
              type="text" 
              placeholder="搜索代理名称..." 
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 pl-10"
            >
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <i class="fas fa-search text-gray-400"></i>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="spinner-border animate-spin inline-block w-8 h-8 border-4 rounded-full text-blue-600" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <span class="ml-2 text-gray-600">加载中...</span>
    </div>
    
    <!-- 空状态 -->
    <div v-else-if="agents.length === 0" class="bg-white rounded-lg shadow p-12 text-center">
      <div class="text-5xl text-gray-300 mb-4">
        <i class="fas fa-robot"></i>
      </div>
      <h3 class="text-xl font-medium text-gray-600 mb-2">{{ hasFilters ? '没有符合条件的代理' : '暂无代理' }}</h3>
      <p class="text-gray-500 mb-6">
        {{ hasFilters 
          ? '尝试调整筛选条件以查看更多代理' 
          : '您还没有创建任何代理，点击下方按钮创建第一个代理' }}
      </p>
      <div class="flex justify-center space-x-4">
        <button 
          v-if="hasFilters"
          @click="clearFilters" 
          class="bg-gray-100 hover:bg-gray-200 text-gray-800 font-medium py-2 px-4 rounded-lg inline-flex items-center"
        >
          <i class="fas fa-times-circle mr-2"></i> 清除筛选
        </button>
        <button 
          @click="openCreateModal" 
          class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg inline-flex items-center"
        >
          <i class="fas fa-plus-circle mr-2"></i> {{ hasFilters ? '创建新代理' : '创建第一个代理' }}
        </button>
      </div>
    </div>
    
    <!-- 代理列表 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
      <!-- 代理卡片 -->
      <div 
        v-for="agent in agents" 
        :key="agent.id" 
        class="bg-white rounded-lg shadow-md overflow-hidden border border-gray-100 hover:shadow-lg transition-shadow duration-300"
      >
        <div class="p-5">
          <!-- 代理标题和图标 -->
          <div class="flex justify-between items-start mb-4">
            <div class="flex items-center">
              <div class="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center text-blue-500 mr-3">
                <i :class="agent.icon || 'fas fa-robot'" class="text-lg"></i>
              </div>
              <div>
                <h3 class="text-lg font-semibold text-gray-800">{{ agent.name }}</h3>
                <p class="text-sm text-gray-500">{{ agent.category }}</p>
              </div>
            </div>
            <div class="flex items-center">
              <span v-if="agent.isPublic" class="px-2 py-1 text-xs rounded bg-green-100 text-green-800">公开</span>
              <span v-else class="px-2 py-1 text-xs rounded bg-gray-100 text-gray-800">私有</span>
            </div>
          </div>
          
          <!-- 代理描述 -->
          <p class="text-gray-600 text-sm mb-4 line-clamp-2">{{ agent.description || '暂无描述' }}</p>
          
          <!-- 代理信息 -->
          <div class="flex items-center text-sm text-gray-500 mb-4">
            <div class="mr-4">
              <i class="fas fa-calendar-alt mr-1"></i>
              {{ formatDate(agent.createdAt) }}
            </div>
            <div>
              <i class="fas fa-play-circle mr-1"></i>
              使用次数: {{ agent.usageCount || 0 }}
            </div>
          </div>
          
          <!-- 操作按钮 -->
          <div class="flex justify-between items-center pt-3 border-t border-gray-100">
            <button 
              @click="viewAgentDetails(agent)" 
              class="text-blue-600 hover:text-blue-800 text-sm font-medium"
            >
              详情
            </button>
            <div v-if="isOwnAgent(agent)">
              <button 
                @click="editAgent(agent)" 
                class="text-gray-600 hover:text-gray-800 mr-3"
              >
                <i class="fas fa-edit"></i>
              </button>
              <button 
                @click="confirmDelete(agent)" 
                class="text-red-600 hover:text-red-800"
              >
                <i class="fas fa-trash-alt"></i>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 分页控制 -->
    <div v-if="!loading && agents.length > 0" class="flex justify-between items-center">
      <div class="text-sm text-gray-500">
        显示 {{ pagination.offset + 1 }} 到 {{ Math.min(pagination.offset + agents.length, pagination.total) }} 条，共 {{ pagination.total }} 条
      </div>
      <div class="flex space-x-2">
        <button 
          @click="prevPage" 
          :disabled="pagination.offset === 0" 
          :class="{'opacity-50 cursor-not-allowed': pagination.offset === 0}"
          class="px-3 py-1 rounded border border-gray-300 bg-white text-gray-700 hover:bg-gray-50"
        >
          上一页
        </button>
        <button 
          @click="nextPage" 
          :disabled="pagination.offset + pagination.limit >= pagination.total" 
          :class="{'opacity-50 cursor-not-allowed': pagination.offset + pagination.limit >= pagination.total}"
          class="px-3 py-1 rounded border border-gray-300 bg-white text-gray-700 hover:bg-gray-50"
        >
          下一页
        </button>
      </div>
    </div>

    <!-- 通知提示 -->
    <div v-if="notification.show" class="fixed top-4 right-4 max-w-md z-50">
      <div 
        :class="[
          'p-4 rounded-lg shadow-lg flex items-start',
          notification.type === 'success' ? 'bg-green-100 text-green-800 border-l-4 border-green-500' : 
          notification.type === 'error' ? 'bg-red-100 text-red-800 border-l-4 border-red-500' : 
          'bg-blue-100 text-blue-800 border-l-4 border-blue-500'
        ]"
      >
        <div class="mr-3 mt-0.5">
          <i 
            :class="[
              notification.type === 'success' ? 'fas fa-check-circle text-green-500' : 
              notification.type === 'error' ? 'fas fa-exclamation-circle text-red-500' : 
              'fas fa-info-circle text-blue-500'
            ]"
          ></i>
        </div>
        <div class="flex-1">
          <p class="font-medium">{{ notification.title }}</p>
          <p class="text-sm">{{ notification.message }}</p>
        </div>
        <button @click="notification.show = false" class="ml-4 text-gray-500 hover:text-gray-700">
          <i class="fas fa-times"></i>
        </button>
      </div>
    </div>

    <!-- 模态框 -->
    <!-- 代理详情模态框 -->
    <div v-if="detailsModal.show" class="fixed inset-0 bg-gray-800 bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-4xl max-h-screen overflow-hidden flex flex-col">
        <!-- 模态框头部 -->
        <div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">代理详情</h3>
          <button @click="detailsModal.show = false" class="text-gray-400 hover:text-gray-500">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <!-- 模态框内容 -->
        <div class="p-6 overflow-y-auto">
          <div v-if="detailsModal.loading" class="flex justify-center py-8">
            <div class="spinner-border animate-spin inline-block w-8 h-8 border-4 rounded-full text-blue-600" role="status">
              <span class="visually-hidden">Loading...</span>
            </div>
          </div>
          
          <div v-else-if="detailsModal.agent">
            <!-- 代理头部信息 -->
            <div class="flex flex-col md:flex-row items-start md:items-center mb-6 pb-4 border-b border-gray-200">
              <div class="w-16 h-16 rounded-full bg-blue-100 flex items-center justify-center text-blue-500 mr-4 mb-4 md:mb-0">
                <i :class="detailsModal.agent.icon || 'fas fa-robot'" class="text-2xl"></i>
              </div>
              <div class="flex-1">
                <div class="flex flex-col md:flex-row md:items-center mb-2">
                  <h2 class="text-2xl font-bold text-gray-800 mr-3">{{ detailsModal.agent.name }}</h2>
                  <div class="flex items-center mt-2 md:mt-0">
                    <span v-if="detailsModal.agent.isPublic" class="px-2 py-1 text-xs rounded bg-green-100 text-green-800 mr-2">公开</span>
                    <span v-else class="px-2 py-1 text-xs rounded bg-gray-100 text-gray-800 mr-2">私有</span>
                    <span class="px-2 py-1 text-xs rounded bg-blue-100 text-blue-800">{{ detailsModal.agent.category }}</span>
                  </div>
                </div>
                <p class="text-gray-600">{{ detailsModal.agent.description || '暂无描述' }}</p>
              </div>
            </div>
            
            <!-- 代理基本信息 -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
              <div class="bg-gray-50 p-3 rounded">
                <div class="text-sm text-gray-500 mb-1">创建时间</div>
                <div class="font-medium">{{ formatDate(detailsModal.agent.createdAt) }}</div>
              </div>
              <div class="bg-gray-50 p-3 rounded">
                <div class="text-sm text-gray-500 mb-1">最后更新</div>
                <div class="font-medium">{{ formatDate(detailsModal.agent.updatedAt) }}</div>
              </div>
              <div class="bg-gray-50 p-3 rounded">
                <div class="text-sm text-gray-500 mb-1">使用次数</div>
                <div class="font-medium">{{ detailsModal.agent.usageCount || 0 }}次</div>
              </div>
              <div class="bg-gray-50 p-3 rounded">
                <div class="text-sm text-gray-500 mb-1">代理类型</div>
                <div class="font-medium">{{ detailsModal.agent.type || '未指定' }}</div>
              </div>
            </div>
            
            <!-- 代理定义 -->
            <div class="mb-6">
              <h3 class="text-lg font-medium text-gray-900 mb-3">代理定义</h3>
              <pre class="bg-gray-800 text-gray-100 p-4 rounded-lg overflow-x-auto whitespace-pre-wrap break-words">{{ formatJson(detailsModal.agent.definition) }}</pre>
            </div>
          </div>
          
          <div v-else class="py-8 text-center text-gray-500">
            未能加载代理详情
          </div>
        </div>
        
        <!-- 模态框底部 -->
        <div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end">
          <button 
            @click="detailsModal.show = false" 
            class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-800 rounded-lg mr-2"
          >
            关闭
          </button>
          <button 
            v-if="detailsModal.agent && isOwnAgent(detailsModal.agent)"
            @click="editAgentFromDetails"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg"
          >
            编辑
          </button>
        </div>
      </div>
    </div>
    
    <!-- 创建/编辑代理模态框 -->
    <div v-if="formModal.show" class="fixed inset-0 bg-gray-800 bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-3xl max-h-screen overflow-hidden flex flex-col">
        <!-- 模态框头部 -->
        <div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">{{ formModal.isEdit ? '编辑代理' : '创建代理' }}</h3>
          <button @click="formModal.show = false" class="text-gray-400 hover:text-gray-500">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <!-- 模态框内容 -->
        <div class="p-6 overflow-y-auto">
          <form @submit.prevent="submitAgentForm">
            <!-- 基本信息 -->
            <div class="mb-6">
              <h4 class="text-base font-medium text-gray-700 mb-4">基本信息</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- 名称 -->
                <div class="col-span-2">
                  <label class="block text-sm font-medium text-gray-700 mb-1">名称 <span class="text-red-500">*</span></label>
                  <input 
                    v-model="formModal.form.name" 
                    type="text" 
                    required
                    class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                    placeholder="输入代理名称"
                  >
                </div>
                
                <!-- 类型 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">类型 <span class="text-red-500">*</span></label>
                  <select 
                    v-model="formModal.form.type" 
                    required
                    class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                  >
                    <option value="">选择类型</option>
                    <option value="chat">对话型</option>
                    <option value="function">函数型</option>
                    <option value="workflow">工作流型</option>
                  </select>
                </div>
                
                <!-- 分类 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">分类 <span class="text-red-500">*</span></label>
                  <select 
                    v-model="formModal.form.category" 
                    required
                    class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                  >
                    <option value="">选择分类</option>
                    <option v-for="category in categories" :key="category" :value="category">
                      {{ category }}
                    </option>
                    <option value="other">其他</option>
                  </select>
                </div>
                
                <!-- 图标 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">图标</label>
                  <input 
                    v-model="formModal.form.icon" 
                    type="text" 
                    class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                    placeholder="输入FontAwesome图标类名"
                  >
                  <div class="text-xs text-gray-500 mt-1">例如: fas fa-robot</div>
                </div>
                
                <!-- 可见性 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">可见性</label>
                  <div class="mt-2">
                    <label class="inline-flex items-center">
                      <input type="checkbox" v-model="formModal.form.isPublic" class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500">
                      <span class="ml-2 text-gray-700">公开（所有人可见）</span>
                    </label>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 描述 -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
              <textarea 
                v-model="formModal.form.description" 
                rows="3"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                placeholder="输入代理描述"
              ></textarea>
            </div>
            
            <!-- 代理定义 -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-1">代理定义 <span class="text-red-500">*</span></label>
              <textarea 
                v-model="formModal.form.definitionText" 
                rows="8"
                required
                class="block w-full font-mono rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                placeholder="输入JSON格式的代理定义"
              ></textarea>
              <div class="flex justify-between text-xs text-gray-500 mt-1">
                <span>请输入有效的JSON格式</span>
                <button 
                  type="button" 
                  @click="formatDefinition" 
                  class="text-blue-600 hover:text-blue-800"
                >
                  格式化JSON
                </button>
              </div>
              <div v-if="formModal.definitionError" class="text-red-500 text-sm mt-1">
                {{ formModal.definitionError }}
              </div>
            </div>
          </form>
        </div>
        
        <!-- 模态框底部 -->
        <div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end">
          <button 
            @click="formModal.show = false" 
            class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-800 rounded-lg mr-2"
            :disabled="formModal.loading"
          >
            取消
          </button>
          <button 
            @click="submitAgentForm"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg flex items-center"
            :disabled="formModal.loading"
          >
            <span v-if="formModal.loading" class="spinner-border spinner-border-sm mr-2" role="status" aria-hidden="true"></span>
            {{ formModal.isEdit ? '保存' : '创建' }}
          </button>
        </div>
      </div>
    </div>
    
    <!-- 删除确认模态框 -->
    <div v-if="deleteModal.show" class="fixed inset-0 bg-gray-800 bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden">
        <!-- 模态框头部 -->
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">确认删除</h3>
        </div>
        
        <!-- 模态框内容 -->
        <div class="p-6">
          <p class="text-gray-700">您确定要删除代理 <strong>{{ deleteModal.agent?.name }}</strong> 吗？此操作不可撤销。</p>
        </div>
        
        <!-- 模态框底部 -->
        <div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end">
          <button 
            @click="deleteModal.show = false" 
            class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-800 rounded-lg mr-2"
            :disabled="deleteModal.loading"
          >
            取消
          </button>
          <button 
            @click="confirmDeleteAgent"
            class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg flex items-center"
            :disabled="deleteModal.loading"
          >
            <span v-if="deleteModal.loading" class="spinner-border spinner-border-sm mr-2" role="status" aria-hidden="true"></span>
            删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { agentApi } from '../api/index'

export default {
  name: 'AgentView',
  setup() {
    // 状态
    const agents = ref([])
    const categories = ref([])
    const loading = ref(false)
    const filters = reactive({
      category: '',
      isPublic: null,
      query: ''
    })
    const pagination = reactive({
      limit: 9,
      offset: 0,
      total: 0
    })

    // 详情模态框状态
    const detailsModal = reactive({
      show: false,
      loading: false,
      agent: null
    })

    // 创建/编辑模态框状态
    const formModal = reactive({
      show: false,
      isEdit: false,
      loading: false,
      definitionError: '',
      form: {
        name: '',
        description: '',
        type: '',
        category: '',
        icon: '',
        isPublic: false,
        definitionText: '{}'
      },
      originalAgent: null
    })

    // 删除确认模态框状态
    const deleteModal = reactive({
      show: false,
      loading: false,
      agent: null
    })

    // 通知提示状态
    const notification = reactive({
      show: false,
      type: 'info', // 'success', 'error', 'info'
      title: '',
      message: '',
      timeout: null
    })

    // 计算属性
    const isOwnAgent = (agent) => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      return agent.userId === user.id
    }

    // 计算是否有筛选条件
    const hasFilters = computed(() => {
      return filters.category !== '' || 
             filters.isPublic !== null || 
             (filters.query && filters.query.trim() !== '')
    })

    // 清除所有筛选条件
    const clearFilters = () => {
      filters.category = ''
      filters.isPublic = null
      filters.query = ''
      // loadAgents() 由于我们使用了 watch，此处不需要手动调用
    }

    // 方法
    const loadAgents = async () => {
      loading.value = true
      try {
        const params = {
          limit: pagination.limit,
          offset: pagination.offset
        }
        
        if (filters.category) {
          params.category = filters.category
        }
        
        if (filters.isPublic !== null) {
          params.is_public = filters.isPublic
        }
        
        // 添加搜索查询参数，如果后端支持
        if (filters.query && filters.query.trim() !== '') {
          params.query = filters.query.trim()
        }
        
        const response = await agentApi.getAgents(params)
        agents.value = response.data || []
        pagination.total = response.total || 0
      } catch (error) {
        console.error('加载代理列表失败:', error)
        showNotification('error', '加载失败', '无法加载代理列表，请稍后重试')
      } finally {
        loading.value = false
      }
    }

    const loadCategories = async () => {
      try {
        const response = await agentApi.getAgentCategories()
        categories.value = response || []
      } catch (error) {
        console.error('加载代理分类失败:', error)
      }
    }

    const formatDate = (dateString) => {
      if (!dateString) return '未知时间'
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      })
    }

    const prevPage = () => {
      if (pagination.offset > 0) {
        pagination.offset = Math.max(0, pagination.offset - pagination.limit)
        loadAgents()
      }
    }

    const nextPage = () => {
      if (pagination.offset + pagination.limit < pagination.total) {
        pagination.offset += pagination.limit
        loadAgents()
      }
    }

    // 重置表单
    const resetForm = () => {
      formModal.form = {
        name: '',
        description: '',
        type: '',
        category: '',
        icon: '',
        isPublic: false,
        definitionText: '{}'
      }
      formModal.definitionError = ''
      formModal.originalAgent = null
    }

    // 打开创建模态框
    const openCreateModal = () => {
      resetForm()
      formModal.isEdit = false
      formModal.show = true
    }

    // 编辑代理
    const editAgent = (agent) => {
      resetForm()
      formModal.isEdit = true
      formModal.originalAgent = agent
      formModal.form = {
        name: agent.name,
        description: agent.description || '',
        type: agent.type,
        category: agent.category,
        icon: agent.icon || '',
        isPublic: agent.isPublic,
        definitionText: formatJson(agent.definition)
      }
      formModal.show = true
    }

    // 确认删除代理
    const confirmDelete = (agent) => {
      deleteModal.agent = agent
      deleteModal.show = true
      deleteModal.loading = false
    }

    // 执行删除代理
    const confirmDeleteAgent = async () => {
      if (!deleteModal.agent) return
      
      deleteModal.loading = true
      try {
        await agentApi.deleteAgent(deleteModal.agent.id)
        showNotification('success', '删除成功', `代理 "${deleteModal.agent.name}" 已删除`)
        loadAgents() // 重新加载代理列表
        deleteModal.show = false
      } catch (error) {
        console.error('删除代理失败:', error)
        showNotification('error', '删除失败', error.message || '无法删除代理，请稍后重试')
      } finally {
        deleteModal.loading = false
      }
    }

    // 格式化代理定义JSON
    const formatDefinition = () => {
      try {
        const parsed = JSON.parse(formModal.form.definitionText)
        formModal.form.definitionText = JSON.stringify(parsed, null, 2)
        formModal.definitionError = ''
        showNotification('success', 'JSON格式化', 'JSON已成功格式化', 2000)
      } catch (error) {
        formModal.definitionError = '无效的JSON格式: ' + error.message
        showNotification('error', 'JSON格式化失败', error.message, 3000)
      }
    }

    // 提交代理表单
    const submitAgentForm = async () => {
      // 验证JSON格式
      try {
        JSON.parse(formModal.form.definitionText)
        formModal.definitionError = ''
      } catch (error) {
        formModal.definitionError = '无效的JSON格式: ' + error.message
        return
      }
      
      formModal.loading = true
      
      try {
        const formData = {
          name: formModal.form.name,
          description: formModal.form.description,
          type: formModal.form.type,
          category: formModal.form.category,
          icon: formModal.form.icon,
          isPublic: formModal.form.isPublic,
          definition: JSON.parse(formModal.form.definitionText)
        }
        
        if (formModal.isEdit && formModal.originalAgent) {
          await agentApi.updateAgent(formModal.originalAgent.id, formData)
          showNotification('success', '更新成功', `代理 "${formData.name}" 已更新`)
        } else {
          await agentApi.createAgent(formData)
          showNotification('success', '创建成功', `代理 "${formData.name}" 已创建`)
        }
        
        loadAgents() // 重新加载代理列表
        formModal.show = false
      } catch (error) {
        console.error(formModal.isEdit ? '更新代理失败:' : '创建代理失败:', error)
        showNotification(
          'error', 
          formModal.isEdit ? '更新失败' : '创建失败', 
          error.message || '操作失败，请稍后重试'
        )
      } finally {
        formModal.loading = false
      }
    }

    // 从详情页编辑代理
    const editAgentFromDetails = () => {
      if (detailsModal.agent) {
        detailsModal.show = false
        editAgent(detailsModal.agent)
      }
    }

    // 查看代理详情
    const viewAgentDetails = async (agent) => {
      detailsModal.show = true
      detailsModal.loading = true
      detailsModal.agent = null
      
      try {
        const response = await agentApi.getAgent(agent.id)
        detailsModal.agent = response
      } catch (error) {
        console.error('获取代理详情失败:', error)
        showNotification('error', '加载失败', '无法加载代理详情，请稍后重试')
      } finally {
        detailsModal.loading = false
      }
    }

    // 格式化JSON显示
    const formatJson = (json) => {
      if (!json) return '{}'
      if (typeof json === 'string') {
        try {
          return JSON.stringify(JSON.parse(json), null, 2)
        } catch (e) {
          return json
        }
      }
      return JSON.stringify(json, null, 2)
    }

    // 监听筛选条件变化
    const setupFilters = () => {
      watch(() => [filters.category, filters.isPublic], () => {
        pagination.offset = 0 // 重置分页
        loadAgents()
      })
      
      // 添加搜索输入防抖处理
      let searchTimeout = null
      watch(() => filters.query, (newValue) => {
        if (searchTimeout) {
          clearTimeout(searchTimeout)
        }
        
        searchTimeout = setTimeout(() => {
          pagination.offset = 0 // 重置分页
          loadAgents()
        }, 500) // 500ms 延迟，用户停止输入后才触发搜索
      })
    }

    // 显示通知
    const showNotification = (type, title, message, duration = 5000) => {
      // 清除之前的定时器
      if (notification.timeout) {
        clearTimeout(notification.timeout)
      }
      
      // 设置新通知
      notification.show = true
      notification.type = type
      notification.title = title
      notification.message = message
      
      // 设置自动关闭
      notification.timeout = setTimeout(() => {
        notification.show = false
      }, duration)
    }

    // 生命周期钩子
    onMounted(() => {
      loadAgents()
      loadCategories()
      setupFilters()
    })

    return {
      agents,
      categories,
      loading,
      filters,
      pagination,
      isOwnAgent,
      loadAgents,
      formatDate,
      prevPage,
      nextPage,
      openCreateModal,
      viewAgentDetails,
      editAgent,
      confirmDelete,
      detailsModal,
      formatJson,
      editAgentFromDetails,
      formModal,
      deleteModal,
      confirmDeleteAgent,
      formatDefinition,
      submitAgentForm,
      notification,
      showNotification,
      hasFilters,
      clearFilters
    }
  }
}
</script>

<style scoped>
/* 添加必要的自定义样式 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.spinner-border {
  border-right-color: transparent;
}

.visually-hidden {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
  border-width: 0.2em;
}
</style> 