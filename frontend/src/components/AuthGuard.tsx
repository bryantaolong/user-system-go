import { useEffect } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import message from '@arco-design/web-react/es/Message'
import { useUserStore } from '@/stores/user'

const AuthGuard = ({ children, requiresAdmin = false }: { children: React.ReactNode; requiresAdmin?: boolean }) => {
  const navigate = useNavigate()
  const location = useLocation()
  const { token, userInfo, isAuthenticated, isAdmin, fetchUserInfo, logout } = useUserStore()

  useEffect(() => {
    const checkAuth = async () => {
      // 游客页面访问控制
      const guestPaths = ['/login', '/register']
      if (guestPaths.includes(location.pathname) && isAuthenticated) {
        navigate('/')
        return
      }

      // 需要认证的页面
      if (requiresAdmin || isAuthenticated) {
        if (!userInfo && token) {
          try {
            const res = await fetchUserInfo()
            const currentUserInfo = useUserStore.getState().userInfo
            if (!res.success || !currentUserInfo) {
              message.error('认证信息失效，请重新登录！')
              logout()
              navigate('/login')
              return
            }
          } catch (error) {
            message.error('网络错误或认证失败，请重新登录！')
            logout()
            navigate('/login')
            return
          }
        } else if (!token) {
          message.error('您尚未登录，请先登录。')
          navigate('/login')
          return
        }
      }

      if (requiresAdmin && !isAdmin) {
        message.error('您没有权限访问此页面！')
        navigate('/')
      }
    }

    checkAuth()
  }, [requiresAdmin, isAuthenticated, isAdmin, token, userInfo, fetchUserInfo, logout, navigate, location.pathname])

  if (!token && !isAuthenticated) {
    return null
  }

  if (requiresAdmin && !isAdmin) {
    return null
  }

  if (!userInfo && token) {
    return null
  }

  return <>{children || <Outlet />}</>
}

export default AuthGuard
