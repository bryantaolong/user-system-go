<template>
  <a-form ref="formRef" :model="form" class="info-form">
    <a-form-item label="用户名">
      <a-input :model-value="username" disabled/>
    </a-form-item>
    <a-form-item label="真实姓名" field="realName">
      <a-input v-model="form.realName" placeholder="请输入真实姓名"/>
    </a-form-item>
    <a-form-item label="性别" field="gender">
      <a-radio-group v-model="form.gender">
        <a-radio :value="1">男</a-radio>
        <a-radio :value="0">女</a-radio>
      </a-radio-group>
    </a-form-item>
    <a-form-item label="生日" field="birthday">
      <a-date-picker v-model="form.birthday" format="YYYY-MM-DD" />
    </a-form-item>
    <a-form-item label="手机号" field="phone">
      <a-input v-model="form.phone" placeholder="请输入手机号"/>
    </a-form-item>
    <a-form-item label="邮箱" field="email">
      <a-input v-model="form.email" placeholder="请输入邮箱"/>
    </a-form-item>
    <a-form-item>
      <a-button type="primary" :loading="loading" @click="handleSave">保存修改</a-button>
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'

const props = defineProps<{
  username?: string
  initialData: {
    realName: string
    gender: 1 | 0
    birthday: string
    phone: string
    email: string
  }
  loading: boolean
}>()

const emit = defineEmits(['save'])

const formRef = ref()
const form = reactive({ ...props.initialData })

watch(() => props.initialData, (newVal) => {
  Object.assign(form, newVal)
}, { deep: true })

const rules = {
  realName: [{ minLength: 2, maxLength: 20, message: '真实姓名长度应在2-20个字符之间', trigger: 'blur' }],
  phone: [{
    match: /^1[3-9]\d{9}$/,
    message: '电话号码格式不正确',
    trigger: 'blur'
  }],
  email: [{type: 'email', message: '邮箱格式不正确', trigger: 'blur'}]
}

const handleSave = async () => {
  if (!formRef.value) return
  try {
    const valid = await formRef.value.validate()
    if (valid) {
      emit('save', { ...form })
    }
  } catch (error) {
    // 验证失败
  }
}
</script>

<style scoped>
.info-form {
  padding: 20px 0;
  max-width: 500px;
}
</style>
