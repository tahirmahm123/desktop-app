<template>
  <button class="w-100 btn-server-header">
    <div class="d-flex align-items-center justify-content-between">
      <div class="d-flex align-items-center">
        <img :src="serverImage" alt="" style="height: 40px; width: 40px" />
        <span class="text-start ms-2">
          <div class="country-name">{{ serverCountry }}</div>
          <div class="server-count">{{ serversCount }} Locations</div>
        </span>
      </div>
      <div v-if="isFastestServerConfig" class="flexRow">
        <!-- CONFIG -->
        <SwitchProgress
          :isChecked="!isSvrExcludedFomFastest(server)"
          :onChecked="
            (value, event) => {
              configFastestSvrClicked(server, event);
            }
          "
        />
      </div>

      <img v-else class="chevron-server-list" :src="chevron" alt="chevron" />
    </div>
  </button>
</template>

<script>
import ArrowUp from "@/assets/img/arrow-up.svg";
import ArrowRight from "@/assets/img/arrow_right.svg";
import SwitchProgress from "@/components/controls/control-switch-small.vue";

export default {
  components: { SwitchProgress },
  props: [
    "server",
    "expanded",
    "configFastestSvrClicked",
    "isFastestServerConfig",
    "isSvrExcludedFomFastest",
  ],
  data: () => ({
    isImgLoadError: false,
  }),
  computed: {
    isShowSingleLine: function () {
      return (
        this.isSingleLine ||
        this.isFullName ||
        this.isFastestServer ||
        this.isRandomServer
      );
    },

    serverImage: function () {
      if (!this.server) return `/flags/unk.svg`;
      try {
        const ccode = this.server.flag.toLowerCase();
        return `/flags/svg/${ccode}.svg`;
      } catch (e) {
        console.log(e);
        return null;
      }
    },
    serversCount: function () {
      return this.server.servers.length;
    },
    serverCountry: function () {
      return this.server.country;
    },
    chevron: function () {
      return this.expanded ? ArrowUp : ArrowRight;
    },
  },

  methods: {
    onImgLoadError() {
      this.isImgLoadError = true;
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

.country-name {
  font-family: "Outfit";
  font-style: normal;
  font-weight: 700;
  font-size: 16px;
  line-height: 18px;
  margin-bottom: 3px;
}

.server-count {
  font-family: "Outfit";
  font-style: normal;
  font-weight: 300;
  font-size: 13px;
  line-height: 15px;
  opacity: 0.7;
}

.btn {
  --bs-btn-padding-y: 10px;
  --bs-btn-padding-x: 10px;
}

.main-server .btn-secondary {
  --bs-btn-color: #000000;
  --bs-btn-bg: rgba(246, 246, 246, 1);
  --bs-btn-border-color: rgba(246, 246, 246, 1);
  --bs-btn-hover-color: #000000;
  --bs-btn-hover-bg: #eeeeee;
  --bs-btn-hover-border-color: #eeeeee;
  --bs-btn-focus-shadow-rgb: 130, 138, 145;
  --bs-btn-active-color: #000000;
  --bs-btn-active-bg: #e3e2e2;
  --bs-btn-active-border-color: #e3e2e2;
  --bs-btn-active-shadow: none;
  --bs-btn-disabled-color: #fff;
  --bs-btn-disabled-bg: #6c757d;
  --bs-btn-disabled-border-color: #6c757d;
}

.dark .main-server .btn-secondary {
  --bs-btn-color: #fff;
  --bs-btn-bg: transparent;
  --bs-btn-border-color: rgba(41, 41, 48, 1);
  --bs-btn-hover-color: #fff;
  --bs-btn-hover-bg: #313131;
  --bs-btn-hover-border-color: #313131;
  --bs-btn-focus-shadow-rgb: 130, 138, 145;
  --bs-btn-active-color: #fff;
  --bs-btn-active-bg: #3d3d3d;
  --bs-btn-active-border-color: #3d3d3d;
  --bs-btn-active-shadow: none;
  --bs-btn-disabled-color: #fff;
  --bs-btn-disabled-bg: #6c757d;
  --bs-btn-disabled-border-color: #6c757d;
}

.chevron-server-list {
  // transform: rotate(90deg);
  transition: 0.2s;
}

.not-collapsed .chevron-server-list {
  transform: rotate(180deg);
  transition: 0.2s;
}

.dark .accordion-button:not(.collapsed)::after {
  background-image: v-bind(ArrowRight);
}
</style>
