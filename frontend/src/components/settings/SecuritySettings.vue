<template>
  <div class="space-y-6">
    <!-- API密钥管理 -->
    <div class="bg-white rounded-lg border border-gray-200 p-6">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-medium">API密钥管理</h3>
        <button
          @click="showCreateApiKeyModal = true"
          class="px-3 py-1.5 bg-primary-600 text-white rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 text-sm"
        >
          创建新密钥
        </button>
      </div>
      
      <!-- API密钥列表 -->
      <div v-if="apiKeys.length > 0" class="space-y-4">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">名称</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建时间</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">上次使用</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="key in apiKeys" :key="key.id">
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ key.name }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(key.createdAt) }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ key.lastUsed ? formatDate(key.lastUsed) : '从未使用' }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span 
                    class="px-2 py-1 text-xs font-semibold rounded-full"
                    :class="key.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                  >
                    {{ key.active ? '活跃' : '已禁用' }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button 
                    @click="toggleApiKeyStatus(key)" 
                    class="text-primary-600 hover:text-primary-900 mr-3"
                  >
                    {{ key.active ? '禁用' : '启用' }}
                  </button>
                  <button 
                    @click="deleteApiKey(key)" 
                    class="text-red-600 hover:text-red-900"
                  >
                    删除
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <!-- 空状态 -->
      <div v-else class="text-center py-8">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">没有API密钥</h3>
        <p class="mt-1 text-sm text-gray-500">创建API密钥以便通过API访问系统</p>
        <div class="mt-6">
          <button
            @click="showCreateApiKeyModal = true"
            class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            <svg class="-ml-1 mr-2 h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
            </svg>
            创建新密钥
          </button>
        </div>
      </div>
    </div>
    
    <!-- 访问令牌管理 -->
    <div class="bg-white rounded-lg border border-gray-200 p-6">
      <h3 class="text-lg font-medium mb-4">访问令牌管理</h3>
      
      <!-- 当前会话 -->
      <div class="mb-6">
        <h4 class="text-base font-medium mb-2">当前会话</h4>
        <div class="bg-gray-50 p-4 rounded-lg">
          <div class="flex justify-between items-center">
            <div>
              <p class="text-sm font-medium text-gray-900">{{ currentSession.deviceName }}</p>
              <p class="text-xs text-gray-500">登录时间: {{ formatDate(currentSession.loginTime) }}</p>
              <p class="text-xs text-gray-500">IP地址: {{ currentSession.ipAddress }}</p>
            </div>
            <span class="px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
              当前会话
            </span>
          </div>
        </div>
      </div>
      
      <!-- 其他活跃会话 -->
      <div>
        <div class="flex justify-between items-center mb-2">
          <h4 class="text-base font-medium">其他活跃会话</h4>
          <button
            v-if="activeSessions.length > 0"
            @click="revokeAllSessions"
            class="text-sm text-red-600 hover:text-red-900"
          >
            注销所有会话
          </button>
        </div>
        
        <div v-if="activeSessions.length > 0" class="space-y-3">
          <div 
            v-for="session in activeSessions" 
            :key="session.id" 
            class="bg-gray-50 p-4 rounded-lg"
          >
            <div class="flex justify-between items-center">
              <div>
                <p class="text-sm font-medium text-gray-900">{{ session.deviceName }}</p>
                <p class="text-xs text-gray-500">登录时间: {{ formatDate(session.loginTime) }}</p>
                <p class="text-xs text-gray-500">IP地址: {{ session.ipAddress }}</p>
              </div>
              <button
                @click="revokeSession(session)"
                class="text-sm text-red-600 hover:text-red-900"
              >
                注销
              </button>
            </div>
          </div>
        </div>
        
        <!-- 空状态 -->
        <div v-else class="text-center py-6 bg-gray-50 rounded-lg">
          <p class="text-sm text-gray-500">没有其他活跃会话</p>
        </div>
      </div>
    </div>
    
    <!-- 安全日志 -->
    <div class="bg-white rounded-lg border border-gray-200 p-6">
      <h3 class="text-lg font-medium mb-4">安全日志</h3>
      
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">事件</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">时间</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">IP地址</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">设备</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="log in securityLogs" :key="log.id">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <span 
                    class="mr-2 flex-shrink-0 h-2 w-2 rounded-full"
                    :class="{
                      'bg-green-400': log.type === 'success',
                      'bg-red-400': log.type === 'error',
                      'bg-yellow-400': log.type === 'warning',
                      'bg-blue-400': log.type === 'info'
                    }"
                  ></span>
                  <span class="text-sm text-gray-900">{{ log.event }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(log.timestamp) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ log.ipAddress }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ log.device }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- 分页 -->
      <div class="flex justify-between items-center mt-4">
        <div class="text-sm text-gray-700">
          显示 <span class="font-medium">1</span> 到 <span class="font-medium">10</span> 条，共 <span class="font-medium">{{ securityLogs.length }}</span> 条记录
        </div>
        <div>
          <button class="px-3 py-1 border border-gray-300 rounded-md mr-2 text-sm">上一页</button>
          <button class="px-3 py-1 border border-gray-300 rounded-md text-sm">下一页</button>
        </div>
      </div>
    </div>
    
    <!-- 创建API密钥对话框 -->
    <div v-if="showCreateApiKeyModal" class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
        <h3 class="text-lg font-medium mb-4">创建新API密钥</h3>
        
        <form @submit.prevent="createApiKey" class="space-y-4">
          <div>
            <label for="keyName" class="block text-sm font-medium text-gray-700 mb-1">密钥名称</label>
            <input
              id="keyName"
              v-model="newApiKey.name"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
              placeholder="例如：开发环境、测试服务器"
              required
            />
          </div>
          
          <div>
            <label for="keyExpiration" class="block text-sm font-medium text-gray-700 mb-1">过期时间</label>
            <select
              id="keyExpiration"
              v-model="newApiKey.expiration"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="never">永不过期</option>
              <option value="30days">30天</option>
              <option value="90days">90天</option>
              <option value="1year">1年</option>
            </select>
          </div>
          
          <div class="pt-4 flex justify-end space-x-3">
            <button
              type="button"
              @click="showCreateApiKeyModal = false"
              class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              取消
            </button>
            <button
              type="submit"
              class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              :disabled="isCreatingKey"
            >
              <span v-if="isCreatingKey" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                创建中...
              </span>
              <span v-else>创建</span>
            </button>
          </div>
        </form>
      </div>
    </div>
    
    <!-- 显示新创建的API密钥对话框 -->
    <div v-if="showApiKeyResult" class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
        <h3 class="text-lg font-medium mb-2">API密钥已创建</h3>
        <p class="text-sm text-gray-500 mb-4">请保存此密钥，它只会显示一次</p>
        
        <div class="bg-gray-50 p-3 rounded-md mb-4">
          <div class="flex justify-between items-center">
            <code class="text-sm break-all">{{ createdApiKey }}</code>
            <button 
              @click="copyApiKey" 
              class="ml-2 text-primary-600 hover:text-primary-900"
              title="复制到剪贴板"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path d="M8 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" />
                <path d="M6 3a2 2 0 00-2 2v11a2 2 0 002 2h8a2 2 0 002-2V5a2 2 0 00-2-2H6zm0 2h8v11H6V5z" />
              </svg>
            </button>
          </div>
        </div>
        
        <div class="flex justify-end">
          <button
            @click="closeApiKeyResult"
            class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            我已保存密钥
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useToast } from 'vue-toastification';

