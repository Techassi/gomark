export default class Home {
    constructor() {
        this.Init = data => {
            let next = data.next.container;

            this.store = {
                dom: {
                    slider: document.querySelector(".slider"),
                    navbutton: document.querySelector(".menu"),
                    nav: document.querySelector("nav.side")
                }
            };

            this.initHandlers();
        };

        this.Kill = () => {
            this.killHandlers();
        };

        this.initHandlers = () => {
            this.store.dom.slider.addEventListener("click", this.switchView);
            this.store.dom.navbutton.addEventListener("click", this.toggleNav);
        };

        this.killHandlers = () => {
            this.store.dom.slider.removeEventListener("click", this.switchView);
            this.store.dom.navbutton.removeEventListener(
                "click",
                this.toggleNav
            );
        };

        this.toggleNav = () => {
            this.store.dom.nav.classList.toggle("open");
        };

        this.switchView = () => {
            console.log("Switch View");
        };
    }
}
