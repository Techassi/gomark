<template>
    <div class="main__most-used">
        <h2>Most used bookmarks</h2>
        <div class="most-used__list"></div>
    </div>
</template>

<style lang="scss" scoped>
.main__most-used {
    h2 {
        margin: 0;
        font-size: 2rem;
        font-weight: 500;
    }

    .most-used__list {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
        gap: 24px;
        margin: 24px 0 0 0;
    }
}
</style>
