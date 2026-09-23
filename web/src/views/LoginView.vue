<script setup lang="ts">
import { LockOutlined, UserOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const route = useRoute()
const { login, register } = useAuth()

const loading = ref(false)
const mode = ref<'login' | 'register'>('login')
const form = reactive({
  username: '',
  password: '',
  nickname: '',
})

function switchMode() {
  mode.value = mode.value === 'login' ? 'register' : 'login'
}

async function onFinish() {
  loading.value = true
  try {
    if (mode.value === 'register') {
      await register(form.username, form.password, form.nickname)
      message.success('注册成功')
      mode.value = 'login'
      return
    }

    await login(form.username, form.password)
    // query 可能是 string[]，只接受单个字符串，避免开放重定向到奇怪的值
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await router.replace(redirect)
  } catch (err) {
    message.error(err instanceof Error ? err.message : mode.value === 'register' ? '注册失败' : '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <a-card :title="mode === 'login' ? '登录 xatu-agent' : '注册 xatu-agent'" class="login-card">
      <!-- @finish：校验通过后才会进来，必填规则不用在 onFinish 里再写一遍 -->
      <a-form layout="vertical" :model="form" @finish="onFinish">
        <a-form-item
          label="用户名"
          name="username"
          :rules="[{ required: true, message: '请输入用户名' }]"
        >
          <a-input v-model:value="form.username" placeholder="用户名" size="large">
            <template #prefix>
              <UserOutlined />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item
          v-if="mode === 'register'"
          label="昵称"
          name="nickname"
        >
          <a-input v-model:value="form.nickname" placeholder="昵称，可不填" size="large" />
        </a-form-item>

        <a-form-item
          label="密码"
          name="password"
          :rules="[{ required: true, message: '请输入密码' }]"
        >
          <a-input-password v-model:value="form.password" placeholder="密码" size="large">
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-button type="primary" html-type="submit" size="large" block :loading="loading">
            {{ mode === 'login' ? '登录' : '注册' }}
          </a-button>
        </a-form-item>
      </a-form>

      <a-button type="link" block @click="switchMode">
        {{ mode === 'login' ? '没有账号？去注册' : '已有账号？去登录' }}
      </a-button>
    </a-card>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}

.login-card {
  width: 400px;
}
</style>
