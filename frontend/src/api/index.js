import axios from 'axios';
import { useToast } from 'vue-toastification';

// 创建axios实例
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// 请求拦截器
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  response => {
    return response.data;
  },
  error => {
    const toast = useToast();
    const { response } = error;
    
    if (response) {
      // 服务器返回错误信息
      const { status, data } = response;
      
      switch (status) {
        case 400:
          toast.error(`请求错误: ${data.message || '参数错误'}`);
          break;
        case 401:
          toast.error('未授权，请重新登录');
          // 清除token并跳转到登录页
          localStorage.removeItem('token');
          window.location.href = '/login';
          break;
        case 403:
          toast.error(`拒绝访问: ${data.message || '权限不足'}`);
          break;
        case 404:
          toast.error(`请求的资源不存在: ${data.message || '未找到'}`);
          break;
        case 500:
          toast.error(`服务器错误: ${data.message || '内部服务器错误'}`);
          break;
        default:
          toast.error(`未知错误: ${data.message || error.message || '请求失败'}`);
      }
    } else {
      // 网络错误或请求被取消
      if (error.message.includes('timeout')) {
        toast.error('请求超时，请检查网络连接');
      } else if (error.message.includes('Network Error')) {
        toast.error('网络错误，请检查网络连接');
      } else {
        toast.error(`请求失败: ${error.message}`);
      }
    }
    
    return Promise.reject(error);
  }
);

// 创建一个用于文件上传的API客户端实例
const uploadApiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 60000, // 文件上传可能需要更长的超时时间
  headers: {
    'Accept': 'application/json'
    // 不设置Content-Type，让axios自动设置为multipart/form-data
  }
});

// 为上传客户端添加请求拦截器
uploadApiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 为上传客户端添加响应拦截器
uploadApiClient.interceptors.response.use(
  response => {
    return response.data;
  },
  error => {
    const toast = useToast();
    const { response } = error;
    
    if (response) {
      const { status, data } = response;
      toast.error(`上传失败: ${data.message || '服务器错误'}`);
    } else {
      toast.error(`上传失败: ${error.message || '未知错误'}`);
    }
    
    return Promise.reject(error);
  }
);

// 认证相关API
export const authApi = {
  login: (credentials) => apiClient.post('/auth/login', credentials),
  register: (userData) => apiClient.post('/auth/register', userData),
  logout: () => apiClient.post('/auth/logout'),
  getProfile: () => apiClient.get('/users/me'),
  updateProfile: (data) => apiClient.put('/users/me', data),
  updateProfileWithAvatar: (formData) => uploadApiClient.post('/users/me/avatar', formData),
  changePassword: (data) => apiClient.put('/users/me/password', data),
  forgotPassword: (email) => apiClient.post('/auth/forgot-password', { email }),
  resetPassword: (data) => apiClient.post('/auth/reset-password', data)
};

// 设置相关API
export const settingsApi = {
  getSettings: () => apiClient.get('/settings'),
  updateSettings: (settings) => apiClient.put('/settings', settings),
  getSecuritySettings: () => apiClient.get('/settings/security'),
  generateApiKey: (description) => apiClient.post('/settings/api-keys', { description }),
  revokeApiKey: (keyId) => apiClient.delete(`/settings/api-keys/${keyId}`),
  getApiKeys: () => apiClient.get('/settings/api-keys'),
  getSessions: () => apiClient.get('/settings/sessions'),
  terminateSession: (sessionId) => apiClient.delete(`/settings/sessions/${sessionId}`)
};

// 工作流相关API
export const workflowApi = {
  getWorkflows: (params) => apiClient.get('/workflows', { params }),
  getWorkflow: (id) => apiClient.get(`/workflows/${id}`),
  createWorkflow: (workflow) => apiClient.post('/workflows', workflow),
  updateWorkflow: (id, workflow) => apiClient.put(`/workflows/${id}`, workflow),
  deleteWorkflow: (id) => apiClient.delete(`/workflows/${id}`),
  executeWorkflow: (id, params) => apiClient.post(`/workflows/${id}/execute`, params)
};

// 执行历史相关API
export const executionApi = {
  getExecutions: (params) => apiClient.get('/executions', { params }),
  getExecution: (id) => apiClient.get(`/executions/${id}`),
  cancelExecution: (id) => apiClient.post(`/executions/${id}/cancel`)
};

// 代理相关API
export const agentApi = {
  getAgents: (params) => apiClient.get('/agents', { params }),
  getAgent: (id) => apiClient.get(`/agents/${id}`),
  createAgent: (agent) => apiClient.post('/agents', agent),
  updateAgent: (id, agent) => apiClient.put(`/agents/${id}`, agent),
  deleteAgent: (id) => apiClient.delete(`/agents/${id}`),
  testAgent: (id) => apiClient.post(`/agents/${id}/test`)
};

export default apiClient; 