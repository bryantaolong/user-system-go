<template>
  <el-dialog
    v-model="visible"
    title="用户详情"
    width="700px"
  >
    <el-descriptions :column="2" border v-if="user">
      <el-descriptions-item label="用户ID">{{ user.id }}</el-descriptions-item>
      <el-descriptions-item label="用户名">{{ user.username }}</el-descriptions-item>
      <el-descriptions-item label="手机号">{{ user.phone || '-' }}</el-descriptions-item>
      <el-descriptions-item label="邮箱">{{ user.email || '-' }}</el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag v-if="user.status === 'NORMAL'" type="success">正常</el-tag>
        <el-tag v-else-if="user.status === 'LOCKED'" type="warning">锁定</el-tag>
        <el-tag v-else type="danger">封禁</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="角色">
        <el-tag v-if="user.roles.includes('ROLE_ADMIN')" type="danger">管理员</el-tag>
        <el-tag v-else type="info">普通用户</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ formatDateTime(user.createdAt) }}</el-descriptions-item>
      <el-descriptions-item label="更新时间">{{ formatDateTime(user.updatedAt) }}</el-descriptions-item>
      <el-descriptions-item label="最后登录时间" :span="2">
        {{ formatDateTime(user.lastLoginAt) || '-' }}
      </el-descriptions-item>
    </el-descriptions>
  </el-dialog>
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
</script>

<style scoped>
/* 可以添加自定义样式 */
</style>

