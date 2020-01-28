/////////////////////////////////////////////
////////////////// IMPORTS //////////////////
/////////////////////////////////////////////

import barba from "@barba/core";
import Hammer from "hammerjs/hammer.js";

import Login from "./namespaces/login";

/////////////////////////////////////////////
////////////////// GENERAL //////////////////
/////////////////////////////////////////////

console.info(
    "%c🗲 gomark @ https://github.com/Techassi/gomark",
    "color: #FFF; background: #011C41; padding: 10px 10px;"
);

/////////////////////////////////////////////
////////////// INIT NAMESPACES //////////////
/////////////////////////////////////////////

let login = new Login();

/////////////////////////////////////////////
/////////////////// BARBA ///////////////////
/////////////////////////////////////////////

barba.init({
    debug: true,
    transitions: [
        {
            name: "default-transition",
            leave(data) {
                return new Promise(resolve => {
                    anime({
                        targets: data.current.container,
                        opacity: [1, 0],
                        duration: 500,
                        easing: "linear",
                        complete: () => {
                            resolve();
                        }
                    });
                });
            },
            enter(data) {
                return new Promise(resolve => {
                    resolve();
                    anime({
                        targets: data.next.container,
                        opacity: [0, 1],
                        duration: 500,
                        easing: "linear"
                    });
                });
            }
        }
    ],
    views: [
        {
            namespace: "login",
            beforeEnter(data) {
                login.Init(data);
            },
            afterEnter() {},
            beforeLeave() {
                login.Kill();
            }
        }
    ]
});
