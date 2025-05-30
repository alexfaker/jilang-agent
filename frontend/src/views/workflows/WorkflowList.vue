<template>
  <div class="flex-1 overflow-y-auto p-6">
    <!-- 工作流过滤选项 -->
    <div class="bg-white rounded-lg p-4 mb-6 flex flex-wrap items-center justify-between gap-4 shadow-sm">
      <div class="flex flex-wrap gap-2">
        <button 
          @click="setFilter('all')"
          :class="[
            'px-4 py-2 rounded-lg font-medium transition-colors',
            currentFilter === 'all' 
              ? 'bg-indigo-600 text-white' 
              : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300'
          ]"
        >
          全部工作流
        </button>
        <button 
          @click="setFilter('active')"
          :class="[
            'px-4 py-2 rounded-lg font-medium transition-colors',
            currentFilter === 'active' 
              ? 'bg-indigo-600 text-white' 
              : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300'
          ]"
        >
          活跃的
        </button>
        <button 
          @click="setFilter('draft')"
          :class="[
            'px-4 py-2 rounded-lg font-medium transition-colors',
            currentFilter === 'draft' 
              ? 'bg-indigo-600 text-white' 
              : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300'
          ]"
        >
          草稿
        </button>
        <button 
          @click="setFilter('archived')"
          :class="[
            'px-4 py-2 rounded-lg font-medium transition-colors',
            currentFilter === 'archived' 
              ? 'bg-indigo-600 text-white' 
              : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300'
          ]"
        >
          已归档
        </button>
      </div>
      
      <div class="flex items-center">
        <label for="sortOrder" class="text-sm text-gray-600 mr-2">排序方式:</label>
        <select 
          id="sortOrder" 
          v-model="sortOrder"
          @change="applySorting"
          class="py-2 px-3 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
        >
          <option value="newest">最新创建</option>
          <option value="name">按名称</option>
          <option value="last_run">最近执行</option>
          <option value="most_run">执行次数</option>
        </select>
      </div>
    </div>
    
    <!-- 工作流列表 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- 工作流卡片 -->
      <div 
        v-for="workflow in filteredWorkflows" 
        :key="workflow.id" 
        class="workflow-card bg-white overflow-hidden rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-2"
      >
        <div class="p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center">
              <div :class="[
                'w-10 h-10 rounded-full flex items-center justify-center mr-3',
                getWorkflowIconBg(getCategoryFromWorkflow(workflow))
              ]">
                <component :is="getWorkflowIcon(getCategoryFromWorkflow(workflow))" :class="getWorkflowIconColor(getCategoryFromWorkflow(workflow))" class="w-5 h-5" />
              </div>
              <div>
                <h3 class="text-lg font-semibold text-gray-800">{{ workflow.name }}</h3>
                <div class="flex items-center mt-1">
                  <span :class="getBadgeClass(workflow.status)">
                    {{ getStatusText(workflow.status) }}
                  </span>
                  <span class="text-xs text-gray-500 ml-2">更新于 {{ formatTimeAgo(workflow.updatedAt || workflow.updated_at) }}</span>
                </div>
              </div>
            </div>
            <div class="dropdown relative">
              <button 
                @click="toggleDropdown(workflow.id)"
                class="text-gray-500 hover:text-gray-700 focus:outline-none p-1"
              >
                <EllipsisVerticalIcon class="w-5 h-5" />
              </button>
              <div 
                v-if="openDropdown === workflow.id"
                class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-10 border border-gray-200"
              >
                <div class="py-1">
                  <button 
                    @click="editWorkflow(workflow.id)"
                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 w-full text-left"
                  >
                    编辑工作流
                  </button>
                  <button 
                    @click="duplicateWorkflow(workflow.id)"
                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 w-full text-left"
                  >
                    复制工作流
                  </button>
                  <button 
                    @click="archiveWorkflow(workflow.id)"
                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 w-full text-left"
                  >
                    归档工作流
                  </button>
                  <hr class="my-1" />
                  <button 
                    @click="deleteWorkflow(workflow.id)"
                    class="block px-4 py-2 text-sm text-red-600 hover:bg-red-50 w-full text-left"
                  >
                    删除工作流
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <p class="text-gray-600 text-sm mb-4 line-clamp-2">{{ workflow.description || '暂无描述' }}</p>
          
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center">
              <ShareIcon class="w-4 h-4 text-gray-400 mr-1" />
              <span class="text-xs text-gray-500">{{ workflow.agentCount || workflow.agent_count || 0 }}个AI代理</span>
            </div>
            <div class="flex items-center">
              <ClockIcon class="w-4 h-4 text-gray-400 mr-1" />
              <span class="text-xs text-gray-500">执行 {{ workflow.runCount || workflow.execution_count || 0 }} 次</span>
            </div>
          </div>
          
          <div class="border-t pt-4 flex justify-between">
            <router-link 
              :to="`/workflows/${workflow.id}`" 
              class="text-indigo-600 hover:text-indigo-800 text-sm font-medium transition-colors"
            >
              查看详情
            </router-link>
            <button 
              @click="executeWorkflow(workflow.id)"
              :disabled="executingWorkflows.includes(workflow.id)"
              class="bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-400 text-white text-sm py-1 px-3 rounded-lg transition-colors flex items-center"
            >
              <PlayIcon v-if="!executingWorkflows.includes(workflow.id)" class="w-4 h-4 mr-1" />
              <ArrowPathIcon v-else class="w-4 h-4 mr-1 animate-spin" />
              执行
            </button>
          </div>
        </div>
      </div>
      
      <!-- 创建新工作流卡片 -->
      <router-link
        to="/workflows/create"
        class="workflow-card bg-white overflow-hidden rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-2 border-2 border-dashed border-gray-300 hover:border-indigo-500"
      >
        <div class="p-6 h-full flex flex-col items-center justify-center text-center">
          <div class="w-12 h-12 rounded-full bg-indigo-100 flex items-center justify-center mb-4">
            <PlusIcon class="w-6 h-6 text-indigo-600" />
          </div>
          <h3 class="text-lg font-semibold text-gray-800 mb-2">创建新工作流</h3>
          <p class="text-sm text-gray-600">开始构建您的自动化工作流</p>
        </div>
      </router-link>
    </div>
    
    <!-- 空状态 -->
    <div v-if="workflows.length === 0 && !loading && !error" class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
        <RectangleStackIcon class="w-8 h-8 text-gray-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">还没有工作流</h3>
      <p class="text-gray-600 mb-6">创建您的第一个工作流来开始自动化任务</p>
      <router-link
        to="/workflows/create"
        class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
      >
        <PlusIcon class="w-5 h-5 mr-2" />
        创建工作流
      </router-link>
    </div>
    
    <!-- 错误状态 -->
    <div v-if="error && !loading" class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 bg-red-100 rounded-full flex items-center justify-center">
        <ExclamationTriangleIcon class="w-8 h-8 text-red-500" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">加载失败</h3>
      <p class="text-gray-600 mb-6">{{ error }}</p>
      <button
        @click="fetchWorkflows"
        class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
      >
        <ArrowPathIcon class="w-5 h-5 mr-2" />
        重试
      </button>
    </div>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="text-center py-12">
      <ArrowPathIcon class="w-8 h-8 mx-auto mb-4 text-indigo-600 animate-spin" />
      <p class="text-gray-600">加载工作流...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import { workflowApi } from '../../api';
