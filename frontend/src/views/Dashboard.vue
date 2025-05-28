<template>
  <div class="flex-1 overflow-y-auto p-6">
    <!-- 欢迎信息 -->
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-gray-800">欢迎回来，{{ userName }}</h2>
      <p class="text-gray-600">这是您的AI工作流平台使用概览</p>
    </div>
    
    <!-- 概览卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <!-- 工作流数量 -->
      <div class="card bg-white p-6 rounded-2xl shadow-lg">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
            <RectangleGroupIcon class="text-blue-600 w-6 h-6" />
          </div>
          <span class="text-sm font-medium text-gray-500">工作流</span>
        </div>
        <div class="flex items-end justify-between">
          <div>
            <p class="text-2xl font-bold text-gray-800">{{ stats.workflows }}</p>
            <p class="text-sm text-gray-500">创建的工作流</p>
          </div>
          <span class="text-green-500 text-sm font-medium flex items-center">
            <ArrowTrendingUpIcon class="w-4 h-4 mr-1" />
            8.2%
          </span>
        </div>
      </div>
      
      <!-- 执行次数 -->
      <div class="card bg-white p-6 rounded-2xl shadow-lg">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 rounded-full bg-purple-100 flex items-center justify-center">
            <PlayCircleIcon class="text-purple-600 w-6 h-6" />
          </div>
          <span class="text-sm font-medium text-gray-500">执行</span>
        </div>
        <div class="flex items-end justify-between">
          <div>
            <p class="text-2xl font-bold text-gray-800">{{ stats.executions }}</p>
            <p class="text-sm text-gray-500">本月执行次数</p>
          </div>
          <span class="text-green-500 text-sm font-medium flex items-center">
            <ArrowTrendingUpIcon class="w-4 h-4 mr-1" />
            12.4%
          </span>
        </div>
      </div>
      
      <!-- 使用的代理 -->
      <div class="card bg-white p-6 rounded-2xl shadow-lg">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 rounded-full bg-amber-100 flex items-center justify-center">
            <CpuChipIcon class="text-amber-600 w-6 h-6" />
          </div>
          <span class="text-sm font-medium text-gray-500">工作流</span>
        </div>
        <div class="flex items-end justify-between">
          <div>
            <p class="text-2xl font-bold text-gray-800">{{ stats.activeWorkflows }}</p>
            <p class="text-sm text-gray-500">使用中的工作流</p>
          </div>
          <span class="text-green-500 text-sm font-medium flex items-center">
            <ArrowTrendingUpIcon class="w-4 h-4 mr-1" />
            4.1%
          </span>
        </div>
      </div>
      
      <!-- 节省时间 -->
      <div class="card bg-white p-6 rounded-2xl shadow-lg">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
            <ClockIcon class="text-green-600 w-6 h-6" />
          </div>
          <span class="text-sm font-medium text-gray-500">节省</span>
        </div>
        <div class="flex items-end justify-between">
          <div>
            <p class="text-2xl font-bold text-gray-800">{{ stats.timeSaved }}</p>
            <p class="text-sm text-gray-500">节省小时（估计）</p>
          </div>
          <span class="text-green-500 text-sm font-medium flex items-center">
            <ArrowTrendingUpIcon class="w-4 h-4 mr-1" />
            15.3%
          </span>
        </div>
      </div>
    </div>

    <!-- 活动列表和最近工作流 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 最近活动 -->
      <div class="card bg-white p-6 rounded-2xl shadow-lg">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-800">最近活动</h3>
          <router-link to="/executions" class="text-sm text-indigo-600 hover:text-indigo-800">查看全部</router-link>
        </div>
        <div class="space-y-4">
          <!-- 活动项 -->
          <div v-for="activity in recentActivities" :key="activity.id" class="flex items-start">
            <div :class="[
              'w-10 h-10 rounded-full flex items-center justify-center mr-4 flex-shrink-0',
              activity.type === 'execute' ? 'bg-blue-100' : 
              activity.type === 'complete' ? 'bg-green-100' : 
              activity.type === 'create' ? 'bg-purple-100' : 'bg-amber-100'
            ]">
              <PlayIcon v-if="activity.type === 'execute'" class="w-5 h-5 text-blue-600" />
              <CheckIcon v-else-if="activity.type === 'complete'" class="w-5 h-5 text-green-600" />
              <PlusIcon v-else-if="activity.type === 'create'" class="w-5 h-5 text-purple-600" />
              <CogIcon v-else class="w-5 h-5 text-amber-600" />
            </div>
            <div class="flex-1">
              <p class="text-sm font-medium text-gray-800">
                {{ activity.title }} 
                <router-link :to="`/workflows/${activity.workflowId}`" class="text-indigo-600 hover:text-indigo-800">{{ activity.workflowName }}</router-link>
              </p>
              <p class="text-xs text-gray-500">{{ formatTimeAgo(activity.timestamp) }}</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 常用工作流 -->
      <div class="card bg-white p-6 rounded-2xl shadow-lg">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-800">常用工作流</h3>
          <router-link to="/workflows" class="text-sm text-indigo-600 hover:text-indigo-800">查看全部</router-link>
        </div>
        <div class="space-y-4">
          <!-- 工作流项 -->
          <div v-for="workflow in frequentWorkflows" :key="workflow.id" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
            <div class="flex items-center">
              <div :class="[
                'w-10 h-10 rounded-full flex items-center justify-center mr-4',
                workflow.category === 'data' ? 'bg-blue-100' :
                workflow.category === 'support' ? 'bg-purple-100' :
                workflow.category === 'content' ? 'bg-green-100' : 'bg-amber-100'
              ]">
                <CpuChipIcon v-if="workflow.category === 'data'" class="w-5 h-5 text-blue-600" />
                <ChatBubbleLeftRightIcon v-else-if="workflow.category === 'support'" class="w-5 h-5 text-purple-600" />
                <DocumentTextIcon v-else-if="workflow.category === 'content'" class="w-5 h-5 text-green-600" />
                <MagnifyingGlassIcon v-else class="w-5 h-5 text-amber-600" />
              </div>
              <div>
                <p class="text-sm font-medium text-gray-800">{{ workflow.name }}</p>
                <p class="text-xs text-gray-500">上次执行: {{ formatTimeAgo(workflow.lastExecuted) }}</p>
              </div>
            </div>
            <button 
              @click="executeWorkflow(workflow.id)"
              class="text-gray-600 hover:text-indigo-600 transition-colors"
              :disabled="executingWorkflows.includes(workflow.id)"
            >
              <PlayIcon v-if="!executingWorkflows.includes(workflow.id)" class="w-5 h-5" />
              <ArrowPathIcon v-else class="w-5 h-5 animate-spin" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import {
  RectangleGroupIcon,
  PlayCircleIcon,
  CpuChipIcon,
  ClockIcon,
  ArrowTrendingUpIcon,
  PlayIcon,
  CheckIcon,
  PlusIcon,
  CogIcon,
  ChatBubbleLeftRightIcon,
  DocumentTextIcon,
  MagnifyingGlassIcon,
  ArrowPathIcon
} from '@heroicons/vue/24/outline';

