<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">仪表盘</h1>
    </div>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- 欢迎卡片 -->
      <div class="mt-6 bg-white dark:bg-gray-800 shadow-sm rounded-lg overflow-hidden">
        <div class="px-4 py-5 sm:p-6">
          <h3 class="text-lg leading-6 font-medium text-gray-900 dark:text-white">
            欢迎使用 JiLang Agent
          </h3>
          <div class="mt-2 max-w-xl text-sm text-gray-500 dark:text-gray-400">
            <p>智能工作流自动化平台 - 让您的团队专注于更有价值的工作</p>
          </div>
          <div class="mt-5">
            <button
              type="button"
              class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              @click="$router.push('/workflows/create')"
            >
              创建新工作流
            </button>
          </div>
        </div>
      </div>

      <!-- 摘要卡片 -->
      <div class="mt-8 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
        <!-- 工作流 -->
        <div class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-primary-500 rounded-md p-3">
                <svg class="h-6 w-6 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    工作流总数
                  </dt>
                  <dd>
                    <div class="text-lg font-medium text-gray-900 dark:text-white">
                      {{ stats.workflows || 0 }}
                    </div>
                  </dd>
                </dl>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 px-5 py-3">
            <div class="text-sm">
              <router-link to="/workflows" class="font-medium text-primary-600 dark:text-primary-400 hover:text-primary-500 dark:hover:text-primary-300">
                查看所有工作流
              </router-link>
            </div>
          </div>
        </div>

        <!-- 执行 -->
        <div class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-green-500 rounded-md p-3">
                <svg class="h-6 w-6 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    成功执行次数
                  </dt>
                  <dd>
                    <div class="text-lg font-medium text-gray-900 dark:text-white">
                      {{ stats.successfulExecutions || 0 }}
                    </div>
                  </dd>
                </dl>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 px-5 py-3">
            <div class="text-sm">
              <router-link to="/executions" class="font-medium text-primary-600 dark:text-primary-400 hover:text-primary-500 dark:hover:text-primary-300">
                查看执行历史
              </router-link>
            </div>
          </div>
        </div>

        <!-- 代理 -->
        <div class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-purple-500 rounded-md p-3">
                <svg class="h-6 w-6 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    活跃代理数
                  </dt>
                  <dd>
                    <div class="text-lg font-medium text-gray-900 dark:text-white">
                      {{ stats.activeAgents || 0 }}
                    </div>
                  </dd>
                </dl>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 px-5 py-3">
            <div class="text-sm">
              <router-link to="/agents" class="font-medium text-primary-600 dark:text-primary-400 hover:text-primary-500 dark:hover:text-primary-300">
                查看所有代理
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- 近期活动 -->
      <div class="mt-8">
        <h2 class="text-lg leading-6 font-medium text-gray-900 dark:text-white">近期活动</h2>
        <div class="mt-2 bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
          <ul role="list" class="divide-y divide-gray-200 dark:divide-gray-700">
            <li v-if="activities.length === 0" class="px-4 py-4 sm:px-6">
              <div class="text-sm text-gray-500 dark:text-gray-400">
                暂无活动记录
              </div>
            </li>
            <li v-for="activity in activities" :key="activity.id" class="px-4 py-4 sm:px-6">
              <div class="flex items-center space-x-4">
                <div class="flex-shrink-0">
                  <svg v-if="activity.type === 'workflow_completed'" class="h-8 w-8 text-green-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <svg v-else-if="activity.type === 'workflow_failed'" class="h-8 w-8 text-red-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                  <svg v-else class="h-8 w-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ activity.title }}
                  </p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">
                    {{ activity.description }}
                  </p>
                </div>
                <div class="text-sm text-gray-500 dark:text-gray-400">
                  {{ formatTime(activity.timestamp) }}
                </div>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

// 统计数据
const stats = ref({
  workflows: 3,
  successfulExecutions: 27,
  activeAgents: 5
});

// 活动记录
const activities = ref([
  {
    id: 1,
    type: 'workflow_completed',
    title: '工作流 "每日数据同步" 执行成功',
    description: '所有任务都已成功完成，耗时 2 分钟 30 秒',
    timestamp: new Date(Date.now() - 30 * 60 * 1000)
  },
  {
    id: 2,
    type: 'workflow_failed',
    title: '工作流 "数据备份" 执行失败',
    description: '任务 "连接远程服务器" 失败: 连接超时',
    timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000)
  },
  {
    id: 3,
    type: 'agent_created',
    title: '新建代理 "数据分析助手"',
    description: '代理已创建并准备就绪',
    timestamp: new Date(Date.now() - 5 * 60 * 60 * 1000)
  }
]);

// 格式化时间
const formatTime = (date) => {
  return format(new Date(date), 'MM月dd日 HH:mm', { locale: zhCN });
};

// 获取统计数据
const fetchStats = async () => {
  // TODO: 从API获取实际数据
  // const response = await statsApi.getDashboardStats();
  // stats.value = response.data;
};

// 获取活动记录
const fetchActivities = async () => {
  // TODO: 从API获取实际数据
  // const response = await activitiesApi.getRecentActivities();
  // activities.value = response.data;
};

onMounted(() => {
  fetchStats();
  fetchActivities();
});
</script> 