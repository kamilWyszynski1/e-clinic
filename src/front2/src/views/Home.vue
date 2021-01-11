<template>
  <div class="home">
  <h1>E-clinic is an application that allows you to create and pay for appointment, manage your account
    but mostly it uses graph algorithms to search for replacements for drugs that interest you.
  </h1>
  <br>
  <h3>
    Application uses Neo4j as graph database and Jaccard index to calculate similarities among drugs.
    <br>
    Drugs are called simillar if they share plenty enough of active substances.
  </h3>
  <img src="https://wikimedia.org/api/rest_v1/media/math/render/svg/eaef5aa86949f49e7dc6b9c8c3dd8b233332c9e7" alt="Jaccard index">
  
  <br><br>

  <h3 id="profile-info">Go to profile page in order to create profile!!!</h3>

  </div>
</template>

<script>

export default {
  data() {
    return {
      email: ""
    }
  },
  name: 'Home',
  methods: {
    setEmailAsCookie() {
     this.$cookies.set("email", this.email, null, null, null, null, "SameSite")
     console.log(this.email)
    }
  }
}
</script>

<style>
#profile-info{
  color: red;
}
</style>
