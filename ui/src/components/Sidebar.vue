<template>
  <div class="flexColumn">

    <div class="sidebar">
      <img src="@/assets/icon-white.svg"  >
      <div class="sidebar-btn-outer">
        <button class="sidebarBtn" v-on:click="routeTo('settings')">
            <img src="@/assets/connect-ico.svg" alt="">
            <p class="sidebar-item-title">Connect</p>
        </button> 
        <button class="sidebarBtn"  v-on:click="routeTo('settings')">
            <img src="@/assets/location-ico.svg" alt="">
            <p class="sidebar-item-title">Location</p>
        </button> 
        <button class="sidebarBtn"  v-ripple v-on:click="routeTo('settings')">
            <img src="@/assets/setting-ico.svg" alt="">
            <p class="sidebar-item-title">Settings</p>
        </button> 
      </div>
    </div>

  </div>
</template>
<script>
import Vue from "vue"
//Register a global custom directive called v-ripple
Vue.directive("ripple", {
  // When the bound element is inserted into the DOM...

  inserted: function(el) {

    // listen for click events to trigger the ripple
    el.addEventListener(
      "click",
      function(e) {
        
        // Setup
    var target = el.getBoundingClientRect();
    var buttonSize = target.width > target.height ? target.width : target.height;
        // remove any previous ripple containers
        var elements = document.getElementsByClassName("ripple");
        while (elements[0]) {
          elements[0].parentNode.removeChild(elements[0]);
        }

        // create the ripple container and append it to the target element
        var ripple = document.createElement("span");
        ripple.setAttribute("class", "ripple");
        el.appendChild(ripple);

        // set the ripple container to the click position and start the animation
        setTimeout(function() {
          ripple.style.width = buttonSize + "px";
          ripple.style.height = buttonSize + "px";
          ripple.style.top = e.offsetY - buttonSize / 2 + "px";
          ripple.style.left = e.offsetX - buttonSize / 2 + "px";
          ripple.setAttribute("class", "ripple ripple-effect");
        }, 100);
      },
      false
    );
  }
});


export default {
  methods: {
    routeTo(value) {
      console.log("route clicked");
      this.$router.push(value);
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
@import "@/components/scss/constants";

.sidebar{
  width: 80px;
  background-image: linear-gradient(#FD2411, #A80004);
  height: 100vh;
  margin: 0 20px;
  border-radius: 10px;
  text-align: center;
  padding: 20px 5px;
}
.sidebar-btn-outer{
  margin-top: 30px;
}
.sidebarBtn{
    background: transparent;
    margin: 5px 0;
    width: 100%;
    border: 2px solid transparent;
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
  border: 2px solid #fff;
}
.sidebarBtn .active{
  background: #6F0305;
}
.sidebarBtn:active {
  width: 0;
  height: 0;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.9);
  transform: scale(0);
  position: absolute;
  opacity: 1;
}

</style>