// Toast通知
const toast = useToast();

// API密钥列表
const apiKeys = ref([]);
// 活跃会话列表
const activeSessions = ref([]);
// 安全日志
const securityLogs = ref([]);
// 当前会话
const currentSession = reactive({
  deviceName: '未知设备',
  loginTime: new Date(),
  ipAddress: '0.0.0.0'
});

// 创建API密钥相关
const showCreateApiKeyModal = ref(false);
const isCreatingKey = ref(false);
const newApiKey = reactive({
  name: '',
  expiration: 'never'
});

// 显示创建的API密钥
const showApiKeyResult = ref(false);
const createdApiKey = ref('');

// 格式化日期
const formatDate = (date) => {
  if (!date) return '';
  
  const d = new Date(date);
  return d.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
};

// 加载数据
onMounted(() => {
  // 加载API密钥
  loadApiKeys();
  
  // 加载会话信息
  loadSessionInfo();
  
  // 加载安全日志
  loadSecurityLogs();
});

// 加载API密钥
const loadApiKeys = async () => {
  // 模拟API调用
  setTimeout(() => {
    apiKeys.value = [
      {
        id: 1,
        name: '开发环境',
        createdAt: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000), // 30天前
        lastUsed: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2天前
        active: true
      },
      {
        id: 2,
        name: '测试服务器',
        createdAt: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000), // 15天前
        lastUsed: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000), // 1天前
        active: true
      },
      {
        id: 3,
        name: '旧密钥',
        createdAt: new Date(Date.now() - 60 * 24 * 60 * 60 * 1000), // 60天前
        lastUsed: null,
        active: false
      }
    ];
  }, 500);
};