const authStore = useAuthStore();

// 用户信息
const userName = ref(authStore.user?.username || '用户');

// 统计数据
const stats = ref({
  workflows: 12,
  executions: 237,
  activeWorkflows: 8,
  timeSaved: 42.5
});

// 最近活动
const recentActivities = ref([
  {
    id: 1,
    type: 'execute',
    title: '执行了工作流',
    workflowName: '数据处理自动化',
    workflowId: 1,
    timestamp: new Date(Date.now() - 10 * 60 * 1000)
  },
  {
    id: 2,
    type: 'complete',
    title: '完成了工作流',
    workflowName: '客户支持回复',
    workflowId: 2,
    timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000)
  },
  {
    id: 3,
    type: 'create',
    title: '创建了新工作流',
    workflowName: '智能内容审核',
    workflowId: 3,
    timestamp: new Date(Date.now() - 24 * 60 * 60 * 1000)
  },
  {
    id: 4,
    type: 'update',
    title: '更新了工作流',
    workflowName: '数据处理自动化',
    workflowId: 1,
    timestamp: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000)
  }
]);

// 常用工作流
const frequentWorkflows = ref([
  {
    id: 1,
    name: '数据处理自动化',
    category: 'data',
    lastExecuted: new Date(Date.now() - 10 * 60 * 1000)
  },
  {
    id: 2,
    name: '客户支持回复',
    category: 'support',
    lastExecuted: new Date(Date.now() - 2 * 60 * 60 * 1000)
  },
  {
    id: 3,
    name: '智能内容审核',
    category: 'content',
    lastExecuted: new Date(Date.now() - 24 * 60 * 60 * 1000)
  },
  {
    id: 4,
    name: '市场数据分析',
    category: 'analysis',
    lastExecuted: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000)
  }
]);

// 正在执行的工作流
const executingWorkflows = ref([]);

// 时间格式化函数
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
  } else {
    return `${days}天前`;
  }
};

// 执行工作流
const executeWorkflow = async (workflowId) => {
  if (executingWorkflows.value.includes(workflowId)) return;
  
  executingWorkflows.value.push(workflowId);
  
  try {
    // 模拟执行
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    // 更新活动记录
    const workflow = frequentWorkflows.value.find(w => w.id === workflowId);
    if (workflow) {
      recentActivities.value.unshift({
        id: Date.now(),
        type: 'execute',
        title: '执行了工作流',
        workflowName: workflow.name,
        workflowId: workflowId,
        timestamp: new Date()
      });
      workflow.lastExecuted = new Date();
    }
    
    // 更新统计
    stats.value.executions++;
  } catch (error) {
    console.error('执行工作流失败:', error);
  } finally {
    executingWorkflows.value = executingWorkflows.value.filter(id => id !== workflowId);
  }
};

onMounted(() => {
  // 获取最新统计数据
  // loadDashboardStats();
});
</script>

<style scoped>
.card {
  transition: all 0.3s ease;
}

.card:hover {
  transform: translateY(-2px);
}
</style> 