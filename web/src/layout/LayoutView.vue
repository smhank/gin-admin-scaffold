<template>
  <el-container class="layout-container">
    <el-aside :width="isCollapse ? '64px' : '220px'" class="sidebar">
      <div class="logo" :class="{ 'logo-collapse': isCollapse }">
        <img src="https://www.gin-vue-admin.com/assets/logo-2f5d2e3f.png" alt="logo" class="logo-img" />
        <span v-show="!isCollapse" class="logo-text">Gin Admin</span>
      </div>
      <el-menu
        :default-active="route.path"
        router
        :collapse="isCollapse"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        class="sidebar-menu"
      >
        <el-menu-item index="/dashboard">
          <el-icon><HomeFilled /></el-icon>
          <template #title>首页</template>
        </el-menu-item>
        <el-sub-menu index="system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/paths">
            <el-icon><Link /></el-icon>
            <template #title>API路径</template>
          </el-menu-item>
          <el-menu-item index="/menus">
            <el-icon><Menu /></el-icon>
            <template #title>菜单管理</template>
          </el-menu-item>
          <el-menu-item index="/roles">
            <el-icon><UserFilled /></el-icon>
            <template #title>角色管理</template>
          </el-menu-item>
          <el-menu-item index="/permissions">
            <el-icon><Lock /></el-icon>
            <template #title>权限管理</template>
          </el-menu-item>
          <el-menu-item index="/users">
            <el-icon><UserFilled /></el-icon>
            <template #title>用户管理</template>
          </el-menu-item>
          <el-menu-item index="/operation-logs">
            <el-icon><Document /></el-icon>
            <template #title>操作历史</template>
          </el-menu-item>
          <el-menu-item index="/migrations">
            <el-icon><List /></el-icon>
            <template #title>迁移记录</template>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="toggleCollapse" :size="20">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/" class="breadcrumb">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="route.meta.title">{{ route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="28" class="user-avatar">A</el-avatar>
              {{ username }}
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <TabsView />
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { HomeFilled, Setting, Lock, ArrowDown, UserFilled, Fold, Expand, Menu, Link, Document, List } from '@element-plus/icons-vue'
import TabsView from './TabsView.vue'

const route = useRoute()
const router = useRouter()
const username = ref(localStorage.getItem('username') || 'admin')
const isCollapse = ref(false)

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const handleCommand = (command: string) => {
  if (command === 'logout') {
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    router.push('/login')
  }
}
</script>

<style scoped>
.layout-container { height: 100vh; }
.sidebar {
  background-color: #001529;
  overflow-y: auto;
  overflow-x: hidden;
  transition: width 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.08);
  z-index: 10;
}
.sidebar::-webkit-scrollbar { width: 0; }
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #002140;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}
.logo-img { width: 32px; height: 32px; margin-right: 8px; flex-shrink: 0; }
.logo-text { font-size: 18px; font-weight: bold; color: #fff; white-space: nowrap; letter-spacing: 1px; }
.logo-collapse .logo-img { margin-right: 0; }
.sidebar-menu { border-right: none; }
.header {
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 50px !important;
}
.header-left { display: flex; align-items: center; gap: 16px; }
.collapse-btn { cursor: pointer; color: #8c8c8c; font-size: 18px; transition: color 0.3s; }
.collapse-btn:hover { color: #1890ff; }
.breadcrumb { font-size: 14px; }
.header-right { cursor: pointer; }
.user-info {
  color: #333;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  padding: 4px 12px;
  border-radius: 6px;
  transition: all 0.3s;
}
.user-info:hover {
  background: #f5f5f5;
}
.user-avatar { background: linear-gradient(135deg, #1890ff, #69c0ff); vertical-align: middle; }
.main {
  background: #f0f2f5;
  padding: 16px;
  min-height: calc(100vh - 50px);
  overflow-y: auto;
}
</style>
