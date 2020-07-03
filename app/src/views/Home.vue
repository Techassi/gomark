<template>
    <div class="home">
        <Recent></Recent>
        <MostUsed></MostUsed>
    </div>
</template>

<script>
import Recent from '@/components/partials/recent.vue'
import MostUsed from '@/components/partials/most-used.vue'
export default {
    name: 'Home',
    components: {
        Recent,
        MostUsed,
    },
}
</script>
