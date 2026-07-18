import { useState, useEffect } from 'react'
import { Form, Input, Button, Divider, Alert } from '@arco-design/web-react'

interface SecuritySettingsProps {
  loading: boolean
  onChangePassword: (data: { oldPassword: string; newPassword: string }) => void
  onDeleteAccount: () => void
}

const SecuritySettings = ({ loading, onChangePassword, onDeleteAccount }: SecuritySettingsProps) => {
  const [form] = Form.useForm()

  useEffect(() => {
    return () => {
      form.resetFields()
    }
  }, [])

  const handlePasswordChange = async () => {
    try {
      const values = await form.validate()
      onChangePassword({
        oldPassword: values.oldPassword,
        newPassword: values.newPassword,
      })
      form.resetFields()
    } catch (error) {
      // 验证失败
    }
  }

  return (
    <div className="security-section">
      <h3>修改密码</h3>
      <Form form={form} layout="vertical" labelWidth={120}>
        <Form.Item label="当前密码" field="oldPassword" rules={[
          { required: true, message: '请输入当前密码' },
        ]}>
          <Input.Password />
        </Form.Item>
        <Form.Item label="新密码" field="newPassword" rules={[
          { required: true, message: '请输入新密码' },
          { minLength: 6, message: '至少6位' },
        ]}>
          <Input.Password />
        </Form.Item>
        <Form.Item label="确认新密码" field="confirmPassword" dependencies={['newPassword']} rules={[
          { required: true, message: '请确认新密码' },
          ({ getFieldValue }) => ({
            validator: (_, value: string) => {
              if (!value || getFieldValue('newPassword') === value) {
                return Promise.resolve()
              }
              return Promise.reject('两次输入不一致')
            },
          }),
        ]}>
          <Input.Password />
        </Form.Item>
        <Form.Item>
          <Button type="primary" loading={loading} onClick={handlePasswordChange}>
            修改密码
          </Button>
        </Form.Item>
      </Form>

      <Divider />

      <div className="danger-section">
        <h3>危险操作</h3>
        <Alert title="注销账号是不可逆的操作，请谨慎操作！" type="error" closable={false} />
        <Button type="primary" status="danger" style={{ marginTop: 16 }} onClick={onDeleteAccount}>
          注销账号
        </Button>
      </div>
    </div>
  )
}

export default SecuritySettings
