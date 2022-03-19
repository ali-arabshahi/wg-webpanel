<template>
  <div id="parentx">
    <vs-sidebar
      parent="#parentx"
      default-index="1"
      color="primary"
      class="sidebarx"
      spacer
      v-model="isSidebarActive"
    >
      <div class="header-sidebar" slot="header">
        <div class="Cuser">
          <p>
            {{ user }}
          </p>
        </div>
        <vs-button icon="reply" color="danger" type="flat" @click="bLogout"
          >log out</vs-button
        >
      </div>

      <vs-sidebar-item index="1" to="home" icon="dashboard"
        >Dashboard</vs-sidebar-item
      >
      <vs-sidebar-item index="2" to="server" icon="dns">
        Server Config
      </vs-sidebar-item>
      <vs-sidebar-item index="3" to="client" icon="manage_accounts">
        Clients Config
      </vs-sidebar-item>
    </vs-sidebar>
  </div>
</template>

<script>
export default {
  name: "AppSideBar",
  emits: ["log-out"],
  props: {
    parent: {
      type: String,
    },
    user: String,
  },
  data: () => ({
    active: false,
  }),
  methods: {
    bLogout() {
      this.$emit("log-out");
    },
  },
  computed: {
    isSidebarActive: {
      get() {
        return this.$store.getters.getSiderState;
      },
      set() {
        this.$store.commit("toggleSidebar");
      },
    },
  },
};
</script>

<style scoped>
.header-sidebar {
  display: flex;
  align-content: center;
  justify-content: space-between;
}
.Cuser {
  padding-top: 15px;
  padding-left: 2px;
}
</style>
