<template>
  <!-- @selected="handleSelected" -->
  <div>
    <vs-table
      v-model="selected"
      @dblSelection.self="doubleSelection"
      max-items="10"
      pagination
      stripe
      :data="tableData"
    >
      <template slot="header">
        <h2>
          {{ title }}
        </h2>
      </template>
      <template slot="thead">
        <vs-th :key="index" v-for="(item, index) in fields"> {{ item }} </vs-th>
      </template>

      <template slot-scope="{ data }">
        <vs-tr :data="tr" :key="indextr" v-for="(tr, indextr) in data">
          <vs-td :data="data[indextr].name">
            {{ data[indextr].name }}
          </vs-td>

          <vs-td :data="data[indextr].allocated_ip">
            {{ data[indextr].allocated_ip }}
          </vs-td>

          <vs-td :data="data[indextr].allowed_ips">
            {{ data[indextr].allowed_ips }}
          </vs-td>
          <vs-td :data="data[indextr]">
            <vs-button
              @click="downloadqrcode(data[indextr])"
              size="small"
              color="primary"
              type="border"
              >QrCode</vs-button
            >
          </vs-td>
          <vs-td :data="data[indextr]">
            <vs-button
              @click="downloadConfig(data[indextr])"
              size="small"
              color="primary"
              type="border"
              >config</vs-button
            >
          </vs-td>
          <vs-td :data="data[indextr].enabled">
            <vs-button
              size="default"
              @click="changeClientStat(indextr)"
              :color="clientStateIconColor(data[indextr].enabled)"
              type="line"
              :icon="clientStateIcon(data[indextr].enabled)"
            ></vs-button>
            <!-- <vs-switch v-model="clientState" :vs-value="data[indextr].name">
              <span slot="on">Enable</span>
              <span slot="off">Disable</span>
            </vs-switch> -->
            <!-- {{ data[indextr].enabled }} -->
          </vs-td>
        </vs-tr>
      </template>
    </vs-table>
  </div>
</template>


<script>
export default {
  name: "AppClientTable",
  emits: ["selected-row", "download-qrcode", "download-config","change-state"],
  props: {
    title: String,
    fields: Array,
    tableData: Array,
  },
  data() {
    return {
      selected: [],
      clientState: [],
    };
  },
  methods: {
    handleSelected() {
      return undefined;
    },
    downloadConfig(client) {
      this.$emit("download-config", client);
    },
    downloadqrcode(client) {
      this.$emit("download-qrcode", client);
    },
    doubleSelection(tr) {
      this.$emit("selected-row", tr);
    },
    clientStateText: function (isEnable) {
      if (isEnable) {
        return "Enable";
      } else {
        return "Disable";
      }
    },
    clientStateIcon: function (isEnable) {
      if (isEnable) {
        return "toggle_on";
      } else {
        return "toggle_off";
      }
    },
    clientStateIconColor: function (isEnable) {
      if (isEnable) {
        return "primary";
      } else {
        return "danger";
      }
    },
    changeClientStat: function (indextr) {
      this.tableData[indextr].enabled = !this.tableData[indextr].enabled;
      this.$emit("change-state", this.tableData[indextr]);
      // this.tableData[indextr].enabled = !this.tableData[indextr].enabled;
    },
  },
  watch: {
  },
  computed: {
  },
};
</script>

<style scoped>
/* .con-expand-users .con-btns-user {
  display: flex;
  padding: 10px;
  padding-bottom: 0px;
  align-items: center;
  justify-content: space-between;
}
.con-expand-users .con-btns-user .con-userx {
  display: flex;
  align-items: center;
  justify-content: flex-start;
}
.con-expand-users .list-icon i {
  font-size: 0.9rem !important;
} */
</style>