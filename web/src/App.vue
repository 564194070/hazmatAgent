<script setup lang="ts">
import { ref } from 'vue'
import { Bubble, Sender, Welcome, useXAgent, useXChat } from 'ant-design-x-vue'

const text = ref('')

const [agent] = useXAgent<string, { message: string }, string>({
  request: async ({ message }, { onSuccess }) => {
    // 第一阶段：假装后端已回复
    await new Promise((r) => setTimeout(r, 800))
    onSuccess([`（mock）你说的是：${message}`])
  },
})

const { messages, onRequest } = useXChat({
  agent: agent.value,
  defaultMessages: [{ message: '你好，我是 xatu-agent。', status: 'success' }],
})

function onSubmit(value: string) {
  onRequest(value)
  text.value = ''
}
</script>

<template>
  <div style="max-width: 720px; margin: 40px auto;">
    <Welcome title="xatu-agent" description="问我任何问题" />
    <Bubble.List
      :items="messages.map(({ id, message, status }) => ({
        key: id,
        role: status === 'local' ? 'user' : 'ai',
        content: message,
        loading: status === 'loading',
      }))"
      style="height: 360px; margin: 16px 0;"
    />
    <Sender v-model:value="text" @submit="onSubmit" />
  </div>
</template>
