<script setup lang="ts">
// xsong 個人短網址服務
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-purple-900 flex flex-col">
    <!-- Header -->
    <header class="bg-black/20 backdrop-blur-sm border-b border-white/10">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-4">
          <div class="flex items-center">
            <a href="https://xsong.us" class="flex items-center hover:opacity-80 transition-opacity">
              <div class="w-10 h-10 bg-gradient-to-r from-purple-500 to-pink-500 rounded-xl flex items-center justify-center shadow-lg">
                <span class="text-white font-bold text-lg">X</span>
              </div>
              <h1 class="ml-3 text-2xl font-bold text-white">xsong</h1>
            </a>
          </div>
          <nav class="hidden md:flex space-x-6">
            <router-link
              to="/"
              class="text-gray-300 hover:text-white px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-white/10"
              :class="{ 'text-white bg-white/20': $route.path === '/' }"
            >
              短網址
            </router-link>
            <router-link
              to="/stats"
              class="text-gray-300 hover:text-white px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-white/10"
              :class="{ 'text-white bg-white/20': $route.path === '/stats' }"
            >
              統計
            </router-link>
          </nav>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1">
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="bg-black/20 backdrop-blur-sm border-t border-white/10 mt-auto">
      <div class="max-w-6xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col md:flex-row justify-between items-center">
          <div class="flex items-center mb-4 md:mb-0">
            <div class="w-8 h-8 bg-gradient-to-r from-purple-500 to-pink-500 rounded-lg flex items-center justify-center">
              <span class="text-white font-bold text-sm">X</span>
            </div>
            <span class="ml-2 text-white font-medium">xsong</span>
          </div>
          <div class="text-gray-400 text-sm">
            © 2025 xsong. 個人短網址服務
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
/* 自定義動畫效果 */
.router-link-active {
  @apply text-white bg-white/20;
}
</style>
