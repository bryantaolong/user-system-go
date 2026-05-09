<template>
  <el-card class="search-card">
    <el-form :model="searchForm" inline class="search-form">
      <el-form-item label="用户名">
        <el-input
            v-model="searchForm.username"
            placeholder="请输入用户名"
            clearable
        />
      </el-form-item>

      <el-form-item label="手机号">
        <el-input
            v-model="searchForm.phone"
            placeholder="请输入手机号"
            clearable
        />
      </el-form-item>

      <el-form-item label="邮箱">
        <el-input
            v-model="searchForm.email"
            placeholder="请输入邮箱"
            clearable
        />
      </el-form-item>

      <el-form-item label="状态">
        <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
        >
          <el-option label="正常" value="NORMAL" />
          <el-option label="锁定" value="LOCKED" />
          <el-option label="封禁" value="BANNED" />
        </el-select>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :icon="Search" @click="handleSearch">
          查询
        </el-button>
        <el-button :icon="Refresh" @click="handleReset">
          重置
        </el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { Search, Refresh } from '@element-plus/icons-vue'

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

.search-form :deep(.el-form-item) {
  margin-right: 20px;
}
</style>