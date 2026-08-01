<template>
  <div class="user-management">
    <a-card class="header-card">
      <div class="header-content">
        <div class="title-section">
          <h2>用户管理</h2>
          <p class="subtitle">管理系统用户信息</p>
        </div>
        <div class="button-group">
          <a-button type="primary" :icon="h(IconPlus)" @click="handleAddUser">新增用户</a-button>
          <a-button type="outline" status="warning" @click="handleExportAllUsers">导出用户数据</a-button>
        </div>
      </div>
    </a-card>

    <a-card class="search-card">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="用户名" field="username">
          <a-input v-model="searchForm.username" placeholder="请输入用户名" allow-clear />
        </a-form-item>
        <a-form-item label="手机号" field="phone">
          <a-input v-model="searchForm.phone" placeholder="请输入手机号" allow-clear />
        </a-form-item>
        <a-form-item label="邮箱" field="email">
          <a-input v-model="searchForm.email" placeholder="请输入邮箱" allow-clear />
        </a-form-item>
        <a-form-item label="状态" field="status">
          <a-select v-model="searchForm.status" placeholder="请选择状态" allow-clear style="width: 120px">
            <a-option :label="'正常'" :value="0" />
            <a-option :label="'锁定'" :value="1" />
            <a-option :label="'封禁'" :value="2" />
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" :icon="h(IconSearch)" @click="handleSearch">查询</a-button>
          <a-button :icon="h(IconRefresh)" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="userList"
        :loading="loading"
        :bordered="true"
        :stripe="true"
        :pagination="false"
      >
        <template #columns="{ record }">
          <template v-for="col in columns" :key="col.dataIndex">
            <template v-if="col.dataIndex === 'roles'">
              <a-space>
                <a-tag v-for="role in record.roles.split(',')" :key="role" :color="role === 'ROLE_ADMIN' ? 'red' : 'gray'">
                  {{ role === 'ROLE_ADMIN' ? '管理员' : role === 'ROLE_USER' ? '普通用户' : role }}
                </a-tag>
              </a-space>
            </template>
            <template v-else-if="col.dataIndex === 'status'">
              <a-tag v-if="record.status === 'NORMAL'" color="green">正常</a-tag>
              <a-tag v-else-if="record.status === 'LOCKED'" color="orange">锁定</a-tag>
              <a-tag v-else color="red">封禁</a-tag>
            </template>
            <template v-else-if="col.dataIndex === 'lastLoginAt'">
              {{ record.lastLoginAt ? new Date(record.lastLoginAt).toLocaleString('zh-CN') : '-' }}
            </template>
            <template v-else-if="col.dataIndex === 'action'">
              <a-button size="small" type="primary" @click="handleEdit(record)">编辑</a-button>
              <a-dropdown trigger="hover" @select="(cmd: string | number | Record<string, any> | undefined) => handleMore(cmd as string, record)">
                <a-button size="small">更多<template #icon><icon-down /></template></a-button>
                <template #content>
                  <a-doption value="view">查看详情</a-doption>
                  <a-doption value="resetPwd">重置密码</a-doption>
                  <a-doption v-if="record.status !== 'BANNED'" value="block">封禁用户</a-doption>
                  <a-doption v-else value="unblock">解封用户</a-doption>
                  <a-doption value="delete">删除用户</a-doption>
                </template>
              </a-dropdown>
            </template>
          </template>
        </template>
      </a-table>

      <div class="pagination">
        <a-pagination
          :current="pageNum"
          :page-size="pageSize"
          :total="total"
          :page-size-options="[10, 20, 50, 100]"
          show-total
          show-page-size
          @change="handlePageChange"
        />
      </div>
    </a-card>

    <!-- User Form Dialog -->
    <a-modal
      v-model:visible="dialogVisible"
      :title="dialogType === 'add' ? '新增用户' : '编辑用户'"
      width="600px"
      @close="handleDialogClose"
    >
      <a-form ref="userFormRef" :model="userForm" :rules="userFormRules" layout="vertical">
        <a-form-item label="用户名" field="username">
          <a-input v-model="userForm.username" :disabled="dialogType === 'edit'" />
        </a-form-item>
        <a-form-item label="手机号" field="phone">
          <a-input v-model="userForm.phone" />
        </a-form-item>
        <a-form-item label="邮箱" field="email">
          <a-input v-model="userForm.email" />
        </a-form-item>
        <a-form-item v-if="dialogType === 'add'" label="密码" field="password">
          <a-input-password v-model="userForm.password" />
        </a-form-item>
        <a-form-item label="角色" field="roleIds">
          <a-select v-model="userForm.roleIds" multiple placeholder="请选择角色" style="width: 100%">
            <a-option v-for="r in roleOptions" :key="r.id" :label="r.roleName" :value="r.id" />
          </a-select>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="dialogVisible = false">取消</a-button>
        <a-button type="primary" :loading="submitting" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>

    <!-- User Detail Dialog -->
    <a-modal v-model:visible="detailDialogVisible" title="用户详情" width="700px">
      <a-descriptions v-if="currentUser" :column="2" bordered>
        <a-descriptions-item label="用户ID">{{ currentUser.id }}</a-descriptions-item>
        <a-descriptions-item label="用户名">{{ currentUser.username }}</a-descriptions-item>
        <a-descriptions-item label="真实姓名">{{ currentUser.realName || '-' }}</a-descriptions-item>
        <a-descriptions-item label="性别">
          <a-tag v-if="currentUser.gender === 0" color="pink">女</a-tag>
          <a-tag v-else-if="currentUser.gender === 1" color="blue">男</a-tag>
          <span v-else>-</span>
        </a-descriptions-item>
        <a-descriptions-item label="生日">{{ formatDateTime(currentUser.birthday) || '-' }}</a-descriptions-item>
        <a-descriptions-item label="头像">
          <a-avatar v-if="currentUser.avatar" :size="40" :src="`/uploads/${currentUser.avatar}`">
            {{ currentUser.username?.charAt(0)?.toUpperCase() }}
          </a-avatar>
          <span v-else>-</span>
        </a-descriptions-item>
        <a-descriptions-item label="手机号">{{ currentUser.phone || '-' }}</a-descriptions-item>
        <a-descriptions-item label="邮箱">{{ currentUser.email || '-' }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag v-if="currentUser.status === 'NORMAL'" color="green">正常</a-tag>
          <a-tag v-else-if="currentUser.status === 'LOCKED'" color="orange">锁定</a-tag>
          <a-tag v-else color="red">封禁</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="角色">
          <a-space>
            <a-tag v-for="role in currentUser.roles.split(',')" :key="role" :color="role === 'ROLE_ADMIN' ? 'red' : 'gray'">
              {{ role === 'ROLE_ADMIN' ? '管理员' : role === 'ROLE_USER' ? '普通用户' : role }}
            </a-tag>
          </a-space>
        </a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ formatDateTime(currentUser.createdAt) }}</a-descriptions-item>
        <a-descriptions-item label="更新时间">{{ formatDateTime(currentUser.updatedAt) }}</a-descriptions-item>
        <a-descriptions-item label="最后登录时间" :span="2">{{ formatDateTime(currentUser.lastLoginAt) || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- Reset Password Prompt -->
    <a-modal v-model:visible="promptVisible" title="重置密码" @ok="handleResetPasswordConfirm">
      <a-input v-model="resetPasswordInput" placeholder="请输入新密码" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { IconPlus, IconSearch, IconRefresh, IconDown } from '@arco-design/web-vue/es/icon'
import * as userApi from '@/api/user'
import * as userExportApi from '@/api/userExport'
import * as userRoleApi from '@/api/userRole'
import * as userProfileApi from '@/api/userProfile'
import type { SysUser } from '@/types/entity'
import type { UserRoleOptionVO } from '@/types/vo'
import type { FormInstance, TableColumnData, FieldRule } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const userList = ref<SysUser[]>([])
const total = ref(0)
const pageNum = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const promptVisible = ref(false)
const resetPasswordInput = ref('')
const resetPasswordUserId = ref<number | null>(null)
const dialogType = ref<'add' | 'edit'>('add')
const currentUser = ref<SysUser | null>(null)
const roleOptions = ref<UserRoleOptionVO[]>([])
const userFormRef = ref<FormInstance>()

const searchForm = reactive({
  username: '',
  phone: '',
  email: '',
  status: undefined as number | undefined,
})

const userForm = reactive({
  username: '',
  phone: '',
  email: '',
  password: '',
  roleIds: [] as number[],
})

const userFormRules: Record<string, FieldRule | FieldRule[]> = {
  username: [
    { required: true, message: '请输入用户名' },
    { minLength: 2, maxLength: 20, message: '用户名长度应在2-20个字符之间' },
  ],
  phone: [{ required: false, match: /^1[3-9]\d{9}$/, message: '电话号码格式不正确' }],
  email: [{ required: false, type: 'email', message: '邮箱格式不正确' }],
  password: [
    { required: true, message: '请输入密码' },
    { minLength: 6, message: '密码至少6位' },
  ],
  roleIds: [{ required: true, message: '请选择角色' }],
}

const columns: TableColumnData[] = [
  { title: 'ID', dataIndex: 'id', width: 80, align: 'center' },
  { title: '用户名', dataIndex: 'username', width: 120, ellipsis: true },
  { title: '手机号', dataIndex: 'phone', width: 140 },
  { title: '邮箱', dataIndex: 'email', ellipsis: true },
  { title: '角色', dataIndex: 'roles', width: 120 },
  { title: '状态', dataIndex: 'status', width: 100, align: 'center' },
  { title: '最后登录时间', dataIndex: 'lastLoginAt', width: 180 },
  { title: '操作', dataIndex: 'action', width: 220, fixed: 'right' },
]

async function loadRoleOptions() {
  if (roleOptions.value.length > 0) return
  const res = await userRoleApi.listRoles()
  if (res.code === 200) {
    roleOptions.value = res.data
  }
}

async function loadUsers() {
  loading.value = true
  try {
    const params: Record<string, any> = {}
    if (searchForm.username) params.username = searchForm.username
    if (searchForm.phone) params.phone = searchForm.phone
    if (searchForm.email) params.email = searchForm.email
    if (searchForm.status) params.status = searchForm.status

    const hasSearch = Object.keys(params).length > 0
    const res = hasSearch
      ? await userApi.queryUsers(params, pageNum.value, pageSize.value)
      : await userApi.listUsers(pageNum.value, pageSize.value)

    if (res.code === 200) {
      userList.value = res.data.rows
      total.value = res.data.total
    } else {
      Message.error(res.message || '加载用户列表失败')
    }
  } catch (error) {
    Message.error('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(current: number) {
  pageNum.value = current
  loadUsers()
}

function handleSearch() {
  pageNum.value = 1
  loadUsers()
}

function handleReset() {
  searchForm.username = ''
  searchForm.phone = ''
  searchForm.email = ''
  searchForm.status = undefined
  pageNum.value = 1
  loadUsers()
}

function handleAddUser() {
  dialogType.value = 'add'
  resetUserForm()
  dialogVisible.value = true
}

async function handleEdit(user: SysUser) {
  dialogType.value = 'edit'
  currentUser.value = user
  resetUserForm()
  userForm.username = user.username || ''
  userForm.phone = user.phone || ''
  userForm.email = user.email || ''

  await loadRoleOptions()
  const roleNames = user.roles ? user.roles.split(',').map((r: string) => r.trim()).filter(Boolean) : []
  userForm.roleIds = roleOptions.value.filter((o: UserRoleOptionVO) => roleNames.includes(o.roleName)).map((o: UserRoleOptionVO) => o.id)
  dialogVisible.value = true
}

async function handleView(user: SysUser) {
  try {
    const res = await userProfileApi.getUserProfileByUserId(user.id)
    if (res.code === 200) {
      currentUser.value = { ...user, ...res.data }
      detailDialogVisible.value = true
    } else {
      Message.error(res.message || '获取用户详情失败')
    }
  } catch (error) {
    Message.error('获取用户详情失败')
  }
}

async function handleResetPasswordConfirm() {
  if (!resetPasswordInput.value) {
    Message.warning('请输入新密码')
    return
  }
  if (!resetPasswordUserId.value) return
  try {
    const res = await userApi.resetPassword(resetPasswordUserId.value, resetPasswordInput.value)
    if (res.code === 200) {
      Message.success('密码重置成功')
      promptVisible.value = false
    } else {
      Message.error(res.message || '密码重置失败')
    }
  } catch (error: any) {
    if (error !== 'cancel' && error !== 'close') {
      console.error('Reset password error:', error)
    }
  }
}

async function handleBlockUser(user: SysUser) {
  try {
    await Modal.confirm({
      title: `确定要封禁用户 "${user.username}" 吗？`,
      content: '警告',
      okText: '确定',
      cancelText: '取消',
      okButtonProps: { status: 'warning' },
    })
    const res = await userApi.blockUser(user.id)
    if (res.code === 200) {
      Message.success('用户已封禁')
      loadUsers()
    } else {
      Message.error(res.message || '封禁失败')
    }
  } catch (error: any) {
    if (error !== 'cancel' && error !== 'close') {
      console.error('Block user error:', error)
    }
  }
}

async function handleUnblockUser(user: SysUser) {
  try {
    await Modal.confirm({
      title: `确定要解封用户 "${user.username}" 吗？`,
      content: '提示',
      okText: '确定',
      cancelText: '取消',
    })
    const res = await userApi.unblockUser(user.id)
    if (res.code === 200) {
      Message.success('用户已解封')
      loadUsers()
    } else {
      Message.error(res.message || '解封失败')
    }
  } catch (error: any) {
    if (error !== 'cancel' && error !== 'close') {
      console.error('Unblock user error:', error)
    }
  }
}

async function handleDeleteUser(user: SysUser) {
  try {
    await Modal.confirm({
      title: `确定要删除用户 "${user.username}" 吗？此操作不可恢复！`,
      content: '危险操作',
      okText: '确定',
      cancelText: '取消',
      okButtonProps: { status: 'danger' },
    })
    const res = await userApi.deleteUser(user.id)
    if (res.code === 200) {
      Message.success('用户已删除')
      if (userList.value.length === 1 && pageNum.value > 1) {
        pageNum.value = pageNum.value - 1
      }
      loadUsers()
    } else {
      Message.error(res.message || '删除失败')
    }
  } catch (error: any) {
    if (error !== 'cancel' && error !== 'close') {
      console.error('Delete user error:', error)
    }
  }
}

function handleMore(command: string, user: SysUser) {
  switch (command) {
    case 'view':
      handleView(user)
      break
    case 'resetPwd':
      resetPasswordUserId.value = user.id
      resetPasswordInput.value = ''
      promptVisible.value = true
      break
    case 'block':
      handleBlockUser(user)
      break
    case 'unblock':
      handleUnblockUser(user)
      break
    case 'delete':
      handleDeleteUser(user)
      break
  }
}

async function handleSubmit() {
  if (!userFormRef.value) return
  try {
    await userFormRef.value.validate()
    submitting.value = true
    if (dialogType.value === 'add') {
      const payload = {
        username: userForm.username || '',
        password: userForm.password || '',
        phone: userForm.phone,
        email: userForm.email,
        roleIds: userForm.roleIds,
      }
      const res = await userApi.createUser(payload)
      if (res.code === 200) {
        Message.success('新增用户成功')
        dialogVisible.value = false
        loadUsers()
      } else {
        Message.error(res.message || '新增失败')
      }
    } else {
      if (!currentUser.value) {
        Message.error('用户信息不存在')
        return
      }
      const updateRes = await userApi.updateUser(currentUser.value.id, {
        phone: userForm.phone,
        email: userForm.email,
        version: currentUser.value.version,
      })
      if (updateRes.code !== 200) {
        Message.error(updateRes.message || '更新失败')
        return
      }
      if (userForm.roleIds && userForm.roleIds.length > 0) {
        const roleRes = await userApi.changeUserRoles(currentUser.value.id, userForm.roleIds)
        if (roleRes.code !== 200) {
          Message.error(roleRes.message || '角色更新失败')
          return
        }
      }
      Message.success('更新成功')
      dialogVisible.value = false
      loadUsers()
    }
  } catch (error) {
    Message.error(dialogType.value === 'add' ? '新增失败' : '更新失败')
  } finally {
    submitting.value = false
  }
}

function resetUserForm() {
  userForm.username = ''
  userForm.phone = ''
  userForm.email = ''
  userForm.password = ''
  userForm.roleIds = []
}

function handleDialogClose() {
  currentUser.value = null
  resetUserForm()
}

async function handleExportAllUsers() {
  try {
    await userExportApi.exportAllUsers('所有用户数据')
    Message.success('所有用户数据已开始导出！')
  } catch (error) {
    Message.error('导出失败，请重试！')
  }
}

function formatDateTime(dateString?: string) {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.header-card {
  margin-bottom: 16px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section h2 {
  font-size: 20px;
  color: #1d2129;
  margin-bottom: 4px;
}

.subtitle {
  font-size: 14px;
  color: #86909c;
}

.button-group {
  display: flex;
  gap: 8px;
}

.search-card {
  margin-bottom: 16px;
}

.table-card {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
