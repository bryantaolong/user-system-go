<template>
  <div class="admin-layout">
    <el-container class="layout-container">
      <el-aside width="240px" class="layout-aside">
        <div class="aside-header">
          <div class="logo">
            <el-icon :size="24"><Platform /></el-icon>
            <span class="logo-text">用户管理系统</span>
          </div>
        </div>

        <el-menu
          router
          :default-active="activeMenu"
          class="aside-menu"
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409eff"
        >
          <el-menu-item index="/">
            <el-icon><House /></el-icon>
            <span>返回首页</span>
          </el-menu-item>
          <el-menu-item index="/admin/users">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/profile">
            <el-icon><Avatar /></el-icon>
            <span>个人中心</span>
          </el-menu-item>
          <el-menu-item index="/admin/settings">
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container class="layout-main">
        <el-header class="layout-header">
          <div class="header-left">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item
                v-for="item in breadcrumbs"
                :key="item.path"
                :to="{ path: item.path }"
              >
                {{ item.meta?.title }}
              </el-breadcrumb-item>
            </el-breadcrumb>
          </div>

          <div class="header-right">
            <el-dropdown trigger="hover" @command="handleCommand">
              <div class="user-info">
                <el-avatar :size="32" :src="getAvatarUrl(userStore.userProfile?.avatar)">
                  {{ userStore.userInfo?.username?.charAt(0).toUpperCase() }}
                </el-avatar>
                <span class="username">{{ userStore.userInfo?.username }}</span>
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">
                    <el-icon><User /></el-icon>
                    个人中心
                  </el-dropdown-item>
                  <el-dropdown-item divided command="logout">
                    <el-icon><SwitchButton /></el-icon>
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <el-main class="layout-content">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Platform,
  User,
  Avatar,
  Setting,
  ArrowDown,
  SwitchButton,
  House
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { getAvatarUrl } from '@/utils/file'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const breadcrumbs = computed(() => route.matched.filter(item => item.meta?.title))

const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/admin/profile')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          type: 'warning'
        })
        await userStore.logout()
        router.push('/login')
        ElMessage.success('已退出登录')
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
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
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

:deep(.el-menu-item.is-active) {
  background-color: #263445 !important;
}
</style>
