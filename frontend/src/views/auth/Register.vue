<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <div class="register-title">用户注册</div>
        <div class="register-subtitle">创建您的账户</div>
      </div>

      <a-form
        ref="formRef"
        :model="registerForm"
        :rules="registerRules"
        class="register-form"
        size="large"
        @submit="handleRegister"
      >
        <a-form-item field="username">
          <a-input
            v-model="registerForm.username"
            placeholder="请输入用户名"
            allow-clear
          >
            <template #add-before><IconUser /></template>
          </a-input>
        </a-form-item>

        <a-form-item field="password">
          <a-input-password
            v-model="registerForm.password"
            placeholder="请输入密码"
            allow-clear
          />
        </a-form-item>

        <a-form-item field="confirmPassword">
          <a-input-password
            v-model="registerForm.confirmPassword"
            placeholder="请确认密码"
            allow-clear
          />
        </a-form-item>

        <a-form-item field="phone">
          <a-input
            v-model="registerForm.phone"
            placeholder="请输入手机号码（可选）"
            allow-clear
          >
            <template #add-before><IconPhone /></template>
          </a-input>
        </a-form-item>

        <a-form-item field="email">
          <a-input
            v-model="registerForm.email"
            placeholder="请输入邮箱地址（可选）"
            allow-clear
          >
            <template #add-before><IconMessage /></template>
          </a-input>
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            size="large"
            :loading="loading"
            class="register-button"
            html-type="submit"
          >
            注册
          </a-button>
        </a-form-item>
      </a-form>

      <div class="register-footer">
        <span>已有账号？</span>
        <a-link type="primary" :underline="false" @click="router.push('/login')">
          立即登录
        </a-link>
      </div>
    </div>

    <div class="register-background">
      <div class="background-circle circle1"></div>
      <div class="background-circle circle2"></div>
      <div class="background-circle circle3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ArcoMessage } from '@/utils/arco-message'
import { IconUser, IconLock, IconPhone, IconMessage } from '@arco-design/web-vue/es/icon'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref()
const loading = ref(false)

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  phone: '',
  email: ''
})

const validateConfirmPassword = (value: string, callback: (error?: string) => void) => {
  if (value !== registerForm.password) {
    callback('两次输入的密码不一致')
  } else {
    callback()
  }
}

const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { minLength: 2, maxLength: 20, message: '用户名长度应在2-20个字符之间', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { minLength: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  phone: [
    {
      match: /^1[3-9]\d{9}$/,
      message: '电话号码格式不正确',
      trigger: 'blur'
    }
  ],
  email: [
    {
      type: 'email',
      message: '邮箱格式不正确',
      trigger: 'blur'
    }
  ]
}

const handleRegister = async () => {
  if (!formRef.value) return

  try {
    const valid = await formRef.value.validate()
    if (!valid) return

    loading.value = true
    const result = await userStore.register({
      username: registerForm.username,
      password: registerForm.password,
      phone: registerForm.phone || undefined,
      email: registerForm.email || undefined
    })

    if (result.success) {
      ArcoMessage.success('注册成功！正在跳转到登录页面...')
      setTimeout(() => {
        router.push('/login')
      }, 1500)
    } else {
      ArcoMessage.error(result.message || '注册失败')
    }
  } catch (error) {
    ArcoMessage.error('注册失败，请稍后重试')
    console.error('Register error:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.register-background {
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
  background: rgba(255, 255, 255, 0.1);
  animation: float 20s infinite linear;
  pointer-events: none;
}

.circle1 {
  width: 300px;
  height: 300px;
  top: -150px;
  right: -150px;
}

.circle2 {
  width: 200px;
  height: 200px;
  bottom: -100px;
  left: -100px;
  animation-delay: -10s;
}

.circle3 {
  width: 150px;
  height: 150px;
  top: 20%;
  left: 10%;
  animation-delay: -15s;
}

@keyframes float {
  0% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(-20px, 20px) rotate(90deg);
  }
  50% {
    transform: translate(0, 40px) rotate(180deg);
  }
  75% {
    transform: translate(20px, 20px) rotate(270deg);
  }
  100% {
    transform: translate(0, 0) rotate(360deg);
  }
}

.register-card {
  position: relative;
  z-index: 10;
  width: 480px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-title {
  font-size: 32px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.register-subtitle {
  font-size: 14px;
  color: #909399;
}

.register-form {
  margin-bottom: 20px;
}

.register-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
}

.register-footer {
  text-align: center;
  font-size: 14px;
  color: #909399;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}

:deep(.arco-input) {
  border-radius: 8px;
}

:deep(.arco-input:hover) {
  border-color: #f5576c;
}

:deep(.arco-input-focused) {
  border-color: #f5576c;
  box-shadow: 0 0 0 2px rgba(245, 87, 108, 0.1);
}

:deep(.arco-form-item-message) {
  color: #f56c6c;
}

:deep(.arco-button-primary) {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border: none;
  transition: all 0.3s;
}

:deep(.arco-button-primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(245, 87, 108, 0.4);
}
</style>
