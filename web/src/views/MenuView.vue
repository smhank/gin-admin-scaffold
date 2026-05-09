<template>
  <div class="menu-container">
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>菜单列表</span>
          <el-button type="primary" size="default" @click="openDialog()">
            <el-icon><Plus /></el-icon>新增菜单
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe row-key="ID" default-expand-all :tree-props="{ children: 'Children', hasChildren: 'hasChildren' }">
        <el-table-column prop="Name" label="菜单名称" min-width="160">
          <template #default="scope">
            <span v-if="scope.row.Icon" style="margin-right: 6px;">
              <el-icon :size="16"><component :is="scope.row.Icon" /></el-icon>
            </span>
            {{ scope.row.Name }}
          </template>
        </el-table-column>
        <el-table-column label="图标" width="80" align="center">
          <template #default="scope">
            <el-icon v-if="scope.row.Icon" :size="20"><component :is="scope.row.Icon" /></el-icon>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="RouteName" label="路由名称" width="140" />
        <el-table-column prop="Path" label="路由路径" width="160" />
        <el-table-column prop="FilePath" label="文件路径" width="200" />
        <el-table-column label="父节点" width="140">
          <template #default="scope">
            {{ getParentName(scope.row.ParentID) }}
          </template>
        </el-table-column>
        <el-table-column prop="Sort" label="排序" width="60" align="center" />
        <el-table-column label="状态" width="80" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.Status === 1 ? 'success' : 'info'" size="small">
              {{ scope.row.Status === 1 ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="scope">
            {{ formatTime(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="openSubMenuDialog(scope.row)">
              <el-icon><Plus /></el-icon>添加子菜单
            </el-button>
            <el-button size="small" type="primary" link @click="handleAssignRole(scope.row)">
              <el-icon><User /></el-icon>分配角色
            </el-button>
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
    <el-drawer v-model="dialogVisible" :title="isEdit ? '编辑菜单' : '新增菜单'" size="600px" direction="rtl">
      <el-form :model="form" label-width="80px" size="default">
        <el-form-item label="上级菜单">
          <el-tree-select
            v-model="form.parentId"
            :data="menuTree"
            :props="{ label: 'Name', value: 'ID', children: 'Children' }"
            placeholder="顶级菜单"
            check-strictly
            clearable
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="菜单名称">
          <el-input v-model="form.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="路由名称">
          <el-input v-model="form.routeName" placeholder="例如: Dashboard" />
        </el-form-item>
        <el-form-item label="路由路径">
          <el-input v-model="form.path" placeholder="例如: /users" />
        </el-form-item>
        <el-form-item label="文件路径">
          <el-input v-model="form.filePath" placeholder="例如: /src/views/DashboardView.vue" />
        </el-form-item>
        <el-form-item label="菜单图标">
          <el-select v-model="form.icon" placeholder="选择图标" filterable clearable style="width: 100%">
            <el-option v-for="icon in iconList" :key="icon" :label="icon" :value="icon">
              <span style="display: flex; align-items: center; gap: 8px;">
                <el-icon :size="16"><component :is="icon" /></el-icon>
                <span>{{ icon }}</span>
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">显示</el-radio>
            <el-radio :value="0">隐藏</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSave">确 定</el-button>
      </template>
    </el-drawer>

    <!-- 分配角色右侧弹出层 -->
    <el-drawer v-model="assignRoleDialogVisible" :title="'分配角色 - ' + currentMenuName" size="500px" direction="rtl">
      <div class="assign-role-config">
        <div class="assign-header">
          <el-button @click="assignRoleDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="handleAssignRoleConfirm">确 定</el-button>
        </div>
        <el-table :data="allRoles" ref="assignRoleTableRef" style="width: 100%" max-height="400" @selection-change="handleRoleSelectionChange" row-key="ID">
          <el-table-column type="selection" width="55" />
          <el-table-column prop="ID" label="ID" width="80" />
          <el-table-column prop="Name" label="角色名称" />
          <el-table-column prop="Code" label="角色编码" />
        </el-table>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, User } from '@element-plus/icons-vue'

const tableData = ref([])
const menuTree = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const form = ref({ parentId: null as number | null, name: '', routeName: '', path: '', filePath: '', icon: '', sort: 0, status: 1 })

const iconList = ['HomeFilled', 'Setting', 'Lock', 'UserFilled', 'Menu', 'Link', 'Fold', 'Expand', 'Plus', 'Edit', 'Delete', 'Search', 'Refresh', 'Upload', 'Download']

// 分配角色相关
const assignRoleDialogVisible = ref(false)
const currentMenuId = ref(0)
const currentMenuName = ref('')
const allRoles = ref<any[]>([])
const selectedRoleIds = ref<number[]>([])
const assignRoleTableRef = ref()

const formatTime = (t: string) => {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/menus')
    tableData.value = res
    menuTree.value = JSON.parse(JSON.stringify(res))
  } catch (e) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

// 获取父节点名称
const getParentName = (parentId: number | null) => {
  if (!parentId) return '-'
  const findName = (items: any[]): string => {
    for (const item of items) {
      if (item.ID === parentId) return item.Name
      if (item.Children && item.Children.length > 0) {
        const name = findName(item.Children)
        if (name) return name
      }
    }
    return '-'
  }
  return findName(tableData.value)
}

const openDialog = (row?: any) => {
  if (row) {
    isEdit.value = true
    editId.value = row.ID
    form.value = {
      parentId: row.ParentID || null,
      name: row.Name,
      routeName: row.RouteName || '',
      path: row.Path || '',
      filePath: row.FilePath || '',
      icon: row.Icon || '',
      sort: row.Sort || 0,
      status: row.Status !== undefined ? row.Status : 1
    }
  } else {
    isEdit.value = false
    editId.value = 0
    form.value = { parentId: null, name: '', routeName: '', path: '', filePath: '', icon: '', sort: 0, status: 1 }
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    if (isEdit.value) {
      await request.put(`/menus/${editId.value}`, form.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/menus', form.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

// 添加子菜单
const openSubMenuDialog = (row: any) => {
  isEdit.value = false
  editId.value = 0
  form.value = { parentId: row.ID, name: '', routeName: '', path: '', filePath: '', icon: '', sort: 0, status: 1 }
  dialogVisible.value = true
}

// 分配角色
const handleAssignRole = async (row: any) => {
  currentMenuId.value = row.ID
  currentMenuName.value = row.Name
  selectedRoleIds.value = []
  
  // 获取所有角色
  try {
    const res = await request.get('/roles')
    allRoles.value = res
  } catch (e) {
    // 错误已在拦截器中处理
  }
  
  // 获取当前菜单已分配的角色
  if (row.Roles && row.Roles.length > 0) {
    selectedRoleIds.value = row.Roles.map((r: any) => r.ID)
  }
  
  assignRoleDialogVisible.value = true
}

const handleRoleSelectionChange = (selection: any[]) => {
  selectedRoleIds.value = selection.map((r: any) => r.ID)
}

const handleAssignRoleConfirm = async () => {
  try {
    await request.put(`/menus/${currentMenuId.value}/roles`, {
      roleIds: selectedRoleIds.value
    })
    ElMessage.success('分配角色成功')
    assignRoleDialogVisible.value = false
    fetchData()
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除菜单「' + row.Name + '」吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.delete(`/menus/${row.ID}`)
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
.menu-container { padding: 0; }
.table-card { border-radius: 4px; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: bold; }
.assign-role-config { padding: 0 20px; }
.assign-header { display: flex; justify-content: flex-end; gap: 10px; margin-bottom: 16px; }
</style>
