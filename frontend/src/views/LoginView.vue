<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1 class="login-title">JiLang Agent</h1>
        <p class="login-subtitle">登录您的账户</p>
      </div>
      
      <div v-if="error" class="alert alert-danger">
        {{ error }}
      </div>
      
      <form @submit.prevent="login" class="login-form">
        <div class="form-group">
          <label for="username">用户名</label>
          <input 
            id="username"
            v-model="credentials.username" 
            type="text" 
            class="form-control" 
            :class="{'is-invalid': validationErrors.username}"
            placeholder="请输入用户名" 
            required 
            autofocus
          />
          <div v-if="validationErrors.username" class="invalid-feedback">
            {{ validationErrors.username }}
          </div>
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <div class="password-input">
            <input 
              id="password"
              v-model="credentials.password" 
              :type="showPassword ? 'text' : 'password'" 
              class="form-control" 
              :class="{'is-invalid': validationErrors.password}"
              placeholder="请输入密码" 
              required
            />
            <button 
              type="button" 
              class="toggle-password-btn"
              @click="showPassword = !showPassword"
            >
              <i :class="showPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
            </button>
          </div>
          <div v-if="validationErrors.password" class="invalid-feedback">
            {{ validationErrors.password }}
          </div>
        </div>
        
        <div class="form-options">
          <div class="form-check">
            <input 
              id="remember"
              v-model="credentials.remember" 
              type="checkbox" 
              class="form-check-input"
            />
            <label for="remember" class="form-check-label">记住我</label>
          </div>
          <a href="#" class="forgot-password">忘记密码?</a>
        </div>
        
        <button 
          type="submit" 
          class="btn btn-primary login-btn" 
          :disabled="loading"
        >
          <span v-if="loading" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
      
      <div class="login-footer">
        <p>
          还没有账户? <router-link to="/register">注册</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LoginView',
  data() {
    return {
      credentials: {
        username: '',
        password: '',
        remember: false
      },
      validationErrors: {},
      error: null,
      loading: false,
      showPassword: false
    }
  },
  methods: {
    validateForm() {
      this.validationErrors = {};
      let isValid = true;
      
      if (!this.credentials.username.trim()) {
        this.validationErrors.username = '用户名不能为空';
        isValid = false;
      }
      
      if (!this.credentials.password) {
        this.validationErrors.password = '密码不能为空';
        isValid = false;
      } else if (this.credentials.password.length < 6) {
        this.validationErrors.password = '密码长度不能少于6个字符';
        isValid = false;
      }
      
      return isValid;
    },
    login() {
      if (!this.validateForm()) {
        return;
      }
      
      this.error = null;
      this.loading = true;
      
      // 模拟API登录请求
      setTimeout(() => {
        try {
          // 这里将替换为实际的API调用
          // const response = await axios.post('/api/auth/login', this.credentials);
          
          // 模拟登录成功
          if (this.credentials.username === 'admin' && this.credentials.password === 'password') {
            // 存储令牌和用户信息到本地存储
            localStorage.setItem('token', 'mock-jwt-token');
            localStorage.setItem('user', JSON.stringify({
              id: 1,
              username: 'admin',
              fullName: '管理员',
              role: 'admin'
            }));
            
            // 重定向到首页或原目标页面
            const redirectPath = this.$route.query.redirect || '/';
            this.$router.push(redirectPath);
          } else {
            // 模拟登录失败
            this.error = '用户名或密码不正确';
          }
          
          this.loading = false;
        } catch (err) {
          this.error = '登录失败: ' + (err.message || '未知错误');
          this.loading = false;
        }
      }, 1500);
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 2rem;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.login-card {
  width: 100%;
  max-width: 420px;
  background-color: #fff;
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  padding: 2.5rem;
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-title {
  margin: 0;
  font-size: 1.8rem;
  font-weight: 700;
  color: #333;
}

.login-subtitle {
  margin-top: 0.5rem;
  color: #666;
  font-size: 1rem;
}

.login-form {
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #333;
}

.form-control {
  width: 100%;
  padding: 0.75rem 1rem;
  font-size: 1rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  transition: border-color 0.2s;
}

.form-control:focus {
  outline: none;
  border-color: #1976d2;
  box-shadow: 0 0 0 3px rgba(25, 118, 210, 0.1);
}

.password-input {
  position: relative;
}

.toggle-password-btn {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #777;
  cursor: pointer;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.form-check {
  display: flex;
  align-items: center;
}

.form-check-input {
  margin-right: 0.5rem;
}

.forgot-password {
  color: #1976d2;
  text-decoration: none;
  font-size: 0.9rem;
}

.forgot-password:hover {
  text-decoration: underline;
}

.login-btn {
  width: 100%;
  padding: 0.75rem;
  font-size: 1rem;
  background-color: #1976d2;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.login-btn:hover {
  background-color: #1565c0;
}

.login-btn:disabled {
  background-color: #90caf9;
  cursor: not-allowed;
}

.login-footer {
  text-align: center;
  margin-top: 1rem;
  font-size: 0.9rem;
  color: #666;
}

.login-footer a {
  color: #1976d2;
  text-decoration: none;
  font-weight: 500;
}

.login-footer a:hover {
  text-decoration: underline;
}

.alert {
  padding: 0.75rem 1rem;
  margin-bottom: 1.5rem;
  border-radius: 6px;
  font-size: 0.9rem;
}

.alert-danger {
  color: #721c24;
  background-color: #f8d7da;
  border: 1px solid #f5c6cb;
}

.is-invalid {
  border-color: #dc3545;
}

.invalid-feedback {
  display: block;
  width: 100%;
  margin-top: 0.25rem;
  font-size: 0.875rem;
  color: #dc3545;
}

@media (max-width: 576px) {
  .login-card {
    padding: 1.5rem;
  }
  
  .form-options {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }
}
</style> 