<template>
  <div>
    <vs-table max-items="10" pagination stripe :data="tableData">
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
            {{ tr.name }}
          </vs-td>

          <vs-td :data="tr.allocated_ip">
            {{ tr.allocated_ip }}
          </vs-td>

          <vs-td :data="tr.send_byte">
            {{ formatBytes(tr.send_byte) }}
          </vs-td>
          <vs-td :data="tr.receive_byte">
            {{ formatBytes(tr.receive_byte) }}
          </vs-td>
          <vs-td :data="tr.handshake">
            <!-- {{tr.handshake}} -->
            {{
              tr.handshake.seen
                ? formatTimeToLocal(tr.handshake.last_handshake_time)
                : "-"
            }}
          </vs-td>
        </vs-tr>
      </template>
    </vs-table>
  </div>
</template>

<script>
import moment from "moment";
export default {
  name: "ApptTafficTable",
  emits: [],
  props: {
    title: String,
    fields: Array,
    tableData: Array,
  },
  data() {
    return {};
  },
  methods: {
    formatBytes(bytes, decimals) {
      if (bytes == 0) return "0 Bytes";
      var k = 1024,
        dm = decimals || 2,
        sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"],
        i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
    },
    formatTimeToLocal(time) {
      const timeToLocal = moment.utc(time).toDate();
      const human = moment(timeToLocal).fromNow();
      return human;
    },
    moment: function () {
      return moment();
    },
  },
  watch: {},
};
</script>

<style scoped>
</style>