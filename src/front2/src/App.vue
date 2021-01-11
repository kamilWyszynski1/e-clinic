<template>
  <div id="app">
    <div id="nav">
      <h1>E-clinic bitches</h1>
      <router-link to="/">Home</router-link> |
      <router-link to="/profile">Profile</router-link> |
      <router-link to="/appointments">Appointments</router-link> |
      <router-link to="/drugs">Drug library</router-link> |
    </div>
    <router-view />
    <div>
      <b-modal id="error" title="BootstrapVue">
        <p class="my-4">
          Something went wrong
        </p>
      </b-modal>
    </div>
  </div>
</template>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
