<template>
  <vs-popup :button-close-hidden="true" :title="title" :active="popupEnable">
    <div class="addForm">
      <vs-row>
        <vs-col>
          <vs-input label="Name" placeholder="" v-model="popData.name" />
        </vs-col>
      </vs-row>
      <vs-row
        type="flex"
        vs-justify="flex-start"
        vs-align="center"
        class="lableGap"
      >
        <vs-input
          :disabled="isIpDisable"
          label="Client IP"
          :success="clientIPValidation.sucess"
          :success-text="clientIPValidation.sucessText"
          val-icon-success="done"
          :danger="clientIPValidation.fail"
          :danger-text="clientIPValidation.failText"
          val-icon-danger="close"
          placeholder=""
          v-model="popData.allocated_ip"
        />
        <vs-tooltip :text="iptooltipText">
          <vs-icon size="15px" icon="help_outline" color="blue"></vs-icon>
        </vs-tooltip>
      </vs-row>
      <vs-row>
        <vs-col>
          <vs-input
            label="DNS Address"
            :success="dnsIPValidation.sucess"
            :success-text="dnsIPValidation.sucessText"
            val-icon-success="done"
            :danger="dnsIPValidation.fail"
            :danger-text="dnsIPValidation.failText"
            placeholder=""
            v-model="popData.dns_address"
          />
        </vs-col>
      </vs-row>
      <vs-row class="rowFlex">
        <vs-row class="lableGap">
          <label for="">Allowed IP</label>
          <vs-tooltip text="add 0.0.0.0/0 to allow all traffic">
            <vs-icon size="15px" icon="help_outline" color="blue"></vs-icon>
          </vs-tooltip>
        </vs-row>

        <vs-col vs-lg="4" vs-sm="4" vs-xs="12">
          <vs-chips
            remove-icon=""
            color="rgb(145, 32, 159)"
            placeholder="press enter to add"
            v-model="popData.allowed_ips"
          >
            <vs-chip
              :key="ips + index"
              @click="remove(popData.allowed_ips, ips)"
              v-for="(ips, index) in popData.allowed_ips"
              closable
            >
              {{ ips }}
            </vs-chip>
          </vs-chips>
        </vs-col>
      </vs-row>
      <vs-row vs-w="12">
        <vs-col vs-type="flex" vs-w="2">
          <vs-button color="success" type="border" @click="buttonReturnData">
            {{ buttonText }}
          </vs-button>
        </vs-col>
        <vs-col v-if="this.popType == 'update'" vs-type="flex" vs-w="2">
          <vs-button color="danger" type="border" @click="removeClient">
            Remove
          </vs-button>
        </vs-col>
        <vs-col vs-type="flex" vs-w="2">
          <vs-button color="#808080" type="border" @click="closePopup">
            Cancel
          </vs-button>
        </vs-col>
      </vs-row>
    </div>
  </vs-popup>
</template>

<script>
export default {
  name: "AppClientPopup",
  emits: ["popup-data", "close-popup", "remove-client"],
  props: {
    title: String,
    availabaleIPs: Array,
    popType: String,
    popupEnable: Boolean,
    popData: Object,
  },
  data() {
    return {
      clientIPValidation: {
        sucess: false,
        fail: false,
        sucessText: "",
        failText: "",
      },
      dnsIPValidation: {
        sucess: false,
        fail: false,
        sucessText: "",
        failText: "",
      },
    };
  },
  computed: {
    isIpDisable(){
      if (this.popType == 'update'){
        return true
      }else{
        return false
      }
    },
    
    buttonText: function () {
      if (this.popType == "add") {
        return "Add Client";
      }
      if (this.popType == "update") {
        return "Update";
      }
      return "";
    },
    clientIP() {
      return this.popData.allocated_ip;
    },
    dnsIPForm(){
      return this.popData.dns_address;
    },
    iptooltipText() {
      return `ip should be in server range E.g. : ${this.availabaleIPs[0]}`;
    },
  },
  methods: {
    removeClient() {
      this.$emit("remove-client", this.popData);
    },
    remove(sourceArr, item) {
      sourceArr.splice(sourceArr.indexOf(item), 1);
    },
    printConsole(val) {
      console.log("on fucose");
      console.log(val);
    },
    closePopup() {
      this.$emit("close-popup");
    },
    buttonReturnData() {
      if (this.clientIPValidation.fail){
        this.$vs.notify({
        title: "client ip",
        text: this.clientIPValidation.failText,
        color: "danger",
        time: 3000,
        });
        return
      }
      if (this.dnsIPValidation.fail){
        this.$vs.notify({
        title: "dns",
        text: this.dnsIPValidation.failText,
        color: "danger",
        time: 3000,
        });
        return
      }
      if (this.popData.allowed_ips.length<=0){
        this.$vs.notify({
        title: "allow ips",
        text: "allow ip list cant be empty.at least add 0.0.0.0/0",
        color: "danger",
        time: 3000,
        });
        return
      }
      this.$emit("popup-data", this.popData);
    },
  },
  watch: {
    clientIP: function (val) {
      if (this.isIpDisable){
        this.clientIPValidation = {
          sucess: false,
          fail: false,
          sucessText: "",
          failText: "",
        };
        return;
      }
      const regex =
        /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
      if (val == "") {
        this.clientIPValidation = {
          sucess: false,
          fail: false,
          sucessText: "",
          failText: "",
        };
        return;
      }
      if (regex.test(val)) {
        if (this.availabaleIPs.includes(val)) {
          this.clientIPValidation = {
            sucess: false,
            fail: true,
            sucessText: "",
            failText: `this ip used already`,
          };
          return;
        }
        this.clientIPValidation = {
          sucess: true,
          fail: false,
          sucessText: "ip is valid",
          failText: "",
        };
        return;
      } else {
        this.clientIPValidation = {
          sucess: false,
          fail: true,
          sucessText: "",
          failText: "not valid ip range",
        };
        return;
      }
    },
    //-------------------
    dnsIPForm:function (val) {
      const regex =
        /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
      if (val!=""){
      if (regex.test(val)) {
        this.dnsIPValidation = {
          sucess: true,
          fail: false,
          sucessText: "ip is valid",
          failText: "",
        };
        return;
      } else {
        this.dnsIPValidation = {
          sucess: false,
          fail: true,
          sucessText: "",
          failText: "not valid ip range",
        };
        return;
      }
      }else{
        this.dnsIPValidation = {
          sucess: false,
          fail: false,
          sucessText: "",
          failText: "",
        };
        return;
      }

    }
  },
};
</script>

<style scoped>
.lableGap {
  gap: 5px;
}
.rowFlex {
  display: flex;
  flex-direction: column;
}
.addForm {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>