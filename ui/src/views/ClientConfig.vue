<template>
  <div v-if="dataIsLoaded">
    <!-- popup to add client----------------------------------------------------------------->
    <app-client-popup
      :title="clientConfigPopUpTitle"
      :availabaleIPs="getClientAvailableIPs"
      :popType="clientConfigPopUpType"
      :popupEnable="showClientConfigPopUp"
      :popData="client"
      @popup-data="processPopUpData"
      @remove-client="removeClient"
      @close-popup="closePopup"
    ></app-client-popup>
    <!---------------------------------------------------------------------------------------->
    <app-qrcode
      :popData="codeData"
      :popupEnable="openCode"
      @close-popup="closeqrcode"
    >
    </app-qrcode>
    <!---------------------------------------------------------------------------------------->
    <vs-row vs-w="12">
      <vs-col>
        <vs-button
          class="holamundo"
          color="primary"
          type="border"
          @click="openAddPopup()"
          >Add Client</vs-button
        >
      </vs-col>
    </vs-row>
    <vs-row vs-w="12">
      <vs-col>
        <vs-divider color="primary"></vs-divider>
      </vs-col>
    </vs-row>
    <vs-row vs-w="12">
      <vs-col
        vs-offset="1"
        vs-justify="center"
        vs-align="center"
        vs-lg="10"
        vs-sm="10"
        vs-xs="10"
      >
        <app-client-table
          v-if="dataIsLoaded"
          :title="tableName"
          @download-qrcode="downloadQrcode"
          @download-config="downloadConfig"
          :fields="tableFiedls"
          :tableData="getAllClient"
          @selected-row="clientFromtableOpenUpdatePopup"
          @change-state="changeClientState"
        ></app-client-table>
      </vs-col>
    </vs-row>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import AppClientTable from "../components/AppClientTable.vue";
