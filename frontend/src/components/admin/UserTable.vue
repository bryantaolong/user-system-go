<template>
  <a-card class="table-card">
    <a-table
      :data="userList"
      :loading="loading"
      border
      stripe
      class="user-table"
      :pagination="false"
    >
      <a-table-column title="ID" data-index="id" :width="80" align="center" />
      <a-table-column title="用户名" data-index="username" :width="120" ellipsis tooltip />
      <a-table-column title="手机号" data-index="phone" :width="140" />
      <a-table-column title="邮箱" data-index="email" ellipsis tooltip />
      <a-table-column title="角色" data-index="roles" :width="140">
        <template #cell="{ record }">
          <a-tag v-if="record.roles.includes('ROLE_ADMIN')" color="red">管理员</a-tag>
          <a-tag v-else color="gray">普通用户</a-tag>
        </template>
      </a-table-column>
      <a-table-column title="状态" data-index="status" :width="100" align="center">
        <template #cell="{ record }">
          <a-tag v-if="record.status === 'NORMAL'" color="green">正常</a-tag>
          <a-tag v-else-if="record.status === 'LOCKED'" color="orange">锁定</a-tag>
          <a-tag v-else color="red">封禁</a-tag>
        </template>
      </a-table-column>
      <a-table-column title="最后登录时间" data-index="lastLoginAt" :width="180">
        <template #cell="{ record }">
          {{ formatDateTime(record.lastLoginAt) }}
        </template>
      </a-table-column>
      <a-table-column title="操作" :width="220" fixed="right">
        <template #cell="{ record }">
          <a-button-group>
            <a-button
              size="small"
              type="primary"
              @click="handleEdit(record)"
            >
              <template #icon><IconEdit /></template>
              编辑
            </a-button>
            <a-dropdown @select="(cmd: string) => handleCommand(cmd, record)">
              <a-button size="small" type="secondary">
                更多<template #icon><IconDown /></template>
              </a-button>
              <template #dropdown>
                <a-doption value="view">
                  <template #icon><IconEye /></template>
                  查看详情
                </a-doption>
                <a-doption value="resetPwd">
                  <template #icon><IconSafe /></template>
                  重置密码
                </a-doption>
                <a-doption v-if="record.status !== 'BANNED'" value="block">
                  <template #icon><IconLock /></template>
                  封禁用户
                </a-doption>
                <a-doption v-else value="unblock">
                  <template #icon><IconUnlock /></template>
                  解封用户
                </a-doption>
                <a-doption value="delete" :style="{ marginTop: '4px' }">
                  <template #icon><IconDelete /></template>
                  删除用户
                </a-doption>
              </template>
            </a-dropdown>
          </a-button-group>
        </template>
      </a-table-column>
    </a-table>

    <a-pagination
      :current="pageNum"
      :page-size="pageSize"
      :total="total"
      :page-size-options="[10, 20, 50, 100]"
      show-total
      show-jumper
      class="pagination"
      @change="handlePageChange"
    />
  </a-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  IconEdit,
  IconEye,
  IconSafe,
  IconLock,
  IconUnlock,
  IconDelete,
  IconDown
} from '@arco-design/web-vue/es/icon'
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

const handlePageChange = (current: number, pageSize: number) => {
  if (current !== pageNum.value) {
    emit('currentChange', current)
  }
  if (pageSize !== props.pageSize) {
    emit('sizeChange', pageSize)
  }
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

:deep(.arco-tag) {
  font-weight: 500;
}

:deep(.arco-button-group) {
  display: inline-flex;
}
</style>
