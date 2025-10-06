import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/login',
            component: () => import('../views/Login.vue')
        },
        {
            path: '/:pathMatch(.*)*', // 匹配所有未匹配的路径
            component: () => import('../views/404.vue'),
        }
    ]
})

export default router