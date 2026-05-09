<template>
  <div class="role-container">
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>角色列表</span>
          <el-button type="primary" size="default" @click="openDialog()">
            <el-icon><Plus /></el-icon>新增角色
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
        <el-table-column prop="ID" label="角色ID" width="100" />
        <el-table-column prop="Name" label="角色名称" min-width="150" />
        <el-table-column label="操作" width="420" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="handleSetPermission(scope.row)">
              <el-icon><Setting /></el-icon>设置权限
            </el-button>
            <el-button size="small" type="primary" link @click="handleAssignUser(scope.row)">
              <el-icon><User /></el-icon>分配给用户
            </el-button>
            <el-button size="small" type="primary" link @click="openSubDialog(scope.row)">
              <el-icon><Plus /></el-icon>新增子角色
            </el-button>
            <el-button size="small" type="primary" link @click="handleCopy(scope.row)">
              <el-icon><CopyDocument /></el-icon>拷贝
            </el-button>
            <el-button size="small" type="primary" link @click="openDialog(scope.row)">
              <el-icon><Edit /></el-icon>编辑
            </el-button>
            <el-button size="small" type="danger" link @click="handleDelete(scope.row)">
              <el-icon><Delete /></el-icon>删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑角色右侧弹出层 -->
    <el-drawer v-model="dialogVisible" :title="isEdit ? '编辑角色' : '新增角色'" size="500px" direction="rtl" destroy-on-close>
      <el-form :model="form" label-width="80px">
        <el-form-item label="角色名称">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色编码">
          <el-input v-model="form.code" placeholder="例如: admin" />
        </el-form-item>
        <el-form-item label="上级角色">
          <el-tree-select
            v-model="form.parentId"
            :data="roleTree"
            :props="{ label: 'Name', value: 'ID', children: 'Children' }"
            placeholder="顶级角色"
            check-strictly
            clearable
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="角色描述">
          <el-input v-model="form.description" type="textarea" placeholder="请输入角色描述" />
        </el-form-item>
        <el-form-item label="角色状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSave">确 定</el-button>
      </template>
    </el-drawer>

    <!-- 设置权限右侧弹出层 -->
    <el-drawer v-model="permissionDialogVisible" title="角色配置" size="650px" direction="rtl" destroy-on-close>
      <div class="permission-config">
        <el-tabs v-model="permissionTab">
          <el-tab-pane label="角色菜单" name="menu">
            <div class="permission-toolbar">
              <el-input v-model="permissionFilter" placeholder="筛选" clearable style="width: 200px" />
              <el-button type="primary" @click="handleSavePermission">确 定</el-button>
            </div>
            <div class="default-home">
              <span>默认首页：</span>
              <el-select v-model="defaultHome" placeholder="请选择" style="width: 150px">
                <el-option label="仪表盘" value="dashboard" />
                <el-option label="官方网站" value="home" />
              </el-select>
            </div>
            <el-tree
              ref="treeRef"
              :data="allPermissions"
              show-checkbox
              node-key="ID"
              :props="{ label: 'Name', children: 'Children' }"
              default-expand-all
              :filter-node-method="filterNode"
            />
          </el-tab-pane>
          <el-tab-pane label="角色api" name="api">
            <div class="permission-toolbar">
              <el-input placeholder="筛选" clearable style="width: 200px" />
              <el-button type="primary" @click="handleSavePermission">确 定</el-button>
            </div>
            <el-table :data="allPaths" style="width: 100%" max-height="400">
              <el-table-column prop="Name" label="接口名称" />
              <el-table-column prop="Path" label="接口路径" />
              <el-table-column prop="Method" label="请求方法" width="100" />
              <el-table-column label="操作" width="80">
                <template #default="scope">
                  <el-checkbox :model-value="scope.row.Checked" @change="(val) => togglePath(val, scope.row)" />
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
          <el-tab-pane label="资源权限" name="resource">
            <div class="permission-toolbar">
              <el-input placeholder="筛选" clearable style="width: 200px" />
              <el-button type="primary" @click="handleSavePermission">确 定</el-button>
            </div>
            <el-tree
              :data="resourceData"
              show-checkbox
              node-key="ID"
              :props="{ label: 'Name', children: 'Children' }"
              default-expand-all
            />
          </el-tab-pane>
        </el-tabs>
      </div>
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSavePermission">确 定</el-button>
      </template>
    </el-drawer>

    <!-- 分配用户右侧弹出层 -->
    <el-drawer v-model="assignDialogVisible" :title="'分配用户 - ' + currentAssignRoleName" size="700px" direction="rtl" @opened="onAssignDialogOpened">
      <div class="assign-user-config">
        <div class="assign-header">
          <el-button @click="assignDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="handleAssignUserConfirm">确 定</el-button>
        </div>
        <el-alert
          title="注：保存时将全量覆盖该角色的用户关联关系；若用户仅剩此一个角色，移除后其主角色保持不变"
          type="warning"
          :closable="false"
          show-icon
          class="assign-tip"
        />
        <div class="assign-search">
          <el-form :inline="true" :model="assignSearchForm">
            <el-form-item label="关键字">
              <el-input v-model="assignSearchForm.keyword" placeholder="用户名/昵称" clearable style="width: 200px" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleAssignSearch">
                <el-icon><Search /></el-icon> 查询
              </el-button>
              <el-button @click="handleAssignReset">
                <el-icon><RefreshRight /></el-icon> 重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>
        <el-table :data="filteredUsers" ref="assignTableRef" style="width: 100%" max-height="400" @selection-change="handleSelectionChange" row-key="ID">
          <el-table-column type="selection" width="55" />
          <el-table-column prop="ID" label="ID" width="80" sortable />
          <el-table-column prop="Username" label="用户名" />
          <el-table-column prop="NickName" label="昵称" />
        </el-table>
        <div class="assign-pagination">
          <span class="total-text">共 {{ filteredUsers.length }} 条</span>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Setting, User, CopyDocument, Search, RefreshRight } from '@element-plus/icons-vue'

