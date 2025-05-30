<template>
  <div class="flex-1 overflow-y-auto">
    <!-- 顶部横幅 -->
    <div class="bg-gradient-to-r from-purple-600 to-indigo-600 text-white p-8">
      <div class="max-w-7xl mx-auto">
        <h2 class="text-3xl font-bold mb-2">探索工作流</h2>
        <p class="text-lg mb-6">发现预构建的工作流，增强您的工作流程，提高效率和智能化水平</p>
        <div class="flex flex-wrap gap-3">
          <button 
            @click="setActiveCategory('')"
            :class="[
              'px-4 py-2 rounded-lg font-medium transition-colors',
              activeCategory === '' 
                ? 'bg-white text-indigo-700' 
                : 'bg-white/20 text-white hover:bg-white/30'
            ]"
          >
            所有工作流
          </button>
          <button 
            v-for="category in categories" 
            :key="category"
            @click="setActiveCategory(category)"
            :class="[
              'px-4 py-2 rounded-lg font-medium transition-colors',
              activeCategory === category 
                ? 'bg-white text-indigo-700' 
                : 'bg-white/20 text-white hover:bg-white/30'
            ]"
          >
            {{ getCategoryDisplayName(category) }}
          </button>
        </div>
      </div>
    </div>
    
    <!-- 搜索和筛选区域 -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="bg-white rounded-lg shadow-sm p-4 mb-6">
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="flex-1">
            <div class="relative">
              <input 
                v-model="searchQuery"
                type="text" 
                placeholder="搜索工作流..." 
                class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                @input="handleSearch"
              >
              <MagnifyingGlassIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            </div>
          </div>
          <div class="flex gap-2">
            <select 
              v-model="sortOrder"
              @change="handleSort"
              class="px-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
              <option value="popular">最受欢迎</option>
              <option value="newest">最新发布</option>
              <option value="rating">评分最高</option>
              <option value="price_low">价格低到高</option>
              <option value="price_high">价格高到低</option>
            </select>
            <button 
              @click="toggleView"
              class="px-4 py-2 rounded-lg border border-gray-300 hover:bg-gray-50 transition-colors"
              :title="viewMode === 'grid' ? '切换到列表视图' : '切换到网格视图'"
            >
              <Bars3Icon v-if="viewMode === 'grid'" class="w-5 h-5" />
              <Squares2X2Icon v-else class="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 工作流列表/网格 -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-8">
      <!-- 加载状态 -->
      <div v-if="loading" class="text-center py-12">
        <ArrowPathIcon class="w-8 h-8 mx-auto mb-4 text-indigo-600 animate-spin" />
        <p class="text-gray-600">加载工作流...</p>
      </div>
      
      <!-- 错误状态 -->
      <div v-else-if="error" class="text-center py-12">
        <ExclamationTriangleIcon class="w-16 h-16 mx-auto mb-4 text-red-500" />
        <h3 class="text-lg font-medium text-gray-900 mb-2">加载失败</h3>
        <p class="text-gray-600 mb-6">{{ error }}</p>
        <button
          @click="fetchAgents"
          class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
        >
          <ArrowPathIcon class="w-5 h-5 mr-2" />
          重试
        </button>
      </div>
      
      <!-- 空状态 -->
      <div v-else-if="agents.length === 0" class="text-center py-12">
        <RectangleStackIcon class="w-16 h-16 mx-auto mb-4 text-gray-400" />
        <h3 class="text-lg font-medium text-gray-900 mb-2">暂无工作流</h3>
        <p class="text-gray-600">{{ searchQuery ? '没有找到匹配的工作流' : '当前分类下暂无工作流' }}</p>
      </div>
      
      <!-- 网格视图 -->
      <div v-else-if="viewMode === 'grid'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="agent in agents" 
          :key="agent.id" 
          class="bg-white rounded-xl shadow-sm hover:shadow-lg transition-all duration-300 hover:-translate-y-1 overflow-hidden border border-gray-100"
        >
          <!-- 工作流图片/图标 -->
          <div class="h-48 bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center relative overflow-hidden">
            <!-- 自定义封面图 -->
            <img 
              v-if="agent.coverImage" 
              :src="agent.coverImage" 
              :alt="agent.name"
              class="w-full h-full object-cover"
              @error="handleImageError"
            >
            <!-- 默认图标背景 -->
            <div v-else class="w-16 h-16 rounded-full bg-white/20 backdrop-blur-sm flex items-center justify-center">
              <component :is="getWorkflowIcon(agent.category)" class="w-8 h-8 text-white" />
            </div>
          </div>
          
          <!-- 工作流信息 -->
          <div class="p-6">
            <div class="flex items-center justify-between mb-4">
              <h4 class="text-lg font-semibold text-gray-800 line-clamp-1">{{ agent.name }}</h4>
              <span :class="getPriceClass(agent.price)">
                {{ agent.price === 0 ? '免费' : `${agent.price} 积分` }}
              </span>
            </div>
            
            <p class="text-gray-600 text-sm mb-4 line-clamp-2">{{ agent.description || '暂无描述' }}</p>
            
            <div class="flex flex-wrap gap-2 mb-4">
              <span class="category-pill" :class="getCategoryPillClass(agent.category)">
                {{ getCategoryDisplayName(agent.category) }}
              </span>
              <span class="category-pill bg-gray-100 text-gray-800">
                {{ agent.type }}
              </span>
            </div>
            
            <div class="flex items-center justify-between text-sm text-gray-500 mb-4">
              <div class="flex items-center">
                <UserGroupIcon class="w-4 h-4 mr-1" />
                <span>{{ agent.purchaseCount || 0 }} 人使用</span>
              </div>
              <div class="flex items-center">
                <StarIcon class="w-4 h-4 mr-1 text-yellow-500" />
                <span>{{ (agent.rating || 0).toFixed(1) }}</span>
              </div>
            </div>
            
            <div class="flex justify-between items-center">
              <button 
                @click="viewAgentDetails(agent)"
                class="text-indigo-600 hover:text-indigo-800 text-sm font-medium transition-colors"
              >
                查看详情
              </button>
              <button 
                @click="purchaseAgent(agent)"
                :disabled="purchasing.includes(agent.id)"
                class="bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-400 text-white text-sm py-2 px-4 rounded-lg transition-colors flex items-center"
              >
                <ArrowPathIcon v-if="purchasing.includes(agent.id)" class="w-4 h-4 mr-1 animate-spin" />
                <ShoppingCartIcon v-else class="w-4 h-4 mr-1" />
                {{ agent.price === 0 ? '免费使用' : '立即购买' }}
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 列表视图 -->
      <div v-else class="bg-white rounded-lg shadow-sm overflow-hidden">
        <div class="divide-y divide-gray-200">
          <div 
            v-for="agent in agents" 
            :key="agent.id" 
            class="p-6 hover:bg-gray-50 transition-colors"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center flex-1">
                <div class="w-12 h-12 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center mr-4 relative overflow-hidden">
                  <!-- 自定义封面图 -->
                  <img 
                    v-if="agent.coverImage" 
                    :src="agent.coverImage" 
                    :alt="agent.name"
                    class="w-full h-full object-cover rounded-lg"
                    @error="handleImageError"
                  >
                  <!-- 默认图标 -->
                  <component v-else :is="getWorkflowIcon(agent.category)" class="w-6 h-6 text-white" />
                </div>
                <div class="flex-1">
                  <div class="flex items-center gap-3 mb-2">
                    <h4 class="text-lg font-semibold text-gray-800">{{ agent.name }}</h4>
                    <span :class="getPriceClass(agent.price)">
                      {{ agent.price === 0 ? '免费' : `${agent.price} 积分` }}
                    </span>
                    <span class="category-pill" :class="getCategoryPillClass(agent.category)">
                      {{ getCategoryDisplayName(agent.category) }}
                    </span>
                  </div>
                  <p class="text-gray-600 text-sm mb-2">{{ agent.description || '暂无描述' }}</p>
                  <div class="flex items-center gap-4 text-sm text-gray-500">
                    <div class="flex items-center">
                      <UserGroupIcon class="w-4 h-4 mr-1" />
                      <span>{{ agent.purchaseCount || 0 }} 人使用</span>
                    </div>
                    <div class="flex items-center">
                      <StarIcon class="w-4 h-4 mr-1 text-yellow-500" />
                      <span>{{ (agent.rating || 0).toFixed(1) }}</span>
                    </div>
                    <div class="flex items-center">
                      <CalendarIcon class="w-4 h-4 mr-1" />
                      <span>{{ formatDate(agent.createdAt) }}</span>
                    </div>
                  </div>
                </div>
              </div>
              <div class="flex items-center gap-3 ml-6">
                <button 
                  @click="viewAgentDetails(agent)"
                  class="text-indigo-600 hover:text-indigo-800 text-sm font-medium transition-colors"
                >
                  查看详情
                </button>
                <button 
                  @click="purchaseAgent(agent)"
                  :disabled="purchasing.includes(agent.id)"
                  class="bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-400 text-white text-sm py-2 px-4 rounded-lg transition-colors flex items-center"
                >
                  <ArrowPathIcon v-if="purchasing.includes(agent.id)" class="w-4 h-4 mr-1 animate-spin" />
                  <ShoppingCartIcon v-else class="w-4 h-4 mr-1" />
                  {{ agent.price === 0 ? '免费使用' : '立即购买' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 分页 -->
      <div v-if="totalPages > 1" class="flex justify-center mt-8">
        <div class="flex items-center space-x-2">
          <button 
            @click="goToPage(currentPage - 1)"
            :disabled="currentPage === 1"
            :class="[
              'px-3 py-2 rounded-lg border border-gray-300 transition-colors',
              currentPage === 1 
                ? 'opacity-50 cursor-not-allowed bg-gray-100' 
                : 'bg-white hover:bg-gray-50'
            ]"
          >
            上一页
          </button>
          
          <button
            v-for="page in getPageNumbers()"
            :key="page"
            @click="goToPage(page)"
            :class="[
              'px-3 py-2 rounded-lg border transition-colors',
              page === currentPage 
                ? 'bg-indigo-600 text-white border-indigo-600' 
                : 'bg-white border-gray-300 hover:bg-gray-50'
            ]"
          >
            {{ page }}
          </button>
          
          <button 
            @click="goToPage(currentPage + 1)"
            :disabled="currentPage === totalPages"
            :class="[
              'px-3 py-2 rounded-lg border border-gray-300 transition-colors',
              currentPage === totalPages 
                ? 'opacity-50 cursor-not-allowed bg-gray-100' 
                : 'bg-white hover:bg-gray-50'
            ]"
          >
            下一页
          </button>
        </div>
      </div>
    </div>
    
    <!-- 工作流详情模态框 -->
    <div v-if="showDetailModal" class="fixed inset-0 bg-gray-800 bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-4xl max-h-screen overflow-hidden flex flex-col">
        <!-- 模态框头部 -->
        <div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">工作流详情</h3>
          <button @click="closeDetailModal" class="text-gray-400 hover:text-gray-500">
            <XMarkIcon class="w-6 h-6" />
          </button>
        </div>
        
        <!-- 模态框内容 -->
        <div class="flex-1 overflow-y-auto p-6" v-if="selectedAgent">
          <div class="flex flex-col md:flex-row gap-6">
            <!-- 左侧：基本信息 -->
            <div class="md:w-1/3">
              <div class="h-48 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center mb-4 relative overflow-hidden">
                <!-- 自定义封面图 -->
                <img 
                  v-if="selectedAgent.coverImage" 
                  :src="selectedAgent.coverImage" 
                  :alt="selectedAgent.name"
                  class="w-full h-full object-cover rounded-lg"
                  @error="handleImageError"
                >
                <!-- 默认图标背景 -->
                <div v-else class="w-16 h-16 rounded-full bg-white/20 backdrop-blur-sm flex items-center justify-center">
                  <component :is="getWorkflowIcon(selectedAgent.category)" class="w-8 h-8 text-white" />
                </div>
              </div>
              
              <div class="space-y-4">
                <div>
                  <h4 class="font-medium text-gray-900 mb-2">基本信息</h4>
                  <div class="space-y-2 text-sm">
                    <div class="flex justify-between">
                      <span class="text-gray-500">价格:</span>
                      <span :class="getPriceClass(selectedAgent.price)">
                        {{ selectedAgent.price === 0 ? '免费' : `${selectedAgent.price} 积分` }}
                      </span>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">分类:</span>
                      <span>{{ getCategoryDisplayName(selectedAgent.category) }}</span>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">类型:</span>
                      <span>{{ selectedAgent.type }}</span>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">使用次数:</span>
                      <span>{{ selectedAgent.purchaseCount || 0 }}</span>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">评分:</span>
                      <div class="flex items-center">
                        <StarIcon class="w-4 h-4 text-yellow-500 mr-1" />
                        <span>{{ (selectedAgent.rating || 0).toFixed(1) }}</span>
                      </div>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">发布时间:</span>
                      <span>{{ formatDate(selectedAgent.createdAt) }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 右侧：详细说明 -->
            <div class="md:w-2/3">
              <h4 class="text-xl font-semibold text-gray-900 mb-4">{{ selectedAgent.name }}</h4>
              <div class="prose max-w-none">
                <p class="text-gray-600 mb-6">{{ selectedAgent.description || '暂无详细描述' }}</p>
                
                <!-- 工作流定义预览 -->
                <div class="bg-gray-50 rounded-lg p-4">
                  <h5 class="font-medium text-gray-900 mb-2">工作流定义</h5>
                  <pre class="text-sm text-gray-600 whitespace-pre-wrap">{{ formatDefinition(selectedAgent.definition) }}</pre>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 模态框底部 -->
        <div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end gap-3">
          <button 
            @click="closeDetailModal" 
            class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-800 rounded-lg transition-colors"
          >
            关闭
          </button>
          <button 
            @click="purchaseAgent(selectedAgent)"
            :disabled="purchasing.includes(selectedAgent.id)"
            class="bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-400 text-white py-2 px-4 rounded-lg transition-colors flex items-center"
          >
            <ArrowPathIcon v-if="purchasing.includes(selectedAgent.id)" class="w-4 h-4 mr-1 animate-spin" />
            <ShoppingCartIcon v-else class="w-4 h-4 mr-1" />
            {{ selectedAgent.price === 0 ? '免费使用' : '立即购买' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import { agentApi } from '../../api';
import {
  MagnifyingGlassIcon,
  Bars3Icon,
  Squares2X2Icon,
  ArrowPathIcon,
  ExclamationTriangleIcon,
  RectangleStackIcon,
  UserGroupIcon,
  StarIcon,
  CalendarIcon,
  ShoppingCartIcon,
  XMarkIcon,
  CpuChipIcon,
  ChatBubbleLeftRightIcon,
  DocumentTextIcon,
  ChartBarIcon,
  BeakerIcon,
  CodeBracketIcon
} from '@heroicons/vue/24/outline';

const router = useRouter();
const toast = useToast();

// 状态变量
const loading = ref(true);
const error = ref(null);
const agents = ref([]);
const categories = ref([]);
const activeCategory = ref('');
const searchQuery = ref('');
const sortOrder = ref('popular');
const viewMode = ref('grid');
const currentPage = ref(1);
const pageSize = ref(12);
const totalItems = ref(0);
const purchasing = ref([]);
const showDetailModal = ref(false);
const selectedAgent = ref(null);
const searchTimeout = ref(null);

// 计算属性
const totalPages = computed(() => Math.ceil(totalItems.value / pageSize.value));

// 获取工作流图标
const getWorkflowIcon = (category) => {
  switch (category) {
    case 'data':
    case '数据处理':
      return CpuChipIcon;
    case 'nlp':
    case '自然语言处理':
      return ChatBubbleLeftRightIcon;
    case 'content':
    case '内容生成':
      return DocumentTextIcon;
    case 'analysis':
    case '数据分析':
      return ChartBarIcon;
    case 'experiment':
    case '实验':
      return BeakerIcon;
    case 'automation':
    case '自动化':
      return CodeBracketIcon;
    default:
      return CpuChipIcon;
  }
};

// 获取分类显示名称
const getCategoryDisplayName = (category) => {
  const categoryMap = {
    'data': '数据处理',
    'nlp': '自然语言处理',
    'content': '内容生成',
    'analysis': '数据分析',
    'experiment': '实验',
    'automation': '自动化'
  };
  return categoryMap[category] || category;
};

// 获取价格样式类
const getPriceClass = (price) => {
  return price === 0 
    ? 'text-sm font-medium text-green-600' 
    : 'text-sm font-medium text-blue-600';
};

// 获取分类标签样式类
const getCategoryPillClass = (category) => {
  const classMap = {
    'data': 'bg-blue-100 text-blue-800',
    'nlp': 'bg-purple-100 text-purple-800',
    'content': 'bg-green-100 text-green-800',
    'analysis': 'bg-amber-100 text-amber-800',
    'experiment': 'bg-pink-100 text-pink-800',
    'automation': 'bg-indigo-100 text-indigo-800'
  };
  return classMap[category] || 'bg-gray-100 text-gray-800';
};

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN');
};

// 格式化工作流定义
const formatDefinition = (definition) => {
  if (!definition) return '暂无定义';
  if (typeof definition === 'string') {
    try {
      return JSON.stringify(JSON.parse(definition), null, 2);
    } catch (e) {
      return definition;
    }
  }
  return JSON.stringify(definition, null, 2);
};

// 获取分页数字
const getPageNumbers = () => {
  const pages = [];
  const maxPages = 5; // 最多显示5个页码
  let start = Math.max(1, currentPage.value - Math.floor(maxPages / 2));
  let end = Math.min(totalPages.value, start + maxPages - 1);
  
  if (end - start < maxPages - 1) {
    start = Math.max(1, end - maxPages + 1);
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  
  return pages;
};

// 方法
const fetchAgents = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    const params = {
      limit: pageSize.value,
      offset: (currentPage.value - 1) * pageSize.value
    };
    
    if (activeCategory.value) {
      params.category = activeCategory.value;
    }
    
    if (searchQuery.value.trim()) {
      params.search = searchQuery.value.trim();
    }
    
    const response = await agentApi.getAgents(params);
    
    if (response.status === 'success') {
      agents.value = response.data.agents || [];
      totalItems.value = response.data.pagination.total || 0;
    } else {
      throw new Error(response.message || '获取工作流列表失败');
    }
  } catch (err) {
    console.error('获取工作流列表失败:', err);
    error.value = '加载工作流失败：' + (err.message || '未知错误');
    agents.value = [];
    totalItems.value = 0;
  } finally {
    loading.value = false;
  }
};

const fetchCategories = async () => {
  try {
    const response = await agentApi.getAgentCategories();
    if (response.status === 'success') {
      categories.value = response.data || [];
    }
  } catch (err) {
    console.error('获取分类列表失败:', err);
  }
};

const setActiveCategory = (category) => {
  activeCategory.value = category;
  currentPage.value = 1;
  fetchAgents();
};

const handleSearch = () => {
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value);
  }
  
  searchTimeout.value = setTimeout(() => {
    currentPage.value = 1;
    fetchAgents();
  }, 500);
};

