const client = {
    state: {
        clients: [],
        availabaleIPs: [],
        clientStat: [],
        loginUser: {
            user: sessionStorage.getItem('user') || "",
            token: sessionStorage.getItem('token') || "",
            isLogiedIn: (sessionStorage.getItem('isLogiedIn') == "true") || false,
        },
    },
    mutations: {
        setClientStat(state, payload) {
            state.clientStat = payload
        },
        setClient(state, payload) {
            state.clients = payload
        },
        setAvailableIPs(state, payload) {
            state.availabaleIPs = payload
        },
        setLoginInfo(state, payload) {
            state.loginUser = payload
        },
    },
    actions: {
        async allClientReq(context) {
            const http_client = this._vm.$http;
            const url = "private/client"
            await http_client.get(url)
                .then(resp => context.commit('setClient', resp.data.Data));
        },
        //------------------------------------------------------------------------
        async addClientReq(context, client) {
            const http_client = this._vm.$http;
            const url = "private/client"
            await http_client.post(url, client)
                .then(() => { });
        },
        //------------------------------------------------------------------------
        async updateClientReq(context, client) {
            const http_client = this._vm.$http;
            const url = "private/client"
            await http_client.put(url, client)
                .then(() => { });
        },
        //------------------------------------------------------------------------
        async changeStatReq(context, client) {
            const http_client = this._vm.$http;
            const url = "private/client/availability"
            await http_client.put(url, client)
                .then(() => { });
        },

        //------------------------------------------------------------------------
        async removeClientReq(context, client) {
            // console.log("REMOVE : ", client)
            const http_client = this._vm.$http;
            const url = "private/client"
            await http_client.delete(url, { data: client })
                .then(() => { });
        },
        //------------------------------------------------------------------------
        async getavailabaleIPsReq(context, client) {
            const http_client = this._vm.$http;
            const url = "private/client/usedIP"
            await http_client.get(url, client)
                .then(resp => context.commit('setAvailableIPs', resp.data.Data));
        },
        //------------------------------------------------------------------------
        async getclientConfig(context, id) {
            const http_client = this._vm.$http;
            const url = "private/client/config/"
            await http_client.get(url + id, {
                responseType: 'arraybuffer',
            })
                .then(resp => {
                    let fileName = resp.headers["x-filename"]
                    var fileURL = window.URL.createObjectURL(new Blob([resp.data]));
                    var fileLink = document.createElement('a');

                    fileLink.href = fileURL;
                    fileLink.setAttribute('download', fileName);
                    document.body.appendChild(fileLink);
                    fileLink.click();
                })
                .catch((error) => {
                    context.commit('setDownloadProgress', 0)
                    console.log(error);
                });
        },
        //------------------------------------------------------------------------
        async getClientStatic(context) {
            const http_client = this._vm.$http;
            const url = "private/client/stat"
            await http_client.get(url)
                .then(resp => context.commit('setClientStat', resp.data.Data));
        },
        //------------------------------------------------------------------------
        sendReqUserLogOUT(context) {
            sessionStorage.removeItem('token')
            sessionStorage.setItem('isLogiedIn', new Boolean(false).toString()),
                sessionStorage.removeItem('user')
            let loginInfo = {
                user: "",
                token: "",
                isLogiedIn: false,
            }
            context.commit('setLoginInfo', loginInfo)
        },
        //------------------------------------------------------------------------
        userLoginRq(context, client) {
            return new Promise((resolve, reject) => {
                const http_client = this._vm.$http;
                const url = "/login"
                http_client.post(url, client)
                    .then(resp => {
                        let loginInfo = {
                            user: resp.data.Data.user,
                            token: resp.data.Data.token,
                            isLogiedIn: true
                        }
                        context.commit('setLoginInfo', loginInfo)
                        sessionStorage.setItem('token', resp.data.Data.token)
                        sessionStorage.setItem('user', resp.data.Data.user)
                        sessionStorage.setItem('isLogiedIn', new Boolean(true).toString()),
                            http_client.defaults.headers.common['token'] = context.state.loginUser.token
                        resolve(resp)

                    })
                    .catch(error => {
                        sessionStorage.removeItem('token')
                        sessionStorage.removeItem('user')
                        sessionStorage.removeItem('isLogiedIn')
                        let loginInfo = {
                            user: "",
                            token: "",
                            isLogiedIn: false,
                        }
                        context.commit('setLoginInfo', loginInfo)
                        delete http_client.defaults.headers.common['token']
                        reject(error)
                    })
            })
        },
    },
    getters: {
        getclientStat(state) {
            return state.clientStat
        },
        getAllClient(state) {
            return state.clients
        },
        getClientAvailableIPs(state) {
            return state.availabaleIPs
        },
        getUserLogin(state) {
            return state.loginUser
        },
    },
}

export default client