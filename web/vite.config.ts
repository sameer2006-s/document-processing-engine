import { defineConfig } from 'vite'
import react, { reactCompilerPreset } from '@vitejs/plugin-react'
import babel from '@rolldown/plugin-babel'

export default defineConfig({
  plugins: [
    react(),
    babel({ presets: [reactCompilerPreset()] }),
  ],
  server: {
    proxy: {
      '/register': { target: 'http://localhost:7070', changeOrigin: true },
      '/login': { target: 'http://localhost:7070', changeOrigin: true },
      '/upload': { target: 'http://localhost:7070', changeOrigin: true },
      '/documents': { target: 'http://localhost:7070', changeOrigin: true },
      '/get-file': { target: 'http://localhost:7070', changeOrigin: true },
      '/search-my-files': { target: 'http://localhost:7070', changeOrigin: true },
      '/chat': { target: 'http://localhost:7070', changeOrigin: true },
    },
  },
  build: { outDir: 'dist' },
})