const handleSort = () => {
  // 在前端对已获取的数据进行排序
  const sortedAgents = [...agents.value];
  
  switch (sortOrder.value) {
    case 'popular':
      sortedAgents.sort((a, b) => (b.purchaseCount || 0) - (a.purchaseCount || 0));
      break;
    case 'newest':
      sortedAgents.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
      break;
    case 'rating':
      sortedAgents.sort((a, b) => (b.rating || 0) - (a.rating || 0));
      break;
    case 'price_low':
      sortedAgents.sort((a, b) => (a.price || 0) - (b.price || 0));
      break;
    case 'price_high':
      sortedAgents.sort((a, b) => (b.price || 0) - (a.price || 0));
      break;
  }
  
  agents.value = sortedAgents;
};

const toggleView = () => {
  viewMode.value = viewMode.value === 'grid' ? 'list' : 'grid';
};

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
    fetchAgents();
  }
};

const viewAgentDetails = async (agent) => {
  try {
    const response = await agentApi.getAgent(agent.id);
    if (response.status === 'success') {
      selectedAgent.value = response.data;
      showDetailModal.value = true;
    } else {
      throw new Error(response.message || '获取工作流详情失败');
    }
  } catch (err) {
    console.error('获取工作流详情失败:', err);
    toast.error('获取工作流详情失败：' + (err.message || '未知错误'));
  }
};

