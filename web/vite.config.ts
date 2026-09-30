import { defineConfig } from 'vite'

export default defineConfig({
  server: {
    allowedHosts: ['.monkeycode-ai.online'],
    proxy: {
      '/api': {
        target: 'http://localhost:3009',
        changeOrigin: true,
      },
    },
  },
})