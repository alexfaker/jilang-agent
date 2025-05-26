import { defineStore } from 'pinia';
import { statsApi } from '../api';

export const useStatsStore = defineStore('stats', {
  state: () => ({
    dashboardStats: {
      totalWorkflows: 0,
      totalExecutions: 0,
      totalAgents: 0,
      successRate: 0,
      recentExecutions: [],
      executionsByStatus: [],
      executionTrend: [],
      topWorkflows: []
    },
    loading: false,
    error: null,
    lastUpdated: null
  }),
  
  getters: {
    formattedSuccessRate: (state) => {
      return `${(state.dashboardStats.successRate * 100).toFixed(1)}%`;
    },
    
    executionStatusColors: () => {
      return {
        'success': '#4caf50',
        'failed': '#f44336',
        'running': '#2196f3',
        'pending': '#ff9800',
        'cancelled': '#9e9e9e'
      };
    },
    
    // 获取最近7天的日期标签
    last7DaysLabels: () => {
      const labels = [];
      for (let i = 6; i >= 0; i--) {
        const date = new Date();
        date.setDate(date.getDate() - i);
        labels.push(date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' }));
      }
      return labels;
    },
    
    // 格式化趋势数据以适应图表
    formattedExecutionTrend: (state) => {
      const labels = [];
      const successData = [];
      const failedData = [];
      
      state.dashboardStats.executionTrend.forEach(item => {
        labels.push(item.date);
        successData.push(item.success || 0);
        failedData.push(item.failed || 0);
      });
      
      return {
        labels,
        datasets: [
          {
            label: '成功',
            data: successData,
            backgroundColor: 'rgba(76, 175, 80, 0.2)',
            borderColor: '#4caf50',
            borderWidth: 2,
            tension: 0.4
          },
          {
            label: '失败',
            data: failedData,
            backgroundColor: 'rgba(244, 67, 54, 0.2)',
            borderColor: '#f44336',
            borderWidth: 2,
            tension: 0.4
          }
        ]
      };
    },
    
    // 格式化状态分布数据以适应饼图
    formattedExecutionsByStatus: (state) => {
      const labels = [];
      const data = [];
      const backgroundColor = [];
      
      state.dashboardStats.executionsByStatus.forEach(item => {
        labels.push(item.status);
        data.push(item.count);
        backgroundColor.push(state.executionStatusColors[item.status] || '#9e9e9e');
      });
      
      return {
        labels,
        datasets: [
          {
            data,
            backgroundColor,
            borderWidth: 1
          }
        ]
      };
    }
  },
  
  actions: {
    async fetchDashboardStats() {
      this.loading = true;
      this.error = null;
      
      try {
        const stats = await statsApi.getDashboardStats();
        this.dashboardStats = stats;
        this.lastUpdated = new Date();
        return stats;
      } catch (error) {
        console.error('获取仪表盘统计数据失败:', error);
        this.error = error.message || '获取仪表盘统计数据失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchExecutionStats(params = {}) {
      this.loading = true;
      this.error = null;
      
      try {
        // 可以传入时间范围等参数
        const stats = await statsApi.getExecutionStats(params);
        
        // 更新相关统计数据
        if (stats.executionsByStatus) {
          this.dashboardStats.executionsByStatus = stats.executionsByStatus;
        }
        
        if (stats.executionTrend) {
          this.dashboardStats.executionTrend = stats.executionTrend;
        }
        
        if (stats.successRate !== undefined) {
          this.dashboardStats.successRate = stats.successRate;
        }
        
        this.lastUpdated = new Date();
        return stats;
      } catch (error) {
        console.error('获取执行统计数据失败:', error);
        this.error = error.message || '获取执行统计数据失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchWorkflowStats() {
      this.loading = true;
      this.error = null;
      
      try {
        const stats = await statsApi.getWorkflowStats();
        
        // 更新相关统计数据
        if (stats.topWorkflows) {
          this.dashboardStats.topWorkflows = stats.topWorkflows;
        }
        
        if (stats.totalWorkflows !== undefined) {
          this.dashboardStats.totalWorkflows = stats.totalWorkflows;
        }
        
        this.lastUpdated = new Date();
        return stats;
      } catch (error) {
        console.error('获取工作流统计数据失败:', error);
        this.error = error.message || '获取工作流统计数据失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchAgentStats() {
      this.loading = true;
      this.error = null;
      
      try {
        const stats = await statsApi.getAgentStats();
        
        // 更新相关统计数据
        if (stats.totalAgents !== undefined) {
          this.dashboardStats.totalAgents = stats.totalAgents;
        }
        
        this.lastUpdated = new Date();
        return stats;
      } catch (error) {
        console.error('获取代理统计数据失败:', error);
        this.error = error.message || '获取代理统计数据失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    // 一次性获取所有仪表盘所需的统计数据
    async fetchAllStats() {
      this.loading = true;
      this.error = null;
      
      try {
        await Promise.all([
          this.fetchDashboardStats(),
          this.fetchExecutionStats(),
          this.fetchWorkflowStats(),
          this.fetchAgentStats()
        ]);
        
        this.lastUpdated = new Date();
      } catch (error) {
        console.error('获取所有统计数据失败:', error);
        this.error = error.message || '获取统计数据失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    // 格式化数字
    formatNumber(number) {
      return new Intl.NumberFormat('zh-CN').format(number);
    },
    
    // 格式化日期时间
    formatDateTime(dateString) {
      if (!dateString) return '';
      
      const date = new Date(dateString);
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
    }
  }
}); 