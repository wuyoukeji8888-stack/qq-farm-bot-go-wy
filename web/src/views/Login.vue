<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { useAccountStore } from '@/stores/account'
import { useAppStore } from '@/stores/app'

const account = useAccountStore()
const app = useAppStore()
const router = useRouter()
const username = ref('')
const password = ref('')
const loading = ref(false)
const isRegister = ref(false) // 切换登录/注册模式

// 登录提交
async function onSubmit() {
  if (!username.value || !password.value) {
    app.error('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    if (isRegister.value) {
      // 注册新用户
      const { data } = await api.post('/api/users/register', {
        username: username.value,
        password: password.value
      })
      if (data.ok) {
        app.success('注册成功，请登录')
        isRegister.value = false
      } else {
        app.error(data.error || '注册失败')
      }
    } else {
      // 用户登录
      const { data } = await api.post('/api/users/login', {
        username: username.value,
        password: password.value
      })
      if (data.ok && data.token) {
        account.setToken(data.token)
        await account.loadAccounts()
        app.success('登录成功')
        router.replace('/')
      } else {
        app.error(data.error || '登录失败')
      }
    }
  } catch (e) {
    app.error(e.response?.data?.error || '操作失败')
  } finally {
    loading.value = false
  }
}

// 切换登录/注册模式
function toggleMode() {
  isRegister.value = !isRegister.value
}
</script>

<template>
  <div class="login-screen">
    <div class="login-card glass">
      <div class="login-logo">🌾</div>
      <h1>QQ 农场</h1>
      <p class="login-title">{{ isRegister ? '用户注册' : '后台管理登录' }}</p>
      <p class="login-sub">{{ isRegister ? '注册新账号' : '请输入用户名和密码' }}</p>
      
      <input
        v-model="username"
        type="text"
        class="ipt"
        placeholder="用户名"
        @keyup.enter="onSubmit"
      />
      <input
        v-model="password"
        type="password"
        class="ipt"
        :placeholder="isRegister ? '设置密码（至少 6 位）' : '密码'"
        @keyup.enter="onSubmit"
      />
      
      <button class="btn primary" :disabled="loading" @click="onSubmit">
        {{ loading ? (isRegister ? '注册中…' : '登录中…') : (isRegister ? '注册' : '登录') }}
      </button>
      
      <p class="toggle-link">
        {{ isRegister ? '已有账号？' : '没有账号？' }}
        <a href="#" @click.prevent="toggleMode">
          {{ isRegister ? '去登录' : '去注册' }}
        </a>
      </p>
    </div>
  </div>
</template>

<style scoped>
.login-screen {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
}
.login-card {
  width: 100%;
  max-width: 360px;
  padding: 32px 28px;
  border-radius: var(--radius-lg);
  text-align: center;
}
.login-title { font-size: 18px; font-weight: 600; margin-bottom: 8px; }
.login-sub { color: var(--muted); font-size: 13px; margin: -4px 0 20px; }
.ipt {
  width: 100%;
  padding: 12px 14px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border);
  background: var(--card-strong);
  color: var(--foreground);
  font-size: 15px;
  margin-bottom: 14px;
  box-sizing: border-box;
}
.btn.primary {
  width: 100%;
  height: 46px;
  font-size: 15px;
  font-weight: 700;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.toggle-link {
  margin-top: 16px;
  font-size: 13px;
  color: var(--muted);
}
.toggle-link a {
  color: var(--primary);
  text-decoration: none;
  margin-left: 4px;
}
.toggle-link a:hover {
  text-decoration: underline;
}
</style>
