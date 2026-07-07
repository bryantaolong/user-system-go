<template>
  <div class="admin-layout">
    <a-layout class="layout-container">
      <a-layout-sider width="240px" class="layout-aside">
        <div class="aside-header">
          <div class="logo">
            <IconDesktop :size="24" />
            <span class="logo-text">用户管理系统</span>
          </div>
        </div>

        <a-menu
          :selected-keys="[activeMenu]"
          @menu-item-click="handleMenuClick"
          class="aside-menu"
        >
          <a-menu-item key="/">
            <template #icon><IconHome /></template>
            <span>返回首页</span>
          </a-menu-item>
          <a-menu-item key="/admin/users">
            <template #icon><IconUser /></template>
            <span>用户管理</span>
          </a-menu-item>
          <a-menu-item key="/admin/profile">
            <template #icon><IconUserGroup /></template>
            <span>个人中心</span>
          </a-menu-item>
          <a-menu-item key="/admin/settings">
            <template #icon><IconSettings /></template>
            <span>系统设置</span>
          </a-menu-item>
        </a-menu>
      </a-layout-sider>

      <a-layout class="layout-main">
        <a-layout-header class="layout-header">
          <div class="header-left">
            <a-breadcrumb>
              <a-breadcrumb-item
                v-for="item in breadcrumbs"
                :key="item.path"
                @click="router.push(item.path)"
                style="cursor: pointer"
              >
                {{ item.meta?.title }}
              </a-breadcrumb-item>
            </a-breadcrumb>
          </div>

          <div class="header-right">
            <a-dropdown trigger="hover" @select="handleCommand">
              <div class="user-info">
                <a-avatar :size="32" :src="getAvatarUrl(userStore.userProfile?.avatar)">
                  {{ userStore.userInfo?.username?.charAt(0).toUpperCase() }}
                </a-avatar>
                <span class="username">{{ userStore.userInfo?.username }}</span>
                <IconDown />
              </div>
              <template #dropdown>
                <a-doption value="profile">
                  <template #icon><IconUser /></template>
                  个人中心
                </a-doption>
                <a-doption value="logout" :style="{ marginTop: '4px' }">
                  <template #icon><IconExport /></template>
                  退出登录
                </a-doption>
              </template>
            </a-dropdown>
          </div>
        </a-layout-header>

        <a-layout-content class="layout-content">
          <router-view />
        </a-layout-content>
      </a-layout>
    </a-layout>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArcoMessage, ArcoMessageBox } from '@/utils/arco-message'
import {
  IconDesktop,
  IconHome,
  IconUser,
  IconUserGroup,
  IconSettings,
  IconDown,
  IconExport,
} from '@arco-design/web-vue/es/icon'
import { useUserStore } from '@/stores/user'
import { getAvatarUrl } from '@/utils/file'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const breadcrumbs = computed(() => route.matched.filter(item => item.meta?.title))

const handleMenuClick = (key: string) => {
  router.push(key)
}

const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/admin/profile')
      break
    case 'logout':
      try {
        await ArcoMessageBox.confirm('确定要退出登录吗？', '提示')
        await userStore.logout()
        router.push('/login')
        ArcoMessage.success('已退出登录')
      } catch {
        // 取消退出
      }
      break
  }
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
}

.layout-container {
  height: 100%;
}

.layout-aside {
  background-color: #304156;
}

.aside-header {
  padding: 20px;
  border-bottom: 1px solid #434a50;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #fff;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.aside-menu {
  border-right: none;
}

.layout-main {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: rgba(255, 255, 255, 0.95);
  border-bottom: 1px solid #e4e7ed;
  backdrop-filter: blur(10px);
}

.header-right .user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 20px;
  transition: all 0.3s;
}

.header-right .user-info:hover {
  background-color: #f5f7fa;
}

.username {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.layout-content {
  padding: 20px;
  background: #fafbfc;
}

:deep(.arco-menu-item) {
  color: #bfcbd9;
}

:deep(.arco-menu-selected) {
  background-color: #263445 !important;
  color: #fff;
}
</style>
