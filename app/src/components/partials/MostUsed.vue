<template>
    <div class="main__most-used">
        <h2 v-on:click="increment">{{ this.$store.state.count }} bookmarks</h2>
    </div>
</template>

<script>
export default {
    name: 'MostUsed',
    methods: {
        increment() {
            this.$store.commit('increment')
            console.log(this.$store.state.count)
        },
    },
    mounted() {
        console.log('Hello')
    },
}
</script>

<style lang="scss" scoped>
.main__most-used {
    h2 {
        margin: 0;
        font-size: 2rem;
        font-weight: 500;
    }
}
</style>
