<template>
  <div class="user-management">
    <a-card class="header-card">
      <div class="header-content">
        <div class="title-section">
          <h2>用户管理</h2>
          <p class="subtitle">管理系统用户信息</p>
        </div>
        <!-- 右侧按钮组 -->
        <div class="button-group">
          <a-button type="primary" @click="handleAddUser">
            <template #icon><IconPlus /></template>
            新增用户
          </a-button>
          <a-button type="warning" @click="handleExportAllUsers">
            导出用户数据
          </a-button>
        </div>
      </div>
    </a-card>

    <UserSearchForm
        @search="handleSearch"
        @reset="handleReset"
    />

    <UserTable
        :loading="loading"
        :user-list="userList"
        :total="total"
        :page-num="pageNum"
        :page-size="pageSize"
        @edit="handleEdit"
        @view="handleView"
        @reset-password="handleResetPassword"
        @block="handleBlockUser"
        @unblock="handleUnblockUser"
        @delete="handleDeleteUser"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
    />

    <UserFormDialog
        v-model="dialogVisible"
        :dialog-type="dialogType"
        :user-form="userForm"
        :submitting="submitting"
        @submit="handleSubmit"
        @close="handleDialogClose"
    />

    <UserDetailDialog
        v-model="detailDialogVisible"
        :user="currentUser"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ArcoMessage, ArcoMessageBox } from '@/utils/arco-message'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import * as userApi from '@/api/user/user'
import * as userExportApi from '@/api/user/userExport'
import * as userRoleApi from '@/api/user/userRole'
import * as userProfileApi from '@/api/user/userProfile'
import type { SysUser } from '@/models/entity'
import UserSearchForm from '@/components/admin/UserSearchForm.vue'
import UserTable from '@/components/admin/UserTable.vue'
import UserFormDialog from '@/components/admin/UserFormDialog.vue'
import UserDetailDialog from '@/components/admin/UserDetailDialog.vue'

import type { UserFormData } from '@/components/admin/UserFormDialog.vue'
import type { UserSearchFormData } from '@/components/admin/UserSearchForm.vue'
import type { UserCreateRequest } from '@/models/request/user'
import type { UserRoleOptionVO } from '@/models/vo'

const loading = ref(false)
const submitting = ref(false)
const userList = ref<SysUser[]>([])
const total = ref(0)
const pageNum = ref(1)
const pageSize = ref(10)

// 搜索条件（保持与后端接口一致）
const searchFormData = reactive<UserSearchFormData>({
  username: '',
  phone: '',
  email: '',
  status: ''
})

const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')

const userForm = reactive<UserFormData>({
  username: '',
  phone: '',
  email: '',
  realName: '',
  gender: '',
  birthday: '',
  avatar: '',
  password: '',
  roleIds: []
})

const currentUser = ref<SysUser | null>(null)

const roleOptions = ref<UserRoleOptionVO[]>([])

const loadRoleOptions = async () => {
  if (roleOptions.value.length > 0) return
  const res = await userRoleApi.listRoles()
  if (res.code === 200) {
    roleOptions.value = res.data
  }
}

const mapStatusToCode = (status: UserSearchFormData['status']): number | undefined => {
  if (!status) return undefined
  switch (status) {
    case 'NORMAL':
      return 0
    case 'BANNED':
      return 1
    case 'LOCKED':
      return 2
  }
}

