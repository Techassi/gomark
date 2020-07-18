<template>
    <div class="gomark">
        <TopNav></TopNav>
        <main>
            <router-view></router-view>
        </main>
    </div>
</template>

<script>
import TopNav from '@/components/nav/nav-top.vue'
export default {
    name: 'App',
    components: {
        TopNav,
    },
}
</script>

<style lang="scss">
body {
    background-color: #ffffff;
}

main {
    width: 100%;
    height: auto;
    max-width: 1280px;
    padding: 0 32px;
    margin: 56px auto 0 auto;
}

* {
    color: $primary-black;
}
</style>
