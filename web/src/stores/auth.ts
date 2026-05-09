import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const permissions = ref<string[]>([])

  function setToken(t: string) {
    token.value = t
    localStorage.setItem('token', t)
  }

  function setPermissions(p: string[]) {
    permissions.value = p
  }

  return { token, permissions, setToken, setPermissions }
})
