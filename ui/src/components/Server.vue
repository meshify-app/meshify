<template>
    <v-container v-if="server">
        <v-row>
            <v-col cols="12">
                <v-card dark>
                    <v-card-title>
                        Mesh Configuration
                    </v-card-title>
                    <div class="d-flex flex-no-wrap justify-space-between">
                        <v-col cols="12">
                            <v-text-field
                                    v-model="server.id"
                                    label="Mesh ID"
                                    disabled
                            />
                            <v-text-field
                                    v-model="server.meshName"
                                    label="Mesh Name"
                                    disabled
                            />
                            <v-combobox
                                    v-model="server.address"
                                    chips
                                    hint="Write IPv4 or IPv6 CIDR and hit enter.  A 100.x.x.0/24 address is recommended."
                                    label="Host interface address pool"
                                    multiple
                                    dark
                            >
                                <template v-slot:selection="{ attrs, item, select, selected }">
                                    <v-chip
                                            v-bind="attrs"
                                            :input-value="selected"
                                            close
                                            @click="select"
                                            @click:close="server.address.splice(server.address.indexOf(item), 1)"
                                    >
                                        <strong>{{ item }}</strong>&nbsp;
                                    </v-chip>
                                </template>
                            </v-combobox>
                        </v-col>
                    </div>
                    <div class="d-flex flex-no-wrap justify-space-between">
                        <v-col cols="12">
                            <v-text-field
                                    v-model="server.listenPort"
                                    type="number"
                                    :rules="[
                          v => !!v || 'Listen port is required',
                        ]"
                                    label="Default Listen port"
                                    required
                            />
                            <v-combobox
                                    v-model="server.dns"
                                    chips
                                    hint="Write IPv4 or IPv6 address and hit enter"
                                    label="DNS servers for clients"
                                    multiple
                                    dark
                            >
                                <template v-slot:selection="{ attrs, item, select, selected }">
                                    <v-chip
                                            v-bind="attrs"
                                            :input-value="selected"
                                            close
                                            @click="select"
                                            @click:close="server.dns.splice(server.dns.indexOf(item), 1)"
                                    >
                                        <strong>{{ item }}</strong>&nbsp;
                                    </v-chip>
                                </template>
                            </v-combobox>
                            <v-combobox
                                    v-model="server.allowedips"
                                    chips
                                    hint="Write IPv4 or IPv6 address and hit enter"
                                    label="Default Allowed IPs for clients"
                                    multiple
                                    dark
                            >
                                <template v-slot:selection="{ attrs, item, select, selected }">
                                    <v-chip
                                            v-bind="attrs"
                                            :input-value="selected"
                                            close
                                            @click="select"
                                            @click:close="server.allowedips.splice(server.allowedips.indexOf(item), 1)"
                                    >
                                        <strong>{{ item }}</strong>&nbsp;
                                    </v-chip>
                                </template>
                            </v-combobox>
                            <v-text-field
                                    type="number"
                                    v-model="server.mtu"
                                    label="Define global MTU"
                                    hint="Leave at 0 and let wg-quick take care of MTU"
                            />
                            <v-text-field
                                    type="number"
                                    v-model="server.persistentKeepalive"
                                    label="Persistent keepalive"
                                    hint="Leave at 0 if you dont want to specify persistent keepalive"
                            />
                        </v-col>
                    </div>
                </v-card>
            </v-col>
        </v-row>
            <v-divider dark/>
            <v-btn
                    class="ma-2"
                    color="success"
                    v-on:click="forceFileDownload()"
            >
                Download server configuration
                <v-icon right dark>mdi-cloud-download-outline</v-icon>
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn
                    class="ma-2"
                    color="warning"
                    @click="update"
            >
                Update server configuration
                <v-icon right dark>mdi-update</v-icon>
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn
                    color="warning"
                    @click="applyGlobals"
            >
                Apply globals to all clients
                <v-icon right dark>mdi-update</v-icon>
            </v-btn>

            <v-divider dark/>
    </v-container>
</template>
<script>
  import {mapActions, mapGetters} from "vuex";

  export default {
    name: 'Server',

    data: () => ({

    }),

    computed:{
      ...mapGetters({
        server: 'server/server',
        config: 'server/config',
        clients: 'client/clients',
      }),
    },

    mounted () {
      this.readServer()
    },

    methods: {
      ...mapActions('server', {
        errorServer: 'error',
        readServer: 'read',
        updateServer: 'update',
      }),

      ...mapActions('client', {
        updateClient: 'update',
      }),

      applyGlobals(){
        this.update()

        this.clients.forEach(client => {
          client.allowedips = this.server.allowedips
          this.updateClient(client)
        })
      },

      update() {
        // convert int values
        this.server.listenPort = parseInt(this.server.listenPort, 10);
        this.server.persistentKeepalive = parseInt(this.server.persistentKeepalive, 10);
        this.server.mtu = parseInt(this.server.mtu, 10);

        // check server addresses
        if (this.server.address.length < 1) {
          this.errorServer('Please provide at least one valid CIDR address for server interface');
          return;
        }
        for (let i = 0; i < this.server.address.length; i++){
          if (this.$isCidr(this.server.address[i]) === 0) {
            this.errorServer(`Invalid CIDR detected, please correct ${this.server.address[i]} before submitting`);
            return
          }
        }

        // check DNS correct
        for (let i = 0; i < this.server.dns.length; i++){
          if (this.$isCidr(this.server.dns[i] + '/32') === 0) {
            this.errorServer(`Invalid IP detected, please correct ${this.server.dns[i]} before submitting`);
            return
          }
        }

        // check client AllowedIPs
        if (this.server.allowedips.length < 0) {
          this.errorServer('Please provide at least one valid CIDR address for client allowed IPs');
          return;
        }
        for (let i = 0; i < this.server.allowedips.length; i++){
          if (this.$isCidr(this.server.allowedips[i]) === 0) {
            this.errorServer('Invalid CIDR detected, please correct before submitting');
            return
          }
        }

        this.updateServer(this.server)
      },

      forceFileDownload(){
        const url = window.URL.createObjectURL(new Blob([this.config]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', 'wg0.conf') //or any other extension
        document.body.appendChild(link)
        link.click()
      },
    }
  };
</script>
