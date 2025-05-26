import { defineStore } from 'pinia';
import { agentApi } from '../api';

export const useAgentStore = defineStore('agent', {
  state: () => ({
    agents: [],
    currentAgent: null,
    loading: false,
    testLoading: false,
    testResult: null,
    error: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0
    },
    filters: {
      type: '',
      status: '',
      search: ''
    }
  }),
  
  getters: {
    getAgentById: (state) => (id) => {
      return state.agents.find(agent => agent.id === id) || null;
    },
    
    filteredAgents: (state) => {
      let result = [...state.agents];
      
      // 按类型筛选
      if (state.filters.type) {
        result = result.filter(agent => agent.type === state.filters.type);
      }
      
      // 按状态筛选
      if (state.filters.status) {
        result = result.filter(agent => agent.status === state.filters.status);
      }
      
      // 按搜索词筛选（名称和描述）
      if (state.filters.search) {
        const searchLower = state.filters.search.toLowerCase();
        result = result.filter(agent => 
          agent.name.toLowerCase().includes(searchLower) || 
          (agent.description && agent.description.toLowerCase().includes(searchLower))
        );
      }
      
      return result;
    },
    
    totalPages: (state) => {
      return Math.ceil(state.pagination.total / state.pagination.limit) || 1;
    },
    
    isFiltered: (state) => {
      return !!state.filters.type || 
             !!state.filters.status || 
             !!state.filters.search;
    },
    
    agentTypes: () => {
      return [
        { value: 'llm', label: '大语言模型' },
        { value: 'function', label: '函数代理' },
        { value: 'tool', label: '工具代理' },
        { value: 'custom', label: '自定义代理' }
      ];
    }
  },
  
  actions: {
    async fetchAgents(params = {}) {
      this.loading = true;
      this.error = null;
      
      try {
        // 合并默认参数和传入的参数
        const queryParams = {
          limit: this.pagination.limit,
          offset: (this.pagination.page - 1) * this.pagination.limit,
          type: this.filters.type,
          status: this.filters.status,
          search: this.filters.search,
          ...params
        };
        
        const response = await agentApi.getAgents(queryParams);
        
        // 假设API返回格式为 { data: [...], total: number }
        this.agents = response.data || [];
        this.pagination.total = response.total || 0;
        
        return response;
      } catch (error) {
        console.error('获取代理列表失败:', error);
        this.error = error.message || '获取代理列表失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchAgentById(id) {
      this.loading = true;
      this.error = null;
      
      try {
        const agent = await agentApi.getAgent(id);
        this.currentAgent = agent;
        
        // 更新列表中的代理
        const index = this.agents.findIndex(a => a.id === id);
        if (index !== -1) {
          this.agents[index] = agent;
        }
        
        return agent;
      } catch (error) {
        console.error(`获取代理 #${id} 失败:`, error);
        this.error = error.message || '获取代理详情失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async createAgent(agentData) {
      this.loading = true;
      this.error = null;
      
      try {
        const newAgent = await agentApi.createAgent(agentData);
        
        // 添加到列表开头
        this.agents.unshift(newAgent);
        
        // 更新总数
        this.pagination.total += 1;
        
        return newAgent;
      } catch (error) {
        console.error('创建代理失败:', error);
        this.error = error.message || '创建代理失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async updateAgent(id, agentData) {
      this.loading = true;
      this.error = null;
      
      try {
        const updatedAgent = await agentApi.updateAgent(id, agentData);
        
        // 更新列表中的代理
        const index = this.agents.findIndex(a => a.id === id);
        if (index !== -1) {
          this.agents[index] = updatedAgent;
        }
        
        // 如果当前查看的是这个代理，也更新currentAgent
        if (this.currentAgent && this.currentAgent.id === id) {
          this.currentAgent = updatedAgent;
        }
        
        return updatedAgent;
      } catch (error) {
        console.error(`更新代理 #${id} 失败:`, error);
        this.error = error.message || '更新代理失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async deleteAgent(id) {
      this.loading = true;
      this.error = null;
      
      try {
        await agentApi.deleteAgent(id);
        
        // 从列表中移除
        const index = this.agents.findIndex(a => a.id === id);
        if (index !== -1) {
          this.agents.splice(index, 1);
          
          // 更新总数
          this.pagination.total -= 1;
        }
        
        // 如果当前查看的是这个代理，清空currentAgent
        if (this.currentAgent && this.currentAgent.id === id) {
          this.currentAgent = null;
        }
      } catch (error) {
        console.error(`删除代理 #${id} 失败:`, error);
        this.error = error.message || '删除代理失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async testAgent(id, testData) {
      this.testLoading = true;
      this.testResult = null;
      this.error = null;
      
      try {
        const result = await agentApi.testAgent(id, testData);
        this.testResult = result;
        return result;
      } catch (error) {
        console.error(`测试代理 #${id} 失败:`, error);
        this.error = error.message || '测试代理失败';
        throw error;
      } finally {
        this.testLoading = false;
      }
    },
    
    setPage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.pagination.page = page;
        this.fetchAgents();
      }
    },
    
    setLimit(limit) {
      this.pagination.limit = limit;
      this.pagination.page = 1; // 重置到第一页
      this.fetchAgents();
    },
    
    setFilter(filters) {
      this.filters = { ...this.filters, ...filters };
      this.pagination.page = 1; // 重置到第一页
      this.fetchAgents();
    },
    
    clearFilters() {
      this.filters = {
        type: '',
        status: '',
        search: ''
      };
      this.pagination.page = 1;
      this.fetchAgents();
    },
    
    // 格式化状态
    formatStatus(status) {
      const statusMap = {
        'active': { text: '活跃', class: 'success' },
        'inactive': { text: '未激活', class: 'warning' },
        'error': { text: '错误', class: 'error' },
        'pending': { text: '待处理', class: 'info' }
      };
      
      return statusMap[status] || { text: status, class: 'default' };
    },
    
    // 格式化代理类型
    formatType(type) {
      const typeMap = {
        'llm': '大语言模型',
        'function': '函数代理',
        'tool': '工具代理',
        'custom': '自定义代理'
      };
      
      return typeMap[type] || type;
    }
  }
}); 