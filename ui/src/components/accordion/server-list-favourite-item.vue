<template>
  <b-card no-body class="w-100 mb-2" v-on:click="onServerChanged(server)">
    <b-card-body class="py-2">
      <div class="d-flex align-items-center justify-content-between">
        <div class="d-flex align-items-center">
          <img :src="serverImage" alt="" style="height: 40px; width: 40px" />
          <span class="text-start ms-2">
            <h6 class="mb-0">{{ serverCountry }}</h6>
            <p class="mb-0">{{ serverName }}</p>
          </span>
        </div>

        <div class="d-flex">
          <serverPingInfoControl
            :isShowPingTime="true"
            :server="server"
            class="pingInfo mx-2"
          />
          <!-- <button class="btn p-1">
            <img
              :src="favoriteImage(server)"
              :alt="server.ip"
              v-on:click="favoriteClicked($event, server)"
            />
          </button> -->
        </div>
      </div>
    </b-card-body>
  </b-card>
</template>

<script>
import ImageFavouritesOn from "@/assets/img/fvrts-on.svg";
import ImageFavouritesOff from "@/assets/img/fvrts-off.svg";
import serverPingInfoControl from "@/components/controls/control-server-ping.vue";

export default {
  components: { serverPingInfoControl },
  props: ["server", "onServerChanged"],
  data: () => ({
    isImgLoadError: false,
  }),
  computed: {
    serverImage: function () {
      if (!this.server) return `/flags/unk.svg`;
      try {
        const code = this.server.country_code.toLowerCase();
        return `/flags/svg/${code}.svg`;
      } catch (e) {
        console.log(e);
        return null;
      }
    },
    serverName: function () {
      return this.server.name;
    },
    serverCountry: function () {
      return this.server.country;
    },
  },

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
      // let serversHashed = this.$store.state.vpnState.serversHashed;
      let gateway = server.id;
      if (favorites.includes(gateway)) {
        // remove
        console.log(`Removing favorite ${gateway}`);
        favorites = favorites.filter((gw) => gw !== gateway);
      } else {
        // add
        console.log(`Adding favorite ${gateway}`);
        // if (serversHashed[gateway] == null) return;
        favorites.push(gateway);
      }

      this.$store.dispatch("settings/serversFavoriteList", favorites);
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
@import "@/components/scss/constants";

.accordion-button {
  align-items: unset;
}

.accordion-button::after {
  flex-shrink: 0;
  width: 1.25rem;
  height: 1.25rem;
  margin-left: auto;
  content: "";
  background-image: v-bind(ArrowUp);
  background-repeat: no-repeat;
  transition: transform 0.2s ease-in-out 0s;
}

.dark .accordion-button:not(.collapsed)::after {
  background-image: v-bind(ArrowRight);
}
</style>
