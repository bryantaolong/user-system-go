<template>
  <div class="page-container">
    <div class="card-box">
      <div class="card-header">
        <h2>用户登录</h2>
        <p>请输入您的账户信息</p>
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

        <a-form-item>
          <a-button type="primary" class="full-width" size="large" :loading="loading" @click="handleSubmit">
            登录
          </a-button>
        </a-form-item>
      </a-form>

      <div class="card-footer">
        <span>还没有账号？</span>
        <a-button type="text" @click="$router.push('/register')">立即注册</a-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { IconUser, IconLock } from '@arco-design/web-vue/es/icon'
import { useUserStore } from '@/stores/user'
import type { FormInstance, FieldRule } from '@arco-design/web-vue'

const router = useRouter()
const store = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
})

const rules: Record<string, FieldRule | FieldRule[]> = {
  username: [
    { required: true, message: '请输入用户名' },
    { minLength: 2, maxLength: 20, message: '用户名长度应在2-20个字符之间' },
  ],
  password: [
    { required: true, message: '请输入密码' },
    { minLength: 6, message: '密码至少6位' },
  ],
}

onMounted(() => {
  if (store.isAuthenticated) {
    router.push('/')
  }
})

async function handleSubmit() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    loading.value = true
    const result = await store.login(form.username, form.password)
    if (result.success) {
      Message.success('登录成功！')
      router.push('/')
    } else {
      Message.error(result.message || '登录失败')
    }
  } catch (error) {
    Message.error('登录失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>