// 加载会话信息
const loadSessionInfo = async () => {
  // 模拟API调用
  setTimeout(() => {
    // 当前会话
    currentSession.deviceName = 'Chrome on macOS';
    currentSession.loginTime = new Date(Date.now() - 3 * 60 * 60 * 1000); // 3小时前
    currentSession.ipAddress = '192.168.1.100';
    
    // 其他活跃会话
    activeSessions.value = [
      {
        id: 1,
        deviceName: 'Firefox on Windows',
        loginTime: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2天前
        ipAddress: '192.168.1.101'
      },
      {
        id: 2,
        deviceName: 'Safari on iPhone',
        loginTime: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000), // 1天前
        ipAddress: '192.168.1.102'
      }
    ];
  }, 500);
};

// 加载安全日志
const loadSecurityLogs = async () => {
  // 模拟API调用
  setTimeout(() => {
    securityLogs.value = [
      {
        id: 1,
        event: '登录成功',
        timestamp: new Date(Date.now() - 3 * 60 * 60 * 1000), // 3小时前
        ipAddress: '192.168.1.100',
        device: 'Chrome on macOS',
        type: 'success'
      },
      {
        id: 2,
        event: '密码修改成功',
        timestamp: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2天前
        ipAddress: '192.168.1.100',
        device: 'Chrome on macOS',
        type: 'success'
      },
      {
        id: 3,
        event: '登录失败 - 密码错误',
        timestamp: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000), // 3天前
        ipAddress: '192.168.1.103',
        device: 'Unknown Device',
        type: 'error'
      },
      {
        id: 4,
        event: '创建API密钥',
        timestamp: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000), // 15天前
        ipAddress: '192.168.1.100',
        device: 'Chrome on macOS',
        type: 'info'
      },
      {
        id: 5,
        event: '异常位置登录',
        timestamp: new Date(Date.now() - 20 * 24 * 60 * 60 * 1000), // 20天前
        ipAddress: '203.0.113.1',
        device: 'Chrome on Windows',
        type: 'warning'
      }
    ];
  }, 500);
};

// 创建API密钥
const createApiKey = async () => {
  isCreatingKey.value = true;
  
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // 生成随机密钥
    const randomKey = Array(40).fill(0).map(() => Math.random().toString(36).charAt(2)).join('');
    createdApiKey.value = randomKey;
    
    // 添加到列表
    const newKey = {
      id: apiKeys.value.length + 1,
      name: newApiKey.name,
      createdAt: new Date(),
      lastUsed: null,
      active: true
    };
    
    apiKeys.value.unshift(newKey);
    
    // 显示结果
    showCreateApiKeyModal.value = false;
    showApiKeyResult.value = true;
    
    // 重置表单
    newApiKey.name = '';
    newApiKey.expiration = 'never';
  } catch (error) {
    toast.error('创建API密钥失败: ' + error.message);
  } finally {
    isCreatingKey.value = false;
  }
};

// 复制API密钥
const copyApiKey = () => {
  navigator.clipboard.writeText(createdApiKey.value)
    .then(() => {
      toast.success('API密钥已复制到剪贴板');
    })
    .catch((error) => {
      toast.error('复制失败: ' + error.message);
    });
};

// 关闭API密钥结果
const closeApiKeyResult = () => {
  showApiKeyResult.value = false;
  createdApiKey.value = '';
};

// 切换API密钥状态
const toggleApiKeyStatus = async (key) => {
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // 更新状态
    key.active = !key.active;
    
    toast.success(`API密钥 "${key.name}" 已${key.active ? '启用' : '禁用'}`);
  } catch (error) {
    toast.error('操作失败: ' + error.message);
  }
};

// 删除API密钥
const deleteApiKey = async (key) => {
  if (!confirm(`确定要删除API密钥 "${key.name}" 吗？此操作不可撤销。`)) {
    return;
  }
  
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // 从列表中移除
    apiKeys.value = apiKeys.value.filter(k => k.id !== key.id);
    
    toast.success(`API密钥 "${key.name}" 已删除`);
  } catch (error) {
    toast.error('删除失败: ' + error.message);
  }
};

// 注销会话
const revokeSession = async (session) => {
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // 从列表中移除
    activeSessions.value = activeSessions.value.filter(s => s.id !== session.id);
    
    toast.success(`会话已注销`);
  } catch (error) {
    toast.error('注销失败: ' + error.message);
  }
};

// 注销所有会话
const revokeAllSessions = async () => {
  if (!confirm('确定要注销所有其他会话吗？这将使所有其他设备上的登录失效。')) {
    return;
  }
  
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // 清空列表
    activeSessions.value = [];
    
    toast.success('所有其他会话已注销');
  } catch (error) {
    toast.error('注销失败: ' + error.message);
  }
};
</script> 