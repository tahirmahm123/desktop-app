<template>
  <div style="cursor: pointer">
    <div
      class="d-flex justify-content-between"
      v-for="host of servers"
      v-bind:key="`Server` + host.id"
    >
      <div class="form-check">
        <input
          class="form-check-input"
          type="radio"
          :id="`host-id-` + host.id"
          :checked="selected(host.id)"
          v-on:change="onServerChanged(host)"
        />
        <label class="form-check-label" :for="`host-id-` + host.id">
          {{ host.name }}
        </label>
      </div>
      <div class="d-flex">
        <serverPingInfoControl
          :isShowPingTime="true"
          :server="host"
          class="pingInfo mx-3"
        />
        <div class="fav-icon checked signals">
          <div class="network-latency signal-3"></div>
          <!-- <div>
            <img
              :src="favoriteImage(host)"
              :alt="host.ip"
              v-on:click="favoriteClicked($event, host)"
            />
          </div> -->
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import serverPingInfoControl from "@/components/controls/control-server-ping.vue";

import ImageFavouritesOn from "@/assets/img/fvrts-on.svg";
import ImageFavouritesOff from "@/assets/img/fvrts-off.svg";

export default {
  components: { serverPingInfoControl },
  props: ["servers", "expanded", "onServerChanged"],
  data: () => ({
    isImgLoadError: false,
  }),
  computed: {},

  methods: {
    onImgLoadError() {
      this.isImgLoadError = true;
    },
    favoriteImage: function (server) {
      if (this.$store.state.settings.serversFavoriteList.includes(server.id))
        return ImageFavouritesOn;
      return ImageFavouritesOff;
    },
    favoriteClicked: function (evt, server) {
      evt.stopPropagation();
      let favorites = this.$store.state.settings.serversFavoriteList.slice();
      let serversHashed = this.$store.state.vpnState.serversHashed;
      let gateway = server.id;
      if (favorites.includes(gateway)) {
        favorites = favorites.filter((gw) => gw !== gateway);
      } else {
        // add
        if (serversHashed[server.ip] == null) return;
        favorites.push(gateway);
      }

      this.$store.dispatch("settings/serversFavoriteList", favorites);
    },

    selected: function (id) {
      let server = this.$store.state.settings.serverEntry;
      return server !== null && server.id === id;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
@import "@/components/scss/constants";

.main {
  display: flex;
}

img.pic {
  width: 22px;
  margin: 1px;
  margin-top: 2.4px;
}

img.flag {
  //border: 1px solid rgba(var(--flag-border-color-rgb), 0.5);

  //border-radius: 2px;
  box-shadow: 0 0 0.4pt 0.4pt rgba(var(--flag-border-color-rgb), 1);
}

.text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

div.textBloack {
  font-size: 14px;
  line-height: 20.8px;
  margin-left: 16px;
}

.text_large {
  font-size: 18px;
  line-height: 21px;
  margin-left: 10px;
}

div.firstLine {
  text-align: left;
  font-size: 16px;
}

div.secondLine {
  text-align: left;
  font-size: 12px;
  color: grey;
}

div.bage {
  font-size: 10px;
  line-height: 14px;
  color: grey;
  display: inline-block;
  padding-left: 4px;
  padding-right: 4px;
  border: 1px solid #777777;
  border-radius: 4px;
}

.flexRow {
  display: flex;
  align-items: center;
}

.marginLeft {
  margin-left: 9px;
}

.pingtext {
  color: var(--text-color-details);
}

.card-body {
  background: rgba(246, 246, 246, 1) !important;
}

.form-check-input:checked {
  background-color: #ebc553;
  border-color: #ebc553;
}

.accordion .collapse {
  background-color: rgba(41, 41, 48, 1);
}
</style>
