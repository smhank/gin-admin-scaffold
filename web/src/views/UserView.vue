<template>
  <div class="user-container">
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <el-button type="primary" size="default" @click="openDialog()">
            <el-icon><Plus /></el-icon>新增用户
          </el-button>
        </div>
      </template>

      <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column prop="Username" label="用户名" min-width="120" />
        <el-table-column prop="NickName" label="昵称" min-width="120" />
        <el-table-column prop="Phone" label="手机号" width="140" />
        <el-table-column prop="Email" label="邮箱" min-width="180" />
        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.Status === 1 ? 'success' : 'danger'" size="small">
              {{ scope.row.Status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="角色" min-width="150">
          <template #default="scope">
            <el-tag v-for="role in scope.row.Roles" :key="role.ID" size="small" style="margin-right: 4px; margin-bottom: 2px;">
              {{ role.Name }}
            </el-tag>
            <span v-if="!scope.row.Roles || scope.row.Roles.length === 0" style="color: #999;">-</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="openDialog(scope.row)">
              <el-icon><Edit /></el-icon>编辑
            </el-button>
            <el-button size="small" type="warning" link @click="handleResetPassword(scope.row)">
              <el-icon><Key /></el-icon>重置密码
            </el-button>
            <el-button size="small" type="primary" link @click="handleAssignRole(scope.row)">
              <el-icon><User /></el-icon>分配角色
            </el-button>
            <el-button size="small" type="danger" link @click="handleDelete(scope.row)">
              <el-icon><Delete /></el-icon>删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑用户右侧弹出层 -->
    <el-drawer v-model="dialogVisible" :title="isEdit ? '编辑用户' : '新增用户'" size="500px" direction="rtl" destroy-on-close>
      <el-form :model="form" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="请输入用户名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="密码" v-if="!isEdit">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickName" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.roleIds" multiple placeholder="请选择角色" style="width: 100%">
            <el-option v-for="role in allRoles" :key="role.ID" :label="role.Name" :value="role.ID" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSave">确 定</el-button>
      </template>
    </el-drawer>

    <!-- 重置密码弹出层 -->
    <el-dialog v-model="passwordDialogVisible" title="重置密码" width="400px" destroy-on-close>
      <el-form :model="passwordForm" label-width="80px">
        <el-form-item label="新密码">
          <el-input v-model="passwordForm.password" type="password" placeholder="请输入新密码" show-password />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请确认新密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleResetPasswordConfirm">确 定</el-button>
      </template>
    </el-dialog>

    <!-- 分配角色右侧弹出层 -->
    <el-drawer v-model="roleDialogVisible" :title="'分配角色 - ' + currentUserName" size="500px" direction="rtl" destroy-on-close>
      <div class="assign-role-config">
        <el-form label-width="80px">
          <el-form-item label="角色">
            <el-select v-model="currentUserRoleIds" multiple placeholder="请选择角色" style="width: 100%">
              <el-option v-for="role in allRoles" :key="role.ID" :label="role.Name" :value="role.ID" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleAssignRoleConfirm">确 定</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Key, User } from '@element-plus/icons-vue'

const tableData = ref([])
const allRoles = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const form = ref({
  username: '',
  password: '',
  nickName: '',
  phone: '',
  email: '',
  status: 1,
  roleIds: [] as number[]
})

// 重置密码
const passwordDialogVisible = ref(false)
const currentUserId = ref(0)
const passwordForm = ref({
  password: '',
  confirmPassword: ''
})

// 分配角色
const roleDialogVisible = ref(false)
const currentUserName = ref('')
const currentUserRoleIds = ref<number[]>([])

const formatTime = (t: string) => {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/users')
    tableData.value = res
  } catch (e) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

const fetchRoles = async () => {
  try {
    const res = await request.get('/roles')
    allRoles.value = res
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

const openDialog = (row?: any) => {
  if (row) {
    isEdit.value = true
    editId.value = row.ID
    form.value = {
      username: row.Username,
      password: '',
      nickName: row.NickName || '',
      phone: row.Phone || '',
      email: row.Email || '',
      status: row.Status !== undefined ? row.Status : 1,
      roleIds: row.Roles ? row.Roles.map((r: any) => r.ID) : []
    }
  } else {
    isEdit.value = false
    editId.value = 0
    form.value = {
      username: '',
      password: '',
      nickName: '',
      phone: '',
      email: '',
      status: 1,
      roleIds: []
    }
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!form.value.username) {
    ElMessage.warning('请输入用户名')
    return
  }
  if (!isEdit.value && !form.value.password) {
    ElMessage.warning('请输入密码')
    return
  }
  try {
    if (isEdit.value) {
      await request.put(`/users/${editId.value}`, {
        nickName: form.value.nickName,
        phone: form.value.phone,
        email: form.value.email,
        status: form.value.status,
        roleIds: form.value.roleIds
      })
      ElMessage.success('更新成功')
    } else {
      await request.post('/users', form.value)
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
    await ElMessageBox.confirm('确定要删除用户「' + row.Username + '」吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.delete(`/users/${row.ID}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e) {
    // 取消或错误已在拦截器中处理
  }
}

const handleResetPassword = (row: any) => {
  currentUserId.value = row.ID
  passwordForm.value = { password: '', confirmPassword: '' }
  passwordDialogVisible.value = true
}

const handleResetPasswordConfirm = async () => {
  if (!passwordForm.value.password) {
    ElMessage.warning('请输入新密码')
    return
  }
  if (passwordForm.value.password !== passwordForm.value.confirmPassword) {
    ElMessage.warning('两次输入的密码不一致')
    return
  }
  try {
    await request.put(`/users/${currentUserId.value}/reset-password`, {
      password: passwordForm.value.password
    })
    ElMessage.success('密码重置成功')
    passwordDialogVisible.value = false
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

const handleAssignRole = (row: any) => {
  currentUserId.value = row.ID
  currentUserName.value = row.NickName || row.Username
  currentUserRoleIds.value = row.Roles ? row.Roles.map((r: any) => r.ID) : []
  roleDialogVisible.value = true
}

const handleAssignRoleConfirm = async () => {
  try {
    await request.put(`/users/${currentUserId.value}/roles`, {
      roleIds: currentUserRoleIds.value
    })
    ElMessage.success('分配角色成功')
    roleDialogVisible.value = false
    fetchData()
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

onMounted(() => {
  fetchData()
  fetchRoles()
})
</script>

<style scoped>
.user-container { padding: 0; }
.table-card { border-radius: 4px; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: bold; }
.assign-role-config { padding: 10px 20px; }
.assign-header { display: flex; justify-content: flex-end; gap: 8px; margin-bottom: 16px; }
</style>
