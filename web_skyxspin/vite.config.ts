import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import imagemin from 'vite-plugin-imagemin'
import viteCompression from 'vite-plugin-compression'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { createStyleImportPlugin, ElementPlusResolve } from 'vite-plugin-style-import'
let timeStamp = new Date().getTime()
// https://vitejs.dev/config/
export default defineConfig({
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
      },
      mangle: true,
    },

    chunkSizeWarningLimit: 1500,
    rollupOptions: {
      output: {
        // 最小化拆分包
        manualChunks(id) {
          if (id.includes('node_modules')) {
            return id.toString().split('node_modules/')[1].split('/')[0].toString()
          }
        },
        // 用于从入口点创建的块的打包输出格式[name]表示文件名,[hash]表示该文件内容hash值
        entryFileNames: `js/[name].[hash]${timeStamp}.js`,
        // 用于命名代码拆分时创建的共享块的输出命名
        // chunkFileNames: `js/[name].[hash]${timeStamp}.js`,
        // 用于输出静态资源的命名，[ext]表示文件扩展名
        assetFileNames: `[ext]/[name].[hash]${timeStamp}.[ext]`,
        // 拆分js到模块文件夹
        chunkFileNames: chunkInfo => {
          const facadeModuleId = chunkInfo.facadeModuleId ? chunkInfo.facadeModuleId.split('/') : []
          const fileName = facadeModuleId[facadeModuleId.length - 2] || '[name]'
          return `js/${fileName}/[name].[hash]${timeStamp}.js`
        }
      }
    }
  },
  base:"/",
  plugins: [
    vue(),
    imagemin({
      gifsicle: {
        interlaced: true,
      },
      optipng: {
        optimizationLevel: 5,
      },
      mozjpeg: {
        quality: 75,
      },
      pngquant: {
        quality: [0.6, 0.8],
        speed: 4,
      },
      webp: {
        quality: 75,
      },
    }),

    AutoImport({
      resolvers: [ElementPlusResolver()]
    }),
    Components({
      resolvers: [ElementPlusResolver()]
    }),
    createStyleImportPlugin({
      resolves: [ElementPlusResolve()]
    })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    host:true,
    proxy: {
      '/api': {
        target: "https://localhost:11100",
        rewrite: (path) => {
          return path.replace("/api/", "/")
        }
      },
    }
  }
})
