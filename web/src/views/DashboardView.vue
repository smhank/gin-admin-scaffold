<template>
  <div class="dashboard-container">
    <el-row :gutter="20" class="dashboard-row">
      <el-col :span="6">
        <el-card shadow="hover" class="dashboard-card">
          <div class="card-content">
            <div class="card-icon" style="background: #1890ff;">
              <el-icon :size="28"><UserFilled /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">1</div>
              <div class="card-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="dashboard-card">
          <div class="card-content">
            <div class="card-icon" style="background: #52c41a;">
              <el-icon :size="28"><UserFilled /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">2</div>
              <div class="card-label">角色总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="dashboard-card">
          <div class="card-content">
            <div class="card-icon" style="background: #faad14;">
              <el-icon :size="28"><Lock /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">12</div>
              <div class="card-label">权限总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="dashboard-card">
          <div class="card-content">
            <div class="card-icon" style="background: #eb2f96;">
              <el-icon :size="28"><Menu /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">4</div>
              <div class="card-label">菜单总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="hover" class="welcome-card">
      <template #header>
        <span>欢迎使用 Gin Admin 后台管理系统</span>
      </template>
      <div class="welcome-content">
        <p>技术栈：Gin + GORM + Redis + Vue3 + Element Plus + TypeScript</p>
        <p>功能特性：RBAC 权限管理、菜单管理、API 路径管理、数据迁移</p>
        <el-divider />
        <p class="welcome-tip">左侧菜单可点击顶部折叠按钮进行伸缩</p>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { UserFilled, Lock, Menu } from '@element-plus/icons-vue'
</script>

<style scoped>
.dashboard-container { padding: 0; }
.dashboard-row { margin-bottom: 20px; }
.dashboard-card { border-radius: 4px; }
.card-content { display: flex; align-items: center; gap: 16px; }
.card-icon {
  width: 56px;
  height: 56px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}
.card-value { font-size: 24px; font-weight: bold; color: #333; }
.card-label { font-size: 14px; color: #999; margin-top: 4px; }
.welcome-card { border-radius: 4px; }
.welcome-content { font-size: 14px; color: #666; line-height: 2; }
.welcome-tip { color: #999; font-size: 12px; }
</style>
