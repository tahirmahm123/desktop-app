<template>
  <div>
    <div class="text-secondary mb-3">GENERAL SETTINGS</div>
    <div class="setting-list-item">
      <div class="d-flex align-items-center">
        <div class="">
          <svg
            width="40"
            height="40"
            viewBox="0 0 40 40"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M18.3334 18.3334V5.00006H8.33339C6.50006 5.00006 5.00006 6.50006 5.00006 8.33339V18.3334H18.3334ZM21.6667 18.3334H35.0001V8.33339C35.0001 6.50006 33.5001 5.00006 31.6667 5.00006H21.6667V18.3334ZM18.3334 21.6667H5.00006V31.6667C5.00006 33.5001 6.50006 35.0001 8.33339 35.0001H18.3334V21.6667ZM21.6667 21.6667V35.0001H31.6667C33.5001 35.0001 35.0001 33.5001 35.0001 31.6667V21.6667H21.6667Z"
              fill="white"
            />
          </svg>
        </div>
        <div class="ms-3">
          <h6 class="mb-0">Lunch at Login</h6>
          <p class="mb-0 text-secondary">
            Auto-start the Golden Guard VPN every time your turn on PC
          </p>
        </div>
      </div>

      <div class="param form-check form-switch">
        <input
          type="checkbox"
          class="form-check-input me-1"
          id="launchAtLogin"
          v-model="isLaunchAtLogin"
          :disabled="isLaunchAtLogin == null"
        />
        <!-- <label class="defColor" for="launchAtLogin">Launch at login</label> -->
      </div>
    </div>
    <div class="setting-list-item">
      <div class="d-flex align-items-center">
        <div class="">
          <svg
            width="40"
            height="40"
            viewBox="0 0 40 40"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M10.0684 23.3734C10.0684 24.516 10.8468 25.2693 12.0396 25.2693H18.7444V31.2208C18.7444 32.966 19.6484 34.7615 20 34.7615C20.339 34.7615 21.243 32.966 21.243 31.2208V25.2693H27.9604C29.1532 25.2693 29.9191 24.516 29.9191 23.3734C29.9191 20.6614 27.772 17.8865 24.2188 16.5807L23.817 10.7297C25.6878 9.67503 27.1694 8.50734 27.8097 7.6661C28.1613 7.22665 28.3245 6.76208 28.3245 6.3603C28.3245 5.53162 27.6967 4.91638 26.7425 4.91638H13.245C12.3033 4.91638 11.6629 5.53162 11.6629 6.3603C11.6629 6.76208 11.8262 7.22665 12.1777 7.6661C12.8181 8.50734 14.2997 9.67503 16.1705 10.7297L15.7687 16.5807C12.2154 17.8865 10.0684 20.6614 10.0684 23.3734Z"
              fill="white"
            />
          </svg>
        </div>
        <div class="ms-3">
          <h6 class="mb-0">Show in Dock</h6>
          <p class="mb-0 text-secondary">
            Show application icon in system dock
          </p>
        </div>
      </div>
      <div class="param form-check form-switch">
        <input
          type="checkbox"
          class="form-check-input me-1"
          id="showAppInSystemDock"
          v-model="showAppInSystemDock"
          :disabled="minimizeToTray !== true"
        />
        <!-- <label class="defColor" for="showAppInSystemDock"
          ></label
        > -->
      </div>
    </div>
    <div class="setting-list-item">
      <div class="d-flex align-items-center">
        <div class="">
          <svg
            width="40"
            height="41"
            viewBox="0 0 40 41"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M20 20.7545C20.3139 20.7545 20.5901 20.6541 20.904 20.3527L25.1102 16.2972C25.3237 16.0838 25.4492 15.8578 25.4492 15.5439C25.4492 14.9538 24.9721 14.5394 24.3946 14.5394C24.1058 14.5394 23.817 14.6524 23.6161 14.8784L21.8583 16.7492L21.0798 17.6784L21.1803 15.8075V5.13511C21.1803 4.50732 20.6529 3.97998 20 3.97998C19.3597 3.97998 18.8323 4.50732 18.8323 5.13511V15.8075L18.9202 17.6784L18.1418 16.7492L16.3965 14.8784C16.1956 14.6524 15.8943 14.5394 15.6055 14.5394C15.0279 14.5394 14.5634 14.9538 14.5634 15.5439C14.5634 15.8578 14.6764 16.0838 14.8898 16.2972L19.1086 20.3527C19.4099 20.6541 19.6987 20.7545 20 20.7545ZM9.25226 33.1346H30.7352C33.4598 33.1346 34.8787 31.7283 34.8787 29.0414V21.8469C34.8787 20.5788 34.7154 19.951 34.1379 19.1599L29.9693 13.5224C28.5129 11.5386 27.4958 10.5843 25.1981 10.5843H23.5031V12.6435H25.3237C26.2026 12.6435 27.0815 13.3717 27.634 14.1125L31.9029 20.064C32.3298 20.629 32.1792 20.8675 31.5388 20.8675H24.043C23.3148 20.8675 22.9883 21.3698 22.9883 21.9097V21.9599C22.9883 23.4289 21.8457 25.011 20 25.011C18.1543 25.011 17.0117 23.4289 17.0117 21.9599V21.9097C17.0117 21.3698 16.6853 20.8675 15.9571 20.8675H8.46124C7.80834 20.8675 7.68278 20.5913 8.08457 20.064L12.3535 14.1125C12.9186 13.3717 13.7975 12.6435 14.6764 12.6435H16.497V10.5843H14.7894C12.4917 10.5843 11.4872 11.5386 10.0307 13.5224L5.84964 19.1599C5.27207 19.951 5.1214 20.5788 5.1214 21.8469V29.0414C5.1214 31.7283 6.5402 33.1346 9.25226 33.1346Z"
              fill="white"
            />
          </svg>
        </div>
        <div class="ms-3">
          <h6 class="mb-0">Minimize to tray</h6>
          <p v-if="isLinux" class="mb-0 text-secondary">
            By enabling this parameter, the application will stay in memory
            after closing the window and it will be accessible only via the tray
            icon.
            <b>Caution:</b> Not all Linux desktop environments support
            displaying the application icon in the system tray.
          </p>
          <p class="mb-0 text-secondary">
            Minimize application to system tray/Dock
          </p>
        </div>
      </div>
      <div class="param form-check form-switch">
        <input
          type="checkbox"
          class="form-check-input me-1"
          id="minimizeToTray"
          v-model="minimizeToTray"
        />
        <!-- <label class="defColor" for="minimizeToTray">Minimize to tray</label> -->
      </div>
    </div>

    <!-- <div class="settingsBoldFont">View:</div>
    <div class="flexRow paramBlock">
      <div class="defColor paramName">Color theme:</div>
      <select
        v-model="colorTheme"
        class="form-select form-select-sm"
        style="margin-left: 30px"
      >
        <option :value="colorThemeEnum.system" :key="colorThemeEnum.system">
          System default
        </option>
        <option :value="colorThemeEnum.light" :key="colorThemeEnum.light">
          Light
        </option>
        <option :value="colorThemeEnum.dark" :key="colorThemeEnum.dark">
          Dark
        </option>
      </select>
    </div> -->

    <!-- <div v-if="isLinux && colorScheme === colorThemeEnum.system">
      <div class="description" style="margin-left: 0px">
        When changing the system color theme, the new application color theme
        will be updated after reopening the application window.
      </div>
    </div> -->
    <div class="setting-list-item">
      <div class="d-flex align-items-center">
        <div class="">
          <svg
            width="40"
            height="41"
            viewBox="0 0 40 41"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M9.23962 24.1478H10.1562C10.0558 23.3568 9.99297 22.503 10.1185 21.7622H9.39029C6.40201 21.7622 4.28008 19.6654 4.28008 16.652C4.28008 13.6637 6.40201 11.5417 9.39029 11.5417H20.4771C23.4653 11.5417 25.5873 13.6637 25.5873 16.652C25.5873 19.6654 23.4653 21.7622 20.4771 21.7622H16.635C16.0323 22.164 16.0825 23.6707 16.861 24.1478H20.6277C24.9846 24.1478 28.0984 21.0716 28.0984 16.652C28.0984 12.2449 24.9846 9.1687 20.6277 9.1687H9.23962C4.88276 9.1687 1.76892 12.2449 1.76892 16.652C1.76892 21.0716 4.88276 24.1478 9.23962 24.1478ZM19.3596 30.0867H30.7352C35.1046 30.0867 38.2184 27.0105 38.2184 22.6034C38.2184 18.1838 35.1046 15.1076 30.7352 15.1076H29.8311C29.9316 15.8986 29.9944 16.7524 29.8688 17.4932H30.597C33.5853 17.4932 35.7073 19.59 35.7073 22.6034C35.7073 25.5917 33.5853 27.7136 30.597 27.7136H19.5103C16.522 27.7136 14.4001 25.5917 14.4001 22.6034C14.4001 19.59 16.522 17.4932 19.5103 17.4932H23.3523C23.955 17.0914 23.9048 15.5847 23.1263 15.1076H19.3596C14.9902 15.1076 11.8889 18.1838 11.8889 22.6034C11.8889 27.0105 14.9902 30.0867 19.3596 30.0867Z"
              fill="white"
            />
          </svg>
        </div>
        <div class="ms-3">
          <h6 class="mb-0">Auto Connect</h6>
          <p class="mb-0 text-secondary">
            Auto-connect VPN when application Launch.
          </p>
        </div>
      </div>
      <div class="param form-check form-switch">
        <input
          type="checkbox"
          id="connectOnLaunch"
          class="form-check-input me-1"
          v-model="autoConnectOnLaunch"
        />
        <!-- <label class="defColor" for="connectOnLaunch">On launch</label> -->
      </div>
    </div>
    <!-- <div class="param d-none" v-if="!isLinux">
      <input
        type="checkbox"
        class="form-check-input me-1"
        id="connectVPNOnInsecureNetwork"
        v-model="connectVPNOnInsecureNetwork"
      />
      <label class="defColor" for="connectVPNOnInsecureNetwork"
        >On joining WiFi networks without encryption</label
      >
    </div> -->
    <!--
    <div class="settingsBoldFont">On exit:</div>
    <div class="param">
      <input
        type="checkbox"
        class="form-check-input me-1"
        id="quitWithoutConfirmation"
        v-model="quitWithoutConfirmation"
      />
      <label class="defColor" for="quitWithoutConfirmation"
        >Quit without confirmation when closing application</label
      >
    </div>
    <div class="param">
      <input
        type="checkbox"
        id="disconnect"
        class="form-check-input me-1"
        v-model="disconnectOnQuit"
        :disabled="quitWithoutConfirmation === false"
      />
      <label class="defColor" for="disconnect"
        >Disconnect when closing application</label
      >
    </div> -->

    <!-- DIAGNOSTIC LOGS-->
    <div class="setting-list-item">
      <div class="d-flex align-items-center">
        <div class="">
          <svg
            width="40"
            height="41"
            viewBox="0 0 40 41"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M13.3078 33.4694H26.6923C29.3918 33.4694 30.7855 32.0506 30.7855 29.3386V18.1011H21.3686C19.7991 18.1011 19.0458 17.3478 19.0458 15.7783V6.26099H13.3078C10.6208 6.26099 9.2146 7.67979 9.2146 10.4044V29.3386C9.2146 32.0506 10.6208 33.4694 13.3078 33.4694ZM21.6072 16.268H30.6599C30.5846 15.6904 30.1702 15.1505 29.5173 14.485L22.6493 7.52912C22.0215 6.88878 21.4565 6.46188 20.8789 6.38654V15.5397C20.8915 16.0294 21.1301 16.268 21.6072 16.268ZM16.3463 29.5394V25.4588C16.3463 24.7431 16.6351 24.3665 17.2252 24.3037V23.136C17.2252 21.2024 18.3929 19.9091 20.1005 19.9091C21.8206 19.9091 22.9883 21.2024 22.9883 23.136V24.3162C23.5659 24.3665 23.8672 24.7431 23.8672 25.4588V29.5394C23.8672 30.343 23.5031 30.7071 22.7623 30.7071H17.4512C16.6979 30.7071 16.3463 30.343 16.3463 29.5394ZM18.4682 24.2911H21.7327V23.023C21.7327 21.8804 21.0798 21.1145 20.1005 21.1145C19.1211 21.1145 18.4682 21.8804 18.4682 23.023V24.2911Z"
              fill="white"
            />
          </svg>
        </div>
        <div class="ms-3">
          <h6 class="mb-0">Allow Logging</h6>
        </div>
      </div>
      <div class="d-flex">
        <div class="param form-check form-switch">
          <input
            type="checkbox"
            id="logging"
            class="form-check-input me-1"
            v-model="logging"
          />
          <!-- <label class="defColor" for="logging">Allow logging</label> -->
        </div>
        <button
          class="btn btn-primary py-0"
          v-on:click="onLogs"
          v-if="logging && isCanSendDiagLogs"
        >
          Diagnostic logs ...
        </button>
      </div>
    </div>
    <div id="diagnosticLogs" v-if="diagnosticLogsShown">
      <ComponentDiagnosticLogs
        :onClose="
          (evtId) => {
            diagnosticLogsShown = false;
          }
        "
      />
    </div>
  </div>
