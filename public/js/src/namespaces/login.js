export default class Login {
    constructor() {
        this.Init = (data) => {
            let next = data.next.container;

            this.store = {
                dom: {
                    username: next.querySelector('input.username'),
                    password: next.querySelector('input.password'),
                    login: next.querySelector('button.login'),
                }
            }

            this.initHandlers();
        }

        this.Kill = () => {
            this.killHandlers();
        }

        this.initHandlers = () => {
            this.store.dom.login.addEventListener('click', this.sendLoginRequest);
        }

        this.killHandlers = () => {
            this.store.dom.login.removeEventListener('click', this.sendLoginRequest);
        }

        this.sendLoginRequest = () => {
            
        }
    }
}