import {
  PlayIcon,
  EllipsisVerticalIcon,
  ShareIcon,
  ClockIcon,
  ArrowPathIcon,
  PlusIcon,
  RectangleStackIcon,
  CpuChipIcon,
  ChatBubbleLeftRightIcon,
  DocumentTextIcon,
  MagnifyingGlassIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline';

const router = useRouter();
const toast = useToast();

// 状态变量
const loading = ref(true);
const error = ref(null);
const currentFilter = ref('all');
const sortOrder = ref('newest');
const openDropdown = ref(null);
const executingWorkflows = ref([]);
const workflows = ref([]);
const totalItems = ref(0);
const currentPage = ref(1);
const itemsPerPage = ref(20);

// 获取工作流数据
const fetchWorkflows = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    const params = {
      limit: itemsPerPage.value,
      offset: (currentPage.value - 1) * itemsPerPage.value
    };
    
    // 如果有状态筛选且不是'all'，添加状态参数
    if (currentFilter.value !== 'all') {
      params.status = currentFilter.value;
    }
    
    const response = await workflowApi.getWorkflows(params);
    
    if (response.status === 'success') {
      workflows.value = response.data.workflows || [];
      totalItems.value = response.data.pagination.total || 0;
    } else {
      throw new Error(response.message || '获取工作流列表失败');
    }
  } catch (err) {
    console.error('获取工作流列表失败:', err);
    error.value = '加载工作流失败：' + (err.message || '未知错误');
    workflows.value = [];
    totalItems.value = 0;
    toast.error('加载工作流失败：' + (err.message || '未知错误'));
  } finally {
    loading.value = false;
  }
};

