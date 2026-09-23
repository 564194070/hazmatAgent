<script setup lang="ts">
import { LogoutOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import {
  Bubble,
  Conversations,
  Sender,
  Welcome,
  useXAgent,
  useXChat,
  type Conversation,
} from 'ant-design-x-vue'
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { createGreeting, useSessions } from '../composables/useSessions'

const router = useRouter()
const { userId, logout } = useAuth()
const sessions = useSessions()

const text = ref('')
const renameOpen = ref(false)
const renameKey = ref('')
const renameText = ref('')

const [agent] = useXAgent<string, { message: string }, string>({
  request: async ({ message }, { onSuccess }) => {
    // 调用发生时再读，不能在 setup 时把 sessionId 抄成常量
    const sessionId = sessions.activeKey.value
    const currentUserId = userId.value ?? ''

    // 以后换成真实接口，Vite 已把 /api 代理到 127.0.0.1:8080：
    // const res = await fetch('/api/chat', {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify({ userId: currentUserId, sessionId, message }),
    // })
    // const data = await res.json()
    // onSuccess([data.reply])
    void currentUserId
    void sessionId

    await new Promise((resolve) => setTimeout(resolve, 800))
    onSuccess([`（mock）你说的是：${message}`])
  },
})

const { messages, onRequest, setMessages } = useXChat<string, string, { message: string }, string>({
  agent: agent.value,
  defaultMessages: [
    { id: 'welcome', message: '你好，我是 xatu-agent。', status: 'success' },
  ],
})

// 输入过程中把气泡同步回当前会话，切换走再回来内容还在
watch(
  messages,
  (list) => {
    sessions.saveMessages(sessions.activeKey.value, list)
  },
  { deep: true },
)

function loadActive() {
  const session = sessions.get(sessions.activeKey.value)
  setMessages(session?.messages ?? createGreeting())
}

function onActiveChange(key: string) {
  if (key === sessions.activeKey.value) return
  sessions.saveMessages(sessions.activeKey.value, messages.value)
  sessions.activeKey.value = key
  loadActive()
}

function onCreate() {
  sessions.saveMessages(sessions.activeKey.value, messages.value)
  sessions.create()
  loadActive()
}

function onSubmit(value: string) {
  const content = value.trim()
  if (!content) return
  sessions.touchTitle(sessions.activeKey.value, content)
  onRequest(content)
  text.value = ''
}

// menu 用函数形式：点哪条会话，菜单回调里就有那条的 key
function sessionMenu(conversation: Conversation) {
  return {
    items: [
      { key: 'rename', label: '重命名' },
      { key: 'delete', label: '删除', danger: true },
    ],
    onClick: ({ key }: { key: string | number }) => {
      if (key === 'rename') {
        renameKey.value = conversation.key
        renameText.value = typeof conversation.label === 'string' ? conversation.label : ''
        renameOpen.value = true
        return
      }
      if (key === 'delete') {
        Modal.confirm({
          title: '删除这个会话？',
          content: '会话目前只在内存里，删除后无法恢复。',
          okText: '删除',
          okType: 'danger',
          cancelText: '取消',
          onOk: () => {
            sessions.saveMessages(sessions.activeKey.value, messages.value)
            sessions.remove(conversation.key)
            loadActive()
          },
        })
      }
    },
  }
}

function applyRename() {
  sessions.rename(renameKey.value, renameText.value)
  renameOpen.value = false
}

async function onLogout() {
  logout()
  // 下一个登录用户不应看到上一个用户的内存会话
  sessions.reset()
  await router.replace('/login')
}
</script>

<template>
  <a-layout class="chat-shell">
    <a-layout-sider class="sider" theme="light" :width="280">
      <div class="sider-head">
        <span class="brand">xatu-agent</span>
        <a-button type="text" @click="onLogout">
          <template #icon>
            <LogoutOutlined />
          </template>
          退出
        </a-button>
      </div>

      <a-button type="primary" block class="new-session" @click="onCreate">
        <template #icon>
          <PlusOutlined />
        </template>
        新建会话
      </a-button>

      <!--
        groupable 打开后按 item.group 分组。
        sort 把分组固定成今天 / 昨天 / 更早。
        active-key 受控：当前会话以 useSessions 为准。
      -->
      <Conversations
        class="session-list"
        :items="sessions.items.value"
        :active-key="sessions.activeKey.value"
        :groupable="{ sort: sessions.groupSort }"
        :menu="sessionMenu"
        @active-change="onActiveChange"
      />
    </a-layout-sider>

    <a-layout-content class="chat-main">
      <Welcome title="xatu-agent" description="问我任何问题" />
      <Bubble.List
        class="bubble-wrap"
        :items="messages.map(({ id, message, status }) => ({
          key: id,
          role: status === 'local' ? 'user' : 'ai',
          content: message,
          loading: status === 'loading',
        }))"
      />
      <Sender v-model:value="text" @submit="onSubmit" />
    </a-layout-content>
  </a-layout>

  <a-modal
    v-model:open="renameOpen"
    title="重命名会话"
    ok-text="保存"
    cancel-text="取消"
    @ok="applyRename"
  >
    <a-input v-model:value="renameText" placeholder="会话名称" @press-enter="applyRename" />
  </a-modal>
</template>

<style scoped>
.chat-shell {
  height: 100vh;
}

.sider {
  height: 100vh;
  border-right: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
}

.sider :deep(.ant-layout-sider-children) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sider-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 16px 8px;
}

.brand {
  font-weight: 600;
}

.new-session {
  margin: 0 16px 12px;
  width: calc(100% - 32px);
}

.session-list {
  flex: 1;
  overflow: auto;
}

.chat-main {
  height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 24px;
  box-sizing: border-box;
}

.bubble-wrap {
  flex: 1;
  min-height: 0;
  margin: 16px 0;
}
</style>