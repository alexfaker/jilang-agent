<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">代理管理</h1>
        <router-link
          to="/agents/create"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          <svg class="-ml-1 mr-2 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          创建代理
        </router-link>
      </div>

      <!-- 搜索和筛选 -->
      <div class="mt-6">
        <div class="flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4">
          <div class="flex-1">
            <label for="search" class="sr-only">搜索代理</label>
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
                placeholder="搜索代理..."
                @input="debounceSearch"
              />
            </div>
          </div>
          <div class="w-full sm:w-48">
            <select
              v-model="statusFilter"
              class="block w-full py-2 pl-3 pr-10 sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
              @change="fetchAgents"
            >
              <option value="">全部状态</option>
              <option value="online">在线</option>
              <option value="offline">离线</option>
              <option value="busy">忙碌</option>
            </select>
          </div>
          <div class="w-full sm:w-48">
            <select
              v-model="typeFilter"
              class="block w-full py-2 pl-3 pr-10 sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500"
              @change="fetchAgents"
            >
              <option value="">全部类型</option>
              <option value="llm">LLM</option>
              <option value="executor">执行器</option>
              <option value="retriever">检索器</option>
            </select>
          </div>
        </div>
      </div>

      <!-- 代理列表 -->
      <div class="mt-6">
        <div class="flex flex-col">
          <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
              <div class="shadow overflow-hidden border-b border-gray-200 dark:border-gray-700 sm:rounded-lg">
                <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                  <thead class="bg-gray-50 dark:bg-gray-800">
                    <tr>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        代理名称
                      </th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        类型
                      </th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        状态
                      </th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        工作流数
                      </th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        最后活动
                      </th>
                      <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        操作
                      </th>
                    </tr>
                  </thead>
                  <tbody v-if="loading" class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-800">
                    <tr>
                      <td colspan="6" class="px-6 py-4 whitespace-nowrap">
                        <div class="flex justify-center">
                          <svg class="animate-spin h-5 w-5 text-primary-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                          </svg>
                          <span class="ml-2">加载中...</span>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                  <tbody v-else-if="agents.length === 0" class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-800">
                    <tr>
                      <td colspan="6" class="px-6 py-4 whitespace-nowrap text-center text-gray-500 dark:text-gray-400">
                        没有找到匹配的代理
                      </td>
                    </tr>
                  </tbody>
                  <tbody v-else class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-800">
                    <tr v-for="agent in agents" :key="agent.id">
                      <td class="px-6 py-4 whitespace-nowrap">
                        <div class="flex items-center">
                          <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-full" :class="agentIconClass(agent.type)">
                            <svg class="h-6 w-6 text-white" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                              <path v-if="agent.type === 'llm'" d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v3h8v-3zM6 8a2 2 0 11-4 0 2 2 0 014 0zM16 18v-3a5.972 5.972 0 00-.75-2.906A3.005 3.005 0 0119 15v3h-3zM4.75 12.094A5.973 5.973 0 004 15v3H1v-3a3 3 0 013.75-2.906z"></path>
                              <path v-if="agent.type === 'executor'" d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z"></path>
                              <path v-if="agent.type === 'retriever'" d="M7 3a1 1 0 000 2h6a1 1 0 100-2H7zM4 7a1 1 0 011-1h10a1 1 0 110 2H5a1 1 0 01-1-1zM2 11a2 2 0 012-2h12a2 2 0 012 2v4a2 2 0 01-2 2H4a2 2 0 01-2-2v-4z"></path>
                            </svg>
                          </div>
                          <div class="ml-4">
                            <div class="text-sm font-medium text-gray-900 dark:text-white">
                              {{ agent.name }}
                            </div>
                            <div class="text-sm text-gray-500 dark:text-gray-400">
                              {{ agent.description }}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap">
                        <div class="text-sm text-gray-900 dark:text-white">{{ agentTypeText(agent.type) }}</div>
                        <div class="text-sm text-gray-500 dark:text-gray-400">{{ agent.model || '-' }}</div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap">
                        <span 
                          :class="statusBadgeClass(agent.status)" 
                          class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                        >
                          {{ statusText(agent.status) }}
                        </span>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                        {{ agent.workflow_count }}
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                        {{ agent.last_active ? formatTime(agent.last_active) : '从未活动' }}
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <router-link 
                          :to="`/agents/${agent.id}`" 
                          class="text-primary-600 hover:text-primary-900 dark:text-primary-400 dark:hover:text-primary-300 mr-4"
                        >
                          查看
                        </router-link>
                        <router-link 
                          :to="`/agents/${agent.id}/edit`" 
                          class="text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300 mr-4"
                        >
                          编辑
                        </router-link>
                        <button 
                          @click="deleteAgent(agent)" 
                          class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
                        >
                          删除
                        </button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="mt-6 flex justify-between items-center">
        <div class="text-sm text-gray-700 dark:text-gray-300">
          显示第 <span class="font-medium">{{ startItem }}</span> 至 <span class="font-medium">{{ endItem }}</span> 项，共 <span class="font-medium">{{ totalItems }}</span> 项
        </div>
        <div>
          <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
            <button
              @click="goToPage(currentPage - 1)"
              :disabled="currentPage === 1"
              :class="currentPage === 1 ? 'opacity-50 cursor-not-allowed' : ''"
              class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
            >
              <span class="sr-only">上一页</span>
              <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
            </button>
            
            <template v-for="page in displayedPages" :key="page">
              <button
                v-if="page !== '...'"
                @click="goToPage(page)"
                :class="page === currentPage ? 'z-10 bg-primary-50 dark:bg-primary-900 border-primary-500 text-primary-600 dark:text-primary-200' : 'bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'"
                class="relative inline-flex items-center px-4 py-2 border text-sm font-medium"
              >
                {{ page }}
              </button>
              <span
                v-else
                class="relative inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                ...
              </span>
            </template>
            
            <button
              @click="goToPage(currentPage + 1)"
              :disabled="currentPage === totalPages"
              :class="currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : ''"
              class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm font-medium text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
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
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import { useToast } from 'vue-toastification';