import AppQrcode from "../components/AppQrcode.vue";
import AppClientPopup from "../components/AppClientPopup.vue";
export default {
  components: {
    AppClientTable,
    AppClientPopup,
    AppQrcode,
  },
  data() {
    return {
      dataIsLoaded: false,
      tableName: "Clients",
      tableFiedls: [
        "Name",
        "IP",
        "AllowedIP",
        "Download QrCode",
        "Download Config",
        "State",
      ],
      client: {},
      codeData: {},
      openCode: false,
      showClientConfigPopUp: false,
      clientConfigPopUpTitle: "",
      clientConfigPopUpType: "",
      refreshing:false,
    };
  },
  computed: {
    ...mapGetters(["getAllClient", "getClientAvailableIPs"]),
  },
  methods: {
    closeqrcode() {
      this.openCode = false;
    },
    closePopup() {
      this.showClientConfigPopUp = false;
    },
    remove(sourceArr, item) {
      sourceArr.splice(sourceArr.indexOf(item), 1);
    },
    removeClient(client) {
      this.showClientConfigPopUp = false;
      this.$vs.loading();
      this.$store.dispatch("removeClientReq", client).then(() => {
        this.clientConfigPopUpTitle = "";
        this.clientConfigPopUpType = "";
        this.$vs.loading.close();
        this.refreshing=true
      });
    },
    clientFromtableOpenUpdatePopup(row) {
      this.clientConfigPopUpTitle = "Update Client";
      this.clientConfigPopUpType = "update";
      this.client = Object.assign({}, row);
      this.showClientConfigPopUp = true;
      // console.log("IN PARENT", row);
    },
    changeClientState(client){
      this.$vs.loading();
      this.$store.dispatch("changeStatReq", client).then(() => {
        this.$vs.loading.close();
        this.refreshing=true
      });
      // console.log(client)
    },
    downloadQrcode(client) {
      this.codeData = {
        clientName: client.name,
        imageData: client.qr_code,
      };
      this.openCode = true;
    },
    downloadConfig(client) {
      this.$store.dispatch("getclientConfig", client.id).then(() => {});
    },
    openAddPopup() {
      this.clientConfigPopUpTitle = "Add Client";
      this.clientConfigPopUpType = "add";
      this.client = {
        name: "",
        allocated_ip: "",
        dns_address: "",
        allowed_ips: [],
      };
      this.showClientConfigPopUp = true;
    },
    processPopUpData(client) {
      this.showClientConfigPopUp = false;
      this.$vs.loading();
      if (this.clientConfigPopUpType == "add") {
        this.$store.dispatch("addClientReq", client).then(() => {
          this.clientConfigPopUpTitle = "";
          this.clientConfigPopUpType = "";
          this.$vs.loading.close();
          this.showClientConfigPopUp = false;
          this.refreshing=true
          this.$vs.notify({
            title: "config change",
            position:"top-center",
            text: "apply change from dashboard in order to change take effect",
            color: "success",
            time: 4000,
          });
        });
      } else {
        this.$store.dispatch("updateClientReq", client).then(() => {
          this.clientConfigPopUpTitle = "";
          this.clientConfigPopUpType = "";
          this.$vs.loading.close();
          this.showClientConfigPopUp = false;
          this.refreshing=true
          this.$vs.notify({
            title: "config change",
            position:"top-center",
            text: "apply change from dashboard in order to change take effect",
            color: "success",
            time: 4000,
          });
        });
      }
    },
  },
  watch: {
    refreshing:function(){
      if (this.refreshing){
        this.$vs.loading();
        this.$store.dispatch("getavailabaleIPsReq").then(() => {
          this.$store.dispatch("allClientReq").then(() => {
            this.dataIsLoaded = true;
            this.$vs.loading.close();
            this.refreshing=false
          });
        });
      }
    },
    allowedIP: function (val) {
      // console.log("change happend",val);
      for (const item of val) {
        // console.log("ITEM",item);
        if (item.includes(" ")) {
          const itemSplit = item.split(" ");
          this.allowedIP.splice(this.allowedIP.indexOf(item), 1);
          for (const sItem of itemSplit) {
            this.allowedIP.push(sItem);
          }
        }
      }
    },
    addClientPopup: function (val) {
      console.log("Popup state", val);
    },
  },
  created() {
    this.$vs.loading();
    this.$store.dispatch("getavailabaleIPsReq").then(() => {
      this.$store.dispatch("allClientReq").then(() => {
        this.dataIsLoaded = true;
        this.$vs.loading.close();
      });
    });
  },
  // updated(){
  //   console.log("update call")
  // },
};
</script>

<style scoped>
/* .rowFlex {
  display: flex;
  flex-direction: column;
}
.addForm {
  display: flex;
  flex-direction: column;
  gap: 10px;
} */
</style>

<!-- <vs-popup title="Add Client" :active.sync="addClientPopup">
      <div class="addForm">
        <vs-row>
          <vs-col>
            <vs-input label="Name" placeholder="" v-model="client.Name" />
          </vs-col>
        </vs-row>
        <vs-row>
          <vs-col>
            <vs-input label="Client IP" placeholder="" v-model="client.IP" />
          </vs-col>
        </vs-row>
        <vs-row>
          <vs-col>
            <vs-input label="DNS Address" placeholder="" v-model="client.DNS" />
          </vs-col>
        </vs-row>
        <vs-row class="rowFlex">
          <label for="">Allowed IP</label>
          <vs-col vs-lg="4" vs-sm="4" vs-xs="12">
            <vs-chips
              remove-icon=""
              color="rgb(145, 32, 159)"
              placeholder="ips"
              v-model="client.AllowedIP"
            >
              <vs-chip
                :key="ips + index"
                @click="remove(client.AllowedIP, ips)"
                v-for="(ips, index) in client.AllowedIP"
                closable
              >
                {{ ips }}
              </vs-chip>
            </vs-chips>
          </vs-col>
        </vs-row>
        <vs-row>
          <vs-col>
            <vs-button color="success" type="border">Save</vs-button>
          </vs-col>
        </vs-row>
      </div>
    </vs-popup> -->