import {createRouter, createWebHistory} from 'vue-router'

import { GetSystemAllData } from "@/util/generatorData";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: (resolve)=> import("@/views/HomeView.vue", resolve)
    }
    // {
    //   path: '/game/:gameId',
    //   name: 'game',
    //   component: ()=> import("../views/game.vue")
    // },
  ]
})


router.beforeEach(async  (to, from, next)=>{

  // if (to.path == "/"){
  //   await GetSystemAllData()
  // }
  next();
})

export default router