const loadUsers = async () => {
  loading.value = true
  try {
    // 构建搜索参数，过滤空值
    const searchParams: Record<string, any> = {}
    if (searchFormData.username) searchParams.username = searchFormData.username
    if (searchFormData.phone) searchParams.phone = searchFormData.phone
    if (searchFormData.email) searchParams.email = searchFormData.email
    if (searchFormData.status) searchParams.status = searchFormData.status

    const hasSearch = Object.keys(searchParams).length > 0

    const res = hasSearch
        ? await userApi.queryUsers(searchParams, pageNum.value, pageSize.value)
        : await userApi.listUsers(pageNum.value, pageSize.value)

    if (res.code === 200) {
      userList.value = res.data.rows
      total.value = res.data.total
    } else {
      ArcoMessage.error(res.message || '加载用户列表失败')
    }
  } catch (error) {
    ArcoMessage.error('加载用户列表失败')
    console.error('Load users error:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = (data: UserSearchFormData) => {
  Object.assign(searchFormData, data)
  pageNum.value = 1
  loadUsers()
}

const handleReset = () => {
  Object.assign(searchFormData, {
    username: '',
    phone: '',
    email: '',
    status: ''
  })
  pageNum.value = 1
  loadUsers()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadUsers()
}

const handleCurrentChange = (val: number) => {
  pageNum.value = val
  loadUsers()
}

const handleAddUser = () => {
  dialogType.value = 'add'
  Object.assign(userForm, {
    username: '',
    phone: '',
    email: '',
    realName: '',
    gender: '',
    birthday: '',
    avatar: '',
    password: '',
    roleIds: []
  })
  dialogVisible.value = true
}

const handleEdit = (user: SysUser) => {
  dialogType.value = 'edit'
  currentUser.value = user

  userForm.username = user.username || ''
  userForm.phone = user.phone || ''
  userForm.email = user.email || ''
  userForm.roleIds = []

  // 回填角色：将用户的 roles 字符串映射为 roleIds
  void (async () => {
    await loadRoleOptions()
    const roleNames = user.roles ? user.roles.split(',').map(r => r.trim()).filter(Boolean) : []
    userForm.roleIds = roleOptions.value
      .filter(o => roleNames.includes(o.roleName))
      .map(o => o.id)
  })()

  dialogVisible.value = true
}

const handleView = async (user: SysUser) => {
  try {
    const res = await userProfileApi.getUserProfileByUserId(user.id)
    if (res.code === 200) {
      // 合并详细信息和列表数据
      currentUser.value = { ...user, ...res.data }
      detailDialogVisible.value = true
    } else {
      ArcoMessage.error(res.message || '获取用户详情失败')
    }
  } catch (error) {
    ArcoMessage.error('获取用户详情失败')
    console.error('View user error:', error)
  }
}

const handleResetPassword = async (user: SysUser) => {
  try {
    const { value } = await ArcoMessageBox.prompt('请输入新密码', '重置密码', {
      confirmText: '确定',
      cancelText: '取消',
      placeholder: '请输入新密码'
    })

    const res = await userApi.resetPassword(user.id, value)
    if (res.code === 200) {
      ArcoMessage.success('密码重置成功')
    } else {
      ArcoMessage.error(res.message || '密码重置失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Reset password error:', error)
    }
  }
}

const handleBlockUser = async (user: SysUser) => {
  try {
    await ArcoMessageBox.confirm(
        `确定要封禁用户 "${user.username}" 吗？`,
        '警告',
        { type: 'warning' }
    )

    const res = await userApi.blockUser(user.id)
    if (res.code === 200) {
      ArcoMessage.success('用户已封禁')
      loadUsers()
    } else {
      ArcoMessage.error(res.message || '封禁失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Block user error:', error)
    }
  }
}

const handleUnblockUser = async (user: SysUser) => {
  try {
    await ArcoMessageBox.confirm(
        `确定要解封用户 "${user.username}" 吗？`,
        '提示',
        { type: 'info' }
    )

    const res = await userApi.unblockUser(user.id)
    if (res.code === 200) {
      ArcoMessage.success('用户已解封')
      loadUsers()
    } else {
      ArcoMessage.error(res.message || '解封失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Unblock user error:', error)
    }
  }
}

const handleDeleteUser = async (user: SysUser) => {
  try {
    await ArcoMessageBox.confirm(
        `确定要删除用户 "${user.username}" 吗？此操作不可恢复！`,
        '危险操作',
        { type: 'error' }
    )

    const res = await userApi.deleteUser(user.id)
    if (res.code === 200) {
      ArcoMessage.success('用户已删除')
      // 如果当前页只有一条数据且不是第一页，返回上一页
      if (userList.value.length === 1 && pageNum.value > 1) {
        pageNum.value--
      }
      loadUsers()
    } else {
      ArcoMessage.error(res.message || '删除失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Delete user error:', error)
    }
  }
}

const handleSubmit = async (formData: UserFormData) => {
  submitting.value = true
  try {
    if (dialogType.value === 'add') {
      const payload: UserCreateRequest = {
        username: formData.username || '',
        password: formData.password || '',
        phone: formData.phone,
        email: formData.email,
        roleIds: formData.roleIds
      }
      const res = await userApi.createUser(payload)
      if (res.code === 200) {
        ArcoMessage.success('新增用户成功')
        dialogVisible.value = false
        loadUsers()
      } else {
        ArcoMessage.error(res.message || '新增失败')
      }
    } else {
      if (!currentUser.value) {
        ArcoMessage.error('用户信息不存在')
        return
      }
      const updateRes = await userApi.updateUser(currentUser.value.id, formData)
      if (updateRes.code !== 200) {
        ArcoMessage.error(updateRes.message || '更新失败')
        return
      }

      if (formData.roleIds && formData.roleIds.length > 0) {
        const roleRes = await userApi.changeUserRoles(currentUser.value.id, formData.roleIds)
        if (roleRes.code !== 200) {
          ArcoMessage.error(roleRes.message || '角色更新失败')
          return
        }
      }

      ArcoMessage.success('更新成功')
      dialogVisible.value = false
      loadUsers()
    }
  } catch (error) {
    ArcoMessage.error(dialogType.value === 'add' ? '新增失败' : '更新失败')
    console.error('Submit user error:', error)
  } finally {
    submitting.value = false
  }
}

const handleDialogClose = () => {
  currentUser.value = null
  Object.assign(userForm, {
    username: '',
    phone: '',
    email: '',
    realName: '',
    gender: '',
    birthday: '',
    avatar: '',
    password: '',
    roleIds: []
  })
}

/* ----------------- 导出功能 ----------------- */
const handleExportAllUsers = async () => {
  try {
    const result = await ArcoMessageBox.prompt(
        '请输入导出文件名:',
        '导出所有用户',
        {
          confirmText: '确定',
          cancelText: '取消',
          placeholder: '请输入导出文件名'
        }
    )

    const fileName = result.value || '所有用户数据'
    // 传递当前筛选的状态
    await userExportApi.exportAllUsers(fileName, mapStatusToCode(searchFormData.status))
    ArcoMessage.success('所有用户数据已开始导出！')
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('导出所有用户失败:', error)
      ArcoMessage.error('导出失败，请重试！')
    }
  }
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.user-management {
  padding: 20px;
}

.header-card {
  margin-bottom: 20px;
  border-radius: 12px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  color: #303133;
}

.subtitle {
  margin: 0;
  font-size: 14px;
  color: #909399;
}

/* 按钮组样式 - 右对齐 */
.button-group {
  display: flex;
  gap: 10px;
  align-items: center;
}
</style>
