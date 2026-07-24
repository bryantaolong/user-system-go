import { useState } from 'react'
import { Table, Pagination, Tag, Button, Dropdown } from '@arco-design/web-react'
import {
  IconEdit,
  IconEye,
  IconSafe,
  IconLock,
  IconUnlock,
  IconDelete,
  IconDown,
} from '@arco-design/web-react/icon'
import type { SysUser } from '@/models/entity'

interface UserTableProps {
  loading: boolean
  userList: SysUser[]
  total: number
  pageNum: number
  pageSize: number
  onEdit: (user: SysUser) => void
  onView: (user: SysUser) => void
  onResetPassword: (user: SysUser) => void
  onBlock: (user: SysUser) => void
  onUnblock: (user: SysUser) => void
  onDelete: (user: SysUser) => void
  onSizeChange: (size: number) => void
  onCurrentChange: (page: number) => void
}

const UserTable = ({
  loading,
  userList,
  total,
  pageNum,
  pageSize,
  onEdit,
  onView,
  onResetPassword,
  onBlock,
  onUnblock,
  onDelete,
  onSizeChange,
  onCurrentChange,
}: UserTableProps) => {
  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80, align: 'center' as const },
    { title: '用户名', dataIndex: 'username', width: 120, ellipsis: true, tooltip: true },
    { title: '手机号', dataIndex: 'phone', width: 140 },
    { title: '邮箱', dataIndex: 'email', ellipsis: true, tooltip: true },
    {
      title: '角色',
      dataIndex: 'roles',
      width: 140,
      render: (record: SysUser) => (
        <Tag color={record.roles.includes('ROLE_ADMIN') ? 'red' : 'gray'}>
          {record.roles.includes('ROLE_ADMIN') ? '管理员' : '普通用户'}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      align: 'center' as const,
      render: (record: SysUser) => {
        if (record.status === 'NORMAL') return <Tag color="green">正常</Tag>
        if (record.status === 'LOCKED') return <Tag color="orange">锁定</Tag>
        return <Tag color="red">封禁</Tag>
      },
    },
    {
      title: '最后登录时间',
      dataIndex: 'lastLoginAt',
      width: 180,
      render: (dateString: string) => {
        if (!dateString) return '-'
        return new Date(dateString).toLocaleString('zh-CN')
      },
    },
    {
      title: '操作',
      width: 220,
      fixed: 'right' as const,
      render: (record: SysUser) => {
        const dropDownList = (
          <>
            <Dropdown.Item key="view" icon={<IconEye />}>查看详情</Dropdown.Item>
            <Dropdown.Item key="resetPwd" icon={<IconSafe />}>重置密码</Dropdown.Item>
            {record.status !== 'BANNED' ? (
              <Dropdown.Item key="block" icon={<IconLock />}>封禁用户</Dropdown.Item>
            ) : (
              <Dropdown.Item key="unblock" icon={<IconUnlock />}>解封用户</Dropdown.Item>
            )}
            <Dropdown.Item key="delete" style={{ marginTop: 4 }} icon={<IconDelete />}>删除用户</Dropdown.Item>
          </>
        )

        return (
          <Button.Group>
            <Button size="small" type="primary" icon={<IconEdit />} onClick={() => onEdit(record)}>
              编辑
            </Button>
            <Dropdown droplist={dropDownList} trigger="hover" onSelect={(key: string) => {
              switch (key) {
                case 'view': onView(record); break
                case 'resetPwd': onResetPassword(record); break
                case 'block': onBlock(record); break
                case 'unblock': onUnblock(record); break
                case 'delete': onDelete(record); break
              }
            }}>
              <Button size="small" secondary icon={<IconDown />}>
                更多
              </Button>
            </Dropdown>
          </Button.Group>
        )
      },
    },
  ]

  return (
    <div className="table-card">
      <Table
        columns={columns}
        data={userList}
        loading={loading}
        border
        stripe
        pagination={false}
        scroll={{ x: 1200 }}
      />
      <div className="pagination">
        <Pagination
          current={pageNum}
          pageSize={pageSize}
          total={total}
          pageSizeOptions={[10, 20, 50, 100]}
          showTotal
          showJumper
          onChange={(current, size) => {
            if (current !== pageNum) {
              onCurrentChange(current)
            }
            if (size !== pageSize) {
              onSizeChange(size)
            }
          }}
        />
      </div>
    </div>
  )
}

export default UserTable
