<template>
    <div class="main__recently">
        <h2>Recent bookmarks</h2>
        <div class="recent__list">
            <BookmarkRecent
                v-for="(entity, index) in entities"
                v-bind:index="index"
                v-bind:entity="entity"
                v-bind:key="entity.id"
            ></BookmarkRecent>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import BookmarkRecent from '@/components/bookmark/bookmark-recent.vue'

export default {
    name: 'MostUsed',
    mounted() {
        axios
            .get(`api/v1/recent`)
            .then(response => {
                this.entities = response.data.entities
            })
            .catch(e => {
                console.error(e)
            })
    },
    data() {
        return {
            entities: [],
        }
    },
    components: {
        BookmarkRecent,
    },
}
</script>

<style lang="scss" scoped>
.main__recently {
    h2 {
        margin: 0;
        font-size: 2rem;
        font-weight: 500;
    }

    .recent__list {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr;
        gap: 16px;
        margin: 24px 0 0 0;
    }
}
</style>
