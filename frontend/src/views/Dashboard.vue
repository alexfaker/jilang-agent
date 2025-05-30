<template>
  <div class="p-6">
    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex items-center justify-center py-12">
      <div class="flex items-center space-x-3">
        <ArrowPathIcon class="w-6 h-6 animate-spin text-indigo-600" />
        <span class="text-lg text-gray-600">加载中...</span>
      </div>
    </div>
    
    <!-- 主要内容 -->
    <div v-else>
      <!-- 欢迎信息 -->
      <div class="mb-8 flex items-center justify-between">
        <div>
          <h2 class="text-2xl font-bold text-gray-800">欢迎回来，{{ userName }}</h2>
          <p class="text-gray-600">这是您的AI工作流平台使用概览</p>
        </div>
        <button 
          @click="refreshData"
          :disabled="isRefreshing"
          class="flex items-center space-x-2 px-4 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-gray-50 rounded-lg transition-colors"
        >
          <ArrowPathIcon :class="['w-4 h-4', isRefreshing ? 'animate-spin' : '']" />
          <span>{{ isRefreshing ? '刷新中...' : '刷新数据' }}</span>
        </button>
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
              <p class="text-sm text-gray-500">总执行次数</p>
            </div>
            <span class="text-green-500 text-sm font-medium flex items-center">
              <ArrowTrendingUpIcon class="w-4 h-4 mr-1" />
              12.4%
            </span>
          </div>
        </div>
        
        <!-- 活跃工作流 -->
        <div class="card bg-white p-6 rounded-2xl shadow-lg">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-full bg-amber-100 flex items-center justify-center">
              <CpuChipIcon class="text-amber-600 w-6 h-6" />
            </div>
            <span class="text-sm font-medium text-gray-500">活跃</span>
          </div>
          <div class="flex items-end justify-between">
            <div>
              <p class="text-2xl font-bold text-gray-800">{{ stats.activeWorkflows }}</p>
              <p class="text-sm text-gray-500">活跃工作流</p>
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
                activity.type === 'create' ? 'bg-purple-100' : 
                activity.type === 'failed' ? 'bg-red-100' : 'bg-amber-100'
              ]">
                <PlayIcon v-if="activity.type === 'execute'" class="w-5 h-5 text-blue-600" />
                <CheckIcon v-else-if="activity.type === 'complete'" class="w-5 h-5 text-green-600" />
                <PlusIcon v-else-if="activity.type === 'create'" class="w-5 h-5 text-purple-600" />
                <ExclamationTriangleIcon v-else-if="activity.type === 'failed'" class="w-5 h-5 text-red-600" />
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
            
            <!-- 空状态 -->
            <div v-if="recentActivities.length === 0" class="text-center py-8 text-gray-500">
              <CogIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
              <p>暂无最近活动</p>
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
            
            <!-- 空状态 -->
            <div v-if="frequentWorkflows.length === 0" class="text-center py-8 text-gray-500">
              <RectangleGroupIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
              <p>暂无工作流</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useUserStore } from '../stores/user';
import { statsApi, workflowApi, executionApi } from '../api/index';
import notify from '../utils/notification';
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
  ArrowPathIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline';

const authStore = useUserStore();

// 加载状态管理
const isLoading = ref(true);
const isRefreshing = ref(false);

// 用户信息
const userName = ref(authStore.user?.username || '用户');

// 统计数据
const stats = ref({
  workflows: 0,
  executions: 0,
  activeWorkflows: 0,
  timeSaved: 0
});

// 最近活动
const recentActivities = ref([]);

// 常用工作流
const frequentWorkflows = ref([]);

// 正在执行的工作流
const executingWorkflows = ref([]);

// 加载仪表盘统计数据
const loadDashboardStats = async () => {
  try {
    const response = await statsApi.getDashboardStats();
    if (response.status === 'success') {
      const data = response.data;
      
      // 更新基础统计
      stats.value.workflows = data.total_workflows || 0;
      stats.value.executions = data.total_executions || 0;
      
      // 计算活跃工作流（有最近执行记录的工作流）
      const recentExecutions = data.recent_executions || [];
      const uniqueWorkflowIds = new Set(recentExecutions.map(e => e.workflow_id));
      stats.value.activeWorkflows = uniqueWorkflowIds.size;
      
      // 估算节省时间（假设每次执行平均节省30分钟）
      stats.value.timeSaved = Math.round((data.total_executions || 0) * 0.5 * 10) / 10;
      
      return data;
    }
  } catch (error) {
    console.error('获取仪表盘统计失败:', error);
    notify.error('获取统计数据失败');
    return null;
  }
};

