<template>
  <div class="flexColumn">
    <div class="sidebar">
      <img src="@/assets/icon-white.svg"  >
      <div class="sidebar-btn-outer">
        <router-link to="/">
          <button :class="['sidebarBtn', active == 'connect' ? 'active' : '']">
              <img src="@/assets/connect-ico.svg" alt="">
              <p class="sidebar-item-title">Connect</p>
          </button> 
        </router-link>
        <router-link to="/servers"  >
          <button :class="['sidebarBtn', active == 'location' ? 'active' : '']" >
              <img src="@/assets/location-ico.svg" alt="">
              <p class="sidebar-item-title">Locations</p>
          </button> 
        </router-link>
        <router-link to="/settings/general">
          <button :class="['sidebarBtn', active == 'settings' ? 'active' : '']">
              <img src="@/assets/setting-ico.svg" alt="">
              <p class="sidebar-item-title">Settings</p>
          </button>
        </router-link>
      </div>
    </div>

  </div>
</template>
<script>


export default {
  props: [
    "active"
  ]
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
@import "@/components/scss/constants";

.sidebar{
  width: 80px;
  background-image: linear-gradient(#FD2411, #A80004);
  height: 100vh;
  margin:  20px;
  border-radius: 10px;
  text-align: center;
  padding: 20px 5px;
  height: 536px;
}
.sidebar-btn-outer{
  margin-top: 30px;
}
.sidebarBtn{
    background: transparent;
    margin: 5px 0;
    width: 100%;
    border: 1px solid transparent;
    padding: 10px 5px;
    border-radius: 5px;
    cursor: pointer;
}
.sidebar-item-title{
    margin: 0px;
    color: #fff;
    margin-top: 5px;
    font-weight: 500;
}
.sidebarBtn:hover{
  border: 1px solid #fff;
}
.sidebarBtn.active{
  background: #ffffff55;
}

</style>
