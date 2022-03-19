<template>
  <div class="main-wrapper">
    <AppNavBar></AppNavBar>
    <AppSideBar
      parent=".main-wrapper"
      :user="getUserLogin.user"
      @log-out="userLogout"
    ></AppSideBar>
    <div class="main-container-fluid">
      <vs-row vs-w="12">
        <vs-col
          vs-offset="1"
          vs-type="flex"
          vs-justify="center"
          vs-align="center"
          vs-w="10"
        >
          <vs-card>
            <!-- <transition name="fade"> -->
            <router-view></router-view>
            <!-- </transition> -->
          </vs-card>
        </vs-col>
      </vs-row>
    </div>
  </div>
</template>


<script>
import { mapGetters } from "vuex";
import AppNavBar from "../components/AppNavBar.vue";
import AppSideBar from "../components/AppSideBar.vue";

export default {
  name: "Home",
  components: {
    AppNavBar,
    AppSideBar,
  },
  computed: {
    ...mapGetters(["getUserLogin"]),
  },
  methods: {
    userLogout() {
      this.$store.dispatch("sendReqUserLogOUT").then(() => {
        console.log("logout done");
        this.$router.push({ name: "LoginPage" }).catch(() => {});
      });
    },
  },
};
</script>
<style scoped>
.main-container-fluid {
  padding: 20px 5px 5px 10px;
}
</style>

