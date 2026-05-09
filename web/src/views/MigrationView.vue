<template>
  <div class="migration-container">
    <el-card shadow="never" class="migration-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">数据迁移记录</span>
          <el-tag type="success" size="small" effect="plain">
            共 {{ migrations.length }} 条记录
          </el-tag>
        </div>
      </template>

      <el-table :data="migrations" stripe style="width: 100%" v-loading="loading">
        <el-table-column type="index" label="序号" width="60" align="center" />
        <el-table-column prop="name" label="迁移名称" min-width="300">
          <template #default="{ row }">
            <div class="migration-name">
              <el-icon :size="16" color="#67c23a"><CircleCheck /></el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="batch" label="批次" width="80" align="center">
          <template #default="{ row }">
            <el-tag size="small" type="primary">#{{ row.batch }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="applied_at" label="应用时间" width="180" align="center" />
        <el-table-column prop="created_at" label="记录时间" width="180" align="center">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && migrations.length === 0" description="暂无迁移记录" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { CircleCheck } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface Migration {
  id: number
  name: string
  batch: number
  applied_at: string
  created_at: string
  updated_at: string
}

const loading = ref(false)
const migrations = ref<Migration[]>([])

const formatTime = (time: string) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 19)
}

const fetchMigrations = async () => {
  loading.value = true
  try {
    const res = await request.get('/migrations')
    migrations.value = res.data || []
  } catch (error) {
    console.error('获取迁移记录失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchMigrations()
})
</script>

<style scoped>
.migration-container {
  padding: 0;
}

.migration-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.migration-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  color: #303133;
}
</style>
