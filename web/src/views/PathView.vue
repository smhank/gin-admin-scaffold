<template>
  <div class="path-container">
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>API路径列表</span>
          <el-button type="primary" size="default" @click="openDialog()">
            <el-icon><Plus /></el-icon>新增API路径
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column prop="Name" label="接口名称" />
        <el-table-column prop="Path" label="接口路径" />
        <el-table-column label="请求方法" width="120">
          <template #default="scope">
            <el-tag :type="methodType(scope.row.Method)" size="small">
              {{ scope.row.Method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="Desc" label="接口描述" />
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="openDialog(scope.row)">
              <el-icon><Edit /></el-icon>
            </el-button>
            <el-button size="small" type="danger" link @click="handleDelete(scope.row)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑右侧弹出层 -->
    <el-drawer v-model="dialogVisible" :title="isEdit ? '编辑API路径' : '新增API路径'" size="500px" direction="rtl">
      <el-form :model="form" label-width="80px" size="default">
        <el-form-item label="接口名称">
          <el-input v-model="form.name" placeholder="接口名称" />
        </el-form-item>
        <el-form-item label="接口路径">
          <el-input v-model="form.path" placeholder="例如: /api/users" />
        </el-form-item>
        <el-form-item label="请求方法">
          <el-select v-model="form.method" placeholder="选择请求方法" style="width: 100%">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
            <el-option label="PATCH" value="PATCH" />
          </el-select>
        </el-form-item>
        <el-form-item label="接口描述">
          <el-input v-model="form.desc" type="textarea" placeholder="接口描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSave">确 定</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'

const tableData = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const form = ref({ name: '', path: '', method: 'GET', desc: '' })

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

const formatTime = (t: string) => {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/paths')
    tableData.value = res
  } catch (e) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

const openDialog = (row?: any) => {
  if (row) {
    isEdit.value = true
    editId.value = row.ID
    form.value = {
      name: row.Name,
      path: row.Path,
      method: row.Method || 'GET',
      desc: row.Desc || ''
    }
  } else {
    isEdit.value = false
    editId.value = 0
    form.value = { name: '', path: '', method: 'GET', desc: '' }
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    if (isEdit.value) {
      await request.put(`/paths/${editId.value}`, form.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/paths', form.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除API路径「' + row.Name + '」吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.delete(`/paths/${row.ID}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e) {
    // 取消或错误已在拦截器中处理
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.path-container { padding: 0; }
.table-card { border-radius: 4px; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: bold; }
</style>
