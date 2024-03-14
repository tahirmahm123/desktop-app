<template>
  <div class="main-server px-3">
    <div class="pt-5 d-flex justify-content-between align-items-center">
      <h4>Choose Locations</h4>
      <button @click="closeServerSection" class="border-0 close-btn">
        <svg
          width="38"
          height="38"
          viewBox="0 0 38 38"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M11.7716 29.6955H26.2164C28.8048 29.6955 30.1407 28.3595 30.1407 25.8069V11.2905C30.1407 8.73793 28.8048 7.40199 26.2164 7.40199H11.7716C9.19516 7.40199 7.84729 8.726 7.84729 11.2905V25.8069C7.84729 28.3595 9.19516 29.6955 11.7716 29.6955ZM14.9683 23.7314C14.3361 23.7314 13.8352 23.2305 13.8352 22.5864C13.8352 22.3001 13.9544 22.0138 14.1811 21.7991L17.4255 18.5428L14.1811 15.2983C13.9544 15.0836 13.8352 14.7974 13.8352 14.5111C13.8352 13.867 14.3361 13.3779 14.9683 13.3779C15.3023 13.3779 15.5647 13.4853 15.7794 13.7L19.0358 16.9444L22.304 13.6881C22.5426 13.4614 22.7931 13.3541 23.1151 13.3541C23.7473 13.3541 24.2483 13.855 24.2483 14.4872C24.2483 14.7854 24.129 15.0479 23.9024 15.2864L20.658 18.5428L23.9024 21.7872C24.1171 22.0138 24.2364 22.2882 24.2364 22.5864C24.2364 23.2305 23.7354 23.7314 23.0913 23.7314C22.7692 23.7314 22.4949 23.6122 22.2683 23.3975L19.0358 20.165L15.8033 23.3975C15.5886 23.6241 15.3023 23.7314 14.9683 23.7314Z"
            fill="white"
          />
        </svg>
      </button>
    </div>
    <div class="serversButtonsSpace" />

    <!-- HEADER -->
    <ul v-if="isFastestServerConfig" class="nav d-flex justify-content-around">
      <li class="nav-item">
        <a
          href="#"
          class="nav-link"
          v-bind:class="{ active: !isFavoritesView }"
          v-on:click.prevent="showAll"
          >Fastest Server Settings</a
        >
      </li>
    </ul>

    <!-- EMPTY FAVORITE SERVERS DESCRIPTION BLOCK -->
    <div v-if="isShowFavoriteDescriptionBlock">
      <div>
        <h4 style="width: 300px; padding-bottom: 5px; margin: 0 auto">
          Your Favourite servers will be displayed here
        </h4>
        <p style="font-size: 12px; color: rgb(191, 191, 191)">
          Save your time by creating your won list of servers.
        </p>
      </div>
    </div>

    <!-- FILTER -->
    <div v-if="!isShowFavoriteDescriptionBlock" class="flexRow mx-2">
      <input
        id="filter"
        v-model="filter"
        class="form-control"
        placeholder="Search for a server"
        v-bind:style="{ backgroundImage: 'url(' + searchImage + ')' }"
      />

      <!--      <div class="buttonWithPopup">
        <button
          v-click-outside="onSortMenuClickedOutside"
          class="noBordersBtn sortBtn sortBtnPlatform"
          v-on:click="onSortMenuClicked()"
        >
          &lt;!&ndash;          <img :src="sortImage" />&ndash;&gt;
        </button>
      </div>-->
    </div>

    <div
      v-if="isFastestServerConfig"
      class="small_text"
      style="margin-bottom: 5px"
    >
      Disable servers you do not want to be choosen as the fastest server
    </div>

    <!-- SERVERS LIST BLOCK -->
    <div
      ref="scrollArea"
      style="height: 440px"
      class="commonMargins m-0 flexColumn scrollableColumnContainer"
      v-on:scroll="recalcScrollButtonVisiblity()"
      v-on:wheel="recalcScrollButtonVisiblity()"
    >
      <!-- FASTEST & RANDOMM SERVER -->
      <div v-if="isFavoritesView === false && isFastestServerConfig === false">
        <div v-if="!isMultihop" class="flexRow">
          <button
            class="serverSelectBtn flexRow"
            v-on:click="onFastestServerClicked()"
          >
            <serverNameControl :isFastestServer="true" class="serverName" />
          </button>
          <button class="noBordersBtn" v-on:click="onFastestServerConfig()">
            <img :src="settingsImage" />
          </button>
        </div>
      </div>

      <!-- SERVERS LIST -->
      <div
        v-for="server of filteredServers"
        v-bind:key="server.flag"
        class="accordion mb-0"
        role="tablist"
      >
        <div class="card">
          <div class="card p-1" role="tab">
            <ServerListAccordionHeader
              :server="server"
              @click="() => onHeaderClicked(server)"
              :expanded="server.flag === expandedHeader"
              :configFastestSvrClicked="configFastestSvrClicked"
              :isFastestServerConfig="isFastestServerConfig"
              :isSvrExcludedFomFastest="isSvrExcludedFomFastest"
            />
          </div>
          <Collapse
            :when="server.flag === expandedHeader"
            v-if="!isFastestServerConfig"
            :id="'accordion-' + server.flag"
            accordion="servers-list"
            role="tabpanel"
          >
            <div class="card-body">
              <ServerListAccordionBody
                :servers="server.servers"
                :expanded="server.flag === expandedHeader"
                :onServerChanged="onServerSelected"
              />
            </div>
          </Collapse>
        </div>
      </div>
    </div>
    <!-- SCROOL DOWN BUTTON -->
    <transition name="fade">
      <button
        style="position: absolute; left: 755px; bottom: 13px"
        v-if="isShowScrollButton"
        class="btnScrollDown"
        v-on:click="onScrollDown()"
      >
        <img src="@/assets/arrow-bottom.svg" />
      </button>
    </transition>
  </div>