const closeDetailModal = () => {
  showDetailModal.value = false;
  selectedAgent.value = null;
};

const purchaseAgent = async (agent) => {
  if (purchasing.value.includes(agent.id)) return;
  
  purchasing.value.push(agent.id);
  
  try {
    // 检查是否免费工作流
    if (agent.price === 0) {
      // 免费工作流直接创建
      const workflowData = {
        name: agent.name,
        description: agent.description,
        definition: agent.definition,
        status: 'active',
        agentId: agent.id
      };
      
      const response = await agentApi.createWorkflow(workflowData);
      if (response.status === 'success') {
        toast.success('工作流已添加到您的工作流列表');
        closeDetailModal();
        router.push('/workflows');
      } else {
        throw new Error(response.message || '添加工作流失败');
      }
    } else {
      // 付费工作流需要购买
      const response = await agentApi.purchaseAgent(agent.id);
      if (response.status === 'success') {
        toast.success('购买成功！工作流已添加到您的工作流列表');
        closeDetailModal();
        router.push('/workflows');
      } else {
        throw new Error(response.message || '购买失败');
      }
    }
  } catch (err) {
    console.error('操作失败:', err);
    toast.error(err.message || '操作失败，请稍后重试');
  } finally {
    purchasing.value = purchasing.value.filter(id => id !== agent.id);
  }
};

// 处理图片加载错误
const handleImageError = (event) => {
  // 图片加载失败时隐藏图片，显示默认背景
  event.target.style.display = 'none';
};

// 生命周期
onMounted(async () => {
  await Promise.all([
    fetchAgents(),
    fetchCategories()
  ]);
});

// 监听排序变化
watch(sortOrder, handleSort);
</script>

<style scoped>
.category-pill {
  border-radius: 9999px;
  padding: 0.25rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 500;
  white-space: nowrap;
}

.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style> 