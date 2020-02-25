const PARAMETERS_INVALID = "parameters_invalid";
const CREDENTIALS_INVALID = "credentials_invalid";
const TWOFA_REQUIRED = "2fa_required";
export default class Login {
    constructor() {
        this.Init = data => {
            let next = data.next.container;

            this.store = {
                dom: {
                    username: next.querySelector("input.username"),
                    password: next.querySelector("input.password"),
                    login: next.querySelector("button.login")
                }
            };

            this.initHandlers();
        };

        this.Kill = () => {
            this.killHandlers();
        };

        this.initHandlers = () => {
            this.store.dom.login.addEventListener(
                "click",
                this.sendLoginRequest
            );
        };

        this.killHandlers = () => {
            this.store.dom.login.removeEventListener(
                "click",
                this.sendLoginRequest
            );
        };

        this.sendLoginRequest = () => {
            var credentials = new FormData();
            credentials.append("username", this.store.dom.username.value);
            credentials.append("password", this.store.dom.password.value);

            fetch("/auth/login", {
                method: "POST",
                body: credentials
            })
                .then(res => {
                    if (res.ok) {
                        return res.json();
                    } else {
                        throw new Error(res.statusText);
                    }
                })
                .then(json => {
                    if (json.status != 200) {
                        if (json.message == PARAMETERS_INVALID) {
                            console.log("Your credentials cannot be empty");
                        } else if (json.message == CREDENTIALS_INVALID) {
                            console.log("Your credentials are invalid");
                        } else {
                            console.log("There was an error loggin you in");
                        }
                        return;
                    }

                    if (json.message == TWOFA_REQUIRED) {
                        this.show2FACodeInput();
                        return;
                    }

                    window.location.replace("/");
                })
                .catch(err => {
                    throw new Error(err);
                });
        };

        this.show2FACodeInput = () => {
            console.log("2FA");
        };
    }
}