// 加载最近活动
const loadRecentActivities = async (dashboardData) => {
  try {
    const recentExecutions = dashboardData?.recent_executions || [];
    
    // 转换执行记录为活动格式
    const activities = recentExecutions.map(execution => ({
      id: execution.id,
      type: execution.status === 'succeeded' ? 'complete' : 
            execution.status === 'failed' ? 'failed' : 'execute',
      title: execution.status === 'succeeded' ? '完成了工作流' :
             execution.status === 'failed' ? '工作流执行失败' : '执行了工作流',
      workflowName: execution.workflow?.name || '未知工作流',
      workflowId: execution.workflow_id,
      timestamp: new Date(execution.created_at)
    }));
    
    recentActivities.value = activities;
  } catch (error) {
    console.error('处理最近活动失败:', error);
    recentActivities.value = [];
  }
};

// 加载常用工作流
const loadFrequentWorkflows = async () => {
  try {
    // 获取用户的工作流列表
    const workflowResponse = await workflowApi.getWorkflows({ limit: 50 });
    if (workflowResponse.status !== 'success') {
      throw new Error('获取工作流列表失败');
    }
    
    const workflows = workflowResponse.data.workflows || [];
    
    // 获取执行记录来计算使用频率
    const executionResponse = await executionApi.getExecutions({ page: 1, page_size: 100 });
    const executions = executionResponse.data?.executions || [];
    
    // 计算每个工作流的执行次数和最后执行时间
    const workflowStats = workflows.map(workflow => {
      const workflowExecutions = executions.filter(e => e.workflow_id === workflow.id);
      const lastExecution = workflowExecutions.length > 0 ? 
        new Date(Math.max(...workflowExecutions.map(e => new Date(e.created_at)))) : 
        new Date(workflow.created_at);
      
      return {
        id: workflow.id,
        name: workflow.name,
        category: getCategoryFromWorkflow(workflow),
        lastExecuted: lastExecution,
        executionCount: workflowExecutions.length
      };
    });
    
    // 按执行次数和最近执行时间排序，取前4个
    const sortedWorkflows = workflowStats
      .filter(w => w.executionCount > 0) // 只显示有执行记录的
      .sort((a, b) => {
        // 先按执行次数排序，然后按最近执行时间
        if (b.executionCount !== a.executionCount) {
          return b.executionCount - a.executionCount;
        }
        return new Date(b.lastExecuted) - new Date(a.lastExecuted);
      })
      .slice(0, 4);
    
    // 如果常用工作流不足4个，用最新创建的补充
    if (sortedWorkflows.length < 4) {
      const remainingWorkflows = workflowStats
        .filter(w => !sortedWorkflows.find(s => s.id === w.id))
        .sort((a, b) => new Date(b.lastExecuted) - new Date(a.lastExecuted))
        .slice(0, 4 - sortedWorkflows.length);
      
      sortedWorkflows.push(...remainingWorkflows);
    }
    
    frequentWorkflows.value = sortedWorkflows;
  } catch (error) {
    console.error('获取常用工作流失败:', error);
    frequentWorkflows.value = [];
  }
};

// 根据工作流信息推断分类
const getCategoryFromWorkflow = (workflow) => {
  const name = (workflow.name || '').toLowerCase();
  const description = (workflow.description || '').toLowerCase();
  const text = name + ' ' + description;
  
  if (text.includes('数据') || text.includes('处理') || text.includes('分析')) {
    return 'data';
  } else if (text.includes('客服') || text.includes('支持') || text.includes('回复')) {
    return 'support';
  } else if (text.includes('内容') || text.includes('文本') || text.includes('审核')) {
    return 'content';
  } else {
    return 'analysis';
  }
};

// 主要数据加载函数
const loadDashboardData = async () => {
  isLoading.value = true;
  try {
    // 并行加载统计数据
    const dashboardData = await loadDashboardStats();
    
    // 加载其他数据
    await Promise.all([
      loadRecentActivities(dashboardData),
      loadFrequentWorkflows()
    ]);
    
  } catch (error) {
    console.error('加载仪表盘数据失败:', error);
    notify.error('加载数据失败，请刷新页面重试');
  } finally {
    isLoading.value = false;
  }
};

// 刷新数据
const refreshData = async () => {
  if (isRefreshing.value) return;
  
  isRefreshing.value = true;
  try {
    await loadDashboardData();
    notify.success('数据已刷新');
  } finally {
    isRefreshing.value = false;
  }
};

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
    // 调用真实的工作流执行API
    const response = await workflowApi.executeWorkflow(workflowId, {});
    
    if (response.status === 'success') {
      notify.success('工作流已开始执行');
      
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
    } else {
      throw new Error(response.message || '执行失败');
    }
  } catch (error) {
    console.error('执行工作流失败:', error);
    notify.error('执行工作流失败: ' + (error.message || '未知错误'));
  } finally {
    executingWorkflows.value = executingWorkflows.value.filter(id => id !== workflowId);
  }
};

onMounted(() => {
  // 加载仪表盘数据
  loadDashboardData();
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