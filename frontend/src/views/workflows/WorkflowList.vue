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
                getWorkflowIconBg(workflow.category)
              ]">
                <component :is="getWorkflowIcon(workflow.category)" :class="getWorkflowIconColor(workflow.category)" class="w-5 h-5" />
              </div>
              <div>
                <h3 class="text-lg font-semibold text-gray-800">{{ workflow.name }}</h3>
                <div class="flex items-center mt-1">
                  <span :class="getBadgeClass(workflow.status)">
                    {{ getStatusText(workflow.status) }}
                  </span>
                  <span class="text-xs text-gray-500 ml-2">更新于 {{ formatTimeAgo(workflow.updated_at) }}</span>
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
              <span class="text-xs text-gray-500">{{ workflow.agent_count || 0 }}个AI代理</span>
            </div>
            <div class="flex items-center">
              <ClockIcon class="w-4 h-4 text-gray-400 mr-1" />
              <span class="text-xs text-gray-500">执行 {{ workflow.execution_count || 0 }} 次</span>
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
    <div v-if="workflows.length === 0 && !loading" class="text-center py-12">
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
  MagnifyingGlassIcon
} from '@heroicons/vue/24/outline';

const router = useRouter();
const toast = useToast();

// 状态变量
const loading = ref(false);
const currentFilter = ref('all');
const sortOrder = ref('newest');
const openDropdown = ref(null);
const executingWorkflows = ref([]);

// 模拟数据
const workflows = ref([
  {
    id: 1,
    name: '数据处理自动化',
    description: '自动处理CSV数据并生成分析报告，定期发送至指定邮箱。',
    status: 'active',
    category: 'data',
    agent_count: 3,
    execution_count: 178,
    created_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 2 * 60 * 60 * 1000)
  },
  {
    id: 2,
    name: '客户支持回复',
    description: '根据客户问题自动生成回复建议，支持多种语言和专业领域。',
    status: 'active',
    category: 'support',
    agent_count: 2,
    execution_count: 245,
    created_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 24 * 60 * 60 * 1000)
  },
  {
    id: 3,
    name: '智能内容审核',
    description: '自动审核文本内容，检测不当言论、敏感信息和违规内容。',
    status: 'active',
    category: 'content',
    agent_count: 4,
    execution_count: 97,
    created_at: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000)
  },
  {
    id: 4,
    name: '市场数据分析',
    description: '收集和分析市场数据，生成趋势报告和投资建议。',
    status: 'active',
    category: 'analysis',
    agent_count: 5,
    execution_count: 156,
    created_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
  },
  {
    id: 5,
    name: '邮件营销自动化',
    description: '根据用户行为自动发送个性化营销邮件。',
    status: 'draft',
    category: 'marketing',
    agent_count: 2,
    execution_count: 0,
    created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000)
  }
]);

// 计算属性
const filteredWorkflows = computed(() => {
  let filtered = workflows.value;
  
  // 按状态过滤
  if (currentFilter.value !== 'all') {
    filtered = filtered.filter(w => w.status === currentFilter.value);
  }
  
  // 排序
  switch (sortOrder.value) {
    case 'newest':
      filtered = filtered.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
      break;
    case 'name':
      filtered = filtered.sort((a, b) => a.name.localeCompare(b.name));
      break;
    case 'last_run':
      filtered = filtered.sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at));
      break;
    case 'most_run':
      filtered = filtered.sort((a, b) => b.execution_count - a.execution_count);
      break;
  }
  
  return filtered;
});

// 方法
const setFilter = (filter) => {
  currentFilter.value = filter;
};

const applySorting = () => {
  // 触发重新计算 filteredWorkflows
};

const formatTimeAgo = (timestamp) => {
  const now = new Date();
  const diff = now - timestamp;
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
    // 模拟执行
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    // 更新执行次数
    const workflow = workflows.value.find(w => w.id === workflowId);
    if (workflow) {
      workflow.execution_count++;
      workflow.updated_at = new Date();
    }
    
    toast.success('工作流执行成功');
  } catch (error) {
    console.error('执行工作流失败:', error);
    toast.error('工作流执行失败');
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
    const original = workflows.value.find(w => w.id === workflowId);
    if (original) {
      const newWorkflow = {
        ...original,
        id: Date.now(),
        name: `${original.name} (副本)`,
        status: 'draft',
        execution_count: 0,
        created_at: new Date(),
        updated_at: new Date()
      };
      workflows.value.push(newWorkflow);
      toast.success('工作流复制成功');
    }
  } catch (error) {
    console.error('复制工作流失败:', error);
    toast.error('工作流复制失败');
  } finally {
    openDropdown.value = null;
  }
};

const archiveWorkflow = async (workflowId) => {
  try {
    const workflow = workflows.value.find(w => w.id === workflowId);
    if (workflow) {
      workflow.status = 'archived';
      workflow.updated_at = new Date();
      toast.success('工作流已归档');
    }
  } catch (error) {
    console.error('归档工作流失败:', error);
    toast.error('工作流归档失败');
  } finally {
    openDropdown.value = null;
  }
};

const deleteWorkflow = async (workflowId) => {
  try {
    workflows.value = workflows.value.filter(w => w.id !== workflowId);
    toast.success('工作流已删除');
  } catch (error) {
    console.error('删除工作流失败:', error);
    toast.error('工作流删除失败');
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

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
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

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style> 