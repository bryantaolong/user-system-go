<template>
  <div class="user-profile">
    <a-card class="profile-header">
      <div class="profile-main">
        <div class="profile-avatar">
          <a-upload
              :show-file-list="false"
              :custom-request="handleUploadAvatar"
              :before-upload="beforeAvatarUpload"
          >
            <a-avatar
                v-if="userStore.userProfile?.avatar"
                :size="120"
                :src="getAvatarUrl(userStore.userProfile.avatar)"
            />
            <a-avatar v-else :size="120">
              {{ userStore.userInfo?.username?.charAt(0).toUpperCase() }}
            </a-avatar>
            <div class="avatar-overlay">
              <IconCamera :size="24" />
              <p>点击更换</p>
            </div>
          </a-upload>
        </div>

        <div class="profile-info">
          <div class="profile-basic">
            <h2 class="profile-username">{{ userStore.userInfo?.username }}</h2>
          </div>
        </div>
      </div>
    </a-card>

    <a-card class="main-content-card">
      <a-tabs v-model:active-key="activeMainTab" class="main-tabs">

        <a-tab-pane key="settings" title="设置">
          <div class="settings-container">
            <a-tabs v-model:active-key="editActiveTab" tab-position="left" class="settings-tabs">
              <a-tab-pane key="basic" title="基本信息">
                <div class="settings-content">
                  <div class="settings-header">
                    <h3>基本信息</h3>
                    <p>管理您的个人信息，包括姓名、联系方式等</p>
                  </div>
                  <BasicInfo
                    :username="userStore.userInfo?.username"
                    :initial-data="basicForm"
                    :loading="updating"
                    @save="handleUpdateBasic"
                  />
                </div>
              </a-tab-pane>

              <a-tab-pane key="security" title="账号安全">
                <div class="settings-content">
                  <div class="settings-header">
                    <h3>账号安全</h3>
                    <p>保护您的账号安全，修改密码或进行账号注销</p>
                  </div>
                  <SecuritySettings
                    ref="securitySettingsRef"
                    :loading="changingPassword"
                    @change-password="handleChangePassword"
                    @delete-account="handleDeleteAccount"
                  />
                </div>
              </a-tab-pane>

              <a-tab-pane key="login-history" title="登录历史">
                <div class="settings-content">
                  <div class="settings-header">
                    <h3>登录历史</h3>
                    <p>查看您最近的账号登录活动</p>
                  </div>
                  <LoginHistory :history="loginHistory" />
                </div>
              </a-tab-pane>
            </a-tabs>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ArcoMessage, ArcoMessageBox } from '@/utils/arco-message'
import { IconCamera } from '@arco-design/web-vue/es/icon'
import { useUserStore } from '@/stores/user'
import * as userApi from '@/api/user/user'
import * as userProfileApi from '@/api/user/userProfile'
import { getAvatarUrl } from '@/utils/file'
import { getLocationFromIp } from '@/utils/ipLocation'
import BasicInfo from '@/components/profile/BasicInfo.vue'
import SecuritySettings from '@/components/profile/SecuritySettings.vue'
import LoginHistory from '@/components/profile/LoginHistory.vue'

/* --- 工具方法 --- */
function genderToNum(g?: string): 1 | 0 {
  return g === 'FEMALE' ? 0 : 1
}

function numToGender(n: 1 | 0): 'MALE' | 'FEMALE' {
  return n === 0 ? 'FEMALE' : 'MALE'
}
/* --- 基础状态 --- */
const userStore = useUserStore()
const activeMainTab = ref('settings')
const editActiveTab = ref('basic')
const updating = ref(false)
const changingPassword = ref(false)
const securitySettingsRef = ref()

const basicForm = reactive({realName: '', gender: 1 as 1 | 0, birthday: '', phone: '', email: ''})

/* --- 操作方法 --- */
const handleUpdateBasic = async (formData: any) => {
  updating.value = true
  try {
    await userStore.updateProfile({
      realName: formData.realName,
      gender: numToGender(formData.gender),
      birthday: formData.birthday ? formData.birthday + 'T00:00:00' : undefined,
      avatar: userStore.userProfile?.avatar
    })
    if (userStore.userInfo?.id) {
      await userApi.updateUser(userStore.userInfo.id, {phone: formData.phone, email: formData.email})
    }
    await userStore.fetchUserInfo()
    ArcoMessage.success('更新成功')
  } catch (e) {
    ArcoMessage.error('更新失败')
  } finally {
    updating.value = false
  }
}

