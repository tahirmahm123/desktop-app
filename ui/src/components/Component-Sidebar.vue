<template>
  <div class="sidebar">
    <div class="text-center mt-4">
      <img src="@/assets/img/main-logo.svg" alt="" class="logo-main" />
      <!-- <h5 class="mt-1">{{ AppName }}</h5> -->
    </div>
    <div class="btn-side-menu">
      <router-link to="/">
        <div :class="['side-menu-item', active === 'connect' ? 'active' : '']">
          <a href="#">
            <span class="">Connect</span>
          </a>
        </div>
      </router-link>
      <!-- <router-link to="/servers">
        <div :class="['side-menu-item', active === 'location' ? 'active' : '']">
          <a href="#">
            <span class="ms-1">Locations</span>
          </a>
        </div>
      </router-link> -->
      <router-link to="/settings/general">
        <div :class="['side-menu-item', active === 'settings' ? 'active' : '']">
          <a href="#">
            <span class="">Settings</span>
          </a>
        </div>
      </router-link>
      <router-link to="/get-pin">
        <div :class="['side-menu-item', active === 'settings' ? 'active' : '']">
          <a href="#">
            <span class="">Get Pin</span>
          </a>
        </div>
      </router-link>
      <router-link to="/get-help">
        <div :class="['side-menu-item', active === 'settings' ? 'active' : '']">
          <a href="#">
            <span class="">Get Help</span>
          </a>
        </div>
      </router-link>
    </div>
  </div>
</template>
<script>
import config from "@/config";

export default {
  data: function () {
    return {
      AppName: config.AppName,
    };
  },
  props: ["active"],
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
@import "@/components/scss/constants";

.sidebar {
  padding: 15px 30px;
  display: flex;
  align-items: center;
  position: relative;
}

.sidebar .app-version {
  font-size: 10px;
  position: absolute;
  top: 560px;
  display: inline-flex;
  left: 50%;
  transform: translateX(-50%);
  color: rgb(174, 174, 174);
}

.sidebar .btn-side-menu {
  margin-top: 30px;
  display: flex;
  margin-left: 30px;
}

.sidebar .btn-side-menu .side-menu-item {
  color: #fff;
  padding: 5px 10px;
  border-radius: 50px;
  border: transparent 1px solid;
  margin: 5px;
  line-height: initial;
}
.logo-main {
  width: 50px;
}
.side-menu-item a {
  color: #000;
  font-size: 14px;
}

.dark .side-menu-item a {
  color: #fff;
}

.sidebar .btn-side-menu .side-menu-item.active {
  background-color: #ebc553;
  border: #ebc553 1px solid;
  font-weight: 600;
  cursor: pointer;
}
.sidebar .btn-side-menu .side-menu-item:hover {
  background-color: transparent;
  border: #ebc553 1px solid;
  cursor: pointer;
}
</style>
