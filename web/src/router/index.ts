import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            redirect: '/login'
        },
        {
            path: '/login',
            component: () => import('../views/Login.vue')
        },
        {
            path: '/register',
            component: () => import('../views/Register.vue')
        },
        {
            path: '/dashboard',
            component: () => import('../views/DashBoard.vue'),
            meta: {
                requiresAuth: true
            }
        },
        {
            path: '/upload',
            component: () => import('../views/Upload.vue'),
            meta: {
                requiresAuth: true
            }
        },
        {
            path: '/gallery',
            component: () => import('../views/Gallery.vue'),
            meta: {
                requiresAuth: true
            }
        },
        {
            path: '/manage',
            component: () => import('../views/Manage.vue'),
            meta: {
                requiresAuth: true,
                requiresAdmin: true
            }
        },
        {
            path: '/:pathMatch(.*)*', // 匹配所有未匹配的路径
            component: () => import('../views/404.vue'),
        }
    ],
})

// 添加全局前置守卫，用于检查是否登录和管理员权限
router.beforeEach((to, _from, next) => {
    if (to.meta.requiresAuth) {
        // 检查本地存储中是否存在 token
        const token = localStorage.getItem('userToken');
        
        // 如果没有 token，说明用户未登录
        if (!token) {
            // 重定向到登录页
            next({
                path: '/login'
            });
        } else {
            // 检查是否需要管理员权限
            if (to.meta.requiresAdmin) {
                const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}');
                if (!userInfo.is_admin) {
                    // 如果不是管理员，重定向到首页
                    next({
                        path: '/dashboard'
                    });
                } else {
                    // 如果是管理员，允许继续跳转
                    next();
                }
            } else {
                // 如果有 token，允许继续跳转
                next();
            }
        }
    } else {
        // 如果目标路由不需要认证，直接允许跳转
        next();
    }
})

export default router