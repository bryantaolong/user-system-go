import { Navigate } from 'react-router-dom'
import AdminLayout from '@/layouts/AdminLayout'
import Login from '@/pages/auth/Login'
import Register from '@/pages/auth/Register'
import UserManagement from '@/pages/admin/UserManagement'
import UserProfile from '@/pages/profile/UserProfile'
import SystemLog from '@/pages/admin/SystemLog'
import NotFound from '@/pages/NotFound'
import AuthGuard from '@/components/AuthGuard'

const routes = [
  {
    path: '/',
    element: <Navigate to="/profile" replace />,
  },
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/register',
    element: <Register />,
  },
  {
    path: '/admin',
    element: (
      <AuthGuard requiresAdmin>
        <AdminLayout />
      </AuthGuard>
    ),
    children: [
      {
        path: 'users',
        element: <UserManagement />,
      },
      {
        path: 'profile',
        element: <UserProfile />,
      },
      {
        path: 'logs',
        element: <SystemLog />,
      },
    ],
  },
  {
    path: '/profile',
    element: (
      <AuthGuard>
        <UserProfile />
      </AuthGuard>
    ),
  },
  {
    path: '/*',
    element: <NotFound />,
  },
]

export default routes