</template>

<script>
import { IsOsDarkColorScheme } from "@/helpers/renderer";

const sender = window.ipcSender;
import serverNameControl from "@/components/controls/control-server-name.vue";
import { isStrNullOrEmpty } from "@/helpers/helpers";
import { Platform, PlatformEnum } from "@/platform/platform";
import { enumValueName } from "@/helpers/helpers";
import { ColorTheme, ServersSortTypeEnum } from "@/store/types";
import ServerListAccordionHeader from "@/components/accordion/server-list-accordion-header";
import ServerListAccordionBody from "@/components/accordion/server-list-accordion-body";
import Image_arrow_left_windows from "@/assets/arrow-left-windows.svg";
import Image_arrow_left_macos from "@/assets/arrow-left-macos.svg";
import Image_arrow_left_linux from "@/assets/arrow-left-linux.svg";
import Image_search_windows from "@/assets/search-windows.svg";
import Image_search_macos from "@/assets/search-macos.svg";
import Image_search_linux from "@/assets/search-linux.svg";
import Image_settings_windows from "@/assets/settings-windows.svg";
import Image_settings_macos from "@/assets/settings-macos.svg";
import Image_settings_linux from "@/assets/settings-linux.svg";
import Image_sort from "@/assets/sort.svg";
import Image_check_thin from "@/assets/check-thin.svg";
import ImageEmptyFavouritesLight from "@/assets/img/empty-favourites-light.svg";
import ImageEmptyFavouritesDark from "@/assets/img/empty-favourites-dark.svg";

import vClickOutside from "click-outside-vue3";
import { Collapse } from "vue-collapsed";

