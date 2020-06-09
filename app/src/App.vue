<template>
    <main>
        <router-view></router-view>
    </main>
</template>

<script>
export default {
    name: 'App',
    components: {},
}
</script>

<style lang="scss">
body {
    background-color: $primary-bg;
}

main {
    width: 100%;
    height: auto;
    padding: 84px 56px 0 84px;
}
</style>