const tableData = ref([])
const allPermissions = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const form = ref({ name: '', code: '', parentId: null as number | null, sort: 0, description: '', status: 1 })
const roleTree = ref([])

// 设置权限相关
const permissionDialogVisible = ref(false)
const permissionTab = ref('menu')
const permissionFilter = ref('')
const defaultHome = ref('dashboard')
const treeRef = ref()
const currentRoleId = ref(0)
const allPaths = ref([])
const resourceData = ref([])

const filterNode = (value: string, data: any) => {
  if (!value) return true
  return data.Name.includes(value)
}

const togglePath = (val: boolean, row: any) => {
  row.Checked = val
}

// 分配给用户相关
const assignDialogVisible = ref(false)
const currentAssignRoleId = ref(0)
const currentAssignRoleName = ref('')
const allUsers = ref<any[]>([])
const selectedUsers = ref<any[]>([])
const assignedUserIds = ref<number[]>([])
const assignSearchForm = ref({ keyword: '' })
const assignTableRef = ref()
const filteredUsers = ref<any[]>([])

// 获取用户列表
const fetchUsers = async () => {
  try {
    const res = await request.get('/users')
    allUsers.value = res
    // 直接设置filteredUsers，不依赖watch
    applyUserFilter()
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

// 应用过滤
const applyUserFilter = () => {
  let filtered = [...allUsers.value]
  const keyword = assignSearchForm.value.keyword
  if (keyword) {
    filtered = filtered.filter(u => 
      u.Username.includes(keyword) || 
      (u.NickName && u.NickName.includes(keyword))
    )
  }
  filteredUsers.value = filtered
}

// 分配用户搜索
const handleAssignSearch = () => {
  applyUserFilter()
}

// 分配用户重置
const handleAssignReset = () => {
  assignSearchForm.value = { keyword: '' }
  applyUserFilter()
}

const handleSelectionChange = (selection: any[]) => {
  selectedUsers.value = selection
}

const isUserAssigned = (userId: number) => {
  return assignedUserIds.value.includes(userId)
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/roles')
    tableData.value = res
  } catch (e) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

const fetchPermissions = async () => {
  try {
    const res = await request.get('/permissions')
    allPermissions.value = res.map((p: any) => ({
      ...p,
      Children: []
    }))
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

const fetchPaths = async () => {
  try {
    const res = await request.get('/paths')
    allPaths.value = res.map((p: any) => ({
      ...p,
      Checked: false
    }))
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

// 新增/编辑角色
const openDialog = (row?: any) => {
  if (row) {
    isEdit.value = true
    editId.value = row.ID
    form.value = {
      name: row.Name,
      code: row.Code || '',
      parentId: row.ParentID || null,
      sort: row.Sort || 0,
      description: row.Description || '',
      status: row.Status !== undefined ? row.Status : 1
    }
  } else {
    isEdit.value = false
    editId.value = 0
    form.value = { name: '', code: '', parentId: null, sort: 0, description: '', status: 1 }
  }
  dialogVisible.value = true
}

// 新增子角色
const openSubDialog = (row: any) => {
  isEdit.value = false
  editId.value = 0
  form.value = { name: '', code: '', parentId: row.ID, sort: 0, description: '', status: 1 }
  dialogVisible.value = true
}

// 保存角色
const handleSave = async () => {
  if (!form.value.name) {
    ElMessage.warning('请输入角色名称')
    return
  }
  try {
    if (isEdit.value) {
      await request.put(`/roles/${editId.value}`, form.value)
      ElMessage.success('更新成功')
    } else {
      await request.post('/roles', form.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

// 设置权限
const handleSetPermission = async (row: any) => {
  currentRoleId.value = row.ID
  permissionTab.value = 'menu'
  defaultHome.value = 'dashboard'
  await fetchPermissions()
  await fetchPaths()
  
  // 获取当前角色的权限ID
  const currentPerms = row.Permissions ? row.Permissions.map((p: any) => p.ID) : []
  
  permissionDialogVisible.value = true
  // 等待树组件渲染完成
  await new Promise(resolve => setTimeout(resolve, 100))
  if (treeRef.value && currentPerms.length > 0) {
    treeRef.value.setCheckedKeys(currentPerms)
  }
}

// 保存权限
const handleSavePermission = async () => {
  const checkedKeys = treeRef.value ? treeRef.value.getCheckedKeys() : []
  try {
    await request.put(`/roles/${currentRoleId.value}/permissions`, { permissions: checkedKeys })
    ElMessage.success('权限设置成功')
    permissionDialogVisible.value = false
    fetchData()
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

// 分配给用户
const handleAssignUser = async (row: any) => {
  currentAssignRoleId.value = row.ID
  currentAssignRoleName.value = row.Name
  selectedUsers.value = []
  assignedUserIds.value = []
  
  // 获取当前已分配给该角色的用户
  if (row.Users && row.Users.length > 0) {
    assignedUserIds.value = row.Users.map((u: any) => u.ID)
  }
  
  // 先清空搜索条件，避免残留筛选条件
  assignSearchForm.value = { keyword: '' }
  
  // 先获取用户数据，再打开弹窗，避免竞态条件
  await fetchUsers()
  
  // 打开弹窗
  assignDialogVisible.value = true
}

// 弹窗完全打开后的回调
const onAssignDialogOpened = async () => {
  // 等待 DOM 渲染完成
  await nextTick()
  // 先清除所有选中状态，避免旧选中残留
  if (assignTableRef.value) {
    assignTableRef.value.clearSelection()
  }
  selectedUsers.value = []
  // 勾选已分配的用户
  if (assignTableRef.value && assignedUserIds.value.length > 0) {
    const selected: any[] = []
    filteredUsers.value.forEach((user: any) => {
      if (assignedUserIds.value.includes(user.ID)) {
        assignTableRef.value.toggleRowSelection(user, true)
        selected.push(user)
      }
    })
    selectedUsers.value = selected
  }
}

// 确认分配用户
const handleAssignUserConfirm = async () => {
  try {
    const userIds = selectedUsers.value.map((u: any) => u.ID)
    await request.post('/roles/assign', {
      roleId: currentAssignRoleId.value,
      userIds: userIds
    })
    ElMessage.success('分配成功')
    assignDialogVisible.value = false
    fetchData()
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

// 拷贝角色
const handleCopy = async (row: any) => {
  try {
    await request.post('/roles', {
      name: row.Name + ' (副本)',
      status: row.Status !== undefined ? row.Status : 1
    })
    ElMessage.success('拷贝成功')
    fetchData()
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

// 删除角色
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除角色「' + row.Name + '」吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.delete(`/roles/${row.ID}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e) {
    // 取消或错误已在拦截器中处理
  }
}

// 构建角色树（用于上级角色选择）
const buildRoleTree = () => {
  const data = JSON.parse(JSON.stringify(tableData.value))
  // 过滤掉当前编辑的角色自身及其子角色，避免循环引用
  const filterSelf = (items: any[], parentId: number | null = null): any[] => {
    return items.filter((item: any) => {
      if (isEdit.value && item.ID === editId.value) return false
      if (item.Children && item.Children.length > 0) {
        item.Children = filterSelf(item.Children, item.ID)
      }
      return true
    })
  }
  roleTree.value = isEdit.value ? filterSelf(data) : data
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.role-container { padding: 0; }
.table-card { border-radius: 4px; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: bold; }
.permission-config { padding: 10px 20px; }
.permission-toolbar { display: flex; justify-content: space-between; margin-bottom: 16px; }
.default-home { margin-bottom: 16px; display: flex; align-items: center; gap: 8px; }
.assign-user-config { padding: 10px 20px; }
.assign-header { display: flex; justify-content: flex-end; gap: 8px; margin-bottom: 16px; }
.assign-tip { margin-bottom: 16px; }
.assign-search { margin-bottom: 16px; }
.assign-search .el-form-item { margin-bottom: 0; }
.assign-pagination { display: flex; justify-content: center; align-items: center; gap: 16px; margin-top: 16px; }
.total-text { font-size: 14px; color: #606266; }
</style>
