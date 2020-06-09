module.exports = {
    css: {
        loaderOptions: {
            sass: {
                prependData: `
                    @import "@/scss/_reset.scss";
                    @import "@/scss/_variables.scss";
                    @import "@/scss/_icons.scss";
                    @import "@/scss/_fonts.scss";
                `,
            },
        },
    },
}
