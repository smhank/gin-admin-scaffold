<template>
  <div class="tabs-container">
    <div class="tabs-scroll" ref="tabsScroll">
      <div class="tabs-wrapper">
        <div
          v-for="tab in visitedTabs"
          :key="tab.path"
          class="tabs-item"
          :class="{ 'tabs-item-active': tab.path === currentPath }"
          @click="switchTab(tab)"
          @contextmenu.prevent="handleContextMenu($event, tab)"
        >
          <span class="tabs-item-dot" :class="{ 'tabs-item-dot-active': tab.path === currentPath }"></span>
          <span class="tabs-item-title">{{ tab.title }}</span>
          <el-icon
            v-if="tab.path !== '/dashboard'"
            class="tabs-item-close"
            @click.stop="closeTab(tab)"
          >
            <Close />
          </el-icon>
        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <teleport to="body">
      <div
        v-show="contextMenu.visible"
        class="context-menu"
        :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
        @click.stop
      >
        <div class="context-menu-item" @click="closeTab(contextMenu.tab)">
          <el-icon><Close /></el-icon>
          关闭当前
        </div>
        <div class="context-menu-item" @click="closeOtherTabs">
          <el-icon><CircleClose /></el-icon>
          关闭其他
        </div>
        <div class="context-menu-item" @click="closeAllTabs">
          <el-icon><Remove /></el-icon>
          关闭全部
        </div>
      </div>
    </teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Close, CircleClose, Remove } from '@element-plus/icons-vue'

interface TabItem {
  path: string
  title: string
  name: string
}

const route = useRoute()
const router = useRouter()
const tabsScroll = ref<HTMLElement | null>(null)

const visitedTabs = ref<TabItem[]>([
  { path: '/dashboard', title: '首页', name: 'dashboard' }
])

const currentPath = computed(() => route.path)

const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  tab: null as TabItem | null
})

// 添加新标签
const addTab = (tab: TabItem) => {
  const exists = visitedTabs.value.find(t => t.path === tab.path)
  if (!exists) {
    visitedTabs.value.push(tab)
  }
}

// 切换标签
const switchTab = (tab: TabItem) => {
  if (tab.path !== route.path) {
    router.push(tab.path)
  }
}

// 关闭标签
const closeTab = (tab: TabItem) => {
  if (tab.path === '/dashboard') return

  const index = visitedTabs.value.findIndex(t => t.path === tab.path)
  if (index === -1) return

  // 如果关闭的是当前标签，需要跳转到其他标签
  if (tab.path === route.path) {
    const nextTab = visitedTabs.value[index - 1] || visitedTabs.value[index + 1]
    if (nextTab) {
      router.push(nextTab.path)
    }
  }

  visitedTabs.value.splice(index, 1)
}

// 关闭其他标签
const closeOtherTabs = () => {
  visitedTabs.value = visitedTabs.value.filter(t => t.path === '/dashboard' || t.path === currentPath.value)
}

// 关闭全部标签
const closeAllTabs = () => {
  visitedTabs.value = visitedTabs.value.filter(t => t.path === '/dashboard')
  if (currentPath.value !== '/dashboard') {
    router.push('/dashboard')
  }
}

// 右键菜单
const handleContextMenu = (e: MouseEvent, tab: TabItem) => {
  contextMenu.value = {
    visible: true,
    x: e.clientX,
    y: e.clientY,
    tab
  }
}

const closeContextMenu = () => {
  contextMenu.value.visible = false
}

// 监听路由变化
watch(
  () => route.path,
  (path) => {
    const title = (route.meta?.title as string) || path.split('/').pop() || ''
    addTab({ path, title, name: route.name as string })
  },
  { immediate: true }
)

onMounted(() => {
  document.addEventListener('click', closeContextMenu)
})

onUnmounted(() => {
  document.removeEventListener('click', closeContextMenu)
})
</script>

<style scoped>
.tabs-container {
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  padding: 0;
  position: relative;
  user-select: none;
}

.tabs-scroll {
  overflow-x: auto;
  overflow-y: hidden;
  white-space: nowrap;
}

.tabs-scroll::-webkit-scrollbar {
  height: 0;
}

.tabs-wrapper {
  display: flex;
  align-items: center;
  padding: 0;
}

.tabs-item {
  display: inline-flex;
  align-items: center;
  padding: 8px 16px;
  font-size: 13px;
  color: #595959;
  cursor: pointer;
  position: relative;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  white-space: nowrap;
  background: #fafafa;
  border-right: 1px solid #f0f0f0;
}

.tabs-item:first-child {
  border-left: none;
}

.tabs-item:hover {
  color: #1890ff;
  background: #e6f7ff;
}

.tabs-item-active {
  color: #1890ff;
  background: #fff;
  font-weight: 500;
}

.tabs-item-active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, #1890ff, #69c0ff);
  border-radius: 2px 2px 0 0;
}

.tabs-item-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #d9d9d9;
  margin-right: 6px;
  transition: all 0.3s;
  flex-shrink: 0;
}

.tabs-item-dot-active {
  background: #1890ff;
  box-shadow: 0 0 6px rgba(24, 144, 255, 0.4);
}

.tabs-item-title {
  margin-right: 6px;
  font-size: 13px;
}

.tabs-item-close {
  font-size: 10px;
  padding: 2px;
  border-radius: 50%;
  transition: all 0.3s;
  opacity: 0;
  color: #8c8c8c;
  flex-shrink: 0;
}

.tabs-item:hover .tabs-item-close {
  opacity: 0.6;
}

.tabs-item-active .tabs-item-close {
  opacity: 0.6;
}

.tabs-item-close:hover {
  background: #bfbfbf;
  color: #fff;
  opacity: 1 !important;
}

.context-menu {
  position: fixed;
  z-index: 9999;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
  padding: 4px 0;
  min-width: 140px;
  animation: contextMenuFadeIn 0.15s ease;
}

@keyframes contextMenuFadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.context-menu-item {
  padding: 8px 16px;
  font-size: 13px;
  color: #333;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.context-menu-item .el-icon {
  font-size: 14px;
  color: #8c8c8c;
}

.context-menu-item:hover {
  background: #e6f7ff;
  color: #1890ff;
}

.context-menu-item:hover .el-icon {
  color: #1890ff;
}
</style>
