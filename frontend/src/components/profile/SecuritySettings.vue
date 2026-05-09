<template>
  <div class="security-section">
    <h3>修改密码</h3>
    <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="120px">
      <el-form-item label="当前密码" prop="oldPassword">
        <el-input v-model="passwordForm.oldPassword" type="password" show-password/>
      </el-form-item>
      <el-form-item label="新密码" prop="newPassword">
        <el-input v-model="passwordForm.newPassword" type="password" show-password/>
      </el-form-item>
      <el-form-item label="确认新密码" prop="confirmPassword">
        <el-input v-model="passwordForm.confirmPassword" type="password" show-password/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="loading" @click="handlePasswordChange">修改密码</el-button>
      </el-form-item>
    </el-form>
  </div>
  <el-divider/>
  <div class="danger-section">
    <h3>危险操作</h3>
    <el-alert title="注销账号是不可逆的操作，请谨慎操作！" type="error" :closable="false" show-icon/>
    <el-button type="danger" plain style="margin-top: 16px" @click="$emit('delete-account')">注销账号</el-button>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { FormInstance } from 'element-plus'

const props = defineProps<{
  loading: boolean
}>()

const emit = defineEmits(['change-password', 'delete-account'])

const passwordFormRef = ref<FormInstance>()
const passwordForm = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' })

const passwordRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [{ required: true, message: '请输入新密码', trigger: 'blur' }, {
    min: 6,
    message: '至少6位',
    trigger: 'blur'
  }],
  confirmPassword: [{ required: true, message: '请确认新密码', trigger: 'blur' }, {
    validator: (_: any, value: string, callback: Function) => {
      value !== passwordForm.newPassword ? callback(new Error('两次输入不一致')) : callback()
    }, trigger: 'blur'
  }]
}

const handlePasswordChange = async () => {
  if (!passwordFormRef.value) return
  await passwordFormRef.value.validate((valid) => {
    if (valid) {
      emit('change-password', {
        oldPassword: passwordForm.oldPassword,
        newPassword: passwordForm.newPassword
      })
    }
  })
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
