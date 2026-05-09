<template>
  <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" class="info-form">
    <el-form-item label="用户名">
      <el-input :value="username" disabled/>
    </el-form-item>
    <el-form-item label="真实姓名" prop="realName">
      <el-input v-model="form.realName" placeholder="请输入真实姓名"/>
    </el-form-item>
    <el-form-item label="性别" prop="gender">
      <el-radio-group v-model="form.gender">
        <el-radio :label="1">男</el-radio>
        <el-radio :label="0">女</el-radio>
      </el-radio-group>
    </el-form-item>
    <el-form-item label="生日" prop="birthday">
      <el-date-picker v-model="form.birthday" type="date" placeholder="选择生日"
                      value-format="YYYY-MM-DD"/>
    </el-form-item>
    <el-form-item label="手机号" prop="phone">
      <el-input v-model="form.phone" placeholder="请输入手机号"/>
    </el-form-item>
    <el-form-item label="邮箱" prop="email">
      <el-input v-model="form.email" placeholder="请输入邮箱"/>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" :loading="loading" @click="handleSave">保存修改</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { FormInstance } from 'element-plus'

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

const formRef = ref<FormInstance>()
const form = reactive({ ...props.initialData })

watch(() => props.initialData, (newVal) => {
  Object.assign(form, newVal)
}, { deep: true })

const rules = {
  realName: [{ min: 2, max: 20, message: '真实姓名长度应在2-20个字符之间', trigger: 'blur' }],
  phone: [{
    validator: (_: any, value: string, callback: Function) => {
      if (!value) return callback()
      const pattern = /^1[3-9]\d{9}$/
      pattern.test(value) ? callback() : callback(new Error('电话号码格式不正确'))
    }, trigger: 'blur'
  }],
  email: [{ type: 'email', message: '邮箱格式不正确', trigger: 'blur' }]
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate((valid) => {
    if (valid) {
      emit('save', { ...form })
    }
  })
}
</script>

<style scoped>
.info-form {
  padding: 20px 0;
  max-width: 500px;
}
</style>
