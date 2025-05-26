<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="sm:flex sm:items-center">
        <div class="sm:flex-auto">
          <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">工作流管理</h1>
          <p class="mt-2 text-sm text-gray-700 dark:text-gray-300">
            创建和管理您的自动化工作流，设置触发条件，监控执行状态。
          </p>
        </div>
        <div class="mt-4 sm:mt-0 sm:ml-16 sm:flex-none">
          <router-link
            to="/workflows/create"
            class="inline-flex items-center justify-center rounded-md border border-transparent bg-primary-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 sm:w-auto"
          >
            创建工作流
          </router-link>
        </div>
      </div>
      
      <!-- 搜索和筛选 -->
      <div class="mt-6 flex flex-col sm:flex-row sm:items-center sm:justify-between search-filter-container">
        <div class="flex-1 min-w-0 max-w-lg mb-3 sm:mb-0">
          <div class="relative rounded-md shadow-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              type="text"
              v-model="searchQuery"
              class="block w-full pl-10 py-2 sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
              placeholder="搜索工作流..."
              @input="debounceSearch"
            />
          </div>
        </div>
        <div class="sm:ml-4">
          <select
            v-model="statusFilter"
            class="block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
            @change="applyFilter"
          >
            <option value="">所有状态</option>
            <option value="active">活跃</option>
            <option value="inactive">暂停</option>
            <option value="draft">草稿</option>
          </select>
        </div>
      </div>
      
      <!-- 工作流列表 -->
      <div class="mt-6">
        <!-- 加载状态 -->
        <div v-if="loading" class="text-center py-12">
          <svg class="animate-spin mx-auto h-8 w-8 text-primary-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">加载工作流...</p>
        </div>
        
        <!-- 错误状态 -->
        <div v-else-if="error" class="text-center py-12">
          <svg class="mx-auto h-12 w-12 text-red-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <p class="mt-2 text-sm text-red-500">{{ error }}</p>
          <button 
            @click="fetchWorkflows" 
            class="mt-3 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            重试
          </button>
        </div>
        
        <!-- 空状态 -->
        <div v-else-if="filteredWorkflows.length === 0" class="text-center py-12 bg-white dark:bg-gray-800 shadow rounded-lg">
          <svg class="mx-auto h-12 w-12 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">没有工作流</h3>
          <p v-if="searchQuery || statusFilter" class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            尝试清除搜索或筛选条件，或者创建新的工作流。
          </p>
          <p v-else class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            开始创建您的第一个工作流以自动化您的任务。
          </p>
          <div class="mt-6">
            <router-link
              to="/workflows/create"
              class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              <svg class="-ml-1 mr-2 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              创建工作流
            </router-link>
          </div>
        </div>
        
        <!-- 工作流列表 -->
        <div v-else class="grid grid-cols-1 gap-4 sm:gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div 
            v-for="workflow in paginatedWorkflows" 
            :key="workflow.id" 
            class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg divide-y divide-gray-200 dark:divide-gray-700 workflow-card"
          >
            <div class="px-4 py-5 sm:px-6">
              <div class="flex items-center justify-between">
                <div class="flex-1 min-w-0">
                  <h3 class="text-lg font-medium text-gray-900 dark:text-white truncate">{{ workflow.name }}</h3>
                </div>
                <div class="ml-2 flex-shrink-0 flex">
                  <span :class="getStatusClass(workflow.status)">
                    {{ getStatusText(workflow.status) }}
                  </span>
                </div>
              </div>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400 line-clamp-2">{{ workflow.description || '无描述' }}</p>
            </div>
            <div class="px-4 py-4 sm:px-6">
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span class="ml-1 text-xs text-gray-500 dark:text-gray-400">{{ formatDate(workflow.created_at) }}</span>
                </div>
                <div class="flex items-center">
                  <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                  </svg>
                  <span class="ml-1 text-xs text-gray-500 dark:text-gray-400">
                    {{ workflow.trigger_type === 'schedule' ? '计划' : workflow.trigger_type === 'webhook' ? 'Webhook' : '事件' }}
                  </span>
                </div>
              </div>
            </div>
            <div class="px-4 py-4 sm:px-6 bg-gray-50 dark:bg-gray-700">
              <div class="flex justify-between space-x-2">
                <button
                  @click="executeWorkflow(workflow.id)"
                  class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded text-primary-700 bg-primary-100 hover:bg-primary-200 dark:bg-primary-900 dark:text-primary-100 dark:hover:bg-primary-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                >
                  <svg class="-ml-0.5 mr-1 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  执行
                </button>
                <button
                  @click="editWorkflow(workflow.id)"
                  class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded text-gray-700 bg-gray-100 hover:bg-gray-200 dark:bg-gray-600 dark:text-gray-100 dark:hover:bg-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
                >
                  <svg class="-ml-0.5 mr-1 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                  编辑
                </button>
                <button
                  @click="confirmDeleteWorkflow(workflow.id)"
                  class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded text-red-700 bg-red-100 hover:bg-red-200 dark:bg-red-900 dark:text-red-100 dark:hover:bg-red-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                >
                  <svg class="-ml-0.5 mr-1 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                  删除
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 分页控件 -->
      <div v-if="filteredWorkflows.length > 0" class="mt-6 flex flex-col sm:flex-row justify-between items-center pagination-container">
        <div class="text-sm text-gray-700 dark:text-gray-300 pagination-info">
          显示第 <span class="font-medium">{{ startItem }}</span> 到 <span class="font-medium">{{ endItem }}</span> 项，共 <span class="font-medium">{{ filteredWorkflows.length }}</span> 项
        </div>
        <div class="flex-1 flex justify-center sm:justify-end mt-3 sm:mt-0">
          <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="分页">
            <button
              @click="goToPage(currentPage - 1)"
              :disabled="currentPage === 1"
              :class="[
                'relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium',
                currentPage === 1 
                  ? 'text-gray-300 dark:text-gray-600 cursor-not-allowed' 
                  : 'text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'
              ]"
            >
              <span class="sr-only">上一页</span>
              <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            
            <button
              v-for="page in displayedPages"
              :key="page"
              @click="goToPage(page)"
              :class="[
                'relative inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium',
                page === currentPage
                  ? 'z-10 bg-primary-50 dark:bg-primary-900 border-primary-500 dark:border-primary-500 text-primary-600 dark:text-primary-200'
                  : 'text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700',
                // 在小屏幕上隐藏部分页码按钮
                (page !== 1 && page !== totalPages && page !== currentPage && 
                 page !== currentPage - 1 && page !== currentPage + 1) ? 'pagination-button-extra' : ''
              ]"
            >
              {{ page }}
            </button>
            
            <button
              @click="goToPage(currentPage + 1)"
              :disabled="currentPage === totalPages"
              :class="[
                'relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium',
                currentPage === totalPages 
                  ? 'text-gray-300 dark:text-gray-600 cursor-not-allowed' 
                  : 'text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'
              ]"
            >
              <span class="sr-only">下一页</span>
              <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </nav>
        </div>
      </div>
      
      <!-- 删除确认对话框 -->
      <div v-if="showDeleteConfirm" class="fixed z-10 inset-0 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
        <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
          <!-- 背景遮罩 -->
          <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="cancelDelete"></div>

          <!-- 使内容居中 -->
          <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

          <!-- 对话框内容 -->
          <div class="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
            <div class="bg-white dark:bg-gray-800 px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <div class="sm:flex sm:items-start">
                <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100 dark:bg-red-900 sm:mx-0 sm:h-10 sm:w-10">
                  <svg class="h-6 w-6 text-red-600 dark:text-red-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                </div>
                <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                  <h3 class="text-lg leading-6 font-medium text-gray-900 dark:text-white" id="modal-title">
                    删除工作流
                  </h3>
                  <div class="mt-2">
                    <p class="text-sm text-gray-500 dark:text-gray-400">
                      您确定要删除这个工作流吗？此操作无法撤销，所有相关的配置和历史记录都将被永久删除。
                    </p>
                  </div>
                </div>
              </div>
            </div>
            <div class="bg-gray-50 dark:bg-gray-700 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button 
                type="button" 
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm"
                @click="confirmDelete"
              >
                删除
              </button>
              <button 
                type="button" 
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 dark:border-gray-600 shadow-sm px-4 py-2 bg-white dark:bg-gray-800 text-base font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                @click="cancelDelete"
              >
                取消
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useWorkflowStore } from '../../stores/workflow';
import { useToast } from 'vue-toastification';

