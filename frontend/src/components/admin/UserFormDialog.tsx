import { useState, useEffect, useRef } from 'react'
import { Modal, Form, Input, Select } from '@arco-design/web-react'
import message from '@arco-design/web-react/es/Message'
import * as userRoleApi from '@/api/user/userRole'
import type { UserRoleOptionVO } from '@shared/models/vo'
import type { UserUpdateRequest } from '@shared/models/request/user'

export interface UserFormData extends UserUpdateRequest {
  username?: string
  password?: string
  roleIds?: number[]
}

const UserFormDialog = ({
  visible,
  dialogType,
  userForm,
  submitting,
  onSubmit,
  onClose,
}: {
  visible: boolean
  dialogType: 'add' | 'edit'
  userForm: UserFormData
  submitting?: boolean
  onSubmit: (formData: UserFormData) => void
  onClose: () => void
}) => {
  const [form] = Form.useForm()
  const [localUserForm, setLocalUserForm] = useState<UserFormData>(userForm)
  const [roleOptions, setRoleOptions] = useState<UserRoleOptionVO[]>([])

  useEffect(() => {
    if (visible) {
      setLocalUserForm(userForm)
    }
  }, [visible, userForm])

  useEffect(() => {
    const loadRoles = async () => {
      if (roleOptions.length > 0) return
      const res = await userRoleApi.listRoles()
      if (res.code === 200) {
        setRoleOptions(res.data)
      }
    }
    if (visible) {
      loadRoles()
    }
  }, [visible])

  const handleSubmit = async () => {
    try {
      const values = await form.validate()
      onSubmit({ ...localUserForm, ...values })
    } catch (error) {
      // 验证失败
    }
  }

  const dialogTitle = dialogType === 'add' ? '新增用户' : '编辑用户'

  return (
    <Modal
      title={dialogTitle}
      visible={visible}
      width={600}
      onCancel={onClose}
      footer={
        <>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" loading={submitting} onClick={handleSubmit}>
            确定
          </Button>
        </>
      }
    >
      <Form
        form={form}
        layout="vertical"
        className="user-form"
        initialValues={localUserForm}
      >
        <Form.Item label="用户名" field="username" rules={dialogType === 'add' ? [
          { required: true, message: '请输入用户名' },
          { minLength: 2, maxLength: 20, message: '用户名长度应在2-20个字符之间' },
        ] : []}>
          <Input disabled={dialogType === 'edit'} />
        </Form.Item>
        <Form.Item label="手机号" field="phone" rules={[
          { match: /^1[3-9]\d{9}$/, message: '电话号码格式不正确' },
        ]}>
          <Input />
        </Form.Item>
        <Form.Item label="邮箱" field="email" rules={[
          { type: 'email', message: '邮箱格式不正确' },
        ]}>
          <Input />
        </Form.Item>
        {dialogType === 'add' && (
          <Form.Item label="密码" field="password" rules={[
            { required: true, message: '请输入密码' },
            { minLength: 6, message: '密码至少6位' },
          ]}>
            <Input.Password />
          </Form.Item>
        )}
        <Form.Item label="角色" field="roleIds" rules={[
          { required: true, message: '请选择角色' },
        ]}>
          <Select placeholder="请选择角色" mode="multiple">
            {roleOptions.map((r) => (
              <Select.Option key={r.id} value={r.id}>
                {r.roleName}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default UserFormDialog