// 计算属性
const filteredWorkflows = computed(() => {
  let filtered = workflows.value;
  
  // API已经处理了状态过滤，这里主要处理排序
  switch (sortOrder.value) {
    case 'newest':
      filtered = filtered.sort((a, b) => new Date(b.createdAt || b.created_at) - new Date(a.createdAt || a.created_at));
      break;
    case 'name':
      filtered = filtered.sort((a, b) => a.name.localeCompare(b.name));
      break;
    case 'last_run':
      filtered = filtered.sort((a, b) => new Date(b.updatedAt || b.updated_at) - new Date(a.updatedAt || a.updated_at));
      break;
    case 'most_run':
      filtered = filtered.sort((a, b) => (b.runCount || b.execution_count || 0) - (a.runCount || a.execution_count || 0));
      break;
  }
  
  return filtered;
});

// 方法
const setFilter = async (filter) => {
  currentFilter.value = filter;
  currentPage.value = 1; // 重置页码
  await fetchWorkflows();
};

const applySorting = () => {
  // 触发重新计算 filteredWorkflows
};

const formatTimeAgo = (timestamp) => {
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now - date;
  const minutes = Math.floor(diff / (1000 * 60));
  const hours = Math.floor(diff / (1000 * 60 * 60));
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (minutes < 60) {
    return `${minutes}分钟前`;
  } else if (hours < 24) {
    return `${hours}小时前`;
  } else if (days === 1) {
    return '昨天';
  } else if (days < 7) {
    return `${days}天前`;
  } else {
    return `${Math.floor(days / 7)}周前`;
  }
};

// 获取工作流分类（从名称或描述推断）
const getCategoryFromWorkflow = (workflow) => {
  const name = workflow.name.toLowerCase();
  const description = (workflow.description || '').toLowerCase();
  const content = name + ' ' + description;
  
  if (content.includes('数据') || content.includes('分析') || content.includes('统计')) {
    return 'data';
  } else if (content.includes('客服') || content.includes('支持') || content.includes('回复')) {
    return 'support';
  } else if (content.includes('内容') || content.includes('审核') || content.includes('文本')) {
    return 'content';
  } else if (content.includes('分析') || content.includes('市场') || content.includes('投资')) {
    return 'analysis';
  } else {
    return 'data'; // 默认分类
  }
};

const getWorkflowIcon = (category) => {
  switch (category) {
    case 'data':
      return CpuChipIcon;
    case 'support':
      return ChatBubbleLeftRightIcon;
    case 'content':
      return DocumentTextIcon;
    case 'analysis':
      return MagnifyingGlassIcon;
    default:
      return CpuChipIcon;
  }
};

const getWorkflowIconBg = (category) => {
  switch (category) {
    case 'data':
      return 'bg-blue-100';
    case 'support':
      return 'bg-purple-100';
    case 'content':
      return 'bg-green-100';
    case 'analysis':
      return 'bg-amber-100';
    default:
      return 'bg-gray-100';
  }
};

const getWorkflowIconColor = (category) => {
  switch (category) {
    case 'data':
      return 'text-blue-600';
    case 'support':
      return 'text-purple-600';
    case 'content':
      return 'text-green-600';
    case 'analysis':
      return 'text-amber-600';
    default:
      return 'text-gray-600';
  }
};

const getBadgeClass = (status) => {
  const baseClasses = 'badge';
  
  switch (status) {
    case 'active':
      return `${baseClasses} badge-active`;
    case 'inactive':
      return `${baseClasses} badge-inactive`;
    case 'draft':
      return `${baseClasses} badge-draft`;
    case 'archived':
      return `${baseClasses} badge-archived`;
    default:
      return `${baseClasses} badge-draft`;
  }
};

