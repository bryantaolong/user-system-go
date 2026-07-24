import { useState, useCallback } from 'react'
import { Card, Button, Modal } from '@arco-design/web-react'
import message from '@arco-design/web-react/es/Message'
import { IconPlus } from '@arco-design/web-react/icon'
import UserSearchForm from '@/components/admin/UserSearchForm'
import UserTable from '@/components/admin/UserTable'
import UserFormDialog from '@/components/admin/UserFormDialog'
import UserDetailDialog from '@/components/admin/UserDetailDialog'
import * as userApi from '@/api/user/user'
import * as userExportApi from '@/api/user/userExport'
import * as userRoleApi from '@/api/user/userRole'
import * as userProfileApi from '@/api/user/userProfile'
import type { SysUser } from '@/models/entity'
import type { UserRoleOptionVO } from '@/models/vo'
import type { UserFormData } from '@/components/admin/UserFormDialog'
import type { UserSearchFormData } from '@/components/admin/UserSearchForm'

const UserManagement = () => {
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [userList, setUserList] = useState<SysUser[]>([])
  const [total, setTotal] = useState(0)
  const [pageNum, setPageNum] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  const [dialogVisible, setDialogVisible] = useState(false)
  const [detailDialogVisible, setDetailDialogVisible] = useState(false)
  const [dialogType, setDialogType] = useState<'add' | 'edit'>('add')
  const [currentUser, setCurrentUser] = useState<SysUser | null>(null)
  const [roleOptions, setRoleOptions] = useState<UserRoleOptionVO[]>([])

  const [searchFormData, setSearchFormData] = useState<UserSearchFormData>({
    username: '',
    phone: '',
    email: '',
    status: '',
  })

  const [userForm, setUserForm] = useState<UserFormData>({
    username: '',
    phone: '',
    email: '',
    realName: '',
    gender: '',
    birthday: '',
    avatar: '',
    password: '',
    roleIds: [],
  })

  const loadRoleOptions = async () => {
    if (roleOptions.length > 0) return
    const res = await userRoleApi.listRoles()
    if (res.code === 200) {
      setRoleOptions(res.data)
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

  const loadUsers = useCallback(async () => {
    setLoading(true)
    try {
      const searchParams: Record<string, any> = {}
      if (searchFormData.username) searchParams.username = searchFormData.username
      if (searchFormData.phone) searchParams.phone = searchFormData.phone
      if (searchFormData.email) searchParams.email = searchFormData.email
      if (searchFormData.status) searchParams.status = searchFormData.status

      const hasSearch = Object.keys(searchParams).length > 0

      const res = hasSearch
        ? await userApi.queryUsers(searchParams, pageNum, pageSize)
        : await userApi.listUsers(pageNum, pageSize)

      if (res.code === 200) {
        setUserList(res.data.rows)
        setTotal(res.data.total)
      } else {
        message.error(res.message || '加载用户列表失败')
      }
    } catch (error) {
      message.error('加载用户列表失败')
      console.error('Load users error:', error)
    } finally {
      setLoading(false)
    }
  }, [searchFormData, pageNum, pageSize])

  const handleSearch = (data: UserSearchFormData) => {
    setSearchFormData(data)
    setPageNum(1)
  }

  const handleReset = () => {
    setSearchFormData({
      username: '',
      phone: '',
      email: '',
      status: '',
    })
    setPageNum(1)
  }

  const handleSizeChange = (val: number) => {
    setPageSize(val)
  }

  const handleCurrentChange = (val: number) => {
    setPageNum(val)
  }

  const handleAddUser = () => {
    setDialogType('add')
    setUserForm({
      username: '',
      phone: '',
      email: '',
      realName: '',
      gender: '',
      birthday: '',
      avatar: '',
      password: '',
      roleIds: [],
    })
    setDialogVisible(true)
  }

  const handleEdit = async (user: SysUser) => {
    setDialogType('edit')
    setCurrentUser(user)
    setUserForm({
      username: user.username || '',
      phone: user.phone || '',
      email: user.email || '',
      realName: '',
      gender: '',
      birthday: '',
      avatar: '',
      password: '',
      roleIds: [],
    })

    await loadRoleOptions()
    const roleNames = user.roles ? user.roles.split(',').map((r) => r.trim()).filter(Boolean) : []
    const roleIds = roleOptions
      .filter((o) => roleNames.includes(o.roleName))
      .map((o) => o.id)
    setUserForm((prev) => ({ ...prev, roleIds }))
    setDialogVisible(true)
  }

  const handleView = async (user: SysUser) => {
    try {
      const res = await userProfileApi.getUserProfileByUserId(user.id)
      if (res.code === 200) {
        setCurrentUser({ ...user, ...res.data })
        setDetailDialogVisible(true)
      } else {
        message.error(res.message || '获取用户详情失败')
      }
    } catch (error) {
      message.error('获取用户详情失败')
      console.error('View user error:', error)
    }
  }

  const handleResetPassword = async (user: SysUser) => {
    try {
      const { value } = await Modal.prompt('请输入新密码', '重置密码', {
        confirmText: '确定',
        cancelText: '取消',
        placeholder: '请输入新密码',
      })

      const res = await userApi.resetPassword(user.id, value)
      if (res.code === 200) {
        message.success('密码重置成功')
      } else {
        message.error(res.message || '密码重置失败')
      }
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('Reset password error:', error)
      }
    }
  }

  const handleBlockUser = async (user: SysUser) => {
    try {
      await Modal.confirm(`确定要封禁用户 "${user.username}" 吗？`, '警告', { type: 'warning' })

      const res = await userApi.blockUser(user.id)
      if (res.code === 200) {
        message.success('用户已封禁')
        loadUsers()
      } else {
        message.error(res.message || '封禁失败')
      }
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('Block user error:', error)
      }
    }
  }

  const handleUnblockUser = async (user: SysUser) => {
    try {
      await Modal.confirm(`确定要解封用户 "${user.username}" 吗？`, '提示', { type: 'info' })

      const res = await userApi.unblockUser(user.id)
      if (res.code === 200) {
        message.success('用户已解封')
        loadUsers()
      } else {
        message.error(res.message || '解封失败')
      }
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('Unblock user error:', error)
      }
    }
  }

  const handleDeleteUser = async (user: SysUser) => {
    try {
      await Modal.confirm(`确定要删除用户 "${user.username}" 吗？此操作不可恢复！`, '危险操作', { type: 'error' })

      const res = await userApi.deleteUser(user.id)
      if (res.code === 200) {
        message.success('用户已删除')
        if (userList.length === 1 && pageNum > 1) {
          setPageNum(pageNum - 1)
        }
        loadUsers()
      } else {
        message.error(res.message || '删除失败')
      }
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('Delete user error:', error)
      }
    }
  }

  const handleSubmit = async (formData: UserFormData) => {
    setSubmitting(true)
    try {
      if (dialogType === 'add') {
        const payload = {
          username: formData.username || '',
          password: formData.password || '',
          phone: formData.phone,
          email: formData.email,
          roleIds: formData.roleIds,
        }
        const res = await userApi.createUser(payload)
        if (res.code === 200) {
          message.success('新增用户成功')
          setDialogVisible(false)
          loadUsers()
        } else {
          message.error(res.message || '新增失败')
        }
      } else {
        if (!currentUser) {
          message.error('用户信息不存在')
          return
        }
        const updateRes = await userApi.updateUser(currentUser.id, formData)
        if (updateRes.code !== 200) {
          message.error(updateRes.message || '更新失败')
          return
        }

        if (formData.roleIds && formData.roleIds.length > 0) {
          const roleRes = await userApi.changeUserRoles(currentUser.id, formData.roleIds)
          if (roleRes.code !== 200) {
            message.error(roleRes.message || '角色更新失败')
            return
          }
        }

        message.success('更新成功')
        setDialogVisible(false)
        loadUsers()
      }
    } catch (error) {
      message.error(dialogType === 'add' ? '新增失败' : '更新失败')
      console.error('Submit user error:', error)
    } finally {
      setSubmitting(false)
    }
  }

  const handleDialogClose = () => {
    setCurrentUser(null)
    setUserForm({
      username: '',
      phone: '',
      email: '',
      realName: '',
      gender: '',
      birthday: '',
      avatar: '',
      password: '',
      roleIds: [],
    })
  }

  // 导出功能
  const handleExportAllUsers = async () => {
    try {
      const fileName = '所有用户数据'
      await userExportApi.exportAllUsers(fileName, mapStatusToCode(searchFormData.status))
      message.success('所有用户数据已开始导出！')
    } catch (error: any) {
      if (error !== 'cancel') {
        console.error('导出所有用户失败:', error)
        message.error('导出失败，请重试！')
      }
    }
  }

  return (
    <div className="user-management">
      <Card className="header-card">
        <div className="header-content">
          <div className="title-section">
            <h2>用户管理</h2>
            <p className="subtitle">管理系统用户信息</p>
          </div>
          <div className="button-group">
            <Button type="primary" icon={<IconPlus />} onClick={handleAddUser}>
              新增用户
            </Button>
            <Button type="warning" onClick={handleExportAllUsers}>
              导出用户数据
            </Button>
          </div>
        </div>
      </Card>

      <UserSearchForm onSearch={handleSearch} onReset={handleReset} />

      <UserTable
        loading={loading}
        userList={userList}
        total={total}
        pageNum={pageNum}
        pageSize={pageSize}
        onEdit={handleEdit}
        onView={handleView}
        onResetPassword={handleResetPassword}
        onBlock={handleBlockUser}
        onUnblock={handleUnblockUser}
        onDelete={handleDeleteUser}
        onSizeChange={handleSizeChange}
        onCurrentChange={handleCurrentChange}
      />

      <UserFormDialog
        visible={dialogVisible}
        dialogType={dialogType}
        userForm={userForm}
        submitting={submitting}
        onSubmit={handleSubmit}
        onClose={handleDialogClose}
      />

      <UserDetailDialog
        visible={detailDialogVisible}
        user={currentUser}
        onClose={() => setDetailDialogVisible(false)}
      />
    </div>
  )
}

export default UserManagement
