<template>
  <div>
    <!-- <div v-if="loggedIn">
      <div class="blocNavbar"><NavBar /></div>
      <div class="blocHeader"><HeaderComp /></div>
      <div class="blocPage"><router-view /></div>
    </div>
    <div v-else>
      <div class="blocPageNotLoggedIn"><router-view /></div>
    </div> -->

    <!-- Pour tester sans avoir à se login :  -->
    <div class="blocNavbar"><NavBar /></div>
    <div class="blocHeader"><HeaderComp /></div>
    <div class="blocPage"><router-view /></div>
    <!--  -->
  </div>
</template>

<script>
import { authStore } from "./store/authStore";
import NavBar from "./components/NavBar.vue";
import HeaderComp from "./components/HeaderComp.vue";

export default {
  data() {
    return {
      loggedIn: authStore.getters.loggedIn,
    };
  },
  components: { NavBar, HeaderComp },
};
</script>

<style lang="scss">
:root {
  --bronze-star-color: #967444;
  --silver-star-color: #c6c6c6;
  --gold-star-color: #ffcf00;
  --main-font: "Roboto", sans-serif;
  font-family: var(--main-font);
  font-weight: 400;
}

.blocNavbar {
  position: fixed;
  top: 10vh;
  left: 0;
  height: 90vh;
  width: 5vw;
  z-index: 1;
}

.blocHeader {
  position: fixed;
  z-index: 10;
  top: 0;
  left: 0;
  height: 10vh;
  width: 100vw;
}

.blocPage {
  position: absolute;
  top: 10vh;
  left: 5vw;
  width: 95vw;
  height: 90vh;
  z-index: 0;
}

.blocPageNotLoggedIn {
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
}
</style>
