<template>
  <div class="hello">
    <nav
      class="navbar navbar-expand-lg navbar-light bg-light justify-content-center"
    >
      <a class="navbar-brand" v-bind:class="{ active: path == '/' }" href="/"
        >Home</a
      >
      <a class="navbar-brand" href="/profile">Profile</a>
      <a class="navbar-brand" href="/appointments">Appointments</a>
      <a class="navbar-brand" href="/drugs">Drug library</a>
    </nav>
  </div>
</template>

<script>
export default {
  name: "NavBar",
  props: ["path"],
};
</script>