export default {
  directives: {
    clickOutside: vClickOutside.directive,
  },
  props: [
    "closeServerSection",
    "onServerChanged",
    "isExitServer",
    "onFastestServer",
    "onRandomServer",
  ],
  components: {
    Collapse,
    serverNameControl,
    ServerListAccordionHeader,
    ServerListAccordionBody,
  },
  data: function () {
    return {
      isDarkTheme: false,
      expandedHeader: "",
      filter: "",
      isFastestServerConfig: false,
      isSortMenu: false,
      isShowScrollButton: false,
    };
  },
  created: function () {},
  mounted() {
    // COLOR SCHEME
    window.matchMedia("(prefers-color-scheme: light)").addListener(() => {
      console.log("media changed");
      this.updateColorScheme();
    });
    this.updateColorScheme();
    this.recalcScrollButtonVisiblity();
    const resizeObserver = new ResizeObserver(this.recalcScrollButtonVisiblity);
    resizeObserver.observe(this.$refs.scrollArea);
  },
  computed: {
    isMinimizedUI: function () {
      return this.$store.state.settings.minimizedUI;
    },
    isFavoritesView: function () {
      return this.$store.state.uiState.serversFavoriteView;
    },
    isMultihop: function () {
      return this.$store.state.settings.isMultiHop;
    },
    isShowFavoriteDescriptionBlock: function () {
      if (!this.isFavoritesView) return false;
      let favSvrs = this.favoriteServers;
      return favSvrs == null || favSvrs.length === 0;
    },
    servers: function () {
      return this.$store.getters["vpnState/activeServers"];
    },

    sortTypeStr: function () {
      return enumValueName(
        ServersSortTypeEnum,
        this.$store.state.settings.serversSortType,
      );
    },

    favoriteServers: function () {
      let favorites = this.$store.state.settings.serversFavoriteList;
      let favouriteServersList = [];
      this.servers.forEach((s) => {
        favouriteServersList = favouriteServersList.concat(
          s.servers.filter((host) => favorites.includes(host.id)),
        );
      });
      return favouriteServersList;
    },

    filteredServers: function () {
      let servers = this.servers;
      if (this.isFavoritesView) servers = this.favoriteServers;

      let filtered = (servers ? servers : []).filter(
        (s) =>
          (s.city && s.city.toLowerCase().includes(this.filter)) ||
          (s.country && s.country.toLowerCase().includes(this.filter)) ||
          (s.country_code &&
            s.country_code.toLowerCase().includes(this.filter)),
      );

      return filtered.slice() /*.sort(compare)*/;
    },

    arrowLeftImagePath: function () {
      switch (Platform()) {
        case PlatformEnum.Windows:
          return Image_arrow_left_windows;
        case PlatformEnum.macOS:
          return Image_arrow_left_macos;
        default:
          return Image_arrow_left_linux;
      }
    },
    searchImage: function () {
      if (!isStrNullOrEmpty(this.filter)) return null;

      switch (Platform()) {
        case PlatformEnum.Windows:
          return Image_search_windows;
        case PlatformEnum.macOS:
          return Image_search_macos;
        default:
          return Image_search_linux;
      }
    },
    settingsImage: function () {
      switch (Platform()) {
        case PlatformEnum.Windows:
          return Image_settings_windows;
        case PlatformEnum.macOS:
          return Image_settings_macos;
        default:
          return Image_settings_linux;
      }
    },
    sortImage: function () {
      return Image_sort;
    },
    selectedImage: function () {
      return Image_check_thin;
    },

    emptyfavouriteImage: function () {
      return this.isDarkTheme
        ? ImageEmptyFavouritesDark
        : ImageEmptyFavouritesLight;
    },
  },

  methods: {
    onBack: function () {
      // this.$router.push("/");
      this.closeServerSection();
    },
    goBack: function () {
      if (this.isFastestServerConfig) {
        this.filter = "";
        this.isFastestServerConfig = false;
        return;
      }
      if (this.onBack != null) this.onBack();
    },
    onHeaderClicked: function (server) {
      if (this.expandedHeader === server.flag) {
        this.expandedHeader = "";
      } else {
        this.expandedHeader = server.flag;
      }
    },
    onServerSelected: function (server) {
      if (this.isInaccessibleServer(server)) {
        sender.showMessageBoxSync({
          type: "info",
          buttons: ["OK"],
          message: "Entry and exit servers cannot be in the same country",
          detail:
            "When using multihop you must select entry and exit servers in different countries. Please select a different entry or exit server.",
        });
        return;
      }
      console.log("trying to update server");
      this.onServerChanged(server, this.isExitServer != null);
      this.onBack();
    },
    onSortMenuClickedOutside: function () {
      this.isSortMenu = false;
    },
    onSortMenuClicked: function () {
      this.isSortMenu = !this.isSortMenu;
    },
    onFastestServerClicked() {
      if (this.onFastestServer != null) this.onFastestServer();
      this.onBack();
    },
    onRandomServerClicked() {
      if (this.onRandomServer != null) this.onRandomServer();
      this.onBack();
    },
    isSvrExcludedFomFastest: function (server) {
      return this.$store.state.settings.serversFastestExcludeList.includes(
        server.flag,
      );
    },

    configFastestSvrClicked(server, event) {
      if (server == null || server.flag == null) return;
      let excludeSvrs =
        this.$store.state.settings.serversFastestExcludeList.slice();

      if (excludeSvrs.includes(server.flag))
        excludeSvrs = excludeSvrs.filter((gw) => gw !== server.flag);
      else excludeSvrs.push(server.flag);

      const activeServers = this.servers.slice();
      const notExcludedActiveServers = activeServers.filter(
        (s) => !excludeSvrs.includes(s.flag),
      );

      if (notExcludedActiveServers.length < 1) {
        sender.showMessageBoxSync({
          type: "info",
          buttons: ["OK"],
          message: "Please, keep at least one server",
          detail: "Not allowed to exclude all servers.",
        });
        event.preventDefault();
        return;
      } else
        this.$store.dispatch("settings/serversFastestExcludeList", excludeSvrs);
    },
    onFastestServerConfig() {
      this.isFastestServerConfig = true;
      this.filter = "";
    },
    isInaccessibleServer: function (server) {
      if (this.$store.state.settings.isMultiHop === false) return false;
      let ccSkip = "";

      let connected = !this.$store.getters["vpnState/isDisconnected"];
      if (
        // ENTRY SERVER
        !this.isExitServer &&
        this.$store.state.settings.serverExit &&
        (connected || !this.$store.state.settings.isRandomExitServer)
      )
        ccSkip = this.$store.state.settings.serverExit.flag;
      else if (
        // EXIT SERVER
        this.isExitServer &&
        this.$store.state.settings.serverEntry &&
        (connected || !this.$store.state.settings.isRandomServer)
      )
        ccSkip = this.$store.state.settings.serverEntry.flag;
      if (server.flag === ccSkip) return true;
      return false;
    },

    showFavorites: function () {
      this.$store.dispatch("uiState/serversFavoriteView", true);
      this.filter = "";
    },
    showAll: function () {
      this.$store.dispatch("uiState/serversFavoriteView", false);
      this.filter = "";
    },
    updateColorScheme() {
      let scheme = sender.ColorScheme();
      console.log(scheme);
      if (scheme === ColorTheme.system) {
        this.isDarkTheme = IsOsDarkColorScheme();
      } else this.isDarkTheme = scheme === ColorTheme.dark;
    },
    recalcScrollButtonVisiblity() {
      console.log("Scroll Btn Added");
      let sa = this.$refs.scrollArea;
      if (sa == null) {
        this.isShowScrollButton = false;
        return;
      }

      const show = sa.scrollHeight > sa.clientHeight + sa.scrollTop;

      // hide - imadiately; show - with 1sec delay
      if (!show) this.isShowScrollButton = false;
      else {
        setTimeout(() => {
          this.isShowScrollButton =
            sa.scrollHeight > sa.clientHeight + sa.scrollTop;
        }, 1000);
      }
    },
    onScrollDown() {
      console.log("Scrolled");
      let sa = this.$refs.scrollArea;
      if (sa == null) return;
      sa.scrollTo({
        top: sa.clientHeight * 0.9 + sa.scrollTop, //sa.scrollHeight,
        behavior: "smooth",
      });
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
@import "@/components/scss/constants";
@import "@/components/scss/popup";

$paddingLeftRight: 20px;

.commonMargins {
  margin-left: $paddingLeftRight;
  margin-right: $paddingLeftRight;
}

input#filter {
  background-position: 97% 50%; //right
  background-repeat: no-repeat;
  margin-top: $paddingLeftRight;
  margin-bottom: $paddingLeftRight;
}

.disabledButton {
  opacity: 0.5;
}

.serverSelectBtn {
  border: none;
  background-color: inherit;
  outline-width: 0;
  border-radius: 8px;
  color: #fff !important;
  cursor: pointer;
  margin: 10px 0;
  height: 48px;
  width: 100%;

  padding: 10px;
}

.serverName {
  max-width: 195px;
  width: 195px;
}

.pingInfo {
  max-width: 72px;
  width: 72px;
}

.pingtext {
  margin-left: 8px;
}

.text {
  margin: $paddingLeftRight;
  margin-top: 60px;
  text-align: center;
}

.small_text {
  margin-left: $paddingLeftRight;
  margin-right: $paddingLeftRight;
  font-size: 11px;
  line-height: 13px;
  color: var(--text-color-details);
}

button.sortBtn {
  margin-left: 5px;
}

div.sortSelectedImg {
  margin-left: 11px;
  position: absolute;
  left: 0px;
  min-width: 13px;
}

.empty-fav-img-bg {
  height: 100vh;
  width: 100%;
  text-align: center;
  padding-top: 30%;
}

.card {
  --bs-card-border-color: none;
  --bs-card-bg: rgba(246, 246, 246, 1);
  --bs-card-spacer-y: 5px;
  --bs-card-spacer-x: 15px;
}

.card-header {
  border-radius: 8px !important;
  background-color: rgba(246, 246, 246, 1) !important;
  padding: 0 !important;
  border-bottom: none;
  border-color: rgba(246, 246, 246, 1);
}

.accordion {
  border-radius: 0 !important;
}

.dark .card {
  --bs-card-border-color: none;
  --bs-card-bg: rgba(41, 41, 48, 0);
}

.dark .card-header {
  background-color: #202325 !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1) !important;
}

.card-body {
  border-bottom: 0.5px solid rgba(255, 255, 255, 0.1) !important;
}

.nav {
  border-bottom: 1px solid #ebc553;
}

.close-btn {
  background: none;
}

a.nav-link.active {
  border-bottom: 3px solid #ebc553;
}
</style>
