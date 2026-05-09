<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="login-title">用户登录</div>
        <div class="login-subtitle">请输入您的账户信息</div>
      </div>

      <el-form
          ref="formRef"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
          size="large"
      >
        <el-form-item prop="username">
          <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              :prefix-icon="User"
              clearable
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              :prefix-icon="Lock"
              show-password
              clearable
              @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item>
          <el-checkbox v-model="rememberMe">记住我</el-checkbox>
          <el-link type="primary" :underline="false" class="forgot-link">
            忘记密码？
          </el-link>
        </el-form-item>

        <el-form-item>
          <el-button
              type="primary"
              size="large"
              :loading="loading"
              class="login-button"
              @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <span>还没有账号？</span>
        <el-link type="primary" :underline="false" @click="$router.push('/register')">
          立即注册
        </el-link>
      </div>
    </div>

    <div class="login-background">
      <div class="background-circle circle1"></div>
      <div class="background-circle circle2"></div>
      <div class="background-circle circle3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref()
const loading = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度应在2-20个字符之间', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    loading.value = true
    try {
      const result = await userStore.login(loginForm.username, loginForm.password)

      if (result.success) {
        ElMessage.success('登录成功！')
        router.push('/')
      } else {
        ElMessage.error(result.message || '登录失败')
      }
    } catch (error) {
      ElMessage.error('登录失败，请稍后重试')
      console.error('Login error:', error)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  /* NVIDIA 标志性绿色渐变 */
  background: linear-gradient(135deg, #76b900 0%, #5b8c00 50%, #3d5a00 100%);
}

.login-card {
  position: relative;
  z-index: 10;
  width: 440px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(10px);
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 1;
  pointer-events: none;
}

.background-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(118, 185, 0, 0.15);
  animation: float 20s infinite linear;
  pointer-events: none;
}

.circle1 {
  width: 300px;
  height: 300px;
  top: -150px;
  left: -150px;
  background: rgba(91, 140, 0, 0.2);
}

.circle2 {
  width: 200px;
  height: 200px;
  bottom: -100px;
  right: -100px;
  animation-delay: -10s;
  background: rgba(118, 185, 0, 0.25);
}

.circle3 {
  width: 150px;
  height: 150px;
  top: 50%;
  right: 10%;
  animation-delay: -15s;
  background: rgba(61, 90, 0, 0.2);
}

@keyframes float {
  0% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(20px, 20px) rotate(90deg);
  }
  50% {
    transform: translate(0, 40px) rotate(180deg);
  }
  75% {
    transform: translate(-20px, 20px) rotate(270deg);
  }
  100% {
    transform: translate(0, 0) rotate(360deg);
  }
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-title {
  font-size: 32px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: #666;
}

.login-form {
  margin-bottom: 20px;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
  /* NVIDIA 绿色按钮 */
  background: linear-gradient(to right, #76b900, #5b8c00);
  border: none;
  transition: all 0.3s ease;
}

.login-button:hover {
  background: linear-gradient(to right, #8bd100, #6ba200);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(118, 185, 0, 0.3);
}

.login-button:active {
  transform: translateY(0);
}

.forgot-link {
  float: right;
  font-size: 14px;
  color: #76b900;
}

.login-footer {
  text-align: center;
  font-size: 14px;
  color: #666;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: none;
  border: 1px solid #dcdfe6;
  transition: all 0.3s;
}

:deep(.el-input__wrapper:hover) {
  border-color: #76b900;
}

:deep(.el-input__wrapper.is-focus) {
  border-color: #76b900;
  box-shadow: 0 0 0 2px rgba(118, 185, 0, 0.1);
}

:deep(.el-form-item.is-error .el-input__wrapper) {
  border-color: #f56c6c;
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: #76b900;
  border-color: #76b900;
}

:deep(.el-link--primary) {
  color: #76b900;
}

:deep(.el-link--primary:hover) {
  color: #5b8c00;
}

/* 输入框图标颜色 */
:deep(.el-input__prefix) {
  color: #76b900;
}
</style>