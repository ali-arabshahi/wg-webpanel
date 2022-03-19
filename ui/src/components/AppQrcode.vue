<template>
  <!-- :button-close-hidden="true"  -->
  <vs-popup title="QrCode" :active.sync="popupToggle">
    <vs-row vs-justify="center" vs-w="12">
      <vs-col type="flex" vs-justify="center" vs-align="center" vs-w="6">
        <vs-card>
          <div slot="header">
            <h3>
              {{ popData.clientName }}
            </h3>
          </div>
          <div>
            <vs-row>
              <vs-col>
                <img v-bind:src="popData.imageData" />
              </vs-col>
            </vs-row>
          </div>
          <div slot="footer">
            <vs-row vs-justify="flex-end">
              <vs-button
                @click="downloadCode"
                type="gradient"
                color="success"
                icon="file_download"
              ></vs-button>
            </vs-row>
          </div>
        </vs-card>
      </vs-col>
    </vs-row>

    <!-- <vs-row vs-type="flex" vs-justify="center" vs-w="12">
      <vs-col vs-type="flex" vs-justify="center" vs-align="center">
        <img v-bind:src="popData.imageData" />
      </vs-col>
    </vs-row> -->
  </vs-popup>
</template>

<script>
export default {
  name: "AppQrcode",
  emits: ["close-popup"],
  props: {
    popupEnable: Boolean,
    popData: Object,
  },
  data() {
    return {};
  },
  computed: {
    popupToggle: {
      get() {
        return this.popupEnable;
      },
      set() {
        this.$emit("close-popup");
      },
    },
  },
  methods: {
    downloadCode() {
      var a = document.createElement("a"); //Create <a>
      a.href = this.popData.imageData;
      // a.href = "data:image/png;base64," + ImageBase64; //Image Base64 Goes here
      a.download = this.popData.clientName + ".png"; //File name Here
      a.click(); //Downloaded file
      this.$emit("close-popup");
    },
    closePopup() {
      this.$emit("close-popup");
    },
  },
  watch: {},
};
</script>

<style scoped>
</style>