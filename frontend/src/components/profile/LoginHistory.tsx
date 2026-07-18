import { Table } from '@arco-design/web-react'

interface LoginHistoryProps {
  history: Array<{
    loginTime: string
    ipAddress: string
    location?: string
    device: string
  }>
}

const LoginHistory = ({ history }: LoginHistoryProps) => {
  const columns = [
    { title: '登录时间', dataIndex: 'loginTime', width: 180 },
    { title: 'IP地址', dataIndex: 'ipAddress', width: 140 },
    { title: '登录地点', dataIndex: 'location' },
    { title: '设备信息', dataIndex: 'device' },
  ]

  return <Table columns={columns} data={history} style={{ width: '100%' }} pagination={false} />
}

export default LoginHistory
