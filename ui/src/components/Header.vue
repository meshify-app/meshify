<template>
    <v-container>
        <v-app-bar app>
            <img class="mr-3" :src="require('../assets/meshify.png')" height="50" alt="meshify"/>
            <v-toolbar-title to="/">meshify.app</v-toolbar-title>

            <v-spacer />

            <v-toolbar-items>
                <v-btn to="/mesh" right>
                    Meshes
                    <img class="ml-1" :src="require('../assets/meshify-bw.png')" height="32" width="32" alt="meshify"/>
                </v-btn>
                <v-btn to="/hosts">
                    Hosts
                    <v-icon right dark>mdi-account-network-outline</v-icon>
                </v-btn>
                <v-btn to="/users">
                    Users
                    <v-icon right dark>mdi-account-multiple</v-icon>
                </v-btn>
<!--                <v-btn to="/server">
                    Server
                    <v-icon right dark>mdi-vpn</v-icon>
                </v-btn> -->
            </v-toolbar-items>

            <v-menu
                    left
                    bottom
            >

                <template v-slot:activator="{ on }">
                    <v-btn icon v-on="on">
                        <v-avatar size="36">
                            <img alt="user.name" :src="user.picture"/>
                        </v-avatar>
                    </v-btn>
                </template>

                <v-card
                        class="mx-auto"
                        max-width="344"
                        outlined
                >
                    <v-list-item three-line>
                        <v-list-item-content>
                            <div class="overline mb-4">connected as</div>
                            <v-list-item-title class="headline mb-1">{{user.name}}</v-list-item-title>
                            <v-list-item-subtitle>Email: {{user.email}}</v-list-item-subtitle>
                            <v-list-item-subtitle>Issuer: {{user.issuer}}</v-list-item-subtitle>
                            <v-list-item-subtitle>Issued at: {{ user.issuedAt | formatDate }}</v-list-item-subtitle>
                            <v-list-item-subtitle>Plan: {{user.plan}}</v-list-item-subtitle>
                        </v-list-item-content>
                    </v-list-item>
                    <v-card-actions>
                        <v-btn small
                                v-on:click="logout()"
                        >
                            logout
                            <v-icon small right dark>mdi-logout</v-icon>
                        </v-btn>
                    </v-card-actions>
                </v-card>
            </v-menu>

        </v-app-bar>
    </v-container>
</template>

<script>
  import {mapActions, mapGetters} from "vuex";

  export default {
    name: 'Header',

    computed:{
      ...mapGetters({
        user: 'auth/user',
      }),
    },

    methods: {
      ...mapActions('auth', {
        logout: 'logout',
      }),
    }
  }
</script>
