<template>
  <div class="profile-page">
    <a-card class="profile-header">
      <div class="profile-main">
        <div class="profile-avatar">
          <a-avatar :size="120" :src="avatarUrl" style="font-size: 48px; background: #165dff;">
            {{ store.user?.username?.charAt(0)?.toUpperCase() }}
          </a-avatar>
        </div>
        <div class="profile-info">
          <h2>{{ store.user?.username }}</h2>
        </div>
      </div>
    </a-card>

    <a-card class="main-content">
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="settings" title="设置">
          <a-tabs v-model:active-key="editTab" tab-position="left">
            <a-tab-pane key="basic" title="基本信息">
              <div class="settings-content">
                <div class="settings-header">
                  <h3>基本信息</h3>
                  <p>管理您的个人信息，包括姓名、联系方式等</p>
                </div>
                <basic-info
                  :username="store.user?.username"
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
                <security-settings
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
                <login-history :history="loginHistory" />
              </div>
            </a-tab-pane>
          </a-tabs>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { useUserStore } from '@/stores/user'
import * as userApi from '@/api/user'
import * as userProfileApi from '@/api/userProfile'
import BasicInfo from './BasicInfo.vue'
import SecuritySettings from './SecuritySettings.vue'
import LoginHistory from './LoginHistory.vue'

const store = useUserStore()
const activeTab = ref('settings')
const editTab = ref('basic')
const updating = ref(false)
const changingPassword = ref(false)

const basicForm = reactive({
  realName: '',
  gender: 1 as 0 | 1,
  birthday: '',
  phone: '',
  email: '',
})

const avatarUrl = computed(() => {
  const avatar = store.profile?.avatar
  if (!avatar) return ''
  if (avatar.startsWith('http')) return avatar
  return `/uploads/${avatar}`
})

const loginHistory = ref<{ loginTime: string; ipAddress: string; location: string; device: string }[]>([])

function loadLoginHistory() {
  if (store.user?.lastLoginAt) {
    loginHistory.value = [
      {
        loginTime: store.user.lastLoginAt.replace('T', ' ').substring(0, 19),
        ipAddress: store.user.lastLoginIp || 'Unknown',
        location: 'Unknown',
        device: store.user.lastLoginDevice || 'Unknown',
      },
    ]
  } else {
    loginHistory.value = []
  }
}

watch(
  () => store.profile,
  (val) => {
    if (val) {
      basicForm.realName = val.realName || ''
      basicForm.gender = val.gender === 0 ? 0 : 1
      basicForm.birthday = val.birthday?.slice(0, 10) || ''
      basicForm.phone = val.phone || ''
      basicForm.email = val.email || ''
    }
    loadLoginHistory()
  },
  { immediate: true }
)

async function handleUpdateBasic(formData: any) {
  updating.value = true
  try {
    await store.updateProfile({
      realName: formData.realName,
      gender: formData.gender,
      birthday: formData.birthday ? formData.birthday + 'T00:00:00' : undefined,
      avatar: store.profile?.avatar,
    })
    if (store.user?.id) {
      await userApi.updateUser(store.user.id, { phone: formData.phone, email: formData.email })
    }
    await store.fetchUserInfo()
    Message.success('更新成功')
  } catch (e) {
    Message.error('更新失败')
  } finally {
    updating.value = false
  }
}

async function handleChangePassword(pwdData: { oldPassword: string; newPassword: string }) {
  changingPassword.value = true
  try {
    const res = await store.changePassword(pwdData.oldPassword, pwdData.newPassword)
    if (res.success) {
      Message.success('密码修改成功')
    } else {
      Message.error(res.message || '修改失败')
    }
  } finally {
    changingPassword.value = false
  }
}

async function handleDeleteAccount() {
  try {
    await Modal.confirm({
      title: '确定注销账号吗？这是不可逆的操作！',
      content: '警告',
      okText: '确定',
      cancelText: '取消',
      okButtonProps: { status: 'warning' },
    })
    const deleteInput = ref('')
    await new Promise<void>((resolve, reject) => {
      // Use a simple inline modal-like approach
      Modal.confirm({
        title: '二次确认',
        content: () => {
          return h('div', {}, [
            h('p', { style: { marginBottom: '12px' } }, '请输入 "DELETE" 确认'),
            h('input', {
              value: deleteInput.value,
              onInput: (e: Event) => deleteInput.value = (e.target as HTMLInputElement).value,
              placeholder: '请输入 DELETE',
              style: { width: '100%' },
            }),
          ])
        },
        okText: '确定',
        cancelText: '取消',
        async onOk() {
          if (deleteInput.value !== 'DELETE') {
            Message.error('请输入 DELETE')
            throw new Error('Invalid input')
          }
          const res = await store.deleteAccount()
          if (res.success) {
            Message.success('注销成功')
            resolve()
          } else {
            reject(new Error('Delete failed'))
          }
        },
        onCancel() {
          reject(new Error('Cancelled'))
        },
      })
    })
  } catch (e) {
    // cancelled
  }
}

import { h } from 'vue'
</script>

<style scoped>
.profile-page {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.profile-header {
  margin-bottom: 24px;
}

.profile-main {
  display: flex;
  align-items: center;
  gap: 24px;
}

.profile-info h2 {
  font-size: 24px;
  color: #1d2129;
}

.settings-content {
  padding: 0 24px;
}

.settings-header {
  margin-bottom: 24px;
}

.settings-header h3 {
  font-size: 18px;
  color: #1d2129;
  margin-bottom: 4px;
}

.settings-header p {
  font-size: 14px;
  color: #86909c;
}
</style>
