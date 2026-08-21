import {defineConfig} from 'vite'
import Vue from '@vitejs/plugin-vue'
import Papa from 'papaparse'
import VueJsx from '@vitejs/plugin-vue-jsx'
import viteCompressoion from 'vite-plugin-compression'
import PrerenderSPAPlugin from 'prerender-spa-plugin';
function CSV() {
    const fileRegex = /\.csv$/
    return {
        name: 'csv',
        transform(src, id) {
            if (fileRegex.test(id)) {
                let v = Papa.parse(src)
                return {
                    code: "export default " + JSON.stringify(v.data),
                    map: null // 如果可行将提供 source map
                }
            }
        }
    }
}

// https://vitejs.dev/config/
export default defineConfig({
    base: "./",
    plugins: [
        Vue(),
        VueJsx(),
        CSV(),
        // imagePresets({
        //     thumbnail: widthPreset({
        //         class: 'img thumb',
        //         loading: 'lazy',
        //         widths: [48, 96],
        //         formats: {
        //             webp: { quality: 50 },
        //             jpg: { quality: 70 },
        //             png: { quality: 70 },
        //         },
        //     }),
        // }),
        // viteCompressoion({
        //     threshold: 1024, // 文件容量大于这个值进行压缩
        //     algorithm: 'gzip', // 压缩方式
        //     ext: 'gz', // 后缀名
        // }),
        PrerenderSPAPlugin({
            // 指定需要预渲染的页面
            routes: [
                '/',
                '/betLog',
            ],
        }),
    ],
    build: {
        assetsInlineLimit: 0,
        reportCompressedSize:false,
        target: ['es2022', /* ... 其他目标 */], // 保留您的目标浏览器
        // 添加以下选项：
        minify: 'esbuild', // 使用 esbuild 进行缩小（推荐）
        rollupOptions: {
            output: {
                entryFileNames: 'js/[hash].js',
                chunkFileNames: 'js/[hash].js',
                assetFileNames: "data/[hash:1]/[hash][extname]",
                manualChunks: (id) => {
                    if (id.includes('node_modules')) {
                        return id.toString().split('node_modules/')[1].split('/')[0].toString();
                    }
                },
            },
            external: ['element-plus/lib/locale/lang/zh-CN'],
        },
        terserOptions: {
            format: {
                comments: false,
            },
            compress: {
                drop_console: true,
                drop_debugger: true,
            },
        },
    },

    resolve: {
        alias: {
            '@': '/src',
            'vue-i18n': 'vue-i18n/dist/vue-i18n.cjs.js',
        },
    },
    css: {
        // CSS 预处理器
        preprocessorOptions: {
            // 定义全局 SCSS 变量
            scss: {
                javascriptEnabled: true,
                additionalData: `
            @use "@/assets/style/variables.scss" as *;
          `,
            },
        },
    },
    server: {
        host: true,
        port: 4568,
        cors: {
            origin: ['https://admin.ppslot.site', 'https://ppslot.site', 'https://h5.ppslot.site', 'https://api.cocodr.xyz', 'https://history.cocodr.xyz', 'https://pg.bxbet.asia', 'https://m-pg.bxbet.asia', 'https://static-pg.bxbet.asia', 'http://localhost:4568', 'http://localhost:3000'],
            credentials: true,
            methods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS'],
            allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With', 'Accept'],
            exposedHeaders: ['Content-Type', 'Authorization'],
            maxAge: 86400
        },
        proxy: {
            '/api': {
                target: 'https://api.cocodr.xyz:11100',
                changeOrigin: true,
                secure: false,
                rewrite: (path) => {
                    return path.replace("/api/", "/")
                },
                headers: {
                    'Access-Control-Allow-Origin': '*',
                    'Access-Control-Allow-Methods': 'GET,POST,PUT,DELETE,PATCH,OPTIONS',
                    'Access-Control-Allow-Headers': 'Content-Type,Authorization,X-Requested-With,Accept'
                }
            },
        }
    }
})
