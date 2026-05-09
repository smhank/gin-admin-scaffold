import { useAuthStore } from '@/stores/auth'
import type { Directive } from 'vue'

export const vPermission: Directive = {
  mounted(el, binding) {
    const { value } = binding
    const authStore = useAuthStore()
    const permissions = authStore.permissions

    if (value && Array.isArray(value)) {
      const hasPermission = permissions.some((p: string) => (value as string[]).includes(p))
      if (!hasPermission) {
        el.parentNode?.removeChild(el)
      }
    }
  }
}
