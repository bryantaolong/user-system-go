import { useState } from 'react'
import { Form, Input, Button } from '@arco-design/web-react'
import message from '@arco-design/web-react/es/Message'
import { IconUser, IconLock, IconPhone, IconMessage } from '@arco-design/web-react/icon'
import { useNavigate } from 'react-router-dom'
import { useUserStore } from '@/stores/user'

const Register = () => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { register } = useUserStore()

  const validateConfirmPassword = (value: string, callback: (error?: string) => void) => {
    if (value !== form.getFieldValue('password')) {
      callback('两次输入的密码不一致')
    } else {
      callback()
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validate()
      setLoading(true)
      const result = await register({
        username: values.username,
        password: values.password,
        phone: values.phone || undefined,
        email: values.email || undefined,
      })

      if (result.success) {
        message.success('注册成功！正在跳转到登录页面...')
        setTimeout(() => {
          navigate('/login')
        }, 1500)
      } else {
        message.error(result.message || '注册失败')
      }
    } catch (error) {
      message.error('注册失败，请稍后重试')
      console.error('Register error:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="register-container">
      <div className="register-card">
        <div className="register-header">
          <div className="register-title">用户注册</div>
          <div className="register-subtitle">创建您的账户</div>
        </div>

        <Form form={form} layout="vertical" className="register-form" size="large" onSubmit={handleSubmit}>
          <Form.Item label="用户名" field="username" rules={[
            { required: true, message: '请输入用户名' },
            { minLength: 2, maxLength: 20, message: '用户名长度应在2-20个字符之间' },
          ]}>
            <Input prefix={<IconUser />} placeholder="请输入用户名" allowClear />
          </Form.Item>

          <Form.Item label="密码" field="password" rules={[
            { required: true, message: '请输入密码' },
            { minLength: 6, message: '密码至少6位' },
          ]}>
            <Input.Password prefix={<IconLock />} placeholder="请输入密码" allowClear />
          </Form.Item>

          <Form.Item label="确认密码" field="confirmPassword" rules={[
            { required: true, message: '请确认密码' },
            { validator: validateConfirmPassword },
          ]}>
            <Input.Password prefix={<IconLock />} placeholder="请确认密码" allowClear />
          </Form.Item>

          <Form.Item label="手机号" field="phone" rules={[
            { match: /^1[3-9]\d{9}$/, message: '电话号码格式不正确' },
          ]}>
            <Input prefix={<IconPhone />} placeholder="请输入手机号码（可选）" allowClear />
          </Form.Item>

          <Form.Item label="邮箱" field="email" rules={[
            { type: 'email', message: '邮箱格式不正确' },
          ]}>
            <Input prefix={<IconMessage />} placeholder="请输入邮箱地址（可选）" allowClear />
          </Form.Item>

          <Form.Item>
            <Button type="primary" size="large" loading={loading} htmlType="submit" className="register-button">
              注册
            </Button>
          </Form.Item>
        </Form>

        <div className="register-footer">
          <span>已有账号？</span>
          <Button type="primary" onClick={() => navigate('/login')}>
            立即登录
          </Button>
        </div>
      </div>

      <div className="register-background">
        <div className="background-circle circle1"></div>
        <div className="background-circle circle2"></div>
        <div className="background-circle circle3"></div>
      </div>
    </div>
  )
}

export default Register
