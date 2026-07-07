<template>
  <a-card class="search-card">
    <a-form :model="searchForm" layout="inline" class="search-form">
      <a-form-item label="用户名">
        <a-input
            v-model="searchForm.username"
            placeholder="请输入用户名"
            allow-clear
        />
      </a-form-item>

      <a-form-item label="手机号">
        <a-input
            v-model="searchForm.phone"
            placeholder="请输入手机号"
            allow-clear
        />
      </a-form-item>

      <a-form-item label="邮箱">
        <a-input
            v-model="searchForm.email"
            placeholder="请输入邮箱"
            allow-clear
        />
      </a-form-item>

      <a-form-item label="状态">
        <a-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            allow-clear
            style="width: 120px"
        >
          <a-option label="正常" value="NORMAL" />
          <a-option label="锁定" value="LOCKED" />
          <a-option label="封禁" value="BANNED" />
        </a-select>
      </a-form-item>

      <a-form-item>
        <a-button type="primary" @click="handleSearch">
          <template #icon><IconSearch /></template>
          查询
        </a-button>
        <a-button @click="handleReset">
          <template #icon><IconRefresh /></template>
          重置
        </a-button>
      </a-form-item>
    </a-form>
  </a-card>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { IconSearch, IconRefresh } from '@arco-design/web-vue/es/icon'

// 状态类型定义
export type UserStatus = 'NORMAL' | 'LOCKED' | 'BANNED' | ''

export interface UserSearchFormData {
  username: string
  phone: string
  email: string
  status: UserStatus
}

const emit = defineEmits<{
  search: [data: UserSearchFormData]
  reset: []
}>()

const searchForm = reactive<UserSearchFormData>({
  username: '',
  phone: '',
  email: '',
  status: ''  // 空字符串表示未选择，与 el-select clearable 配合
})

const handleSearch = () => {
  emit('search', { ...searchForm })
}

const handleReset = () => {
  Object.assign(searchForm, {
    username: '',
    phone: '',
    email: '',
    status: ''
  })
  emit('reset')
}

defineExpose({
  searchForm
})
</script>

<style scoped>
.search-card {
  margin-bottom: 20px;
  border-radius: 12px;
}

.search-form {
  margin-bottom: -18px;
}

.search-form :deep(.arco-form-item) {
  margin-right: 20px;
}
</style>
