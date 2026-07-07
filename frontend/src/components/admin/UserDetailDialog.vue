<template>
  <a-modal
    v-model:visible="visible"
    title="用户详情"
    width="700px"
    @cancel="handleClose"
  >
    <a-descriptions :column="2" bordered v-if="user">
      <a-descriptions-item label="用户ID">{{ user.id }}</a-descriptions-item>
      <a-descriptions-item label="用户名">{{ user.username }}</a-descriptions-item>
      <a-descriptions-item label="手机号">{{ user.phone || '-' }}</a-descriptions-item>
      <a-descriptions-item label="邮箱">{{ user.email || '-' }}</a-descriptions-item>
      <a-descriptions-item label="状态">
        <a-tag v-if="user.status === 'NORMAL'" color="green">正常</a-tag>
        <a-tag v-else-if="user.status === 'LOCKED'" color="orange">锁定</a-tag>
        <a-tag v-else color="red">封禁</a-tag>
      </a-descriptions-item>
      <a-descriptions-item label="角色">
        <a-tag v-if="user.roles.includes('ROLE_ADMIN')" color="red">管理员</a-tag>
        <a-tag v-else color="gray">普通用户</a-tag>
      </a-descriptions-item>
      <a-descriptions-item label="创建时间">{{ formatDateTime(user.createdAt) }}</a-descriptions-item>
      <a-descriptions-item label="更新时间">{{ formatDateTime(user.updatedAt) }}</a-descriptions-item>
      <a-descriptions-item label="最后登录时间" :span="2">
        {{ formatDateTime(user.lastLoginAt) || '-' }}
      </a-descriptions-item>
    </a-descriptions>
  </a-modal>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { SysUser } from '@/models/entity'

const props = defineProps<{
  modelValue: boolean
  user: SysUser | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const formatDateTime = (dateString: string | undefined | null) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

const handleClose = () => {
  visible.value = false
}
</script>

<style scoped>
/* 可以添加自定义样式 */
</style>
