<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/api'
import { useAppStore } from '@/stores/app'
import { useAccountStore } from '@/stores/account'
import { useRouter } from 'vue-router'

const app = useAppStore()
const account = useAccountStore()
const router = useRouter()
const oldPwd = ref(''); const newPwd = ref(''); const newPwd2 = ref('')
const busy = ref(false)
const cardCode = ref('')
const renewBusy = ref(false)

const limitText = computed(() => {
  const n = account.userInfo?.accountLimit
  if (n === -1) return '无限'
  return n || 2
})
const expireText = computed(() => {
  const u = account.userInfo
  if (!u) return '-'
  if (u.isPermanent) return '永久'
  if (!u.expiresAt) return '-'
  return new Date(u.expiresAt).toLocaleString('zh-CN')
})

async function change() {
  if (!oldPwd.value) { app.error('请输入原密码'); return }
  if (!newPwd.value || newPwd.value.length < 6) { app.error('新密码至少 6 位'); return }
  if (newPwd.value !== newPwd2.value) { app.error('两次新密码不一致'); return }
  busy.value = true
  try {
    const username = account.userInfo?.username
    if (!username) { app.error('用户未登录'); return }
    const { data } = await api.post(`/api/users/change-password?username=${username}`, { 
      oldPassword: oldPwd.value, 
      newPassword: newPwd.value 
    })
    if (data?.ok) { app.success('密码修改成功'); oldPwd.value = ''; newPwd.value = ''; newPwd2.value = '' }
    else app.error(data?.error || '修改失败')
  } catch (e) { app.error(e.response?.data?.error || '网络错误') } finally { busy.value = false }
}

async function renew() {
  const code = cardCode.value.trim()
  if (!code) { app.error('请输入卡密'); return }
  renewBusy.value = true
  try {
    const { data } = await api.post('/api/users/renew', { cardCode: code })
    if (data?.ok) {
      app.success(data.message || '续费成功')
      cardCode.value = ''
      await account.loadUserInfo()
    } else app.error(data?.error || '续费失败')
  } catch (e) { app.error(e.response?.data?.error || '续费失败') } finally { renewBusy.value = false }
}

onMounted(() => { account.loadUserInfo() })
</script>

<template>
  <div>
    <div class="subbar"><button class="icon-btn" @click="router.push('/more')">‹</button><h3>用户设置</h3></div>
    <div style="padding:12px;">
      <div class="sec-title" style="margin:2px 0 12px"><span>修改登录密码</span></div>
      <div style="display:flex;flex-direction:column;gap:12px;">
        <input v-model="oldPwd" class="field" type="password" placeholder="原密码" autocomplete="current-password">
        <input v-model="newPwd" class="field" type="password" placeholder="新密码（至少 6 位）" autocomplete="new-password">
        <input v-model="newPwd2" class="field" type="password" placeholder="确认新密码" autocomplete="new-password">
        <button :disabled="busy" style="padding:12px;border-radius:10px;background:var(--primary,#3b82f6);color:#fff;border:none;font-size:15px;font-weight:700;cursor:pointer;" @click="change">{{ busy ? '提交中…' : '修改密码' }}</button>
      </div>

      <div class="sec-title" style="margin:20px 0 12px"><span>续费 / 提升上限</span></div>
      <div style="font-size:12px;color:var(--muted);margin-bottom:10px;">
        账号上限: {{ limitText }} · 到期时间: {{ expireText }}
      </div>
      <div style="display:flex;flex-direction:column;gap:12px;">
        <input v-model="cardCode" class="field" placeholder="输入时间卡密或额度卡密">
        <button :disabled="renewBusy" style="padding:12px;border-radius:10px;background:var(--primary,#3b82f6);color:#fff;border:none;font-size:15px;font-weight:700;cursor:pointer;" @click="renew">{{ renewBusy ? '提交中…' : '使用卡密' }}</button>
      </div>
    </div>
  </div>
</template>
