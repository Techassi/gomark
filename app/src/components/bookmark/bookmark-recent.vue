<template>
    <div class="recent__list--item">
        <div class="image__wrapper">
            <img :src="imageURL(entity)" v-bind:alt="entity.name" />
        </div>
        <h3>{{ entity.name }}</h3>
        <a
            v-bind:href="entity.bookmark.url"
            target="_blank"
            rel="noopener noreferrer"
        >{{ entity.bookmark.url }}</a>
        <span>{{ this.trim(entity.bookmark.description, 100) }}</span>
    </div>
</template>

<script>
export default {
    name: 'BookmarkRecent',
    methods: {
        trim(s, l) {
            if (s.length == 0) return ''
            return s.length > l ? s.substring(0, l - 3) + '...' : s
        },
        imageURL(entity) {
            return entity.image_url == ''
                ? `image/fallback/image.jpg`
                : `image/${entity.hash}/${entity.image_url}`
        },
    },
    props: {
        entity: Object,
    },
}
</script>

<style lang="scss" scoped>
.recent__list--item {
    padding: 16px;
    border: 1px solid #e1e4e8;
    border-radius: 6px;

    .image__wrapper {
        width: 100%;
        height: 0;
        padding-bottom: 50%;
        position: relative;
        border-radius: 6px;
        overflow: hidden;

        img {
            position: absolute;
            width: 100%;
            height: 100%;
            object-fit: cover;
        }
    }

    h3 {
        font-size: 1.4rem;
        line-height: 1.5;
        max-height: 21px;
        margin: 16px 0 0 0;
        overflow: hidden;
    }

    a {
        display: inline-block;
        font-size: 1.2rem;
        line-height: 1.5;
    }

    span {
        display: block;
        font-size: 1.2rem;
        line-height: 1.5;
        margin: 16px 0 0 0;
    }
}
</style>
