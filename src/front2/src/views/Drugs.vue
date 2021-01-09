<template>
  <div class="drugs">
    <div>
      <div class="drug">
        <table class="tg" id="customers">
          <thead>
            <tr>
              <th class="tg-0lax">Nazwa</th>
              <th class="tg-0lax">Rodzaj preparatu</th>
              <th class="tg-0lax">Nazwa powszechna</th>
              <th class="tg-0lax">Moc</th>
              <th class="tg-0lax">Typ preparatu</th>
            </tr>
          </thead>
          <tbody>
            <!-- <tr v-bind:key="drug.id" v-for="drug in drugs">
              <td class="tg-0lax">{{ drug.name }}</td>
              <td class="tg-0lax">{{ drug.type_of_preparation }}</td>
              <td class="tg-0lax">{{ drug.common_name }}</td>
              <td class="tg-0lax">{{ drug.strength }}</td>
              <td class="tg-0lax">{{ drug.shape }}</td>
            </tr> -->
            <router-link v-bind:key="drug.id" v-for="drug in drugs" tag="tr"  :to="{name: 'drug', params: {id: drug.id}}">
              <td class="tg-0lax"><b>{{ drug.name }}</b></td>
              <td class="tg-0lax">{{ drug.type_of_preparation }}</td>
              <td class="tg-0lax">{{ drug.common_name }}</td>
              <td class="tg-0lax">{{ drug.strength }}</td>
              <td class="tg-0lax">{{ drug.shape }}</td>
            </router-link>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Drugs",
  data() {
    return {
      drugs: [],
    };
  },
  mounted() {
    axios
      .get(
        "http://localhost:8081/api/v1/Assistant/GetDrugs?prefix=&offset=0&limit=100"
      )
      .then((response) => (this.drugs = response.data.drugs));
  },
};
</script>

<style>
.tg {
  border-collapse: collapse;
  border-spacing: 0;
}
.tg td {
  border-color: black;
  border-style: solid;
  border-width: 1px;
  font-family: Arial, sans-serif;
  font-size: 20px;
  overflow: hidden;
  padding: 10px 5px;
  word-break: normal;
}
.tg th {
  border-color: black;
  border-style: solid;
  border-width: 1px;
  font-family: Arial, sans-serif;
  font-size: 25px;
  font-weight: normal;
  overflow: hidden;
  padding: 10px 5px;
  word-break: normal;
}
.tg .tg-0lax {
  text-align: left;
  vertical-align: top;
}

#customers {
  font-family: Arial, Helvetica, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

#customers td,
#customers th {
  border: 1px solid #ddd;
  padding: 8px;
}

#customers tr:nth-child(even) {
  background-color: #f2f2f2;
}

#customers tr:hover {
  background-color: #ddd;
}

#customers th {
  padding-top: 12px;
  padding-bottom: 12px;
  text-align: left;
  background-color: #42b983;
  color: white;
}
</style>