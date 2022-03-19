<template>
  <div v-if="dataIsLoaded" class="form">
    <!-- {{ getServer }}
    {{ localCopy }}
    <hr />
    {{ getInterfaces }} -->
    <!-- {{ selectInt }} -->

    <vs-row class="inputCheckBoxRow">
      <!-- <vs-col> -->
      <vs-select label="Interfaces" v-model="localCopy.interface">
        <vs-select-item
          :is-selected.sync="item.IsSelected"
          :key="index"
          :value="item.Interface"
          :text="item.Text"
          v-for="(item, index) in getInterfaces"
        />
      </vs-select>
      <!-- </vs-col> -->
      <!-- <vs-tooltip text="Tooltip Default">
          <vs-icon size="15px" icon="help_outline" color="blue"></vs-icon>
        </vs-tooltip> -->
    </vs-row>
    <vs-row class="inputCheckBoxRow">
      <vs-input
        :disabled="!manualIP"
        label="IP Address"
        placeholder="ip"
        v-model="localCopy.address"
      />
      <vs-checkbox class="inputCheckBox" v-model="manualIP"
        >Manual IP</vs-checkbox
      >
    </vs-row>
    <vs-row>
      <vs-input
        label="Server Port"
        type="number"
        placeholder="port"
        v-model="localCopy.listen_port"
      />
    </vs-row>
    <vs-row class="inputCheckBoxRow">
      <vs-input
        label="Tunel IP"
        :success="ipTunnelValidation.sucess"
        :success-text="ipTunnelValidation.sucessText"
        val-icon-success="done"
        :danger="ipTunnelValidation.fail"
        :danger-text="ipTunnelValidation.failText"
        val-icon-danger="close"
        placeholder="ip"
        v-model="localCopy.tunnel_address"
      />
      <vs-tooltip
        text="server tunnel ip with mask.
      shuould be in this format : x.x.x.x/x. 
      this is a private range"
      >
        <vs-icon size="15px" icon="help_outline" color="blue"></vs-icon>
      </vs-tooltip>
    </vs-row>
    <vs-row>
      <vs-col>
        <label for="">Auto Generate public / private key</label>
        <vs-switch class="lableIp" v-model="autoKey" />
      </vs-col>
    </vs-row>
    <vs-row class="inputCheckBoxRow">
      <vs-input
        :disabled="autoKey"
        label="Public Key"
        placeholder="Public Key"
        v-model="localCopy.public_key"
      />
      <vs-icon
        class="copyIcon"
        size="15px"
        icon="content_copy"
        @click="copyKey('Public')"
      ></vs-icon>
    </vs-row>
    <vs-row class="inputCheckBoxRow">
      <vs-input
        :disabled="autoKey"
        label="Private Key"
        placeholder="Private Key"
        v-model="localCopy.private_key"
      />
      <vs-icon
        class="copyIcon"
        size="15px"
        icon="content_copy"
        @click="copyKey('Private')"
      ></vs-icon>
    </vs-row>
    <!-- <vs-row class="rowFlex">
      <label for="">Allowed IP</label>
      <vs-col vs-lg="4" vs-sm="6" vs-xs="12">
        <vs-chips
          remove-icon=""
          color="rgb(145, 32, 159)"
          placeholder="ips"
          v-model="allowedIP"
        >
          <vs-chip
            :key="ips + index"
            @click="remove(allowedIP, ips)"
            v-for="(ips, index) in allowedIP"
            closable
          >
            {{ ips }}
          </vs-chip>
        </vs-chips>
      </vs-col>
    </vs-row> -->
    <vs-divider color="primary"></vs-divider>
    <!-- <vs-row>
      <vs-col>
        <label for="">Enable IP Forwarding</label>
        <vs-switch class="lableIp" v-model="localCopy['ip-forwarding']" />
      </vs-col>
    </vs-row> -->
    <vs-row>
      <vs-col vs-w="3">
        <label for="">Automate firewall Rull Generate</label>
        <vs-switch class="lableIp" v-model="enableNatRule" />
      </vs-col>
      <vs-col vs-w="1">
        <vs-tooltip text="add nat rule for tunnel">
          <vs-icon size="15px" icon="help_outline" color="blue"></vs-icon>
        </vs-tooltip>
      </vs-col>
    </vs-row>
    <vs-row class="rowFlex">
      <label for="">Pre Script</label>
      <vs-col
        vs-lg="4"
        vs-sm="6"
        vs-xs="12"
        :class="{ disable: enableNatRule }"
      >
        <vs-chips
          remove-icon=""
          color="rgb(145, 32, 159)"
          placeholder="Pre Script"
          v-model="localCopy.post_up"
        >
          <vs-chip
            :key="pre + index"
            @click="remove(localCopy.post_up, pre)"
            v-for="(pre, index) in localCopy.post_up"
            closable
          >
            {{ pre }}
          </vs-chip>
        </vs-chips>
      </vs-col>
    </vs-row>
    <vs-row class="rowFlex">
      <label for="">Post Script</label>
      <vs-col
        vs-lg="4"
        vs-sm="6"
        vs-xs="12"
        :class="{ disable: enableNatRule }"
      >
        <vs-chips
          remove-icon=""
          color="rgb(145, 32, 159)"
          placeholder="Post script"
          v-model="localCopy.post_down"
        >
          <vs-chip
            :key="postScr + index"
            @click="remove(localCopy.post_down, postScr)"
            v-for="(postScr, index) in localCopy.post_down"
            closable
          >
            {{ postScr }}
          </vs-chip>
        </vs-chips>
      </vs-col>
    </vs-row>
    <vs-divider color="primary"></vs-divider>
    <vs-row>
      <vs-button color="success" type="filled" @click="saveConfigPromt"
        >Save Config</vs-button
      >
    </vs-row>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
