const server = {
    state: {
        serverConfig: {},
        intefaces: [],
        serverStatus: {},
    },
    mutations: {
        setServerStatusData(state, payload) {
            state.serverStatus = payload
        },
        setServerData(state, payload) {
            state.serverConfig = payload
        },
        setInterfacesData(state, payload) {
            state.intefaces = payload
        },
    },
    actions: {
        async serverConfigReq(context) {
            const http_client = this._vm.$http;
            const url = "private/server"
            await http_client.get(url)
                .then(resp => context.commit('setServerData', resp.data.Data));
        },
        //------------------------------------------------------------------------
        async saveServerConfigReq(context, serverConf) {
            context.state.serverConfig = { ...serverConf }
            const http_client = this._vm.$http;
            const url = "private/server"
            await http_client.post(url, context.state.serverConfig)
                .then(() => { });
        },
        //------------------------------------------------------------------------
        async serverInterfacesReq(context) {
            const http_client = this._vm.$http;
            const url = "private/interfaces"
            await http_client.get(url)
                .then(resp => context.commit('setInterfacesData', resp.data.Data));
        },
        //------------------------------------------------------------------------
        async serverStatusReq(context) {
            const http_client = this._vm.$http;
            const url = "private/server/status"
            await http_client.get(url)
                .then(resp => context.commit('setServerStatusData', resp.data.Data));
        },
        //------------------------------------------------------------------------
        async serverReloadReq() {
            const http_client = this._vm.$http;
            const url = "private/server/reload"
            await http_client.post(url)
                .then(() => { });
        },
        //------------------------------------------------------------------------
        async serverOperationReq(context, action) {
            const http_client = this._vm.$http;
            const url = "private/server/" + action
            await http_client.post(url)
                .then(() => { });
        },
        //------------------------------------------------------------------------
    },
    getters: {
        getServer(state) {
            return state.serverConfig
        },
        getInterfaces(state) {
            return state.intefaces
        },
        getServerStatus(state) {
            return state.serverStatus
        }
    },
}

export default server