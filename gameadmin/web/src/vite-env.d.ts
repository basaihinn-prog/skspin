/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}


declare module "*.csv" {
  const value: string
  export default value
}


type JsMap<K, V> = { [key in K]: V }