export default {
  data() {
    return {
      ipTunnelValidation: {
        sucess: false,
        fail: false,
        sucessText: "",
        failText: "",
      },
      localCopy: {},
      // selectInt: "",
      dataIsLoaded: false,
      manualIP: false,
      // enableIpForward: false,
      enableNatRule: false,
      autoKey: false,
    };
  },
  computed: {
    ...mapGetters(["getServer", "getInterfaces"]),
    iptunnel() {
      return this.localCopy.tunnel_address;
    },
  },
  methods: {
    copyKey(key) {
      let ckey = "";
      if (key === "Private") {
        ckey = this.localCopy.private_key;
      }
      if (key === "Public") {
        ckey = this.localCopy.public_key;
      }
      navigator.clipboard.writeText(ckey);
      this.$vs.notify({
        title: "",
        text: `copy ${key.toLowerCase()} key`,
        color: "success",
        time: 1000,
      });
      console.log("click copy");
    },
    remove(sourceArr, item) {
      sourceArr.splice(sourceArr.indexOf(item), 1);
    },
    saveConfigPromt() {
      // if (this.getServer.public_key !="" && this.getServer.private_key !=""){
      this.$vs.dialog({
        color: "danger",
        type: "confirm",
        title: `change server Config`,
        text: "Any change to server config might result in broken clients.this means you need to create config and apply to all client again.please do in carefully",
        accept: this.sendSaveReq,
      });
      // }

      // this.localCopy.manualIpAdd = this.selectInt;
      // console.log("save");
    },
    sendSaveReq() {
      this.$vs.loading();
      this.$store.dispatch("saveServerConfigReq", this.localCopy).then(() => {
        this.$vs.notify({
          title: "save",
          text: "server config save",
          color: "success",
          position: "top-right",
        });
        // this.$vs.loading.close();
        this.$vs.loading();
        this.$store.dispatch("serverInterfacesReq").then(() => {
          this.$store.dispatch("serverConfigReq").then(() => {
            this.localCopy = Object.assign({}, this.getServer);
            this.manualIP = this.localCopy.manual_ip;
            this.autoKey = this.localCopy.auto_generate_key;
            this.enableNatRule = this.localCopy.auto_generate_script;
            this.dataIsLoaded = true;
            this.$vs.loading.close();
          });
        });
      });
    },
  },
  watch: {
    iptunnel: function (val) {
      const regex =
        /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/[0-9][0-9]$/;
      if (val == "") {
        this.ipTunnelValidation = {
          sucess: false,
          fail: false,
          sucessText: "",
          failText: "",
        };
        return;
      }
      if (regex.test(val)) {
        this.ipTunnelValidation = {
          sucess: true,
          fail: false,
          sucessText: "ip is valid",
          failText: "",
        };
        return;
      } else {
        this.ipTunnelValidation = {
          sucess: false,
          fail: true,
          sucessText: "",
          failText: "ip format should be x.x.x.x/x",
        };
        return;
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
    getInterfaces: {
      deep: true,
      handler: function (newval) {
        let selectedInterface = newval.find((item) => {
          return item.IsSelected == true;
        });
        if (selectedInterface === undefined) {
          return;
        }
        if (this.localCopy != undefined) {
          if (this.localCopy.manual_ip == false) {
            this.localCopy.address = selectedInterface.IP;
          }
        }
      },
    },
    autoKey: function (val) {
      this.localCopy.auto_generate_key = val;
      if (val == false) {
        this.localCopy.public_key = "";
        this.localCopy.private_key = "";
      }
      if (val == true) {
        this.localCopy.public_key = this.getServer.public_key;
        this.localCopy.private_key = this.getServer.private_key;
      }
    },
    enableNatRule: function (val) {
      this.localCopy.auto_generate_script = val;
      if (val == false) {
        this.localCopy.post_down = this.getServer.post_down;
        this.localCopy.post_up = this.getServer.post_up;
      }
      if (val == true) {
        this.localCopy.post_down = [];
        this.localCopy.post_up = [];
      }
    },
    manualIP: function (val) {
      this.localCopy.manual_ip = val;
      if (val === false) {
        let selectedInterface = this.getInterfaces.find((item) => {
          return item.IsSelected == true;
        });
        if (selectedInterface != undefined) {
          this.localCopy.address = selectedInterface.IP;
        }
      }
      if (val === true) {
        this.localCopy.address = this.getServer.address;
      }
    },
  },
  created() {
    this.$vs.loading();
    this.$store.dispatch("serverInterfacesReq").then(() => {
      this.$store.dispatch("serverConfigReq").then(() => {
        this.localCopy = Object.assign({}, this.getServer);
        this.manualIP = this.localCopy.manual_ip;
        this.autoKey = this.localCopy.auto_generate_key;
        this.enableNatRule = this.localCopy.auto_generate_script;
        this.dataIsLoaded = true;
        this.$vs.loading.close();
      });
    });
  },
};
</script>

<style scoped>
.copyIcon {
  cursor: pointer;
}
.disable {
  pointer-events: none;
  cursor: not-allowed !important;
  /* visibility: hidden; */
  /* filter:alpha(opacity=25); */
  /* -moz-opacity:.25; */
  opacity: 0.6;
  /* background-color: black; */
  /* height: 100%; width: 100%; */
  /* position: absolute; */
  /* top: 0px; */
  /* left: 0px; */
  /* z-index: 500; just preventing  */
}
/* .con-chips {
  justify-content: flex-start!important;
} */
.lableIp {
  margin: 10px 0 0 0;
}
.rowFlex {
  display: flex;
  flex-direction: column;
}
.form {
  padding-left: 10px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}
.inputCheckBoxRow {
  /* display: flex;
  flex-direction: row;
  align-items: center; */
  align-items: center;
  gap: 10px;
  /* justify-content: space-between; */
}
.inputCheckBox {
  /* padding-top: 10px; */
  margin: 15px 0 0 0;
  /* align-self: baseline; */
  /* background-color: brown; */
}
</style>