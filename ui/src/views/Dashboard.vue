<template>
  <div v-if="dataIsLoaded">
    
    <vs-row vs-justify="space-around">
      <vs-col
        type="flex"
        vs-justify="center"
        vs-align="center"
        vs-lg="5"
        vs-sm="7"
        vs-xs="7"
        vs-w="5"
      >
        <vs-card fixed-height>
          <div slot="header">
            <h2>Server</h2>
          </div>
          <div>
            <vs-row vs-justify="space-around" vs-align="center" vs-w="12">
              <div>
                <vs-row
                  vs-type="flex"
                  vs-justify="space-between"
                  vs-align="center"
                >
                  <vs-col
                    vs-type="flex"
                    vs-justify="center"
                    vs-align="center"
                    vs-w="5"
                  >
                    {{ getServerStatus.interface_name }}
                  </vs-col>
                  <vs-col
                    vs-type="flex"
                    vs-justify="center"
                    vs-align="center"
                    vs-w="5"
                  >
                    <vs-icon
                      size="15px"
                      :color="wgStatusColor"
                      icon="circle"
                    ></vs-icon>
                    <!-- <vs-radio
                      v-model="wgStatus"
                      disabled="true"
                      color="success"
                      vs-name="server-status"
                      vs-value="up"
                    ></vs-radio> -->
                  </vs-col>
                </vs-row>
              </div>
              <div>
                {{ formatBytes(getServerStatus.send) }}
                <vs-icon size="15px" color="red" icon="call_made"></vs-icon>
              </div>
              <div>
                {{ formatBytes(getServerStatus.receive) }}
                <vs-icon
                  size="15px"
                  color="rgb(70, 150, 0)"
                  icon="call_received"
                ></vs-icon>
              </div>
            </vs-row>
            <!-- <vs-divider color="primary"></vs-divider> -->
            <vs-row vs-w="12">
              <vs-list>
                <vs-list-item title="Apply change">
                  <vs-button radius @click="reload" size="small" color="success" type="flat" icon="sync"></vs-button>
                  <!-- <vs-button @click="reload" size="small" color="danger">apply change</vs-button> -->
                </vs-list-item>
                <vs-list-item title="Start/Stop">
                <!-- <vs-button
                  @click="serverOperation(wgStartStopBText)"
                  size="small"
                  color="primary"
                  type="border"
                  :icon="wgStartStopBIcon"
                  >{{ wgStartStopBText }}</vs-button
                > -->
                <vs-button radius @click="serverOperation(wgStartStopBText)" size="small" :color="wgStartStopBColor" type="flat" :icon="wgStartStopBIcon"></vs-button>
                  <!-- <vs-button size="small" color="danger">One action</vs-button> -->
                </vs-list-item>
              </vs-list>
              <!-- <vs-col
                vs-type="flex"
                vs-justify="flex-start"
                vs-align="center"
                vs-w="3"
              >
                <vs-button
                  @click="reload"
                  size="small"
                  color="rgb(127, 170, 240)"
                  type="filled"
                  icon="sync"
                  >Apply change</vs-button
                >
              </vs-col>
              <vs-col
                vs-type="flex"
                vs-justify="flex-start"
                vs-align="center"
                vs-w="3"
              >
                <vs-button
                  @click="serverOperation(wgStartStopBText)"
                  size="small"
                  color="primary"
                  type="border"
                  :icon="wgStartStopBIcon"
                  >{{ wgStartStopBText }}</vs-button
                >
              </vs-col> -->


              <!-- <div> -->
              <!-- <vs-button radius color="primary" type="border" icon="search"></vs-button> -->
              <!-- <vs-button size="small" color="primary" type="border" icon="sync">Apply change</vs-button> -->

              <!-- <vs-icon
                  size="small"
                  color="rgb(70, 150, 0)"
                  icon="sync"
                ></vs-icon> -->

              <!-- </div> -->
              <!-- <div> -->
              <!-- <vs-button size="small" color="primary" type="border" icon="sync">Apply change</vs-button> -->
              <!-- <vs-icon
                  size="small"
                  color="rgb(70, 150, 0)"
                  icon="sync"
                ></vs-icon> -->
              <!-- </div> -->
            </vs-row>
          </div>
          <!-- <div slot="footer">
        </div> -->
        </vs-card>
      </vs-col>
      <vs-col
        type="flex"
        vs-justify="center"
        vs-align="center"
        vs-lg="5"
        vs-sm="7"
        vs-xs="7"
        vs-w="5"
      >
        <vs-card fixed-height>
          <div slot="header">
            <h2>Clients</h2>
          </div>
          <vs-row  vs-type="flex" vs-justify="space-between">
            <vs-col vs-w="6" vs-type="flex" vs-justify="center" vs-align="center" >
              <vs-list>
                    <vs-list-header icon="supervisor_account" title="All Client"></vs-list-header>
                    <vs-list-item icon="" :title="String(getServerStatus.all_client)" subtitle=""></vs-list-item>
                </vs-list>
               
              </vs-col>
              <vs-col vs-w="6" vs-type="flex" vs-justify="center" vs-align="center" >
                <vs-list>
                    <vs-list-header icon="supervisor_account" title="Enabled Client" color="success"></vs-list-header>
                    <vs-list-item icon="" :title="String(getServerStatus.enable_client)" subtitle=""></vs-list-item>
                </vs-list>

              </vs-col>
          </vs-row>
          <!-- <vs-row   vs-w="12">
            <vs-col vs-align="center" vs-w="6">40</vs-col>
            <vs-col vs-align="center" vs-w="6">5</vs-col>
          </vs-row> -->
          <!-- <div class="mainContainer">
            <div class="clientItem">{{ getServerStatus.all_client }}</div>
            <div class="clientItem">{{ getServerStatus.enable_client }}</div>  
         </div>  -->
        </vs-card>
      </vs-col>
    </vs-row>
    <vs-divider color="success"></vs-divider>
    <!-- vs-justify="flex-start" vs-align="center" vs-lg ="5" vs-sm="7" vs-xs="3" vs-w="12" -->
    <vs-row vs-w="12">
      <vs-col
        vs-offset="1"
        vs-justify="center"
        vs-align="center"
        vs-lg="10"
        vs-sm="10"
        vs-xs="10"
      >
        <ApptTafficTable
          :fields="trafficUsageTablehedear"
          title="Traffic Usage"
          :tableData="getclientStat"
        >
        </ApptTafficTable>
      </vs-col>
    </vs-row>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import ApptTafficTable from "../components/AppTrafficTable.vue";

