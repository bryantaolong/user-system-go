<template>
  <a-layout class="admin-layout">
    <a-layout-sider class="layout-aside" :width="240">
      <div class="aside-header">
        <div class="logo">
          <icon-monitor />
          <span class="logo-text">用户管理系统</span>
        </div>
      </div>
      <a-menu :selected-keys="[activeMenu]" class="aside-menu" @menu-item-click="handleMenuClick">
        <a-menu-item key="/">
          <template #icon><icon-home /></template>
          <span>返回首页</span>
        </a-menu-item>
        <a-menu-item key="/admin/users">
          <template #icon><icon-user /></template>
          <span>用户管理</span>
        </a-menu-item>
        <a-menu-item key="/admin/profile">
          <template #icon><icon-user-group /></template>
          <span>个人中心</span>
        </a-menu-item>
        <a-menu-item key="/admin/logs">
          <template #icon><icon-settings /></template>
          <span>系统日志</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout class="layout-main">
      <a-layout-header class="layout-header">
        <a-breadcrumb>
          <a-breadcrumb-item v-for="(item, idx) in breadcrumbItems" :key="idx">
            <router-link v-if="idx < breadcrumbItems.length - 1" :to="item.path">{{ item.title }}</router-link>
            <span v-else>{{ item.title }}</span>
          </a-breadcrumb-item>
        </a-breadcrumb>
        <a-dropdown trigger="hover" @select="(val: string | number | Record<string, any> | undefined, ev: Event) => handleCommand(val as string, ev)">
          <div class="user-info">
            <a-avatar :size="32" style="margin-right: 8px; background: #165dff;">
              {{ store.user?.username?.charAt(0)?.toUpperCase() }}
            </a-avatar>
            <span class="username">{{ store.user?.username }}</span>
            <icon-down />
          </div>
          <template #content>
            <a-doption value="profile">个人中心</a-doption>
            <a-doption value="logout">退出登录</a-doption>
          </template>
        </a-dropdown>
      </a-layout-header>

      <a-layout-content class="layout-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Message, Modal } from '@arco-design/web-vue'
import { IconHome, IconUser, IconUserGroup, IconSettings, IconDown } from '@arco-design/web-vue/es/icon'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const store = useUserStore()

const activeMenu = computed(() => route.path)

const breadcrumbItems = computed(() => {
  const segments = route.path.split('/').filter(Boolean)
  const items: { title: string; path: string }[] = []
  let path = ''
  for (const seg of segments) {
    path += '/' + seg
    const title = seg === 'admin' ? '管理后台' : seg === 'users' ? '用户管理' : seg === 'profile' ? '个人中心' : seg === 'logs' ? '系统日志' : seg
    items.push({ title, path })
  }
  return items
})

async function handleMenuClick(key: string) {
  router.push(key)
}

async function handleCommand(command: string | number | Record<string, any> | undefined, ev: Event) {
  if (command === 'profile') {
    router.push('/admin/profile')
  } else if (command === 'logout') {
    try {
      await Modal.confirm({
        title: '确定要退出登录吗？',
        content: '提示',
        okText: '确定',
        cancelText: '取消',
      })
      await store.logout()
      router.push('/login')
      Message.success('已退出登录')
    } catch {
      // cancelled
    }
  }
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
}

.layout-aside {
  background: #001529;
  color: #fff;
}

.aside-header {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 16px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #fff;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
}

.aside-menu {
  border-right: none;
  background: #001529;
  color: #fff;
}

:deep(.arco-menu-inner) {
  background: #001529;
}

:deep(.arco-menu-item) {
  color: rgba(255, 255, 255, 0.65);
}

:deep(.arco-menu-selected) {
  background: #165dff !important;
  color: #fff !important;
}

:deep(.arco-menu-item:hover) {
  color: #fff;
}

.layout-header {
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e5e8ef;
  padding: 0 24px;
  height: 56px;
  line-height: 56px;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.username {
  margin-right: 4px;
  font-size: 14px;
  color: #1d2129;
}

.layout-content {
  background: #f2f3f5;
  padding: 24px;
  overflow-y: auto;
}
</style>
