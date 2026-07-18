import { useState, useEffect, useCallback } from 'react'
import { Card, Tabs, Modal } from '@arco-design/web-react'
import message from '@arco-design/web-react/es/Message'
import { IconCamera } from '@arco-design/web-react/icon'
import { useUserStore } from '@/stores/user'
import * as userApi from '@/api/user/user'
import * as userProfileApi from '@/api/user/userProfile'
import { getAvatarUrl } from '@/utils/file'
import { getLocationFromIp } from '@/utils/ipLocation'
import BasicInfo from '@/components/profile/BasicInfo'
import SecuritySettings from '@/components/profile/SecuritySettings'
import LoginHistory from '@/components/profile/LoginHistory'

function genderToNum(g?: string): 1 | 0 {
  return g === 'FEMALE' ? 0 : 1
}

function numToGender(n: 1 | 0): 'MALE' | 'FEMALE' {
  return n === 0 ? 'FEMALE' : 'MALE'
}

const UserProfile = () => {
  const userStore = useUserStore()
  const [activeMainTab, setActiveMainTab] = useState('settings')
  const [editActiveTab, setEditActiveTab] = useState('basic')
  const [updating, setUpdating] = useState(false)
  const [changingPassword, setChangingPassword] = useState(false)

  const [basicForm, setBasicForm] = useState({
    realName: '',
    gender: 1 as 1 | 0,
    birthday: '',
    phone: '',
    email: '',
  })

  const handleUpdateBasic = async (formData: any) => {
    setUpdating(true)
    try {
      await userStore.updateProfile({
        realName: formData.realName,
        gender: numToGender(formData.gender),
        birthday: formData.birthday ? formData.birthday + 'T00:00:00' : undefined,
        avatar: userStore.userProfile?.avatar,
      })
      if (userStore.userInfo?.id) {
        await userApi.updateUser(userStore.userInfo.id, { phone: formData.phone, email: formData.email })
      }
      await userStore.fetchUserInfo()
      message.success('更新成功')
    } catch (e) {
      message.error('更新失败')
    } finally {
      setUpdating(false)
    }
  }

  const handleChangePassword = async (pwdData: any) => {
    setChangingPassword(true)
    try {
      const res = await userStore.changePassword(pwdData.oldPassword, pwdData.newPassword)
      if (res.success) {
        message.success('密码修改成功')
      } else {
        message.error(res.message)
      }
    } finally {
      setChangingPassword(false)
    }
  }

  const handleDeleteAccount = async () => {
    try {
      await Modal.confirm('确定注销账号吗？这是不可逆的操作！', '警告')
      const { value } = await Modal.prompt('请输入 "DELETE" 确认', '二次确认', {
        confirmText: '确定',
        cancelText: '取消',
        placeholder: '请输入 DELETE',
      })
      if (value === 'DELETE') {
        const res = await userStore.deleteAccount()
        if (res.success) {
          message.success('注销成功')
        }
      }
    } catch (e) {
      // 取消或错误
    }
  }

  const handleUploadAvatar = async (options: any) => {
    try {
      const res = await userProfileApi.uploadAvatar(options.file)
      if (res.code === 200) {
        message.success('头像上传成功')
        if (userStore.userProfile) {
          userStore.userProfile.avatar = res.data
        }
      } else {
        message.error(res.message || '上传失败')
      }
    } catch (e: any) {
      message.error(e.message || '上传失败')
    }
  }

  const beforeAvatarUpload = (file: File) => {
    const isLt2M = file.size / 1024 / 1024 < 2
    if (!isLt2M) message.error('大小不能超过 2MB!')
    return isLt2M
  }

  const [loginHistory, setLoginHistory] = useState<
    { loginTime: string; ipAddress: string; location: string; device: string }[]
  >([])

  const loadLoginHistory = useCallback(async () => {
    if (userStore.userInfo?.lastLoginAt) {
      try {
        const ipAddress = userStore.userInfo.lastLoginIp || 'Unknown'
        const location = ipAddress !== 'Unknown' ? await getLocationFromIp(ipAddress) : 'Unknown'

        setLoginHistory([
          {
            loginTime: userStore.userInfo.lastLoginAt.replace('T', ' ').substring(0, 19),
            ipAddress: ipAddress,
            location: location,
            device: userStore.userInfo.lastLoginDevice || 'Unknown',
          },
        ])
      } catch (error) {
        console.error('Failed to load login history:', error)
        setLoginHistory([
          {
            loginTime: userStore.userInfo.lastLoginAt.replace('T', ' ').substring(0, 19),
            ipAddress: 'Unknown',
            location: 'Unknown',
            device: 'Unknown',
          },
        ])
      }
    } else {
      setLoginHistory([])
    }
  }, [userStore.userInfo])

  useEffect(() => {
    if (userStore.userProfile) {
      setBasicForm({
        realName: userStore.userProfile.realName || '',
        gender: genderToNum(userStore.userProfile.gender),
        birthday: userStore.userProfile.birthday?.slice(0, 10) || '',
        phone: userStore.userProfile.phone || '',
        email: userStore.userProfile.email || '',
      })
    }
    loadLoginHistory()
  }, [userStore.userProfile, loadLoginHistory])

  return (
    <div className="user-profile">
      <Card className="profile-header">
        <div className="profile-main">
          <div className="profile-avatar">
            {userStore.userProfile?.avatar ? (
              <img
                src={getAvatarUrl(userStore.userProfile.avatar)}
                alt="avatar"
                style={{ width: 120, height: 120, borderRadius: '50%', objectFit: 'cover' }}
              />
            ) : (
              <div
                style={{
                  width: 120,
                  height: 120,
                  borderRadius: '50%',
                  backgroundColor: '#165DFF',
                  color: '#fff',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: 48,
                  fontWeight: 600,
                }}
              >
                {userStore.userInfo?.username?.charAt(0).toUpperCase()}
              </div>
            )}
          </div>
          <div className="profile-info">
            <div className="profile-basic">
              <h2 className="profile-username">{userStore.userInfo?.username}</h2>
            </div>
          </div>
        </div>
      </Card>

      <Card className="main-content-card">
        <Tabs activeTab={activeMainTab} onChange={setActiveMainTab} className="main-tabs">
          <Tabs.TabPane key="settings" title="设置">
            <div className="settings-container">
              <Tabs activeTab={editActiveTab} tabPosition="left" className="settings-tabs">
                <Tabs.TabPane key="basic" title="基本信息">
                  <div className="settings-content">
                    <div className="settings-header">
                      <h3>基本信息</h3>
                      <p>管理您的个人信息，包括姓名、联系方式等</p>
                    </div>
                    <BasicInfo
                      username={userStore.userInfo?.username}
                      initialData={basicForm}
                      loading={updating}
                      onSave={handleUpdateBasic}
                    />
                  </div>
                </Tabs.TabPane>

                <Tabs.TabPane key="security" title="账号安全">
                  <div className="settings-content">
                    <div className="settings-header">
                      <h3>账号安全</h3>
                      <p>保护您的账号安全，修改密码或进行账号注销</p>
                    </div>
                    <SecuritySettings
                      loading={changingPassword}
                      onChangePassword={handleChangePassword}
                      onDeleteAccount={handleDeleteAccount}
                    />
                  </div>
                </Tabs.TabPane>

                <Tabs.TabPane key="login-history" title="登录历史">
                  <div className="settings-content">
                    <div className="settings-header">
                      <h3>登录历史</h3>
                      <p>查看您最近的账号登录活动</p>
                    </div>
                    <LoginHistory history={loginHistory} />
                  </div>
                </Tabs.TabPane>
              </Tabs>
            </div>
          </Tabs.TabPane>
        </Tabs>
      </Card>
    </div>
  )
}

export default UserProfile
