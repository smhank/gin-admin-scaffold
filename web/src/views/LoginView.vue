<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-left">
        <div class="login-left-content">
          <div class="login-logo">
            <img src="https://www.gin-vue-admin.com/assets/logo-2f5d2e3f.png" alt="logo" class="logo-img" />
            <span class="logo-text">Gin Admin</span>
          </div>
          <div class="login-desc">
            <p>基于 Gin + Vue3 + Element Plus 的后台管理系统</p>
          </div>
        </div>
      </div>
      <div class="login-right">
        <div class="login-form-wrapper">
          <h2 class="login-title">欢迎登录</h2>
          <p class="login-subtitle">请输入您的账户信息</p>
          <el-form :model="loginForm" class="login-form" size="large">
            <el-form-item>
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名"
                :prefix-icon="User"
              />
            </el-form-item>
            <el-form-item>
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                :prefix-icon="Lock"
                show-password
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleLogin" style="width: 100%" :loading="loading">
                登 录
              </el-button>
            </el-form-item>
          </el-form>
          <div class="login-tips">
            <p>默认账号: admin / admin123</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import request from '@/utils/request'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const loginForm = ref({ username: 'admin', password: 'admin123' })

const handleLogin = async () => {
  if (!loginForm.value.username || !loginForm.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const res = await request.post('/login', loginForm.value)
    authStore.setToken(res.token)
    authStore.setPermissions(res.permissions)
    localStorage.setItem('username', loginForm.value.username)
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  display: flex;
  width: 900px;
  height: 500px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.login-left {
  flex: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.login-left-content {
  text-align: center;
  padding: 40px;
}

.login-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.logo-img {
  width: 48px;
  height: 48px;
  margin-right: 12px;
}

.logo-text {
  font-size: 28px;
  font-weight: bold;
  color: #fff;
}

.login-desc {
  font-size: 14px;
  opacity: 0.9;
  line-height: 1.8;
}

.login-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-form-wrapper {
  width: 320px;
  padding: 40px;
}

.login-title {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: #999;
  margin-bottom: 30px;
}

.login-form {
  width: 100%;
}

.login-tips {
  margin-top: 20px;
  text-align: center;
  font-size: 12px;
  color: #999;
}
</style>
