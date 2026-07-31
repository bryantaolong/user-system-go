<template>
  <div class="basic-info">
    <div class="avatar-section">
      <a-avatar :size="80" :src="avatarUrl" style="font-size: 32px; background: #165dff;">
        {{ username?.charAt(0)?.toUpperCase() }}
      </a-avatar>
      <div class="avatar-actions">
        <a-button size="small" @click="triggerUpload">上传头像</a-button>
        <input ref="fileInput" type="file" accept="image/png,image/jpeg,image/gif,image/webp" style="display: none" @change="handleFileChange" />
      </div>
    </div>

    <a-form ref="formRef" :model="form" layout="horizontal" style="max-width: 500px">
      <a-form-item label="用户名" field="username">
        <a-input :model-value="username" disabled />
      </a-form-item>
      <a-form-item label="真实姓名" field="realName">
        <a-input v-model="form.realName" placeholder="请输入真实姓名" />
      </a-form-item>
    <a-form-item label="性别" field="gender">
      <a-radio-group v-model="form.gender">
        <a-radio :value="1">男</a-radio>
        <a-radio :value="0">女</a-radio>
      </a-radio-group>
    </a-form-item>
    <a-form-item label="生日" field="birthday">
      <a-date-picker v-model="form.birthday" format="YYYY-MM-DD" value-format="YYYY-MM-DD" placeholder="选择日期" style="width: 100%" />
    </a-form-item>
    <a-form-item label="手机号" field="phone">
      <a-input v-model="form.phone" placeholder="请输入手机号" />
    </a-form-item>
    <a-form-item label="邮箱" field="email">
      <a-input v-model="form.email" placeholder="请输入邮箱" />
    </a-form-item>
    <a-form-item>
      <a-button type="primary" :loading="loading" @click="handleSave">保存修改</a-button>
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { Message } from '@arco-design/web-vue'
import type { FormInstance, FieldRule } from '@arco-design/web-vue'
import { uploadAvatar } from '@/api/userProfile'

const props = defineProps<{
  username?: string
  initialData: { realName: string; gender: 0 | 1; birthday: string; phone: string; email: string; avatar?: string }
  loading: boolean
}>()

const emit = defineEmits<{
  save: [data: any]
}>()

const formRef = ref<FormInstance>()
const fileInput = ref<HTMLInputElement | null>(null)
const avatarUrl = ref('')

const form = reactive({
  realName: '',
  gender: 1 as 0 | 1,
  birthday: '',
  phone: '',
  email: '',
})

const rules: Record<string, FieldRule | FieldRule[]> = {
  realName: [{ required: false, minLength: 2, maxLength: 20, message: '真实姓名长度应在2-20个字符之间' }],
  phone: [{ required: false, match: /^1[3-9]\d{9}$/, message: '电话号码格式不正确' }],
  email: [{ required: false, type: 'email', message: '邮箱格式不正确' }],
}

watch(
  () => props.initialData,
  (val) => {
    if (val) {
      form.realName = val.realName
      form.gender = val.gender
      form.birthday = val.birthday
      form.phone = val.phone
      form.email = val.email
      avatarUrl.value = val.avatar || ''
    }
  },
  { immediate: true, deep: true }
)

function triggerUpload() {
  fileInput.value?.click()
}

async function handleFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  try {
    const res = await uploadAvatar(file)
    if (res.code === 200 && res.data) {
      avatarUrl.value = `/uploads/${res.data}`
      Message.success('头像上传成功')
      emit('save', { ...form, avatar: res.data })
    } else {
      Message.error(res.message || '头像上传失败')
    }
  } catch (err: any) {
    Message.error(err.message || '头像上传失败')
  } finally {
    target.value = ''
  }
}

async function handleSave() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    emit('save', { ...form, avatar: avatarUrl.value })
  } catch (e) {
    // validation failed
  }
}
</script>

<style scoped>
.basic-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar-actions {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
</style>
