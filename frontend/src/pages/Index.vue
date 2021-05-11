<template>
  <q-page class="flex flex-center">
    <q-card class="ipInfo-card">
      <q-card-section>
        <div class="text-h6 text-black text-center">IP 查詢</div>
        <q-input
          v-model="ipQuery"
          rounded
          outlined
          bg-color="white"
          class="col-10 q-pa-md"
          placeholder="輸入 IP"
          @keydown.enter.prevent="getIPInfo"
          :dense="false"
        >
          <template #append>
            <q-icon
              v-if="ipQuery !== ''"
              name="close"
              class="cursor-pointer"
              @click="ipQuery = ''"
            />
            <q-icon name="search" glossy @click="getIPInfo" />
          </template>
        </q-input>
      </q-card-section>

      <q-separator dark inset />

      <q-card-section class="text-black flex-center">
        <div class="column">
          <div class="row col-10">
            <div class="">IP 位址:</div>
            <div class="">{{ ipReturn }}</div>
          </div>
          <div class="row">
            <div>國家:</div>
            <div>{{ contry }}</div>
          </div>
          <div class="row">
            <div>城市:</div>
            <div>{{ city }}</div>
          </div>
          <div class="row">
            <div>地理位置:</div>
            <div>{{ coordinates }}</div>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script>
import { axios } from "boot/axios";

export default {
  name: "PageIndex",
  components: {},
  data() {
    return {
      ipQuery: "",
      ipReturn: "",
      contry: "",
      city: "",
      coordinates: ""
    };
  },
  created() {
    this.getIPInfo();
  },
  methods: {
    getIPInfo() {
      var path = "/ipinfo/" + this.ipQuery;
      axios.get(path).then(res => {
        var data = res.data;
        this.ipReturn = data["IP"];
        this.contry = data["Country"];
        this.city = data["City"];
        this.coordinates = data["Coordinates"];
      });
    }
  }
};
</script>

<style lang="sass" scoped>
.ipInfo-card
  width: 95%
  max-width: 480px
.search-bar
 width: 95%
</style>
