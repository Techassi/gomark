<template>
    <nav class="nav__side">
        <div class="nav__side--icon add">
            <i class="gg-add-r"></i>
        </div>
        <div class="nav__side--icon notes">
            <i class="gg-notes"></i>
        </div>
        <div class="nav__side--icon album">
            <i class="gg-album"></i>
        </div>
        <div class="nav__side--icon share">
            <i class="gg-share"></i>
        </div>
    </nav>
</template>

<script>
export default {
    name: 'SideNav',
}
</script>

<style lang="scss" scoped>
.nav__side {
    position: fixed;
    display: grid;
    grid-template-rows: 56px 56px 56px 56px;
    width: 56px;
    height: 100%;
    padding: 84px 0 0 0;
    background-color: #fff;
    border-right: 2px solid $secondary-grey;
    z-index: 10;

    &--icon {
        display: flex;
        align-items: center;
        justify-content: center;
        transition: 0.2s;
        cursor: pointer;

        --ggs: 0.8;

        &:hover {
            color: $primary-accent;
        }

        &.share {
            i {
                margin: 0 0 0 -6px;
            }
        }
    }
}
</style>
