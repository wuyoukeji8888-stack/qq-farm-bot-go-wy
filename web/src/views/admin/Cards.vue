<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { useAppStore } from '@/stores/app'
import { useRouter } from 'vue-router'

const app = useAppStore()
const router = useRouter()

const cards = ref([])
const loading = ref(false)
const generating = ref(false)
const genCount = ref(10)
const genDays = ref(30)
const genType = ref('time')
const genValue = ref(1)
const genDesc = ref('')
const selected = ref([])

async function loadCards() {
  loading.value = true
  try {
    const { data } = await api.get('/api/admin/cards')
    cards.value = data.cards || []
  } catch (e) {
    app.error(e.response?.data?.error || '加载失败')
  } finally {
    loading.value = false
  }
}

async function generate() {
  if (genCount.value <= 0 || genCount.value > 1000) {
    app.error('生成数量需在 1-1000 之间')
    return
  }
  if (genType.value === 'time' && genDays.value <= 0) {
    app.error('天数必须大于 0')
    return
  }
  if (genType.value === 'quota' && genValue.value <= 0) {
    app.error('额度必须大于 0')
    return
  }
  generating.value = true
  try {
    const payload = {
      action: 'generate',
      cardType: genType.value,
      count: genCount.value,
      description: genDesc.value
    }
    if (genType.value === 'quota') {
      payload.value = genValue.value
      payload.description = genDesc.value || `额度卡密 +${genValue.value} 账号`
    } else {
      payload.days = genDays.value
      payload.description = genDesc.value || `${genDays.value}天卡密`
    }
    const { data } = await api.post('/api/admin/cards', payload)
    if (data.ok) {
      app.success(`成功生成 ${data.codes.length} 个卡密`)
      genDesc.value = ''
      await loadCards()
    }
  } catch (e) {
    app.error(e.response?.data?.error || '生成失败')
  } finally {
    generating.value = false
  }
}

async function toggleCards(enabled) {
  if (selected.value.length === 0) {
    app.error('请先选择卡密')
    return
  }
  try {
    await api.post('/api/admin/cards', {
      action: 'toggle',
      codes: selected.value,
      enabled
    })
    app.success(enabled ? '已启用' : '已禁用')
    selected.value = []
    await loadCards()
  } catch (e) {
    app.error(e.response?.data?.error || '操作失败')
  }
}

async function deleteCards() {
  if (selected.value.length === 0) {
    app.error('请先选择卡密')
    return
  }
  if (!confirm(`确定删除选中的 ${selected.value.length} 个卡密？`)) return
  try {
    const { data } = await api.post('/api/admin/cards', {
      action: 'delete',
      codes: selected.value
    })
    const n = data.deleted
    if (n > 0) {
      app.success(`已删除 ${n} 个卡密`)
    } else {
      app.error('未删除任何卡密，请重新勾选后再试')
    }
    selected.value = []
    await loadCards()
  } catch (e) {
    app.error(e.response?.data?.error || '删除失败')
  }
}

function toggleSelect(code) {
  const i = selected.value.indexOf(code)
  if (i >= 0) {
    selected.value.splice(i, 1)
  } else {
    selected.value.push(code)
  }
}

function selectAll() {
  if (selected.value.length === cards.value.length) {
    selected.value = []
  } else {
    selected.value = cards.value.map(c => c.code)
  }
}

function formatTime(ms) {
  if (!ms) return '-'
  const d = new Date(ms)
  return d.toLocaleString('zh-CN')
}

function copyCode(code) {
  navigator.clipboard.writeText(code)
  app.success('已复制')
}

function typeBtnStyle(on) {
  return {
    flex: 1,
    padding: '8px 12px',
    borderRadius: '8px',
    border: '1px solid var(--border)',
    cursor: 'pointer',
    fontSize: '13px',
    background: on ? 'var(--primary,#3b82f6)' : 'var(--card-strong)',
    color: on ? '#fff' : 'inherit'
  }
}

function cardTypeLabel(card) {
  if (card.type === 'quota') return `额度 +${card.value || 1}`
  if (card.days) return `${card.days}天`
  return '时间卡'
}

