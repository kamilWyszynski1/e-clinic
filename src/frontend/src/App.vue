<template>
  <div id="app">
    <router-view>
      <router-link to="/"></router-link>
      <NavBar v-bind:path="path" />
      <div v-if="isPatient">
        <Patient
          v-bind:path="path"
          v-bind:id="'147fc868-efa7-4f37-b85c-6a3f05e78b44'"
        />
      </div>
      <div v-else>
        <Specialist v-bind:path="path" />
      </div>
    </router-view>
  </div>
</template>


<script>
import NavBar from "./components/layout/NavBar.vue";
import Patient from "./components/patient/Patient.vue";
import Specialist from "./components/specialist/Specialist.vue";

export default {
  name: "App",
  components: {
    NavBar,
    Patient,
    Specialist,
  },

  data() {
    return {
      path: window.location.pathname,
      isPatient: true,
    };
  },
};
</script>

<style>
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
body {
  font-family: Arial, Helvetica, sans-serif;
  line-height: 1.4;
}
</style>
