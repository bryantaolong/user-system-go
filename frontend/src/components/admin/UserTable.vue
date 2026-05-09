<template>
  <el-card class="table-card">
    <el-table
      v-loading="loading"
      :data="userList"
      border
      stripe
      class="user-table"
    >
      <el-table-column prop="id" label="ID" width="80" align="center" />
      <el-table-column prop="username" label="用户名" min-width="120" show-overflow-tooltip />
      <el-table-column prop="phone" label="手机号" width="140" />
      <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
      <el-table-column prop="roles" label="角色" width="140">
        <template #default="{ row }">
          <el-tag v-if="row.roles.includes('ROLE_ADMIN')" type="danger">管理员</el-tag>
          <el-tag v-else type="info">普通用户</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.status === 'NORMAL'" type="success">正常</el-tag>
          <el-tag v-else-if="row.status === 'LOCKED'" type="warning">锁定</el-tag>
          <el-tag v-else type="danger">封禁</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="lastLoginAt" label="最后登录时间" width="180">
        <template #default="{ row }">
          {{ formatDateTime(row.lastLoginAt) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button-group>
            <el-button
              size="small"
              type="primary"
              :icon="Edit"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-dropdown @command="(cmd: string) => handleCommand(cmd, row)">
              <el-button size="small" type="info">
                更多<el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="view" :icon="View">
                    查看详情
                  </el-dropdown-item>
                  <el-dropdown-item command="resetPwd" :icon="Key">
                    重置密码
                  </el-dropdown-item>
                  <el-dropdown-item v-if="row.status !== 'BANNED'" command="block" :icon="Lock">
                    封禁用户
                  </el-dropdown-item>
                  <el-dropdown-item v-else command="unblock" :icon="Unlock">
                    解封用户
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" :icon="Delete" divided>
                    删除用户
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pageNum"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      class="pagination"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  Edit,
  View,
  Key,
  Lock,
  Unlock,
  Delete,
  ArrowDown
} from '@element-plus/icons-vue'
import type { SysUser } from '@/models/entity'

const props = defineProps<{
  loading: boolean
  userList: SysUser[]
  total: number
  pageNum: number
  pageSize: number
}>()

const emit = defineEmits<{
  edit: [user: SysUser]
  view: [user: SysUser]
  resetPassword: [user: SysUser]
  block: [user: SysUser]
  unblock: [user: SysUser]
  delete: [user: SysUser]
  sizeChange: [size: number]
  currentChange: [page: number]
}>()

const pageNum = computed({
  get: () => props.pageNum,
  set: (val) => emit('currentChange', val)
})

const pageSize = computed({
  get: () => props.pageSize,
  set: (val) => emit('sizeChange', val)
})

const formatDateTime = (dateString: string | undefined | null) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

const handleEdit = (user: SysUser) => {
  emit('edit', user)
}

const handleCommand = (command: string, user: SysUser) => {
  switch (command) {
    case 'view':
      emit('view', user)
      break
    case 'resetPwd':
      emit('resetPassword', user)
      break
    case 'block':
      emit('block', user)
      break
    case 'unblock':
      emit('unblock', user)
      break
    case 'delete':
      emit('delete', user)
      break
  }
}

const handleSizeChange = (val: number) => {
  emit('sizeChange', val)
}

const handleCurrentChange = (val: number) => {
  emit('currentChange', val)
}
</script>

<style scoped>
.table-card {
  border-radius: 12px;
}

.user-table {
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

:deep(.el-tag) {
  font-weight: 500;
}

:deep(.el-button-group) {
  display: inline-flex;
}
</style>

