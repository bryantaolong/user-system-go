<template>
  <div class="page-container">
    <div class="card-box register-box">
      <div class="card-header">
        <h2>用户注册</h2>
        <p>创建您的账户</p>
      </div>

      <a-form ref="formRef" :model="form" :rules="rules" size="large" @keyup.enter="handleSubmit">
        <a-form-item field="username">
          <a-input v-model="form.username" placeholder="请输入用户名">
            <template #add-before>
              <icon-user />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item field="password">
          <a-input-password v-model="form.password" placeholder="请输入密码">
            <template #add-before>
              <icon-lock />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item field="confirmPassword">
          <a-input-password v-model="form.confirmPassword" placeholder="请确认密码">
            <template #add-before>
              <icon-lock />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item field="phone">
          <a-input v-model="form.phone" placeholder="请输入手机号码（可选）">
            <template #add-before>
              <icon-phone />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item field="email">
          <a-input v-model="form.email" placeholder="请输入邮箱地址（可选）">
            <template #add-before>
              <icon-email />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item>
          <a-button type="primary" class="full-width" size="large" :loading="loading" @click="handleSubmit">
            注册
          </a-button>
        </a-form-item>
      </a-form>

      <div class="card-footer">
        <span>已有账号？</span>
        <a-button type="text" @click="$router.push('/login')">立即登录</a-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { IconUser, IconLock, IconPhone, IconEmail } from '@arco-design/web-vue/es/icon'
import { useUserStore } from '@/stores/user'
import type { FormInstance, FieldRule } from '@arco-design/web-vue'

const router = useRouter()
const store = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  phone: '',
  email: '',
})

const validateConfirmPassword = (value: string) => {
  if (value !== form.password) {
    return new Error('两次输入的密码不一致')
  }
  return true
}

const rules: Record<string, FieldRule | FieldRule[]> = {
  username: [
    { required: true, message: '请输入用户名' },
    { minLength: 2, maxLength: 20, message: '用户名长度应在2-20个字符之间' },
  ],
  password: [
    { required: true, message: '请输入密码' },
    { minLength: 6, message: '密码至少6位' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码' },
    { validator: validateConfirmPassword },
  ],
  phone: [
    { match: /^1[3-9]\d{9}$/, message: '电话号码格式不正确' },
  ],
  email: [
    { type: 'email', message: '邮箱格式不正确' },
  ],
}

async function handleSubmit() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    loading.value = true
    const result = await store.register({
      username: form.username,
      password: form.password,
      phone: form.phone || undefined,
      email: form.email || undefined,
    })
    if (result.success) {
      Message.success('注册成功！正在跳转到登录页面...')
      setTimeout(() => router.push('/login'), 1500)
    } else {
      Message.error(result.message || '注册失败')
    }
  } catch (error) {
    Message.error('注册失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>
