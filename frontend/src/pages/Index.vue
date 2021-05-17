<template>
  <q-page class="column items-center">
    <div class="row">
      <div class="column ipinfoLayout">
        <div class="col">
          <div class="column">
            <div class="text-h6 text-black text-center">
              IP 查詢
            </div>
            <q-input
              v-model="ipQuery"
              rounded
              outlined
              bg-color="white"
              class="q-pa-md"
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
          </div>
        </div>

        <div class="col">
          <div class="column">
            <div class="row col-10">
              <div class="row">IP 位址:</div>
              <div class="row">{{ ipReturn }}</div>
            </div>
            <div class="row col-10">
              <div class="row">國家:</div>
              <div class="row">{{ contry }}</div>
            </div>
            <div class="row  col-10">
              <div class="row">城市:</div>
              <div class="row">{{ city }}</div>
            </div>
            <div class="row  col-10">
              <div class="row">經度:</div>
              <div class="row">{{ longitude }}</div>
            </div>
            <div class="row  col-10">
              <div class="row">緯度:</div>
              <div class="row">{{ latitude }}</div>
            </div>
          </div>
        </div>
      </div>
      <div class="column  mapLayout items-center">
        <l-map
          ref="map"
          :zoom="mapZoom"
          :center="mapCenter"
          :options="mapOptions"
        >
          <l-tile-layer :url="mapUrl" :attribution="mapAttribution" />
          <l-marker v-if="mapShowMark" :lat-lng="mapMark"
            ><l-icon :icon-url="iconUrl" />
          </l-marker>
        </l-map>
      </div>
    </div>
  </q-page>
</template>

<script>
import L from "leaflet";
import { LMap, LTileLayer, LMarker, LIcon } from "vue2-leaflet";
import "leaflet/dist/leaflet.css";
import numeral from "numeral";

import { axios } from "boot/axios";

