import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import store from '../store'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
    },
    {
        path: '/all',
        name: 'All',
        component: () =>
            import(/* webpackChunkName: "all" */ '../views/All.vue'),
    },
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes,
})

router.beforeEach((to, from, next) => {
    // redirect to login page if not logged in and trying to access a restricted page
    const publicPages = ['/login']
    const authRequired = !publicPages.includes(to.path)
    const loggedIn = store.getters.getLoginState

    if (authRequired && !loggedIn) {
        return next('/login')
    }

    next()
})

export default router
