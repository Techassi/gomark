import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
// import store from '../store'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
        meta: {
            title: 'Home | Gomark',
            // metaTags: [
            //   {
            //     name: 'description',
            //     content: 'The home page of our example app.'
            //   },
            //   {
            //     property: 'og:description',
            //     content: 'The home page of our example app.'
            //   }
            // ]
        },
    },
    {
        path: '/all',
        name: 'All',
        component: () =>
            import(/* webpackChunkName: "all" */ '../views/All.vue'),
        meta: {
            title: 'All | Gomark',
            // metaTags: [
            //   {
            //     name: 'description',
            //     content: 'The home page of our example app.'
            //   },
            //   {
            //     property: 'og:description',
            //     content: 'The home page of our example app.'
            //   }
            // ]
        },
    },
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes,
})

router.beforeEach((to, from, next) => {
    // redirect to login page if not logged in and trying to access a restricted page
    // const publicPages = ['/login']
    // const authRequired = !publicPages.includes(to.path)
    // const loggedIn = store.getters.getLoginState

    // if (authRequired && !loggedIn) {
    //     return next('/login')
    // }
    document.title = to.meta.title
    next()
})

export default router