const handleChangePassword = async (pwdData: any) => {
  changingPassword.value = true
  try {
    const res = await userStore.changePassword(pwdData.oldPassword, pwdData.newPassword)
    if (res.success) {
      ArcoMessage.success('密码修改成功')
      securitySettingsRef.value?.resetPasswordForm()
    } else ArcoMessage.error(res.message)
  } finally {
    changingPassword.value = false
  }
}

const handleDeleteAccount = async () => {
  try {
    await ArcoMessageBox.confirm('确定注销账号吗？这是不可逆的操作！', '警告')
    const { value } = await ArcoMessageBox.prompt('请输入 "DELETE" 确认', '二次确认', {
      confirmText: '确定',
      cancelText: '取消',
      placeholder: '请输入 DELETE'
    })
    if (value === 'DELETE') {
      const res = await userStore.deleteAccount()
      if (res.success) ArcoMessage.success('注销成功')
    }
  } catch (e) {
  }
}

const handleUploadAvatar = async (options: any) => {
  try {
    const res = await userProfileApi.uploadAvatar(options.file)
    if (res.code === 200) {
      ArcoMessage.success('头像上传成功')
      if (userStore.userProfile) {
        userStore.userProfile.avatar = res.data
      }
    } else {
      ArcoMessage.error(res.message || '上传失败')
    }
  } catch (e: any) {
    ArcoMessage.error(e.message || '上传失败')
  }
}

const beforeAvatarUpload = (file: File) => {
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isLt2M) ArcoMessage.error('大小不能超过 2MB!')
  return isLt2M
}

const loginHistory = ref<any[]>([])

const loadLoginHistory = async () => {
  if (userStore.userInfo?.lastLoginAt) {
    try {
      const ipAddress = userStore.userInfo.lastLoginIp || 'Unknown'
      const location = ipAddress !== 'Unknown' 
        ? await getLocationFromIp(ipAddress) 
        : 'Unknown'
      
      loginHistory.value = [{
        loginTime: userStore.userInfo.lastLoginAt.replace('T', ' ').substring(0, 19),
        ipAddress: ipAddress,
        location: location,
        device: userStore.userInfo.lastLoginDevice || 'Unknown'
      }]
    } catch (error) {
      console.error('Failed to load login history:', error)
      loginHistory.value = [{
        loginTime: userStore.userInfo.lastLoginAt.replace('T', ' ').substring(0, 19),
        ipAddress: 'Unknown',
        location: 'Unknown',
        device: 'Unknown'
      }]
    }
  } else {
    loginHistory.value = []
  }
}

onMounted(() => {
  if (userStore.userProfile) {
    Object.assign(basicForm, {
      realName: userStore.userProfile.realName || '',
      gender: genderToNum(userStore.userProfile.gender),
      birthday: userStore.userProfile.birthday?.slice(0, 10) || '',
      phone: userStore.userProfile.phone || '',
      email: userStore.userProfile.email || ''
    })
  }
  loadLoginHistory()
})
</script>

<style scoped>
.user-profile {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.profile-header, .main-content-card {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.profile-main {
  display: flex;
  align-items: center;
  gap: 40px;
  padding: 20px;
}

.profile-avatar {
  position: relative;
  cursor: pointer;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  opacity: 0;
  transition: 0.3s;
}

.profile-avatar:hover .avatar-overlay {
  opacity: 1;
}

.profile-info {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.profile-username {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 15px;
}

.profile-stats {
  display: flex;
  gap: 30px;
}

.stat-item {
  text-align: center;
  cursor: pointer;
}

.stat-number {
  display: block;
  font-size: 18px;
  font-weight: 600;
}

.stat-label {
  font-size: 13px;
  color: #909399;
}

.main-tabs {
  padding: 0 10px;
}

.settings-container {
  min-height: 400px;
  padding: 20px 0;
}

.settings-tabs {
  height: 100%;
}

:deep(.arco-tabs-tab-left) {
  margin-right: 30px;
  width: 160px;
  border-right: 1px solid var(--color-border-2);
}

:deep(.arco-tabs-tab-title) {
  text-align: left;
  height: 45px;
  line-height: 45px;
  font-size: 15px;
  padding: 0 20px;
}

.settings-content {
  padding-left: 10px;
}

.settings-header {
  margin-bottom: 25px;
  padding-bottom: 15px;
  border-bottom: 1px solid var(--color-border-2);
}

.settings-header h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-1);
}

.settings-header p {
  margin: 0;
  font-size: 14px;
  color: var(--color-text-3);
}
</style>
