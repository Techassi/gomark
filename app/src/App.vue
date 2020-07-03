<template>
    <div class="gomark">
        <main>
            <router-view></router-view>
        </main>
    </div>
</template>

<script>
export default {
    name: 'App',
    components: {},
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
    margin: 0 auto;
}

* {
    color: $primary-black;
}
</style>