const router = useRouter();
const workflowStore = useWorkflowStore();
const toast = useToast();

// 状态变量
const loading = ref(false);
const error = ref(null);
const searchQuery = ref('');
const statusFilter = ref('');
const currentPage = ref(1);
const pageSize = ref(9);
const showDeleteConfirm = ref(false);
const workflowToDelete = ref(null);

// 计算属性
const filteredWorkflows = computed(() => {
  return workflowStore.filteredWorkflows;
});

const totalPages = computed(() => {
  return Math.ceil(filteredWorkflows.value.length / pageSize.value);
});

const paginatedWorkflows = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return filteredWorkflows.value.slice(start, end);
});

const startItem = computed(() => {
  if (filteredWorkflows.value.length === 0) return 0;
  return (currentPage.value - 1) * pageSize.value + 1;
});

const endItem = computed(() => {
  if (filteredWorkflows.value.length === 0) return 0;
  return Math.min(currentPage.value * pageSize.value, filteredWorkflows.value.length);
});

const displayedPages = computed(() => {
  const maxVisiblePages = 5;
  const halfVisible = Math.floor(maxVisiblePages / 2);
  
  if (totalPages.value <= maxVisiblePages) {
    // 如果总页数小于等于最大可见页数，则显示所有页码
    return Array.from({ length: totalPages.value }, (_, i) => i + 1);
  }
  
  let startPage = Math.max(currentPage.value - halfVisible, 1);
  let endPage = startPage + maxVisiblePages - 1;
  
  if (endPage > totalPages.value) {
    endPage = totalPages.value;
    startPage = Math.max(endPage - maxVisiblePages + 1, 1);
  }
  
  return Array.from({ length: endPage - startPage + 1 }, (_, i) => startPage + i);
});

