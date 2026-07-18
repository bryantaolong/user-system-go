import { Layout, Menu, Breadcrumb, Dropdown, Avatar } from '@arco-design/web-react'
import { IconDesktop, IconHome, IconUser, IconUserGroup, IconSettings, IconDown, IconExport } from '@arco-design/web-react/icon'
import { useNavigate, useLocation } from 'react-router-dom'
import { useUserStore } from '@/stores/user'
import { Modal } from '@arco-design/web-react'
import message from '@arco-design/web-react/es/Message'
import { getAvatarUrl } from '@/utils/file'
import { Outlet } from 'react-router-dom'

const { Sider, Header, Content } = Layout

const AdminLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { userInfo, userProfile, logout } = useUserStore()

  const activeMenu = location.pathname

  const handleMenuClick = (key: string) => {
    navigate(key)
  }

  const handleCommand = async (command: string) => {
    switch (command) {
      case 'profile':
        navigate('/admin/profile')
        break
      case 'logout':
        try {
          await Modal.confirm('确定要退出登录吗？', '提示')
          await logout()
          navigate('/login')
          message.success('已退出登录')
        } catch {
          // 取消退出
        }
        break
    }
  }

  const menuItems = [
    { key: '/', icon: <IconHome />, label: '返回首页' },
    { key: '/admin/users', icon: <IconUser />, label: '用户管理' },
    { key: '/admin/profile', icon: <IconUserGroup />, label: '个人中心' },
    { key: '/admin/settings', icon: <IconSettings />, label: '系统设置' },
  ]

  const dropDownList = [
    { key: 'profile', icon: <IconUser />, label: '个人中心' },
    { key: 'logout', icon: <IconExport />, label: '退出登录' },
  ]

  const pathSegments = location.pathname.split('/').filter(Boolean)

  return (
    <Layout className="admin-layout">
      <Sider width={240} className="layout-aside">
        <div className="aside-header">
          <div className="logo">
            <IconDesktop style={{ fontSize: 24 }} />
            <span className="logo-text">用户管理系统</span>
          </div>
        </div>
        <Menu selectedKeys={[activeMenu]} onClickMenuItem={handleMenuClick} className="aside-menu">
          {menuItems.map((item) => (
            <Menu.Item key={item.key} icon={item.icon}>
              {item.label}
            </Menu.Item>
          ))}
        </Menu>
      </Sider>
      <Layout className="layout-main">
        <Header className="layout-header">
          <div className="header-left">
            <Breadcrumb>
              {pathSegments.map((segment, index) => {
                const path = '/' + pathSegments.slice(0, index + 1).join('/')
                const title = segment === 'admin' ? '管理后台' : segment === 'users' ? '用户管理' : segment === 'profile' ? '个人中心' : segment === 'logs' ? '系统日志' : segment
                return (
                  <Breadcrumb.Item key={path} onClick={() => navigate(path)} style={{ cursor: 'pointer' }}>
                    {title}
                  </Breadcrumb.Item>
                )
              })}
            </Breadcrumb>
          </div>
          <div className="header-right">
            <Dropdown droplist={dropDownList} onSelect={handleCommand}>
              <div className="user-info">
                <Avatar size={32} style={{ marginRight: 8 }}>
                  {userInfo?.username?.charAt(0).toUpperCase()}
                </Avatar>
                <span className="username">{userInfo?.username}</span>
                <IconDown />
              </div>
            </Dropdown>
          </div>
        </Header>
        <Content className="layout-content">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default AdminLayout
