<template>
  <div class="text-center">
    <div class="pt-3 pb-5 align-items-center connect-status-bg">
      <img v-bind:src="protectedImage" alt="" />
      <h6
        v-bind:class="[`text-${protectedColor}`, 'mb-0']"
        style="font-size: 16px; font-weight: 700"
      >
        {{ protectedText }}
      </h6>
    </div>
    <div class="connect-action my-3 text-center [`text-${protectedColor}`]">
      <div class="line-btn-bg"></div>
      <SwitchProgress
        v-bind:class="{ lowOpacity: isCanResume }"
        :onChecked="onChecked"
        :isChecked="isChecked"
        :isProgress="isProgress"
      />
    </div>
  </div>
</template>

<script>
import SwitchProgress from "@/components/controls/control-switch.vue";
// import imgPause from "@/components/images/pause-btn.vue";
import { PauseStateEnum } from "@/store/types";
import { GetTimeLeftText } from "@/helpers/renderer";
import vClickOutside from "click-outside-vue3";
import Connected from "@/assets/img/shield-tick.svg";
import Processing from "@/assets/img/shield.svg";
import Disconnected from "@/assets/img/shield-cross.svg";

export default {
  directives: {
    clickOutside: vClickOutside.directive,
  },
  components: {
    SwitchProgress,
    // imgPause,
  },
  props: [
    "onChecked",
    "onPauseResume",
    "pauseState",
    "isChecked",
    "isProgress",
  ],
  data: () => ({
    isPauseMenuAllowed: false,
    isPauseExtendMenuShow: false,
    pauseTimeUpdateTimer: null,
    pauseTimeLeftText: "",
  }),
  mounted() {
    this.startPauseTimer();
  },
  created: function () {
    let self = this;
    window.addEventListener("click", function (e) {
      // close dropdown when clicked outside
      if (!self.$el.contains(e.target)) {
        self.isPauseMenuAllowed = false;
      }
    });
  },
  computed: {
    protectedText: function () {
      if (this.$store.state.vpnState.pauseState === PauseStateEnum.Paused)
        return "Paused";
      if (this.isProgress) return "You are about to secure.";
      if (this.isChecked !== true || this.isCanResume)
        return "You’re not secured";
      return "You’re secured till:";
    },
    protectedImage: function () {
      if (this.isProgress) return Processing;
      if (this.isChecked !== true || this.isCanResume) return Disconnected;
      return Connected;
    },
    protectedColor: function () {
      if (this.$store.state.vpnState.pauseState === PauseStateEnum.Paused)
        return "Paused";
      if (this.isProgress) return "warning";
      if (this.isChecked !== true || this.isCanResume) return "danger";
      return "success";
    },
    isConnected: function () {
      return this.$store.getters["vpnState/isConnected"];
    },
    pauseConnectionTill: function () {
      return this.$store.state.uiState.pauseConnectionTill;
    },
    isPaused: function () {
      if (this.$store.state.vpnState.pauseState !== PauseStateEnum.Paused)
        return false;
      return this.pauseConnectionTill != null;
    },
    isCanPause: function () {
      if (!this.isConnected) return false;
      if (this.isProgress === true) return false;

      var connInfo = this.$store.state.vpnState.connectionInfo;
      if (connInfo === null || connInfo.IsCanPause === false) return false;
      if (this.$store.state.vpnState.pauseState === PauseStateEnum.Resumed)
        return true;
      return false;
    },
    isCanResume: function () {
      if (this.isCanPause) return false;
      if (!this.isConnected) return false;
      if (this.isProgress === true) return false;
      if (this.$store.state.vpnState.pauseState === PauseStateEnum.Paused)
        return true;
      return false;
    },
    isCanShowPauseMenu: function () {
      return this.isCanPause && this.isPauseMenuAllowed;
    },
  },
  watch: {
    isPaused() {
      this.startPauseTimer();
    },
  },
  methods: {
    onPauseMenuClickOutside() {
      this.isPauseExtendMenuShow = false;
      this.isPauseMenuAllowed = false;
    },
    onPauseMenu() {
      if (this.isPauseMenuAllowed != true) this.onPauseResume(null);
      this.isPauseMenuAllowed = !this.isPauseMenuAllowed;
    },
    onAddPauseTimeMenu() {
      if (this.isCanResume != true) this.isPauseExtendMenuShow = false;
      else this.isPauseExtendMenuShow = !this.isPauseExtendMenuShow;
    },
    onPauseMenuItem(seconds) {
      this.isPauseMenuAllowed = false;
      this.isPauseExtendMenuShow = false;
      if (this.onPauseResume != null) this.onPauseResume(seconds);
    },
    startPauseTimer() {
      if (this.pauseTimeUpdateTimer) return;
      if (!this.$store.state.uiState.pauseConnectionTill) return;

      this.pauseTimeUpdateTimer = setInterval(() => {
        this.pauseTimeLeftText = GetTimeLeftText(
          this.$store.state.uiState.pauseConnectionTill,
        );

        if (this.$store.state.vpnState.pauseState !== PauseStateEnum.Paused) {
          clearInterval(this.pauseTimeUpdateTimer);
          this.pauseTimeUpdateTimer = null;
        }
      }, 1000);

      this.pauseTimeLeftText = GetTimeLeftText(
        this.$store.state.uiState.pauseConnectionTill,
      );
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
@import "@/components/scss/constants";
@import "@/components/scss/popup";

$shadow:
  0px 3px 1px rgba(0, 0, 0, 0.06),
  0px 3px 8px rgba(0, 0, 0, var(--shadow-opacity-koef));

.main {
  @extend .left_panel_block;

  align-items: center;
}
.connect-action {
  position: relative;
}
.connect-action::after {
  content: "";
  background: rgba(255, 255, 255, 0.15);
  height: 1px;
  width: 100%;
  z-index: 0;
  position: absolute;
  top: 50%;
  left: 0;
}
.buttons {
  align-items: center;
  min-height: 92px;
}
.lowOpacity {
  opacity: 0.5;
}

.large_text {
  font-style: normal;
  font-weight: 600;
  letter-spacing: -0.3px;
  font-size: 24px;
  line-height: 29px;
}

.small_text {
  font-size: 14px;
  line-height: 17px;
  letter-spacing: -0.3px;
  color: var(--text-color-details);
}

.settingsBtn {
  float: right;
  padding: 8px;
  border: none;
  background-color: #ffffff;
  outline-width: 0;
  cursor: pointer;
  box-shadow: $shadow;

  // centering content
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 5px;
}

.settingsBtn:hover {
  background-color: #f0f0f0;
}

.settingsBtnResume {
  @extend .settingsBtn;
  background-color: #fd2411;
}

.settingsBtnResume:hover {
  background-color: #a80004;
}
</style>