const toast = useToast();
const loading = ref(false);
const agents = ref([]);
const searchQuery = ref('');
const statusFilter = ref('');
const typeFilter = ref('');
const totalItems = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);

// 模拟数据
const mockAgents = [
  {
    id: 1,
    name: 'GPT 助手',
    description: 'OpenAI GPT-4 模型助手',
    type: 'llm',
    model: 'gpt-4',
    status: 'online',
    workflow_count: 5,
    last_active: new Date(Date.now() - 15 * 60 * 1000)
  },
  {
    id: 2,
    name: '脚本执行器',
    description: '执行 Python 和 Shell 脚本的代理',
    type: 'executor',
    model: null,
    status: 'online',
    workflow_count: 3,
    last_active: new Date(Date.now() - 2 * 60 * 60 * 1000)
  },
  {
    id: 3,
    name: '数据检索器',
    description: '从数据库和文档中检索信息',
    type: 'retriever',
    model: null,
    status: 'offline',
    workflow_count: 2,
    last_active: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000)
  },
  {
    id: 4,
    name: 'Claude 助手',
    description: 'Anthropic Claude 模型助手',
    type: 'llm',
    model: 'claude-3-opus',
    status: 'busy',
    workflow_count: 8,
    last_active: new Date(Date.now() - 5 * 60 * 1000)
  },
  {
    id: 5,
    name: '文件处理器',
    description: '处理和转换各种文件格式',
    type: 'executor',
    model: null,
    status: 'online',
    workflow_count: 4,
    last_active: new Date(Date.now() - 45 * 60 * 1000)
  }
];

// 计算属性
const totalPages = computed(() => Math.ceil(totalItems.value / pageSize.value));
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
const fetchAgents = async () => {
  loading.value = true;
  try {
    // TODO: 替换为实际API调用
    // const response = await agentsApi.getAgents({
    //   page: currentPage.value,
    //   pageSize: pageSize.value,
    //   search: searchQuery.value,
    //   status: statusFilter.value,
    //   type: typeFilter.value
    // });
    
    // 模拟API响应
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // 过滤模拟数据
    let filteredAgents = [...mockAgents];
    
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      filteredAgents = filteredAgents.filter(a => 
        a.name.toLowerCase().includes(query) || 
        a.description.toLowerCase().includes(query)
      );
    }
    
    if (statusFilter.value) {
      filteredAgents = filteredAgents.filter(a => a.status === statusFilter.value);
    }
    
    if (typeFilter.value) {
      filteredAgents = filteredAgents.filter(a => a.type === typeFilter.value);
    }
    
    totalItems.value = filteredAgents.length;
    
    // 应用分页
    const start = (currentPage.value - 1) * pageSize.value;
    const end = start + pageSize.value;
    agents.value = filteredAgents.slice(start, end);
    
  } catch (error) {
    console.error('获取代理失败', error);
    toast.error('获取代理列表失败');
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
      fetchAgents();
    }, 300);
  };
})();

const formatTime = (date) => {
  return format(new Date(date), 'MM月dd日 HH:mm', { locale: zhCN });
};

const goToPage = (page) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  fetchAgents();
};

const statusText = (status) => {
  const statusMap = {
    online: '在线',
    offline: '离线',
    busy: '忙碌'
  };
  return statusMap[status] || status;
};

const statusBadgeClass = (status) => {
  const badgeMap = {
    online: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
    offline: 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200',
    busy: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
  };
  return badgeMap[status] || 'bg-gray-100 text-gray-800';
};

const agentTypeText = (type) => {
  const typeMap = {
    llm: '语言模型',
    executor: '执行器',
    retriever: '检索器'
  };
  return typeMap[type] || type;
};

const agentIconClass = (type) => {
  const colorMap = {
    llm: 'bg-blue-500',
    executor: 'bg-purple-500',
    retriever: 'bg-green-500'
  };
  return colorMap[type] || 'bg-gray-500';
};

const deleteAgent = async (agent) => {
  if (!confirm(`确定要删除代理 "${agent.name}" 吗？`)) {
    return;
  }
  
  try {
    // TODO: 替换为实际API调用
    // await agentsApi.deleteAgent(agent.id);
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 300));
    
    // 从本地列表中移除
    agents.value = agents.value.filter(a => a.id !== agent.id);
    totalItems.value--;
    
    toast.success(`代理 "${agent.name}" 已删除`);
    
    // 如果当前页面为空且不是第一页，则跳转到上一页
    if (agents.value.length === 0 && currentPage.value > 1) {
      goToPage(currentPage.value - 1);
    } else {
      // 重新获取当前页数据
      fetchAgents();
    }
  } catch (error) {
    console.error('删除代理失败', error);
    toast.error('删除代理失败');
  }
};

onMounted(() => {
  fetchAgents();
});
</script> 