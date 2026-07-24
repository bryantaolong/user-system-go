import { Modal, Descriptions, Tag } from '@arco-design/web-react'
import type { SysUser } from '@/models/entity'

const UserDetailDialog = ({ visible, user, onClose }: { visible: boolean; user: SysUser | null; onClose: () => void }) => {
  const formatDateTime = (dateString: string | undefined | null) => {
    if (!dateString) return '-'
    return new Date(dateString).toLocaleString('zh-CN')
  }

  return (
    <Modal title="用户详情" visible={visible} width={700} onCancel={onClose}>
      {user && (
        <Descriptions column={2} bordered>
          <Descriptions.Item label="用户ID">{user.id}</Descriptions.Item>
          <Descriptions.Item label="用户名">{user.username}</Descriptions.Item>
          <Descriptions.Item label="手机号">{user.phone || '-'}</Descriptions.Item>
          <Descriptions.Item label="邮箱">{user.email || '-'}</Descriptions.Item>
          <Descriptions.Item label="状态">
            {user.status === 'NORMAL' && <Tag color="green">正常</Tag>}
            {user.status === 'LOCKED' && <Tag color="orange">锁定</Tag>}
            {user.status === 'BANNED' && <Tag color="red">封禁</Tag>}
          </Descriptions.Item>
          <Descriptions.Item label="角色">
            <Tag color={user.roles.includes('ROLE_ADMIN') ? 'red' : 'gray'}>
              {user.roles.includes('ROLE_ADMIN') ? '管理员' : '普通用户'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">{formatDateTime(user.createdAt)}</Descriptions.Item>
          <Descriptions.Item label="更新时间">{formatDateTime(user.updatedAt)}</Descriptions.Item>
          <Descriptions.Item label="最后登录时间" span={2}>
            {formatDateTime(user.lastLoginAt) || '-'}
          </Descriptions.Item>
        </Descriptions>
      )}
    </Modal>
  )
}

export default UserDetailDialog
