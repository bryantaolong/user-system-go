<template>
  <a-modal
      v-model:visible="visible"
      :title="dialogTitle"
      width="600px"
      @cancel="handleClose"
  >
    <a-form
        ref="formRef"
        :model="localUserForm"
        :rules="userRules"
        label-width="100px"
        class="user-form"
    >
      <a-form-item label="用户名" field="username">
        <a-input v-model="localUserForm.username" :disabled="dialogType === 'edit'"/>
      </a-form-item>
      <a-form-item label="手机号" field="phone">
        <a-input v-model="localUserForm.phone"/>
      </a-form-item>
      <a-form-item label="邮箱" field="email">
        <a-input v-model="localUserForm.email"/>
      </a-form-item>
      <a-form-item v-if="dialogType === 'add'" label="密码" field="password">
        <a-input-password v-model="localUserForm.password" />
      </a-form-item>
      <a-form-item label="角色" field="roleIds">
        <a-select v-model="localUserForm.roleIds" multiple placeholder="请选择角色">
          <a-option
              v-for="r in roleOptions"
              :key="r.id"
              :value="r.id"
          >
            {{ r.roleName }}
          </a-option>
        </a-select>
      </a-form-item>
    </a-form>
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" :loading="submitting" @click="handleSubmit">
        确定
      </a-button>
    </template>
  </a-modal>
</template>

<script setup lang="ts">
import {ref, computed, watch, nextTick} from 'vue'
import type { FormInstance } from '@arco-design/web-vue'
import type {UserUpdateRequest} from '@/models/request/user'
import * as userRoleApi from '@/api/user/userRole'
import type {UserRoleOptionVO} from '@/models/vo'

export interface UserFormData extends UserUpdateRequest {
  username?: string
  password?: string
  roleIds?: number[]
}

const props = defineProps<{
  modelValue: boolean
  dialogType: 'add' | 'edit'
  userForm: UserFormData
  submitting?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: [formData: UserFormData]
  close: []
  'update:userForm': [formData: UserFormData]
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const dialogTitle = computed(() =>
    props.dialogType === 'add' ? '新增用户' : '编辑用户'
)

const formRef = ref<FormInstance>()

const roleOptions = ref<UserRoleOptionVO[]>([])

const loadRoleOptions = async () => {
  if (roleOptions.value.length > 0) return
  const res = await userRoleApi.listRoles()
  if (res.code === 200) {
    roleOptions.value = res.data
  }
}

// 初始化 localUserForm，确保所有字段有默认值
const localUserForm = ref<UserFormData>({
  username: props.userForm.username || '',
  phone: props.userForm.phone || '',
  email: props.userForm.email || '',
  roleIds: props.userForm.roleIds ? [...props.userForm.roleIds] : [],
  password: props.userForm.password || ''
})

// 仅监听 visible 变化，在打开时同步最新 props 数据
watch(visible, async (isVisible) => {
  if (isVisible) {
    await loadRoleOptions()
    // 同步最新数据，同时防止 undefined 覆盖
    localUserForm.value = {
      username: props.userForm.username || '',
      phone: props.userForm.phone || '',
      email: props.userForm.email || '',
      roleIds: props.userForm.roleIds ? [...props.userForm.roleIds] : [],
      password: props.userForm.password || ''
    }

    await nextTick()
    formRef.value?.clearValidate()
  }
}, {immediate: true})

// 动态验证规则
const userRules = computed(() => {
  const rules: Record<string, any> = {
    phone: [{match: /^1[3-9]\d{9}$/, message: '电话号码格式不正确', trigger: 'blur'}],
    email: [{type: 'email', message: '邮箱格式不正确', trigger: 'blur'}],
    roleIds: [{required: true, message: '请选择角色', trigger: 'change'}]
  }

  if (props.dialogType === 'add') {
    rules.username = [
      {required: true, message: '请输入用户名', trigger: 'blur'},
      {minLength: 2, maxLength: 20, message: '用户名长度应在2-20个字符之间', trigger: 'blur'}
    ]
    rules.password = [
      {required: true, message: '请输入密码', trigger: 'blur'},
      {minLength: 6, message: '密码至少6位', trigger: 'blur'}
    ]
  }

  return rules
})

const handleSubmit = async () => {
  if (!formRef.value) return

  formRef.value.clearValidate()

  await nextTick()

  try {
    const valid = await formRef.value.validate()
    if (!valid) return
    emit('submit', {...localUserForm.value})
  } catch (error) {
    // 验证失败
  }
}

const handleCancel = () => {
  visible.value = false
}

const handleClose = () => {
  formRef.value?.resetFields()
  formRef.value?.clearValidate()
  // 重置为当前 props 值（安全）
  localUserForm.value = {
    username: props.userForm.username || '',
    phone: props.userForm.phone || '',
    email: props.userForm.email || '',
    roleIds: props.userForm.roleIds ? [...props.userForm.roleIds] : [],
    password: props.userForm.password || ''
  }
  emit('close')
}

defineExpose({
  userFormRef: formRef
})
</script>

<style scoped>
.user-form {
  padding: 20px 20px 0 20px;
}
</style>