export default {
  name: "Dashboard",
  components: {
    ApptTafficTable,
  },
  data() {
    return {
      dataIsLoaded: true,
      serverStatus: "",
      trafficUsageTablehedear: ["Name", "IP", "Send", "Recieve", "Last Seen"],
      trafficUssage: [
        {
          id: "",
          ip: "",
          name: "ali",
          private_key: "",
          public_key: "z1OsMX1m10F35/dfgRSDQ7nVoRwe7V9pEiy4h2K9OV0=",
          last_handshake_time: "2022-02-22T07:58:14.440686993Z",
          allocated_ip: "192.168.90.2",
          "send-byte": 47722,
          "recieve-byte": 102048,
        },
        {
          id: "",
          ip: "",
          name: "sh-phone",
          private_key: "",
          public_key: "cFoLN2Nl2MtKrcVLfJty5mA3ksg3pNjFextOlH+eWiA=",
          last_handshake_time: "0001-01-01T00:00:00Z",
          allocated_ip: "192.168.90.3",
          "send-byte": 0,
          "recieve-byte": 0,
        },
      ],
    };
  },
  computed: {
    ...mapGetters(["getclientStat", "getServerStatus"]),
    wgStatusColor() {
      if (this.getServerStatus.is_enable == true) {
        return "#12EE1690";
      } else {
        return "#f54242";
      }
    },
    wgStartStopBText() {
      if (this.getServerStatus.is_enable) {
        return "Stop";
      } else {
        return "Start";
      }
    },
    wgStartStopBIcon() {
      if (this.getServerStatus.is_enable) {
        return "stop";
      } else {
        return "play_arrow";
      }
    },
    wgStartStopBColor(){
      if (this.getServerStatus.is_enable) {
        return "danger";
      } else {
        return "success";
      }
    },
  },
  methods: {
    serverOperation(action) {
      action = action.toLowerCase();
      this.$vs.loading();
      this.$store.dispatch("serverOperationReq", action).then(() => {
        this.$store.dispatch("serverStatusReq").then(() => {
          this.$store.dispatch("getClientStatic").then(() => {
            this.dataIsLoaded = true;
            this.$vs.loading.close();
          });
        });
        // this.$vs.loading.close();
      });
    },
    reload() {
      this.$vs.loading();
      this.$store.dispatch("serverReloadReq").then(() => {
        this.$store.dispatch("serverStatusReq").then(() => {
          this.$store.dispatch("getClientStatic").then(() => {
            this.dataIsLoaded = true;
            this.$vs.loading.close();
          });
        });
        // this.$vs.loading.close();
      });
    },
    formatBytes(bytes, decimals) {
      if (bytes == 0) return "0 Bytes";
      var k = 1024,
        dm = decimals || 2,
        sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"],
        i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
    },
  },
  watch: {},
  created() {
    this.$vs.loading();
    this.$store.dispatch("serverStatusReq").then(() => {
      this.$store.dispatch("getClientStatic").then(() => {
        this.dataIsLoaded = true;
        this.$vs.loading.close();
      });
    });
  },
};
</script>

<style scoped>
.clientInfo{
  padding: 50px;
}
</style>