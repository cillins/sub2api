<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page -->
  <div
    v-else
    class="relative flex min-h-screen flex-col overflow-hidden bg-gray-50 dark:bg-dark-950"
  >
    <!-- Background Decorations -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div
        class="absolute -right-20 -top-20 h-80 w-80 rounded-full bg-primary-400/10 blur-3xl"
      ></div>
      <div
        class="absolute -bottom-20 -left-20 h-72 w-72 rounded-full bg-primary-500/8 blur-3xl"
      ></div>
    </div>

    <!-- Header -->
    <header class="relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <!-- Logo -->
        <div class="flex items-center">
          <div class="h-10 w-10 overflow-hidden rounded-xl shadow-md">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
        </div>

        <!-- Nav Actions -->
        <div class="flex items-center gap-3">
          <!-- Language Switcher -->
          <LocaleSwitcher />

          <!-- Doc Link -->
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>

          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>

          <!-- Login / Dashboard Button -->
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="inline-flex items-center gap-1.5 rounded-full bg-gray-900 py-1 pl-1 pr-2.5 transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            <span
              class="flex h-5 w-5 items-center justify-center rounded-full bg-primary-600 text-[10px] font-semibold text-white"
            >
              {{ userInitial }}
            </span>
            <span class="text-xs font-medium text-white">{{ t('home.dashboard') }}</span>
            <svg
              class="h-3 w-3 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25"
              />
            </svg>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="inline-flex items-center rounded-full bg-gray-900 px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="relative z-10 flex-1 px-6 py-16">
      <div class="mx-auto max-w-6xl">
        <!-- Hero Section - Left/Right Layout -->
        <div class="mb-12 flex flex-col items-center justify-between gap-12 lg:flex-row lg:gap-16">
          <!-- Left: Text Content -->
          <div class="flex-1 text-center lg:text-left">
            <h1
              class="mb-4 text-4xl font-bold text-gray-900 dark:text-white md:text-5xl lg:text-6xl"
            >
              {{ siteName }}
            </h1>
            <p class="mb-8 text-lg text-gray-600 dark:text-dark-300 md:text-xl">
              {{ siteSubtitle }}
            </p>

            <!-- CTA Button -->
            <div>
              <router-link
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="btn btn-primary px-6 py-2.5 text-sm"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
              </router-link>
            </div>
          </div>

          <!-- Right: Code Example -->
          <div class="flex flex-1 justify-center lg:justify-end">
            <div class="w-full max-w-md rounded-xl border border-gray-200/60 bg-gray-950 p-5 font-mono text-[13px] leading-relaxed shadow-sm dark:border-dark-700/40">
              <div class="mb-3 flex items-center gap-2 text-gray-500">
                <span class="inline-block h-2 w-2 rounded-full bg-green-500"></span>
                <span class="text-xs">API Request</span>
              </div>
              <div class="space-y-1 text-gray-300">
                <div><span class="text-blue-400">POST</span> <span class="text-gray-400">/v1/chat/completions</span></div>
                <div class="text-gray-500">{</div>
                <div class="pl-4"><span class="text-amber-400">"model"</span>: <span class="text-green-400">"claude-sonnet-4-20250514"</span>,</div>
                <div class="pl-4"><span class="text-amber-400">"messages"</span>: [</div>
                <div class="pl-8">{ <span class="text-amber-400">"role"</span>: <span class="text-green-400">"user"</span>,</div>
                <div class="pl-10"><span class="text-amber-400">"content"</span>: <span class="text-green-400">"Hello"</span> }</div>
                <div class="pl-4">]</div>
                <div class="text-gray-500">}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Feature Tags -->
        <div class="mb-10 flex flex-wrap items-center justify-center gap-3 text-sm text-gray-500 dark:text-dark-400">
          <span>{{ t('home.tags.subscriptionToApi') }}</span>
          <span class="text-gray-300 dark:text-dark-600">·</span>
          <span>{{ t('home.tags.stickySession') }}</span>
          <span class="text-gray-300 dark:text-dark-600">·</span>
          <span>{{ t('home.tags.realtimeBilling') }}</span>
        </div>

        <!-- Features -->
        <div class="mb-12 grid gap-5 md:grid-cols-2">
          <!-- Feature 1: Unified Gateway (larger) -->
          <div
            class="rounded-xl border border-gray-200/60 bg-white/70 p-6 transition-shadow hover:shadow-md dark:border-dark-700/40 dark:bg-dark-800/50 md:col-span-2 md:p-8"
          >
            <div class="flex flex-col gap-4 md:flex-row md:items-start md:gap-8">
              <div class="flex-1">
                <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
                  {{ t('home.features.unifiedGateway') }}
                </h3>
                <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
                  {{ t('home.features.unifiedGatewayDesc') }}
                </p>
              </div>
              <div class="flex flex-wrap shrink-0 items-center gap-2 text-sm text-gray-500 dark:text-dark-400">
                <span class="rounded-md bg-blue-50 px-2 py-1 text-xs font-medium text-blue-700 dark:bg-blue-900/20 dark:text-blue-400">Claude</span>
                <span class="rounded-md bg-green-50 px-2 py-1 text-xs font-medium text-green-700 dark:bg-green-900/20 dark:text-green-400">GPT</span>
                <span class="rounded-md bg-sky-50 px-2 py-1 text-xs font-medium text-sky-700 dark:bg-sky-900/20 dark:text-sky-400">Gemini</span>
                <span class="rounded-md bg-purple-50 px-2 py-1 text-xs font-medium text-purple-700 dark:bg-purple-900/20 dark:text-purple-400">DeepSeek</span>
              </div>
            </div>
          </div>

          <!-- Feature 2: Account Pool -->
          <div
            class="rounded-xl border border-gray-200/60 bg-white/70 p-6 transition-shadow hover:shadow-md dark:border-dark-700/40 dark:bg-dark-800/50"
          >
            <h3 class="mb-2 text-base font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.multiAccount') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.multiAccountDesc') }}
            </p>
          </div>

          <!-- Feature 3: Billing & Quota -->
          <div
            class="rounded-xl border border-gray-200/60 bg-white/70 p-6 transition-shadow hover:shadow-md dark:border-dark-700/40 dark:bg-dark-800/50"
          >
            <h3 class="mb-2 text-base font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.balanceQuota') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.balanceQuotaDesc') }}
            </p>
          </div>
        </div>

        <!-- Supported Providers -->
        <div class="mb-8 text-center">
          <h2 class="mb-3 text-2xl font-bold text-gray-900 dark:text-white">
            {{ t('home.providers.title') }}
          </h2>
          <p class="text-sm text-gray-600 dark:text-dark-400">
            {{ t('home.providers.description') }}
          </p>
        </div>

        <div class="mb-16 flex flex-wrap items-center justify-center gap-3">
          <span class="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200">
            <span class="h-1.5 w-1.5 rounded-full bg-orange-500"></span>
            {{ t('home.providers.claude') }}
          </span>
          <span class="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200">
            <span class="h-1.5 w-1.5 rounded-full bg-green-500"></span>
            {{ t('home.providers.gpt') }}
          </span>
          <span class="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200">
            <span class="h-1.5 w-1.5 rounded-full bg-blue-500"></span>
            {{ t('home.providers.gemini') }}
          </span>
          <span class="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200">
            <span class="h-1.5 w-1.5 rounded-full bg-rose-500"></span>
            {{ t('home.providers.antigravity') }}
          </span>
          <span class="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200">
            <span class="h-1.5 w-1.5 rounded-full bg-purple-500"></span>
            {{ t('home.providers.deepseek') }}
          </span>
          <span class="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200">
            <span class="h-1.5 w-1.5 rounded-full bg-amber-500"></span>
            {{ t('home.providers.qwen') }}
          </span>
          <span class="inline-flex items-center gap-1.5 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200">
            <span class="h-1.5 w-1.5 rounded-full bg-cyan-500"></span>
            {{ t('home.providers.seedance') }}
          </span>
          <span class="inline-flex items-center gap-1.5 rounded-lg border border-gray-200/50 bg-white/50 px-4 py-2 text-sm text-gray-400 dark:border-dark-700/50 dark:bg-dark-800/50 dark:text-dark-500">
            {{ t('home.providers.more') }}
            <span class="text-xs">{{ t('home.providers.soon') }}</span>
          </span>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-8 dark:border-dark-800/50">
      <div
        class="mx-auto flex max-w-6xl flex-col items-center justify-center gap-4 text-center sm:flex-row sm:text-left"
      >
        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
        <div class="flex items-center gap-4">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
          <a
            :href="githubUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            GitHub
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// GitHub URL
const githubUrl = 'https://github.com/Wei-Shaw/sub2api'

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Initialize theme
function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()

  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
</style>
