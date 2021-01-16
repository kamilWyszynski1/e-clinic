<template>
  <div>
    <h3>Create profile, you need to do it before creating appointment</h3>
    <form class="form" v-on:submit.prevent="submit">
      <div class="form-group row">
        <label for="symptoms" class="col-4 col-form-label">E-mail</label>
        <div class="col-4">
          <input
            id="symptoms"
            name="symptoms"
            type="text"
            class="form-control"
            required="required"
            v-model="email"
          />
           <span id="symptomsHelpBlock" class="form-text text-muted"
            >Pass your real e-mail. Server will send you email with payment link.</span
          >
        </div>
      </div>

      <div class="form-group row">
        <label for="symptoms" class="col-4 col-form-label">Name</label>
        <div class="col-4">
          <input
            id="symptoms"
            name="symptoms"
            type="text"
            class="form-control"
            required="required"
            v-model="name"
          />
        </div>
      </div>

      <div class="form-group row">
        <label for="symptoms" class="col-4 col-form-label">Surname</label>
        <div class="col-4">
          <input
            id="symptoms"
            name="symptoms"
            type="text"
            class="form-control"
            required="required"
            v-model="surname"
          />
        </div>
      </div>

      <div class="form-group row">
        <label for="select" class="col-4 col-form-label">Gender</label>
        <div class="col-4">
          <select
            id="select"
            name="select"
            class="custom-select"
            required="required"
            v-model="gender"
          >
            <option
              v-bind:value="g"
              v-bind:key="g"
              v-for="g in gedners"
            >
              {{ g }}
            </option>
          </select>
        </div>
      </div>

      <div class="form-group row">
        <div class="offset-4 col-4">
          <button name="submit" class="btn btn-primary">Submit</button>
        </div>
      </div>
    </form>

    Profile uuid: {{this.uuid}}

     <div>
      <b-modal id="modal-1" title="BootstrapVue">
        <p class="my-4">Can't create appointment without account. Go to Profile!</p>
      </b-modal>
    </div>
  </div>
</template>

<script>
import axios from "axios";


export default {
  name: "Profile",
  data() {
    return {
      email: "",
      name: "",
      surname: "",
      gender: "",
      gedners: ["MALE", "FEMALE"],
      uuid: this.$cookies.get("uuid")
    };
  },
  methods: {
    submit() {
      axios.
        post("http://localhost:8081/api/v1/Assistant/CreatePatient", {
          email: this.email,
          name: this.name, 
          surname: this.surname, 
          gender: this.gender
        }).
        then( res => {
          console.log(res)
          this.$cookies.set("uuid", res.data.id, null, null, null, null, "SameSite")
          this.uuid = res.data.id
        })
    },
  },
};
</script>