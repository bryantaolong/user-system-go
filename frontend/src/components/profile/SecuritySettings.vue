<template>
  <div class="security-section">
    <h3>修改密码</h3>
    <a-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="120px">
      <a-form-item label="当前密码" field="oldPassword">
        <a-input-password v-model="passwordForm.oldPassword" />
      </a-form-item>
      <a-form-item label="新密码" field="newPassword">
        <a-input-password v-model="passwordForm.newPassword" />
      </a-form-item>
      <a-form-item label="确认新密码" field="confirmPassword">
        <a-input-password v-model="passwordForm.confirmPassword" />
      </a-form-item>
      <a-form-item>
        <a-button type="primary" :loading="loading" @click="handlePasswordChange">修改密码</a-button>
      </a-form-item>
    </a-form>
  </div>
  <a-divider/>
  <div class="danger-section">
    <h3>危险操作</h3>
    <a-alert title="注销账号是不可逆的操作，请谨慎操作！" type="error" :closable="false" />
    <a-button type="primary" status="danger" style="margin-top: 16px" @click="$emit('delete-account')">注销账号</a-button>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'

const props = defineProps<{
  loading: boolean
}>()

const emit = defineEmits(['change-password', 'delete-account'])

const passwordFormRef = ref()
const passwordForm = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' })

const passwordRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [{ required: true, message: '请输入新密码', trigger: 'blur' }, {
    minLength: 6,
    message: '至少6位',
    trigger: 'blur'
  }],
  confirmPassword: [{ required: true, message: '请确认新密码', trigger: 'blur' }, {
    validator: (_: any, value: string) => {
      return value === passwordForm.newPassword || '两次输入不一致'
    },
    trigger: 'blur'
  }]
}

const handlePasswordChange = async () => {
  if (!passwordFormRef.value) return
  try {
    const valid = await passwordFormRef.value.validate()
    if (valid) {
      emit('change-password', {
        oldPassword: passwordForm.oldPassword,
        newPassword: passwordForm.newPassword
      })
    }
  } catch (error) {
    // 验证失败
  }
}

defineExpose({
  resetPasswordForm: () => passwordFormRef.value?.resetFields()
})
</script>

<style scoped>
.security-section h3 {
  margin-bottom: 20px;
}

.danger-section {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}
</style>
