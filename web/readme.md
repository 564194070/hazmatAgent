# 用 Vite 生成 Vue 3 + TypeScript 工程
npm create vite@latest . -- --template vue-ts

    index.html
    浏览器入口，里面挂一个 <div id="app">
    src/main.ts
    创建 Vue 应用，挂到 #app
    src/App.vue
    根组件，聊天页就写在这里
    vite.config.ts
    开发服务器、代理、打包配置
    package.json
    依赖和 dev / build 脚本

# 安装UI依赖

npm install ant-design-vue ant-design-x-vue @ant-design/icons-vue


# 在入口里接入 Ant Design Vue
    src/main.ts：

    import { createApp } from 'vue'
    import Antd from 'ant-design-vue'
    import 'ant-design-vue/dist/reset.css'
    import 'ant-design-x-vue/dist/style.css'
    import App from './App.vue'
    createApp(App).use(Antd).mount('#app')


# 配置开发代理

    vite.config.ts

    import { defineConfig } from 'vite'
    import vue from '@vitejs/plugin-vue'
    export default defineConfig({
    plugins: [vue()],
    server: {
        proxy: {
        '/api': 'http://127.0.0.1:8080',
        },
    },
    })
    浏览器访问的是 localhost:5173，Go 以后会在 8080。直接 fetch('http://127.0.0.1:8080/api/chat') 会跨域。代理的意思是：前端写 fetch('/api/chat')，Vite 帮你转到 8080，浏览器以为还在同源。


# 小聊天外壳

src/App.vue


<script setup lang="ts">
import { ref } from 'vue'
import { Bubble, Sender, Welcome } from 'ant-design-x-vue'
const text = ref('')
const items = ref([
  { key: '1', role: 'ai', content: '你好，我是 xatu-agent。' },
])
</script>
<template>
  <div style="max-width: 720px; margin: 40px auto;">
    <Welcome title="xatu-agent" description="问我任何问题" />
    <Bubble.List :items="items" style="height: 360px; margin: 16px 0;" />
    <Sender v-model:value="text" />
  </div>
</template>


# 测试
npm run dev

# 安装路由

npm install vue-router