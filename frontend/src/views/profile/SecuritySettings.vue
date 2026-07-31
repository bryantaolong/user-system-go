<template>
  <div class="security-section">
    <h3>修改密码</h3>
    <a-form ref="formRef" :model="form" :rules="rules" layout="horizontal" style="max-width: 500px">
      <a-form-item label="当前密码" field="oldPassword">
        <a-input-password v-model="form.oldPassword" placeholder="请输入当前密码" />
      </a-form-item>
      <a-form-item label="新密码" field="newPassword">
        <a-input-password v-model="form.newPassword" placeholder="请输入新密码" />
      </a-form-item>
      <a-form-item label="确认新密码" field="confirmPassword">
        <a-input-password v-model="form.confirmPassword" placeholder="请确认新密码" />
      </a-form-item>
      <a-form-item>
        <a-button type="primary" :loading="loading" @click="handlePasswordChange">修改密码</a-button>
      </a-form-item>
    </a-form>

    <a-divider />

    <div class="danger-section">
      <h3>危险操作</h3>
      <a-alert title="注销账号是不可逆的操作，请谨慎操作！" type="error" :closable="false" />
      <a-button type="primary" status="danger" style="margin-top: 16px" @click="$emit('delete-account')">注销账号</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { FormInstance, FieldRule } from '@arco-design/web-vue'

const props = defineProps<{ loading: boolean }>()

const emit = defineEmits<{
  'change-password': [data: { oldPassword: string; newPassword: string }]
  'delete-account': []
}>()

const formRef = ref<FormInstance>()

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const validateConfirm = (value: string) => {
  if (!value || form.newPassword !== value) {
    return new Error('两次输入不一致')
  }
  return true
}

const rules: Record<string, FieldRule | FieldRule[]> = {
  oldPassword: [{ required: true, message: '请输入当前密码' }],
  newPassword: [
    { required: true, message: '请输入新密码' },
    { minLength: 6, message: '至少6位' },
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码' },
    { validator: validateConfirm },
  ],
}

async function handlePasswordChange() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    emit('change-password', { oldPassword: form.oldPassword, newPassword: form.newPassword })
    form.oldPassword = ''
    form.newPassword = ''
    form.confirmPassword = ''
  } catch (e) {
    // validation failed
  }
}
</script>

<style scoped>
.security-section {
  max-width: 600px;
}

.security-section h3,
.danger-section h3 {
  font-size: 16px;
  color: #1d2129;
  margin-bottom: 16px;
}

.danger-section {
  margin-top: 24px;
}
</style>