onMounted(loadCards)
</script>

<template>
  <div>
    <div class="subbar">
      <button class="icon-btn" @click="router.push('/more')">‹</button>
      <h3>卡密管理</h3>
    </div>
    <div style="padding:12px;">
      <!-- 生成卡密 -->
      <div class="sec-title" style="margin:2px 0 12px"><span>批量生成卡密</span></div>
      <div style="display:flex;flex-direction:column;gap:10px;">
        <div style="display:flex;gap:8px;">
          <button :style="typeBtnStyle(genType==='time')" @click="genType='time'">时间卡</button>
          <button :style="typeBtnStyle(genType==='quota')" @click="genType='quota'">额度卡</button>
        </div>
        <div style="display:flex;gap:8px;">
          <input v-model.number="genCount" type="number" class="field" placeholder="数量" style="flex:1">
          <input v-if="genType==='time'" v-model.number="genDays" type="number" class="field" placeholder="天数" style="flex:1">
          <input v-else v-model.number="genValue" type="number" class="field" placeholder="增加账号数" style="flex:1">
        </div>
        <input v-model="genDesc" class="field" placeholder="描述（可选）">
        <button :disabled="generating" style="padding:12px;border-radius:10px;background:var(--primary,#3b82f6);color:#fff;border:none;font-size:15px;font-weight:700;cursor:pointer;" @click="generate">
          {{ generating ? '生成中...' : (genType === 'quota' ? '生成额度卡' : '生成时间卡') }}
        </button>
      </div>

      <!-- 操作栏 -->
      <div class="sec-title" style="margin:16px 0 12px">
        <span>卡密列表（{{ cards.length }}）</span>
      </div>
      <div style="display:flex;gap:8px;margin-bottom:12px;flex-wrap:wrap;">
        <button style="padding:8px 12px;border-radius:8px;background:var(--card-strong);border:1px solid var(--border);cursor:pointer;font-size:13px;" @click="selectAll">全选</button>
        <button style="padding:8px 12px;border-radius:8px;background:#22c55e;color:#fff;border:none;cursor:pointer;font-size:13px;" @click="toggleCards(true)">启用</button>
        <button style="padding:8px 12px;border-radius:8px;background:#f59e0b;color:#fff;border:none;cursor:pointer;font-size:13px;" @click="toggleCards(false)">禁用</button>
        <button style="padding:8px 12px;border-radius:8px;background:#ef4444;color:#fff;border:none;cursor:pointer;font-size:13px;" @click="deleteCards">删除</button>
        <button style="padding:8px 12px;border-radius:8px;background:var(--card-strong);border:1px solid var(--border);cursor:pointer;font-size:13px;" @click="loadCards">刷新</button>
      </div>

      <!-- 卡密列表 -->
      <div v-if="loading" style="text-align:center;padding:24px;color:var(--muted)">加载中...</div>
      <div v-else-if="cards.length === 0" style="text-align:center;padding:24px;color:var(--muted)">暂无卡密</div>
      <div v-else style="display:flex;flex-direction:column;gap:8px;">
        <div v-for="card in cards" :key="card.code" style="padding:12px;border-radius:10px;background:var(--card-strong);border:1px solid var(--border);">
          <div style="display:flex;align-items:center;gap:8px;">
            <input type="checkbox" :checked="selected.includes(card.code)" @change="toggleSelect(card.code)">
            <div style="flex:1;min-width:0;">
              <div style="font-family:monospace;font-size:13px;word-break:break-all;">{{ card.code }}</div>
              <div style="font-size:12px;color:var(--muted);margin-top:4px;">
                {{ cardTypeLabel(card) }} | {{ card.description || '-' }} | {{ card.usedBy ? '已使用' : '未使用' }} | {{ card.enabled ? '启用' : '禁用' }}
              </div>
            </div>
            <button style="padding:4px 8px;border-radius:6px;background:transparent;border:1px solid var(--border);cursor:pointer;font-size:12px;" @click="copyCode(card.code)">复制</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
