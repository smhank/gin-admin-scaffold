<template>
  <div class="permission-container">
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>权限列表</span>
          <el-button type="primary" size="default" @click="openDialog()">
            <el-icon><Plus /></el-icon>新增权限
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column prop="Name" label="名称" />
        <el-table-column prop="Code" label="权限标识" />
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
    <el-drawer v-model="dialogVisible" :title="isEdit ? '编辑权限' : '新增权限'" size="500px" direction="rtl">
      <el-form :model="form" label-width="80px" size="default">
        <el-form-item label="名称">
          <el-input v-model="form.name" placeholder="请输入权限名称" />
        </el-form-item>
        <el-form-item label="权限标识">
          <el-input v-model="form.code" placeholder="例如: user:list" />
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
const form = ref({ name: '', code: '' })

const formatTime = (t: string) => {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/permissions')
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
    form.value = { name: row.Name, code: row.Code }
  } else {
    isEdit.value = false
    editId.value = 0
    form.value = { name: '', code: '' }
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    if (isEdit.value) {
      await request.put(`/permissions/${editId.value}`, form.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/permissions', form.value)
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
    await ElMessageBox.confirm('确定要删除权限「' + row.Name + '」吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.delete(`/permissions/${row.ID}`)
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
.permission-container { padding: 0; }
.table-card { border-radius: 4px; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: bold; }
</style>
