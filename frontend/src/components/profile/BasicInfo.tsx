import { useState, useEffect } from 'react'
import { Form, Input, Radio, DatePicker, Button } from '@arco-design/web-react'
import message from '@arco-design/web-react/es/Message'
import dayjs from 'dayjs'

interface BasicInfoProps {
  username?: string
  initialData: {
    realName: string
    gender: 1 | 0
    birthday: string
    phone: string
    email: string
  }
  loading: boolean
  onSave: (data: any) => void
}

const BasicInfo = ({ username, initialData, loading, onSave }: BasicInfoProps) => {
  const [form] = Form.useForm()

  useEffect(() => {
    form.setFieldsValue({
      ...initialData,
      birthday: initialData.birthday ? dayjs(initialData.birthday) : undefined,
    })
  }, [initialData, form])

  const handleSave = async () => {
    try {
      const values = await form.validate()
      const birthdayStr = values.birthday ? values.birthday.format('YYYY-MM-DD') : ''
      onSave({
        realName: values.realName,
        gender: values.gender,
        birthday: birthdayStr,
        phone: values.phone,
        email: values.email,
      })
    } catch (error) {
      // 验证失败
    }
  }

  return (
    <Form form={form} className="info-form" style={{ maxWidth: 500 }}>
      <Form.Item label="用户名">
        <Input value={username} disabled />
      </Form.Item>
      <Form.Item label="真实姓名" field="realName" rules={[
        { minLength: 2, maxLength: 20, message: '真实姓名长度应在2-20个字符之间' },
      ]}>
        <Input placeholder="请输入真实姓名" />
      </Form.Item>
      <Form.Item label="性别" field="gender">
        <Radio.Group>
          <Radio value={1}>男</Radio>
          <Radio value={0}>女</Radio>
        </Radio.Group>
      </Form.Item>
      <Form.Item label="生日" field="birthday">
        <DatePicker format="YYYY-MM-DD" />
      </Form.Item>
      <Form.Item label="手机号" field="phone" rules={[
        { match: /^1[3-9]\d{9}$/, message: '电话号码格式不正确' },
      ]}>
        <Input placeholder="请输入手机号" />
      </Form.Item>
      <Form.Item label="邮箱" field="email" rules={[
        { type: 'email', message: '邮箱格式不正确' },
      ]}>
        <Input placeholder="请输入邮箱" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" loading={loading} onClick={handleSave}>
          保存修改
        </Button>
      </Form.Item>
    </Form>
  )
}

export default BasicInfo
