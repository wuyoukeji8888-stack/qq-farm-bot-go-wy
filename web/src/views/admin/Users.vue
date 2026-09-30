<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { useAppStore } from '@/stores/app'
import { useRouter } from 'vue-router'

const app = useAppStore()
const router = useRouter()

const users = ref([])
const loading = ref(false)
const selected = ref(new Set())
const editing = ref(null)
const editLimit = ref(2)
const editDays = ref(30)
const editPermanent = ref(false)
const editPassword = ref('')
const saving = ref(false)

async function loadUsers() {
  loading.value = true
  try {
    const { data } = await api.get('/api/admin/users')
    users.value = data.users || []
  } catch (e) {
    app.error(e.response?.data?.error || '加载失败')
  } finally {
    loading.value = false
  }
}

async function deleteUsers() {
  if (selected.value.size === 0) {
    app.error('请先选择用户')
    return
  }
  if (!confirm(`确定删除选中的 ${selected.value.size} 个用户？`)) return
  try {
    await api.delete('/api/admin/users', {
      data: { usernames: Array.from(selected.value) }
    })
    app.success('删除成功')
    selected.value.clear()
    await loadUsers()
  } catch (e) {
    app.error(e.response?.data?.error || '删除失败')
  }
}

function toggleSelect(username) {
  if (selected.value.has(username)) {
    selected.value.delete(username)
  } else {
    selected.value.add(username)
  }
}

function selectAll() {
  if (selected.value.size === users.value.length) {
    selected.value.clear()
  } else {
    users.value.forEach(u => selected.value.add(u.username))
  }
}

function formatTime(ms) {
  if (!ms) return '-'
  const d = new Date(ms)
  return d.toLocaleString('zh-CN')
}

function remainingDays(user) {
  if (user.isPermanent) return -1
  const exp = user.expiresAt || (user.card && user.card.expiresAt) || 0
  if (!exp) return 0
  const left = exp - Date.now()
  if (left <= 0) return 0
  return Math.ceil(left / 86400000)
}

function startEdit(user) {
  editing.value = user.username
  editLimit.value = user.accountLimit === -1 ? -1 : (user.accountLimit || 2)
  editPermanent.value = !!user.isPermanent
  const days = remainingDays(user)
  editDays.value = days > 0 ? days : 30
  editPassword.value = ''
}

function cancelEdit() {
  editing.value = null
}

async function saveEdit(username) {
  saving.value = true
  try {
    const payload = {
      username,
      accountLimit: Number(editLimit.value)
    }
    if (editPermanent.value) {
      payload.permanent = true
    } else {
      payload.days = Number(editDays.value)
    }
    if (editPassword.value) {
      payload.password = editPassword.value
    }
    const { data } = await api.post('/api/admin/users', payload)
    if (data.ok) {
      app.success('已保存')
      editing.value = null
      await loadUsers()
    } else {
      app.error(data.error || '保存失败')
    }
  } catch (e) {
    app.error(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(loadUsers)
</script>

<template>
  <div>
    <div class="subbar">
      <button class="icon-btn" @click="router.push('/more')">‹</button>
      <h3>用户管理</h3>
    </div>
    <div style="padding:12px;">
      <div class="sec-title" style="margin:2px 0 12px">
        <span>用户列表（{{ users.length }}）</span>
      </div>
      <div style="display:flex;gap:8px;margin-bottom:12px;flex-wrap:wrap;">
        <button style="padding:8px 12px;border-radius:8px;background:var(--card-strong);border:1px solid var(--border);cursor:pointer;font-size:13px;" @click="selectAll">全选</button>
        <button style="padding:8px 12px;border-radius:8px;background:#ef4444;color:#fff;border:none;cursor:pointer;font-size:13px;" @click="deleteUsers">删除</button>
        <button style="padding:8px 12px;border-radius:8px;background:var(--card-strong);border:1px solid var(--border);cursor:pointer;font-size:13px;" @click="loadUsers">刷新</button>
      </div>

      <div v-if="loading" style="text-align:center;padding:24px;color:var(--muted)">加载中...</div>
      <div v-else-if="users.length === 0" style="text-align:center;padding:24px;color:var(--muted)">暂无用户</div>
      <div v-else style="display:flex;flex-direction:column;gap:8px;">
        <div v-for="user in users" :key="user.username" style="padding:12px;border-radius:10px;background:var(--card-strong);border:1px solid var(--border);">
          <div style="display:flex;align-items:center;gap:8px;">
            <input type="checkbox" :checked="selected.has(user.username)" @change="toggleSelect(user.username)">
            <div style="flex:1;min-width:0;">
              <div style="font-weight:600;">{{ user.username }}</div>
              <div style="font-size:12px;color:var(--muted);margin-top:4px;">
                角色: {{ user.role === 'admin' ? '管理员' : '普通用户' }} |
                账号限制: {{ user.accountLimit === -1 ? '无限' : user.accountLimit }} |
                到期: {{ user.isPermanent ? '永久' : formatTime(user.expiresAt || (user.card && user.card.expiresAt)) }}
              </div>
              <div v-if="editing === user.username && user.role !== 'admin'" style="display:flex;flex-direction:column;gap:8px;margin-top:10px;">
                <div style="display:flex;gap:8px;">
                  <input v-model.number="editLimit" type="number" class="field" placeholder="账号上限(-1无限)" style="flex:1">
                  <input v-if="!editPermanent" v-model.number="editDays" type="number" class="field" placeholder="剩余天数" style="flex:1">
                </div>
                <label style="font-size:12px;display:flex;align-items:center;gap:6px;">
                  <input type="checkbox" v-model="editPermanent"> 永久时长
                </label>
                <input v-model="editPassword" type="password" class="field" placeholder="新密码（留空不修改）" style="margin-top:4px">
                <div style="font-size:11px;color:var(--muted);">密码需至少 8 位，并包含大小写字母、数字或符号</div>
                <div style="display:flex;gap:8px;">
                  <button :disabled="saving" style="flex:1;padding:8px;border-radius:8px;background:var(--primary,#3b82f6);color:#fff;border:none;cursor:pointer;" @click="saveEdit(user.username)">{{ saving ? '保存中...' : '保存' }}</button>
                  <button style="flex:1;padding:8px;border-radius:8px;background:var(--card-strong);border:1px solid var(--border);cursor:pointer;" @click="cancelEdit">取消</button>
                </div>
              </div>
            </div>
            <button v-if="user.role !== 'admin' && editing !== user.username" style="padding:4px 8px;border-radius:6px;background:transparent;border:1px solid var(--border);cursor:pointer;font-size:12px;" @click="startEdit(user)">编辑</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
