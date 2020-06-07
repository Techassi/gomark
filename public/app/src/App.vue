<template>
    <main>
        <TopNav></TopNav>
        <SideNav></SideNav>
    </main>
</template>

<script>
import TopNav from '@/components/nav/TopNav.vue'
import SideNav from '@/components/nav/SideNav.vue'

export default {
    name: 'App',
    components: {
        TopNav,
        SideNav,
    },
}
</script>

<style lang="scss">
body {
    background-color: $primary-bg;
}
</style>
