<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="sm:flex sm:items-center">
        <div class="sm:flex-auto">
          <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">执行历史</h1>
          <p class="mt-2 text-sm text-gray-700 dark:text-gray-300">
            所有工作流的执行记录列表，包括成功和失败的执行。
          </p>
        </div>
      </div>
      
      <!-- 搜索和筛选 -->
      <div class="mt-6">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-4">
          <div class="sm:col-span-2">
            <label for="search" class="sr-only">搜索</label>
            <div class="relative rounded-md shadow-sm">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <input
                type="text"
                name="search"
                id="search"
                v-model="searchQuery"
                class="block w-full pl-10 py-2 sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
                placeholder="搜索工作流名称或ID..."
                @input="debounceSearch"
              />
            </div>
          </div>
          <div>
            <select
              v-model="statusFilter"
              class="block w-full py-2 pl-3 pr-10 sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
              @change="fetchExecutions"
            >
              <option value="">所有状态</option>
              <option value="success">成功</option>
              <option value="failed">失败</option>
              <option value="running">运行中</option>
              <option value="pending">等待中</option>
              <option value="cancelled">已取消</option>
            </select>
          </div>
          <div>
            <select
              v-model="timeFilter"
              class="block w-full py-2 pl-3 pr-10 sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
              @change="fetchExecutions"
            >
              <option value="all">所有时间</option>
              <option value="today">今天</option>
              <option value="yesterday">昨天</option>
              <option value="week">本周</option>
              <option value="month">本月</option>
            </select>
          </div>
        </div>
      </div>
      
      <!-- 加载状态 -->
      <div v-if="loading" class="flex justify-center items-center py-12">
        <div class="spinner"></div>
        <span class="ml-3">加载执行历史...</span>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="bg-red-50 dark:bg-red-900/20 p-4 rounded-md">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800 dark:text-red-200">加载失败</h3>
            <div class="mt-2 text-sm text-red-700 dark:text-red-300">
              <p>{{ error }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 执行记录为空 -->
      <div v-else-if="executions.length === 0" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">暂无执行记录</h3>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ searchQuery || statusFilter || timeFilter !== 'all' ? '没有符合条件的执行记录，请尝试修改筛选条件。' : '暂无任何工作流执行记录，请先执行一个工作流。' }}
        </p>
      </div>
      
      <!-- 执行列表 -->
      <div v-else class="mt-6 flex flex-col">
        <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
            <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
              <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-700">
                <thead class="bg-gray-50 dark:bg-gray-800">
                  <tr>
                    <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 dark:text-white sm:pl-6">ID</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900 dark:text-white">工作流</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900 dark:text-white">触发方式</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900 dark:text-white">状态</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900 dark:text-white">开始时间</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900 dark:text-white">执行时长</th>
                    <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                      <span class="sr-only">操作</span>
                    </th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200 dark:divide-gray-700 bg-white dark:bg-gray-900">
                  <tr v-for="execution in executions" :key="execution.id" class="hover:bg-gray-50 dark:hover:bg-gray-800">
                    <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm sm:pl-6">
                      <div class="flex items-center">
                        <div>
                          <div class="font-medium text-gray-900 dark:text-white">
                            {{ execution.id.substring(0, 8) }}
                          </div>
                          <div class="text-gray-500 dark:text-gray-400">
                            {{ execution.workflow_name || execution.workflowName }}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500 dark:text-gray-400">
                      {{ triggerText(execution.trigger_type || execution.triggerType) }}
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                      <span 
                        :class="statusBadgeClass(execution.status)" 
                        class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                      >
                        {{ statusText(execution.status) }}
                      </span>
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500 dark:text-gray-400">
                      {{ formatDate(execution.start_time || execution.startedAt) }}
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500 dark:text-gray-400">
                      {{ formatDuration(execution.duration) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div class="flex justify-end space-x-2">
                        <router-link
                          :to="`/executions/${execution.id}`"
                          class="text-primary-600 dark:text-primary-400 hover:text-primary-900 dark:hover:text-primary-300"
                        >
                          查看详情
                        </router-link>
                        <button
                          v-if="execution.status === 'running' || execution.status === 'pending'"
                          @click="cancelExecution(execution.id)"
                          class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300"
                        >
                          取消
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 分页 -->
      <div class="py-4 flex items-center justify-between border-t border-gray-200 dark:border-gray-700 mt-4">
        <div class="flex-1 flex justify-between sm:hidden">
          <button
            @click="goToPage(currentPage - 1)"
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700"
            :class="currentPage === 1 ? 'opacity-50 cursor-not-allowed' : ''"
          >
            上一页
          </button>
          <button
            @click="goToPage(currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700"
            :class="currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : ''"
          >
            下一页
          </button>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700 dark:text-gray-300">
              显示第 <span class="font-medium">{{ startItem }}</span> 至 <span class="font-medium">{{ endItem }}</span> 项，共 <span class="font-medium">{{ totalItems }}</span> 项
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
              <button
                @click="goToPage(currentPage - 1)"
                :disabled="currentPage === 1"
                class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700"
                :class="currentPage === 1 ? 'opacity-50 cursor-not-allowed' : ''"
              >
                <span class="sr-only">上一页</span>
                <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
              
              <template v-for="page in displayedPages">
                <button
                  v-if="page !== '...'"
                  @click="goToPage(page)"
                  :key="page"
                  :class="page === currentPage ? 'z-10 bg-primary-50 dark:bg-primary-900 border-primary-500 text-primary-600 dark:text-primary-200' : 'bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'"
                  class="relative inline-flex items-center px-4 py-2 border text-sm font-medium"
                >
                  {{ page }}
                </button>
                <span
                  v-else
                  :key="'ellipsis-' + page"
                  class="relative inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium text-gray-700 dark:text-gray-300"
                >
                  ...
                </span>
              </template>
              
              <button
                @click="goToPage(currentPage + 1)"
                :disabled="currentPage === totalPages"
                class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700"
                :class="currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : ''"
              >
                <span class="sr-only">下一页</span>
                <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { format, formatDistanceStrict } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import { useToast } from 'vue-toastification';
import { useExecutionStore } from '../../stores/execution';
import { useRouter } from 'vue-router';

const toast = useToast();
const executionStore = useExecutionStore();
const router = useRouter();

const loading = ref(false);
const executions = ref([]);
const searchQuery = ref('');
const statusFilter = ref('');
const timeFilter = ref('all');
const totalItems = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const error = ref(null);

// 计算属性
const totalPages = computed(() => Math.ceil(totalItems.value / pageSize.value) || 1);
const startItem = computed(() => ((currentPage.value - 1) * pageSize.value) + 1);
const endItem = computed(() => Math.min(currentPage.value * pageSize.value, totalItems.value));

const displayedPages = computed(() => {
  const pages = [];
  if (totalPages.value <= 7) {
    for (let i = 1; i <= totalPages.value; i++) {
      pages.push(i);
    }
  } else {
    pages.push(1);
    if (currentPage.value > 3) {
      pages.push('...');
    }
    
    let startPage = Math.max(2, currentPage.value - 1);
    let endPage = Math.min(totalPages.value - 1, currentPage.value + 1);
    
    if (currentPage.value <= 3) {
      endPage = Math.min(5, totalPages.value - 1);
    }
    
    if (currentPage.value >= totalPages.value - 2) {
      startPage = Math.max(2, totalPages.value - 4);
    }
    
    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }
    
    if (currentPage.value < totalPages.value - 2) {
      pages.push('...');
    }
    
    pages.push(totalPages.value);
  }
  return pages;
});

// 方法
const fetchExecutions = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    // 准备查询参数
    const params = {
      limit: pageSize.value,
      offset: (currentPage.value - 1) * pageSize.value,
    };
    
    // 添加搜索参数
    if (searchQuery.value) {
      params.search = searchQuery.value;
    }
    
    // 添加状态过滤
    if (statusFilter.value) {
      params.status = statusFilter.value;
    }
    
    // 添加时间过滤
    if (timeFilter.value) {
      const now = new Date();
      let startDate = null;
      
      switch (timeFilter.value) {
        case 'today':
          startDate = new Date(now.setHours(0, 0, 0, 0));
          break;
        case 'yesterday':
          startDate = new Date(now.setDate(now.getDate() - 1));
          startDate.setHours(0, 0, 0, 0);
          break;
        case 'week':
          startDate = new Date(now.setDate(now.getDate() - 7));
          break;
        case 'month':
          startDate = new Date(now.setMonth(now.getMonth() - 1));
          break;
      }
      
      if (startDate) {
        params.start_date = startDate.toISOString();
      }
    }
    
    // 调用 store 的 fetchExecutions 方法
    const response = await executionStore.fetchExecutions(params);
    
    // 更新本地状态
    executions.value = response.data || [];
    totalItems.value = response.total || 0;
  } catch (err) {
    console.error('获取执行历史失败:', err);
    error.value = err.message || '获取执行历史失败';
    toast.error(`获取执行历史失败: ${err.message || '未知错误'}`);
  } finally {
    loading.value = false;
  }
};

