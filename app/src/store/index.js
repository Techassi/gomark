import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const store = new Vuex.Store({
    state: {
        count: 0,
        loginState: localStorage.getItem('gomark-user'),
    },
    mutations: {
        increment(state) {
            state.count++
        },
        updateLoginState(state) {
            state.loginState = localStorage.getItem('gomark-user')
        },
    },
    getters: {
        getLoginState(state) {
            return state.loginState
        },
    },
})

export default store