// 方法
const fetchWorkflows = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    await workflowStore.fetchWorkflows();
  } catch (err) {
    error.value = err.message || '获取工作流失败';
    console.error('获取工作流失败:', err);
  } finally {
    loading.value = false;
  }
};

const debounceSearch = (() => {
  let timeout;
  return () => {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      applyFilter();
    }, 300);
  };
})();

const applyFilter = () => {
  workflowStore.setFilter({
    search: searchQuery.value,
    status: statusFilter.value
  });
  currentPage.value = 1; // 重置到第一页
};

const formatDate = (dateString) => {
  if (!dateString) return '未知';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('zh-CN', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date);
};

const getStatusClass = (status) => {
  const baseClasses = 'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium';
  
  switch (status) {
    case 'active':
      return `${baseClasses} bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100`;
    case 'inactive':
      return `${baseClasses} bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100`;
    case 'draft':
      return `${baseClasses} bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-100`;
    default:
      return `${baseClasses} bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-100`;
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

const editWorkflow = (id) => {
  router.push(`/workflows/edit/${id}`);
};

const executeWorkflow = async (id) => {
  try {
    loading.value = true;
    await workflowStore.executeWorkflow(id);
    toast.success('工作流执行已开始');
    router.push('/executions');
  } catch (err) {
    toast.error(`执行失败: ${err.message || '未知错误'}`);
    console.error('执行工作流失败:', err);
  } finally {
    loading.value = false;
  }
};

const confirmDeleteWorkflow = (id) => {
  workflowToDelete.value = id;
  showDeleteConfirm.value = true;
};

const confirmDelete = async () => {
  if (workflowToDelete.value) {
    await deleteWorkflow(workflowToDelete.value);
  }
};

const cancelDelete = () => {
  showDeleteConfirm.value = false;
  workflowToDelete.value = null;
};

const deleteWorkflow = async (id) => {
  try {
    loading.value = true;
    await workflowStore.deleteWorkflow(id);
    toast.success('工作流已删除');
    showDeleteConfirm.value = false;
    workflowToDelete.value = null;
  } catch (err) {
    toast.error(`删除失败: ${err.message || '未知错误'}`);
    console.error('删除工作流失败:', err);
  } finally {
    loading.value = false;
  }
};

const goToPage = (page) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  // 如果需要，这里可以添加滚动到页面顶部的逻辑
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

// 生命周期钩子
onMounted(() => {
  fetchWorkflows();
});
</script>

<style scoped>
/* 添加行高限制，防止描述过长 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 响应式调整 */
@media (max-width: 640px) {
  /* 在小屏幕上调整搜索和筛选布局 */
  .search-filter-container {
    flex-direction: column;
  }
  
  /* 在小屏幕上调整分页控件 */
  .pagination-container {
    flex-direction: column;
    align-items: center;
  }
  
  .pagination-info {
    margin-bottom: 1rem;
  }
  
  /* 在小屏幕上隐藏部分页码按钮 */
  .pagination-button-extra {
    display: none;
  }
}

/* 暗色模式适配 */
@media (prefers-color-scheme: dark) {
  .workflow-card {
    background-color: var(--color-gray-800);
    border-color: var(--color-gray-700);
  }
}
</style> 