const debounceSearch = (() => {
  let timeout;
  return () => {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      currentPage.value = 1;
      fetchExecutions();
    }, 300);
  };
})();

const formatDate = (date) => {
  if (!date) return '未知';
  return format(new Date(date), 'MM月dd日 HH:mm:ss', { locale: zhCN });
};

const formatDuration = (seconds) => {
  if (seconds === null || seconds === undefined) return '进行中';
  
  if (seconds < 60) {
    return `${seconds}秒`;
  }
  
  const durationDate = new Date(seconds * 1000);
  return formatDistanceStrict(new Date(0), durationDate, { locale: zhCN });
};

const goToPage = (page) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  fetchExecutions();
};

const statusText = (status) => {
  const statusMap = {
    success: '成功',
    failed: '失败',
    running: '运行中',
    pending: '等待中',
    cancelled: '已取消'
  };
  return statusMap[status] || status;
};

const statusBadgeClass = (status) => {
  const badgeMap = {
    success: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
    failed: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
    running: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
    pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
    cancelled: 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200'
  };
  return badgeMap[status] || 'bg-gray-100 text-gray-800';
};

const triggerText = (triggerType) => {
  const triggerMap = {
    schedule: '定时计划',
    manual: '手动触发',
    webhook: 'Webhook',
    event: '事件触发'
  };
  return triggerMap[triggerType] || triggerType || '未知';
};

const cancelExecution = async (executionId) => {
  if (!confirm(`确定要取消执行吗？`)) {
    return;
  }
  
  try {
    await executionStore.cancelExecution(executionId);
    toast.success(`已取消执行`);
    
    // 刷新执行列表
    await fetchExecutions();
  } catch (error) {
    console.error('取消执行失败', error);
    toast.error(`取消执行失败: ${error.message || '未知错误'}`);
  }
};

onMounted(() => {
  fetchExecutions();
});
</script> 