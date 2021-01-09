<template>
  <div v-if="data" class="drug">
    <h2>Lek</h2>
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
        <tr>
          <td class="tg-0lax">
            <b>{{ data.drug.name }}</b>
          </td>
          <td class="tg-0lax">{{ data.drug.type_of_preparation }}</td>
          <td class="tg-0lax">{{ data.drug.common_name }}</td>
          <td class="tg-0lax">{{ data.drug.strength }}</td>
          <td class="tg-0lax">{{ data.drug.shape }}</td>
        </tr>
      </tbody>
    </table>

    <h2>Substancje czynne</h2>
    <table class="tg" id="customers">
      <thead>
        <tr>
          <th class="tg-0lax">Nazwa</th>
        </tr>
      </thead>
      <tbody>
        <tr v-bind:key="s.id" v-for="s in data.substances">
          <td class="tg-0lax">
            <b>{{ s.name }}</b>
          </td>
        </tr>
      </tbody>
    </table>

    <h2>Zamienniki</h2>
    <table v-if="replacements" class="tg" id="customers">
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
        <router-link
          v-bind:key="drug.id"
          v-for="drug in replacements"
          tag="tr"
          :to="{ name: 'drug', params: { id: drug.id } }"
        >
          <td class="tg-0lax">
            <b>{{ drug.name }}</b>
          </td>
          <td class="tg-0lax">{{ drug.type_of_preparation }}</td>
          <td class="tg-0lax">{{ drug.common_name }}</td>
          <td class="tg-0lax">{{ drug.strength }}</td>
          <td class="tg-0lax">{{ drug.shape }}</td>
        </router-link>
      </tbody>
    </table>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Drug",
  data() {
    return {
      data: null,
      replacements: null,
    };
  },
  methods: {
    getContent(uid) {
      axios
        .get(`http://localhost:8081/api/v1/Assistant/GetDrug?drugID=${uid}`)
        .then((response) => {
          console.log(response.data);
          this.data = response.data;
        });

      axios
        .get(
          `http://localhost:8081/api/v1/Assistant/GetReplacement?drugID=${uid}&minSimilarity=0.5`
        )
        .then((response) => {
          console.log(response.data);
          this.replacements = response.data.drugs;
        });
    },
  },
  created() {
      this.getContent(this.$route.params.id);
  },
  beforeRouteUpdate(to, from, next) {
    console.log(to, from, next);
    this.getContent(to.params.id);
    next();
  },
};
</script>