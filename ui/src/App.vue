<template>
  <div id="app" class="dark">
    <!--
    For no-bordered windows: print border manually -
    show transparent but bordered div top of the window
    Div should be 'transparent' for mouse events -->
    <div v-if="!isWindowHasFrame" class="border"></div>

    <!-- ability to move by mouse when no title for window -->
    <div class="title">
      <CustomTitleBar />
    </div>

    <router-view />
  </div>
</template>

<script>
import { IsWindowHasFrame } from "@/platform/platform";

const sender = window.ipcSender;
import { InitDefaultCopyMenus } from "@/context-menu/renderer";
import CustomTitleBar from "@/views/CustomTitleBar.vue";

export default {
  components: {
    CustomTitleBar,
  },
  beforeCreate() {
    // function using to re-apply all mutations
    // This is required to send to renderer processes current storage state
    sender.RefreshStorage();
  },
  mounted() {},
  data: function () {
    return {};
  },
  computed: {
    isWindowHasFrame: function () {
      return IsWindowHasFrame();
    },
    isLoggedIn: function () {
      return this.$store.getters["account/isLoggedIn"];
    },
  },
  watch: {
    isLoggedIn() {
      if (this.isLoggedIn === false) this.$router.push("/");
    },
  },
};

InitDefaultCopyMenus();
</script>

<style lang="scss">
@import "@/components/scss/constants";

html * {
  // disable elements\text selection
  -webkit-user-select: none;

  // assign default properties globally for all elements
}

@font-face {
  font-family: "Outfit";
  src:
    local("Outfit"),
    url("~@/assets/fonts/outfit/Outfit.woff") format("woff");
}

body {
  font-family: "Outfit", serif;
}

#app {
  position: absolute;
  left: 0;
  top: 0;
  width: 100vw;
  height: 100vh;

  // disable scrollbars
  overflow-y: hidden;
  overflow-x: hidden;
}

.title {
  // Panel can be draggable by mouse
  // (we need this because using no title style for main window (for  macOS))
  -webkit-app-region: drag;
  height: 24px;
  width: 100%;

  position: absolute;
}

.border {
  // For no-bordered windows: print border manually -
  // show transparent but bordered div top of the window
  // Div should be 'transparent' for mouse events

  // top of all other elements
  z-index: 100;
  // full window
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  // allow elements located under this div to receive mouse events
  pointer-events: none;
  // border style
  //border: 1px solid rgba(128, 128, 128, 0.5);
}

::-webkit-scrollbar {
  background: transparent;
  width: 8px;
  margin-left: 10pc;
}

::-webkit-scrollbar-thumb {
  background: #d2d2d2;
  border-radius: 10px;
}

.dark ::-webkit-scrollbar-thumb {
  background: #232323;
  border-radius: 10px;
}
</style>
