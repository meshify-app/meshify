<template>
    <v-container>
        <v-card>
            <v-card-title>
                Join
            </v-card-title>
            <p>Joining {{ $route.params.id }}</p>
            <p>{{ this.result }}</p>
        </v-card>
    </v-container>
</template>
<script>
  import { mapActions, mapGetters } from 'vuex'

  export default {
    name: 'Join',

    data: () => ({
      result : "",
      panel: 1,
      valid: false,
      search: '',
    }),

    computed:{
      ...mapGetters({
        user: 'auth/user',
        accounts: 'account/accounts',

      }),
    },

    mounted () {
        this.result = this.activate(this.$route.query.id)
    },

    methods: {
      ...mapActions('join', {
        activate: 'activate',
      }),

    }
  };
</script>