const getStatusText = (status) => {
  switch (status) {
    case 'active':
      return '活跃';
    case 'inactive':
      return '暂停';
    case 'draft':
      return '草稿';
    case 'archived':
      return '已归档';
    default:
      return '未知';
  }
};

const toggleDropdown = (workflowId) => {
  openDropdown.value = openDropdown.value === workflowId ? null : workflowId;
};

const executeWorkflow = async (workflowId) => {
  if (executingWorkflows.value.includes(workflowId)) return;
  
  executingWorkflows.value.push(workflowId);
  
  try {
    const response = await workflowApi.executeWorkflow(workflowId, {});
    
    if (response.status === 'success') {
      toast.success('工作流执行成功');
      await fetchWorkflows(); // 重新加载数据以更新执行次数
    } else {
      throw new Error(response.message || '执行失败');
    }
  } catch (error) {
    console.error('执行工作流失败:', error);
    toast.error('工作流执行失败：' + (error.message || '未知错误'));
  } finally {
    executingWorkflows.value = executingWorkflows.value.filter(id => id !== workflowId);
    openDropdown.value = null;
  }
};

const editWorkflow = (workflowId) => {
  router.push(`/workflows/edit/${workflowId}`);
  openDropdown.value = null;
};

const duplicateWorkflow = async (workflowId) => {
  try {
    // 先获取原始工作流
    const originalResponse = await workflowApi.getWorkflow(workflowId);
    if (originalResponse.status !== 'success') {
      throw new Error('获取原始工作流失败');
    }
    
    const original = originalResponse.data;
    const newWorkflowData = {
      name: `${original.name} (副本)`,
      description: original.description,
      definition: original.definition,
      status: 'draft'
    };
    
    const response = await workflowApi.createWorkflow(newWorkflowData);
    if (response.status === 'success') {
      toast.success('工作流复制成功');
      await fetchWorkflows(); // 重新加载列表
    } else {
      throw new Error(response.message || '复制失败');
    }
  } catch (error) {
    console.error('复制工作流失败:', error);
    toast.error('工作流复制失败：' + (error.message || '未知错误'));
  } finally {
    openDropdown.value = null;
  }
};

const archiveWorkflow = async (workflowId) => {
  try {
    const response = await workflowApi.updateWorkflow(workflowId, { status: 'archived' });
    if (response.status === 'success') {
      toast.success('工作流已归档');
      await fetchWorkflows(); // 重新加载列表
    } else {
      throw new Error(response.message || '归档失败');
    }
  } catch (error) {
    console.error('归档工作流失败:', error);
    toast.error('工作流归档失败：' + (error.message || '未知错误'));
  } finally {
    openDropdown.value = null;
  }
};

const deleteWorkflow = async (workflowId) => {
  if (!confirm('确定要删除这个工作流吗？此操作不可撤销。')) {
    openDropdown.value = null;
    return;
  }
  
  try {
    const response = await workflowApi.deleteWorkflow(workflowId);
    if (response.status === 'success') {
      toast.success('工作流已删除');
      await fetchWorkflows(); // 重新加载列表
    } else {
      throw new Error(response.message || '删除失败');
    }
  } catch (error) {
    console.error('删除工作流失败:', error);
    toast.error('工作流删除失败：' + (error.message || '未知错误'));
  } finally {
    openDropdown.value = null;
  }
};

// 点击外部关闭下拉菜单
const handleClickOutside = (event) => {
  if (!event.target.closest('.dropdown')) {
    openDropdown.value = null;
  }
};

onMounted(async () => {
  document.addEventListener('click', handleClickOutside);
  await fetchWorkflows();
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped>
.workflow-card {
  transition: all 0.3s ease;
}

.workflow-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08);
}

.badge {
  padding: 0.125rem 0.5rem;
  border-radius: 0.75rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.badge-active {
  background-color: #dcfce7;
  color: #166534;
}

.badge-inactive {
  background-color: #fee2e2;
  color: #991b1b;
}

.badge-draft {
  background-color: #e5e7eb;
  color: #374151;
}

.badge-archived {
  background-color: #fef2f2;
  color: #b91c1c;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style> 