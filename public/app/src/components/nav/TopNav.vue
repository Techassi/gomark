<template>
    <nav class="nav__top">
        <router-link to="/" class="nav__top--logo">
            <i class="gg-bookmark"></i>
        </router-link>
        <div class="nav__top--links">
            <router-link to="/">Home</router-link>
            <router-link to="/all">All</router-link>
        </div>
    </nav>
</template>

<script>
export default {
    name: 'TopNav',
}
</script>

<style lang="scss" scoped>
.nav__top {
    position: fixed;
    display: grid;
    grid-template-columns: 56px 2fr 1fr;
    column-gap: 28px;
    width: 100%;
    height: 56px;
    background-color: $primary-white;
    border-bottom: 2px solid $secondary-grey;
    z-index: 20;

    &--logo {
        display: flex;
        align-items: center;
        justify-content: center;

        i {
            color: $primary-accent;
        }
    }

    &--links {
        display: flex;
        align-items: center;

        a {
            font-size: 1.6rem;
            font-weight: 500;
            padding: 0 28px 0 0;
            color: $primary-black;

            &.router-link-exact-active {
                color: $primary-accent;
            }
        }
    }
}
</style>