export default {
  name: "PageIndex",
  components: {
    LMap,
    LTileLayer,
    LMarker,
    LIcon
  },
  data() {
    return {
      ipQuery: "",
      ipReturn: "",
      contry: "",
      city: "",
      longitude: 0,
      latitude: 0,
      mapShowMark: false,
      iconUrl:
        "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABkAAAApCAYAAADAk4LOAAAFgUlEQVR4Aa1XA5BjWRTN2oW17d3YaZtr2962HUzbDNpjszW24mRt28p47v7zq/bXZtrp/lWnXr337j3nPCe85NcypgSFdugCpW5YoDAMRaIMqRi6aKq5E3YqDQO3qAwjVWrD8Ncq/RBpykd8oZUb/kaJutow8r1aP9II0WmLKLIsJyv1w/kqw9Ch2MYdB++12Onxee/QMwvf4/Dk/Lfp/i4nxTXtOoQ4pW5Aj7wpici1A9erdAN2OH64x8OSP9j3Ft3b7aWkTg/Fm91siTra0f9on5sQr9INejH6CUUUpavjFNq1B+Oadhxmnfa8RfEmN8VNAsQhPqF55xHkMzz3jSmChWU6f7/XZKNH+9+hBLOHYozuKQPxyMPUKkrX/K0uWnfFaJGS1QPRtZsOPtr3NsW0uyh6NNCOkU3Yz+bXbT3I8G3xE5EXLXtCXbbqwCO9zPQYPRTZ5vIDXD7U+w7rFDEoUUf7ibHIR4y6bLVPXrz8JVZEql13trxwue/uDivd3fkWRbS6/IA2bID4uk0UpF1N8qLlbBlXs4Ee7HLTfV1j54APvODnSfOWBqtKVvjgLKzF5YdEk5ewRkGlK0i33Eofffc7HT56jD7/6U+qH3Cx7SBLNntH5YIPvODnyfIXZYRVDPqgHtLs5ABHD3YzLuespb7t79FY34DjMwrVrcTuwlT55YMPvOBnRrJ4VXTdNnYug5ucHLBjEpt30701A3Ts+HEa73u6dT3FNWwflY86eMHPk+Yu+i6pzUpRrW7SNDg5JHR4KapmM5Wv2E8Tfcb1HoqqHMHU+uWDD7zg54mz5/2BSnizi9T1Dg4QQXLToGNCkb6tb1NU+QAlGr1++eADrzhn/u8Q2YZhQVlZ5+CAOtqfbhmaUCS1ezNFVm2imDbPmPng5wmz+gwh+oHDce0eUtQ6OGDIyR0uUhUsoO3vfDmmgOezH0mZN59x7MBi++WDL1g/eEiU3avlidO671bkLfwbw5XV2P8Pzo0ydy4t2/0eu33xYSOMOD8hTf4CrBtGMSoXfPLchX+J0ruSePw3LZeK0juPJbYzrhkH0io7B3k164hiGvawhOKMLkrQLyVpZg8rHFW7E2uHOL888IBPlNZ1FPzstSJM694fWr6RwpvcJK60+0HCILTBzZLFNdtAzJaohze60T8qBzyh5ZuOg5e7uwQppofEmf2++DYvmySqGBuKaicF1blQjhuHdvCIMvp8whTTfZzI7RldpwtSzL+F1+wkdZ2TBOW2gIF88PBTzD/gpeREAMEbxnJcaJHNHrpzji0gQCS6hdkEeYt9DF/2qPcEC8RM28Hwmr3sdNyht00byAut2k3gufWNtgtOEOFGUwcXWNDbdNbpgBGxEvKkOQsxivJx33iow0Vw5S6SVTrpVq11ysA2Rp7gTfPfktc6zhtXBBC+adRLshf6sG2RfHPZ5EAc4sVZ83yCN00Fk/4kggu40ZTvIEm5g24qtU4KjBrx/BTTH8ifVASAG7gKrnWxJDcU7x8X6Ecczhm3o6YicvsLXWfh3Ch1W0k8x0nXF+0fFxgt4phz8QvypiwCCFKMqXCnqXExjq10beH+UUA7+nG6mdG/Pu0f3LgFcGrl2s0kNNjpmoJ9o4B29CMO8dMT4Q5ox8uitF6fqsrJOr8qnwNbRzv6hSnG5wP+64C7h9lp30hKNtKdWjtdkbuPA19nJ7Tz3zR/ibgARbhb4AlhavcBebmTHcFl2fvYEnW0ox9xMxKBS8btJ+KiEbq9zA4RthQXDhPa0T9TEe69gWupwc6uBUphquXgf+/FrIjweHQS4/pduMe5ERUMHUd9xv8ZR98CxkS4F2n3EUrUZ10EYNw7BWm9x1GiPssi3GgiGRDKWRYZfXlON+dfNbM+GgIwYdwAAAAASUVORK5CYII=",
      mapZoom: 5.5,
      mapCenter: L.latLng(23.5, 121),
      mapMark: L.latLng(23.5, 121),
      mapOptions: {
        zoomSnap: 0.5
      },
      mapUrl: "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
      mapAttribution:
        '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
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
        var orgLongitude = data["Longitude"];
        var orgLatitude = data["Latitude"];
        if (this.contry == "") {
          this.mapShowMark = false;
        } else {
          this.mapShowMark = true;
          this.mapCenter = new L.LatLng(orgLatitude, orgLongitude);
          this.mapMark = new L.LatLng(orgLatitude, orgLongitude);
        }
        this.longitude = numeral(orgLongitude).format("0.00");
        this.latitude = numeral(orgLatitude).format("0.00");
      });
    }
  }
};
</script>

<style lang="sass" scoped>
.column > div
  padding: 5px 5px
.mapLayout
  width: 325px
  height: 450px
.ipinfoLayout
  width: 350px
  height: 450px
</style>