</template>

<script>
import { ColorTheme } from "@/store/types";
import ComponentDiagnosticLogs from "@/components/DiagnosticLogs.vue";
import { Platform, PlatformEnum } from "@/platform/platform";

const sender = window.ipcSender;

// VUE component
export default {
  components: {
    ComponentDiagnosticLogs,
  },
  data: function () {
    return {
      diagnosticLogsShown: false,
      isLaunchAtLoginValue: null,
      colorScheme: null,
    };
  },
  mounted() {
    this.colorScheme = sender.ColorScheme();
    this.doUpdateIsLaunchAtLogin();
  },
  methods: {
    async onLogs() {
      this.diagnosticLogsShown = true;
    },
    async doUpdateIsLaunchAtLogin() {
      try {
        this.isLaunchAtLoginValue = await sender.AutoLaunchIsEnabled();
      } catch (err) {
        console.error("Error obtaining 'LaunchAtLogin' value: ", err);
        this.isLaunchAtLoginValue = null;
      }
    },
  },
  computed: {
    isLinux() {
      return Platform() === PlatformEnum.Linux;
    },
    isLaunchAtLogin: {
      get() {
        return this.isLaunchAtLoginValue;
      },
      set(value) {
        this.isLaunchAtLoginValue = value;
        let theThis = this;
        (async function () {
          try {
            await sender.AutoLaunchSet(theThis.isLaunchAtLoginValue);
          } catch (err) {
            console.error("Error changing 'LaunchAtLogin' value: ", err);
            theThis.isLaunchAtLoginValue = null;
          }
        })();
      },
    },
    autoConnectOnLaunch: {
      get() {
        return this.$store.state.settings.autoConnectOnLaunch;
      },
      set(value) {
        this.$store.dispatch("settings/autoConnectOnLaunch", value);
      },
    },
    connectVPNOnInsecureNetwork: {
      get() {
        return this.$store.state.settings.wifi?.connectVPNOnInsecureNetwork;
      },
      set(value) {
        let wifi = Object.assign({}, this.$store.state.settings.wifi);
        wifi.connectVPNOnInsecureNetwork = value;
        this.$store.dispatch("settings/wifi", wifi);
      },
    },

    minimizeToTray: {
      get() {
        return this.$store.state.settings.minimizeToTray;
      },
      set(value) {
        this.$store.dispatch("settings/minimizeToTray", value);
        if (value !== true)
          this.$store.dispatch("settings/showAppInSystemDock", true);
      },
    },
    connectSelectedMapLocation: {
      get() {
        return this.$store.state.settings.connectSelectedMapLocation;
      },
      set(value) {
        this.$store.dispatch("settings/connectSelectedMapLocation", value);
      },
    },
    showAppInSystemDock: {
      get() {
        return this.$store.state.settings.showAppInSystemDock;
      },
      set(value) {
        this.$store.dispatch("settings/showAppInSystemDock", value);
      },
    },
    disconnectOnQuit: {
      get() {
        return this.$store.state.settings.disconnectOnQuit;
      },
      set(value) {
        this.$store.dispatch("settings/disconnectOnQuit", value);
      },
    },
    quitWithoutConfirmation: {
      get() {
        return this.$store.state.settings.quitWithoutConfirmation;
      },
      set(value) {
        this.$store.dispatch("settings/quitWithoutConfirmation", value);
      },
    },
    logging: {
      get() {
        return this.$store.state.settings.logging;
      },
      set(value) {
        this.$store.dispatch("settings/logging", value);
        sender.SetLogging();
      },
    },

    isCanSendDiagLogs() {
      return sender.IsAbleToSendDiagnosticReport();
    },

    colorThemeEnum() {
      return ColorTheme;
    },
    colorTheme: {
      get() {
        return this.colorScheme;
      },
      set(value) {
        sender.ColorSchemeSet(value);
        this.colorScheme = value;
      },
    },
  },
};
</script>

<style scoped lang="scss">
@import "@/components/scss/constants";

.defColor {
  @extend .settingsDefaultTextColor;
}

div.param {
  @extend .flexRow;
  margin-top: 3px;
}

input:disabled {
  opacity: 0.5;
}

input:disabled + label {
  opacity: 0.5;
}

label {
  margin-left: 1px;
}

div.description {
  @extend .settingsGrayLongDescriptionFont;
  margin-left: 22px;
  max-width: 490px;
}

#diagnosticLogs {
  background: white;
  z-index: 99;
  position: absolute;
  left: 0%;
  top: 0%;
  width: 100%;
  height: 100%;
}

select {
  background: linear-gradient(180deg, #ffffff 0%, #ffffff 100%);
  border: 0.5px solid rgba(0, 0, 0, 0.2);
  border-radius: 3.5px;
  width: 186px;
}
</style>
