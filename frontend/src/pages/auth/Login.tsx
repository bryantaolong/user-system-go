import { useState, useEffect } from 'react'
import { Form, Input, Button, Checkbox } from '@arco-design/web-react'
import message from '@arco-design/web-react/es/Message'
import { IconUser, IconLock } from '@arco-design/web-react/icon'
import { useNavigate } from 'react-router-dom'
import { useUserStore } from '@/stores/user'

const Login = () => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { login } = useUserStore()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      navigate('/')
    }
  }, [navigate])

  const handleSubmit = async () => {
    try {
      const values = await form.validate()
      setLoading(true)
      const result = await login(values.username, values.password)
      if (result.success) {
        message.success('登录成功！')
        navigate('/')
      } else {
        message.error(result.message || '登录失败')
      }
    } catch (error) {
      message.error('登录失败，请稍后重试')
      console.error('Login error:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-container">
      <div className="login-card">
        <div className="login-header">
          <div className="login-title">用户登录</div>
          <div className="login-subtitle">请输入您的账户信息</div>
        </div>

        <Form form={form} layout="vertical" className="login-form" size="large" onSubmit={handleSubmit}>
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

          <Form.Item>
            <Checkbox>记住我</Checkbox>
          </Form.Item>

          <Form.Item>
            <Button type="primary" size="large" loading={loading} htmlType="submit" className="login-button">
              登录
            </Button>
          </Form.Item>
        </Form>

        <div className="login-footer">
          <span>还没有账号？</span>
          <Button type="primary" onClick={() => navigate('/register')}>
            立即注册
          </Button>
        </div>
      </div>

      <div className="login-background">
        <div className="background-circle circle1"></div>
        <div className="background-circle circle2"></div>
        <div className="background-circle circle3"></div>
      </div>
    </div>
  )
}

export default Login
