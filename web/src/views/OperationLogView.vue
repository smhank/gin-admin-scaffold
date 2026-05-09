<template>
  <div class="operation-log-container">
    <!-- 搜索区域 -->
    <el-card shadow="hover" class="search-card">
      <el-form :model="queryParams" inline>
        <el-form-item label="操作人">
          <el-input v-model="queryParams.username" placeholder="请输入操作人" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item label="操作动作">
          <el-input v-model="queryParams.action" placeholder="请输入操作动作" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item label="操作结果">
          <el-select v-model="queryParams.result" placeholder="请选择" clearable style="width: 120px">
            <el-option label="成功" value="success" />
            <el-option label="失败" value="fail" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 320px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card shadow="hover" class="table-card">
      <template #header>
        <div class="table-header">
          <span>操作日志列表</span>
          <el-button type="danger" size="small" @click="handleClearAll">
            <el-icon><Delete /></el-icon>
            清空日志
          </el-button>
        </div>
      </template>

      <el-table :data="logs" stripe style="width: 100%" v-loading="loading">
        <el-table-column type="index" label="序号" width="60" align="center" />
        <el-table-column prop="username" label="操作人" width="120" />
        <el-table-column prop="action" label="操作动作" width="160" />
        <el-table-column prop="method" label="请求方法" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="methodType(row.method)" size="small" effect="plain">
              {{ row.method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="请求路径" min-width="200" show-overflow-tooltip />
        <el-table-column prop="result" label="操作结果" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.result === 'success' ? 'success' : 'danger'" size="small">
              {{ row.result === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="耗时(ms)" width="100" align="center">
          <template #default="{ row }">
            <span :class="durationClass(row.duration)">{{ row.duration }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="CreatedAt" label="操作时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="danger" size="small" link @click="handleDelete(row)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Refresh, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

interface OperationLog {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  username: string
  action: string
  method: string
  path: string
  params: string
  result: string
  duration: number
  ip: string
  user_agent: string
}

const loading = ref(false)
const logs = ref<OperationLog[]>([])
const dateRange = ref<string[]>([])

const queryParams = reactive({
  username: '',
  action: '',
  result: '',
  startTime: '',
  endTime: ''
})

const fetchLogs = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (queryParams.username) params.username = queryParams.username
    if (queryParams.action) params.action = queryParams.action
    if (queryParams.result) params.result = queryParams.result
    if (queryParams.startTime) params.startTime = queryParams.startTime
    if (queryParams.endTime) params.endTime = queryParams.endTime

    const res = await request.get('/operation-logs', { params })
    logs.value = res
  } catch (err) {
    ElMessage.error('获取日志失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  if (dateRange.value && dateRange.value.length === 2) {
    queryParams.startTime = dateRange.value[0]
    queryParams.endTime = dateRange.value[1]
  } else {
    queryParams.startTime = ''
    queryParams.endTime = ''
  }
  fetchLogs()
}

const handleReset = () => {
  queryParams.username = ''
  queryParams.action = ''
  queryParams.result = ''
  queryParams.startTime = ''
  queryParams.endTime = ''
  dateRange.value = []
  fetchLogs()
}

const handleDelete = (row: OperationLog) => {
  ElMessageBox.confirm('确定要删除该日志吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete(`/operation-logs/${row.ID}`)
      ElMessage.success('删除成功')
      fetchLogs()
    } catch (err) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const handleClearAll = () => {
  ElMessageBox.confirm('确定要清空所有操作日志吗？此操作不可恢复！', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete('/operation-logs/all')
      ElMessage.success('清空成功')
      fetchLogs()
    } catch (err) {
      ElMessage.error('清空失败')
    }
  }).catch(() => {})
}

const methodType = (method: string) => {
  const map: Record<string, string> = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'info'
  }
  return map[method] || 'info'
}

const durationClass = (duration: number) => {
  if (duration > 1000) return 'duration-slow'
  if (duration > 500) return 'duration-medium'
  return 'duration-fast'
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 19)
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.operation-log-container { padding: 0; }
.search-card { margin-bottom: 16px; border-radius: 4px; }
.table-card { border-radius: 4px; }
.table-header { display: flex; justify-content: space-between; align-items: center; }
.duration-fast { color: #67c23a; }
.duration-medium { color: #e6a23c; }
.duration-slow { color: #f56c6c; font-weight: bold; }
</style>
