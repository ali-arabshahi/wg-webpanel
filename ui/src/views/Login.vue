<template>
  <div class="container">
    <!-- <div class="loginBox"> -->
    <vs-row vs-justify="center" vs-w="12">
      <vs-col type="flex" vs-justify="center" vs-align="center" vs-w="4">
        <vs-card>
          <div slot="header">
            <h3>Login</h3>
          </div>
          <div class="rowFlex">
            <vs-row vs-justify="center" vs-w="12">
              <vs-input label="User Name" placeholder="" v-model="userName" />
            </vs-row>
            <vs-row vs-justify="center" vs-w="12">
              <vs-input
                type="password"
                label="Password"
                placeholder=""
                v-model="password"
              />
            </vs-row>
          </div>
          <div slot="footer">
            <vs-row vs-justify="center">
              <!-- <vs-button color="success" type="filled">Login</vs-button> -->
              <vs-button
                icon-after
                :color="btcolor"
                type="gradient"
                :icon="loginButtonIcon"
                @click="loginReq"
                >Login</vs-button
              >
            </vs-row>
          </div>
        </vs-card>
      </vs-col>
    </vs-row>
    <!-- </div> -->
  </div>
</template>

<script>
export default {
  name: "Login",
  components: {},
  data() {
    return {
      userName: "",
      password: "",
      isLogin: false,
      btcolor: "success",
    };
  },
  computed: {
    loginButtonIcon() {
      if (!this.isLogin) {
        return "lock";
      } else {
        return "lock_open";
      }
    },
  },
  methods: {
    loginReq() {
      let user = {
        UserName: this.userName,
        Password: this.password,
      };
      this.$vs.loading();
      this.$store
        .dispatch("userLoginRq", user)
        .then(() => {
          console.log("login sucessful");
          this.isLogin = true;
          this.$vs.loading.close();
          setTimeout(() => {
            this.$router.push("/home");
          }, 500);
        })
        .catch(() => {
          this.btcolor = "danger";
          setTimeout(() => {
            this.btcolor = "success";
          }, 1000);
          console.log(" login fail from login page");
          this.$vs.loading.close();
        });
    },
  },
  watch: {},
  created() {},
};
</script>

<style scoped>
.CbuttonNon {
}
.CbuttonFail {
  animation-name: shake;
  animation-duration: 5s;
  animation-iteration-count: infinite;
  animation-timing-function: ease-in;
}
.rowFlex {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}
.loginBox {
  /* width: 100px; */
  /* height: 100px; */
  background-color: aquamarine;
}
.main-container-fluid {
  height: 100vh;
}
</style>


