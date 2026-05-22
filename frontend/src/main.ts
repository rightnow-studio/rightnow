import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import router from './router'
import App from './App.vue'
import './style.css'

import zh from './locales/zh.json'
import en from './locales/en.json'

const messages = { 'zh-CN': zh, 'zh': zh, 'en-US': en, 'en': en }

const browserLang = navigator.language || 'zh-CN'
const fallbackLocale = 'zh-CN'

const i18n = createI18n({
  legacy: false,
  locale: browserLang,
  fallbackLocale,
  messages,
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(i18n)
app.mount('#app')
