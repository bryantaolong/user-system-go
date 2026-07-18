import { useState } from 'react'
import { Form, Input, Button, Select } from '@arco-design/web-react'
import { IconSearch, IconRefresh } from '@arco-design/web-react/icon'

export type UserStatus = 'NORMAL' | 'LOCKED' | 'BANNED' | ''

export interface UserSearchFormData {
  username: string
  phone: string
  email: string
  status: UserStatus
}

const UserSearchForm = ({ onSearch, onReset }: { onSearch: (data: UserSearchFormData) => void; onReset: () => void }) => {
  const [form] = Form.useForm()

  const handleSearch = () => {
    const values = form.getFieldsValue() as UserSearchFormData
    onSearch(values)
  }

  const handleResetClick = () => {
    form.resetFields()
    onReset()
  }

  return (
    <div className="search-card">
      <Form form={form} layout="inline" className="search-form">
        <Form.Item label="用户名" field="username">
          <Input placeholder="请输入用户名" allowClear />
        </Form.Item>

        <Form.Item label="手机号" field="phone">
          <Input placeholder="请输入手机号" allowClear />
        </Form.Item>

        <Form.Item label="邮箱" field="email">
          <Input placeholder="请输入邮箱" allowClear />
        </Form.Item>

        <Form.Item label="状态" field="status">
          <Select placeholder="请选择状态" allowClear style={{ width: 120 }}>
            <Select.Option label="正常" value="NORMAL" />
            <Select.Option label="锁定" value="LOCKED" />
            <Select.Option label="封禁" value="BANNED" />
          </Select>
        </Form.Item>

        <Form.Item>
          <Button type="primary" icon={<IconSearch />} onClick={handleSearch}>
            查询
          </Button>
          <Button icon={<IconRefresh />} onClick={handleResetClick}>
            重置
          </Button>
        </Form.Item>
      </Form>
    </div>
  )
}

export default UserSearchForm
