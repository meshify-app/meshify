<template>
  <v-content>
    <Hosts/>
  </v-content>
</template>

<script>
  import Hosts from '../components/Hosts'

  export default {
    name: 'hosts',
    components: {
      Hosts
    }
  }
</script>
