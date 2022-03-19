const sideBar = {
    state: {
        sidebarMinimize: false
    },
    mutations: {
        toggleSidebar(state) {
            state.sidebarMinimize = !state.sidebarMinimize
        },
        closeSidebar(state) {
            state.sidebarMinimize = false
        }
    },
    actions: {
        actionSidebare(context) {
            context.commit('toggleSidebar')
        }
    },
    getters: {
        getSiderState(state) {
            return state.sidebarMinimize
        }
    },
}

export default sideBar