<template>
    <v-container>
        <v-card>
            <v-card-title>
                Hosts
                <v-spacer></v-spacer>
                <v-text-field
                        v-if="listView"
                        v-model="search"
                        append-icon="mdi-magnify"
                        label="Search"
                        single-line
                        hide-details
                ></v-text-field>
                <v-spacer></v-spacer>
                <v-btn
                        color="success"
                        @click="startCreate"
                >
                    Add host manually
                    <v-icon right dark>mdi-account-network-outline</v-icon>
                </v-btn>
            </v-card-title>
            <v-data-table
                    v-if="listView"
                    :headers="headers"
                    :items="hosts"
                    :search="search"
            >
                <template v-slot:item.address="{ item }">
                    <v-chip
                            v-for="(ip, i) in item.address"
                            :key="i"
                            color="#336699"
                            text-color="white"
                    >
                        <v-icon left>mdi-ip-network</v-icon>
                        {{ ip }}
                    </v-chip>
                </template>
                <template v-slot:item.tags="{ item }">
                    <v-chip
                            v-for="(tag, i) in item.tags"
                            :key="i"
                            color="blue-grey"
                            text-color="white"
                    >
                        <v-icon left>mdi-tag</v-icon>
                        {{ tag }}
                    </v-chip>
                </template>
                <template v-slot:item.created="{ item }">
                    <v-row>
                        <p>{{ item.createdBy }} at {{ item.created | formatDate }}</p>
                    </v-row>
                </template>
                <template v-slot:item.updated="{ item }">
                    <v-row>
                        <p>At {{ item.updated | formatDate }} by {{ item.updatedBy }}</p>
                    </v-row>
                </template>
                <template v-slot:item.action="{ item }">
                    <v-row>
                        <v-icon
                                class="pr-1 pl-1"
                                @click.stop="startUpdate(item)"
                        >
                            mdi-square-edit-outline
                        </v-icon>
                        <v-icon
                                class="pr-1 pl-1"
                                @click.stop="forceFileDownload(item)"
                        >
                            mdi-cloud-download-outline
                        </v-icon>
                        <v-icon
                                class="pr-1 pl-1"
                                @click.stop="email(item)"
                        >
                            mdi-email-send-outline
                        </v-icon>                        
                        <v-icon
                                class="pr-1 pl-1"
                                @click="remove(item)"
                        >
                            mdi-trash-can-outline
                        </v-icon>
                        <v-switch
                                dark
                                class="pr-1 pl-1"
                                color="success"
                                v-model="item.enable"
                                v-on:change="update(item)"
                        />
                    </v-row>
                </template>

            </v-data-table>
            <v-card-text v-else>
                <v-row>
                    <v-col
                            v-for="(host, i) in hosts"
                            :key="i"
                            sm12 lg6
                    >
                        <v-card
                                :color="host.enable ? '#1F7087' : 'warning'"
                                class="mx-auto"
                                raised
                                shaped
                        >
                            <v-list-item>
                                <v-list-item-content>
                                    <v-list-item-title class="headline">{{ host.name }}</v-list-item-title>
                                    <v-list-item-subtitle>{{ host.email }}</v-list-item-subtitle>
                                    <v-list-item-subtitle>Created: {{ host.created | formatDate }} by {{ host.createdBy }}</v-list-item-subtitle>
                                    <v-list-item-subtitle>Updated: {{ host.updated | formatDate }} by {{ host.updatedBy }}</v-list-item-subtitle>
                                </v-list-item-content>

                                <v-list-item-avatar
                                        tile
                                        size="150"
                                >
<!--                                    <v-img :src="'data:image/png;base64, ' + gethostQrcode(host.id)"/> -->
                                </v-list-item-avatar>
                            </v-list-item>

                            <v-card-text class="text--primary">
                                <v-chip
                                        v-for="(ip, i) in host.address"
                                        :key="i"
                                        color="indigo"
                                        text-color="white"
                                >
                                    <v-icon left>mdi-ip-network</v-icon>
                                    {{ ip }}
                                </v-chip>
                            </v-card-text>
                            <v-card-text class="text--primary">
                                <v-chip
                                        v-for="(tag, i) in host.tags"
                                        :key="i"
                                        color="blue-grey"
                                        text-color="white"
                                >
                                    <v-icon left>mdi-tag</v-icon>
                                    {{ tag }}
                                </v-chip>
                            </v-card-text>
                            <v-card-actions>
                                <v-tooltip bottom>
                                    <template v-slot:activator="{ on }">
                                        <v-btn
                                                text
                                                v-on:click="forceFileDownload(host)"
                                                v-on="on"
                                        >
                                            <span class="d-none d-lg-flex">Download</span>
                                            <v-icon right dark>mdi-cloud-download-outline</v-icon>
                                        </v-btn>
                                    </template>
                                    <span>Download</span>
                                </v-tooltip>

                                <v-tooltip bottom>
                                    <template v-slot:activator="{ on }">
                                        <v-btn
                                                text
                                                @click.stop="startUpdate(host)"
                                                v-on="on"
                                        >
                                            <span class="d-none d-lg-flex">Edit</span>
                                            <v-icon right dark>mdi-square-edit-outline</v-icon>
                                        </v-btn>
                                    </template>
                                    <span>Edit</span>
                                </v-tooltip>

                                <v-tooltip bottom>
                                    <template v-slot:activator="{ on }">
                                        <v-btn
                                                text
                                                @click="remove(host)"
                                                v-on="on"
                                        >
                                            <span class="d-none d-lg-flex">Delete</span>
                                            <v-icon right dark>mdi-trash-can-outline</v-icon>
                                        </v-btn>
                                    </template>
                                    <span>Delete</span>
                                </v-tooltip>

                                <v-tooltip bottom>
                                    <template v-slot:activator="{ on }">
                                        <v-btn
                                                text
                                                @click="email(host)"
                                                v-on="on"
                                        >
                                            <span class="d-none d-lg-flex">Send Email</span>
                                            <v-icon right dark>mdi-email-send-outline</v-icon>
                                        </v-btn>
                                    </template>
                                    <span>Send Email</span>
                                </v-tooltip>
                                <v-spacer/>
                                <v-tooltip right>
                                    <template v-slot:activator="{ on }">
                                        <v-switch
                                                dark
                                                v-on="on"
                                                color="success"
                                                v-model="host.enable"
                                                v-on:change="update(host)"
                                        />
                                    </template>
                                    <span> {{host.enable ? 'Disable' : 'Enable'}} this host</span>
                                </v-tooltip>

                            </v-card-actions>
                        </v-card>
                    </v-col>
                </v-row>
            </v-card-text>
        </v-card>
        <v-dialog
                v-if="host"
                v-model="dialogCreate"
                max-width="550"
        >
            <v-card>
                <v-card-title class="headline">Add new host</v-card-title>
                <v-card-text>
                    <v-row>
                        <v-col
                                cols="12"
                        >
                            <v-form
                                    ref="form"
                                    v-model="valid"
                            >
                                <v-text-field
                                        v-model="host.name"
                                        label="host friendly name"
                                        :rules="[ v => !!v || 'host name is required', ]"
                                        required
                                />
                                <v-select
                                        v-model="meshList"
                                        :items="meshList"
                                        label="Join this mesh"
                                        :rules="[ v => !!v || 'Mesh is required', ]"
                                        multiple
                                        chips
                                        persistent-hint
                                        required
                                />
                                <v-combobox
                                        v-model="host.allowedIPs"
                                        chips
                                        hint="Write IPv4 or IPv6 CIDR and hit enter"
                                        label="Allowed IPs"
                                        multiple
                                        dark
                                >
                                    <template v-slot:selection="{ attrs, item, select, selected }">
                                        <v-chip
                                                v-bind="attrs"
                                                :input-value="selected"
                                                close
                                                @click="select"
                                                @click:close="host.allowedIPs.splice(host.allowedIPs.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-combobox
                                        v-model="host.tags"
                                        chips
                                        hint="Enter a tag, hit tab, hit enter."
                                        label="Tags"
                                        multiple
                                        dark
                                >
                                    <template v-slot:selection="{ attrs, item, select, selected }">
                                        <v-chip
                                                v-bind="attrs"
                                                :input-value="selected"
                                                close
                                                @click="select"
                                                @click:close="host.tags.splice(host.tags.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-switch
                                        v-model="host.enable"
                                        color="success"
                                        inset
                                        :label="host.enable ? 'Enable host after creation': 'Disable host after creation'"
                                />
                                <v-switch
                                        v-model="host.ignorePersistentKeepalive"
                                        color="red"
                                        inset
                                        :label="'Ignore global persistent keepalive: ' + (host.ignorePersistentKeepalive ? 'Yes': 'NO')"
                                />
                            </v-form>
                        </v-col>
                    </v-row>
                </v-card-text>
                <v-card-actions>
                    <v-spacer/>
                    <v-btn
                            :disabled="!valid"
                            color="success"
                            @click="create(host)"
                    >
                        Submit
                        <v-icon right dark>mdi-check-outline</v-icon>
                    </v-btn>
                    <v-btn
                            color="primary"
                            @click="dialogCreate = false"
                    >
                        Cancel
                        <v-icon right dark>mdi-close-circle-outline</v-icon>
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
        <v-dialog
                v-if="host"
                v-model="dialogUpdate"
                max-width="550"
        >
            <v-card>
                <v-card-title class="headline">Edit Host</v-card-title>
                <v-card-text>

                    <v-row>
                        <v-col
                                cols="12"
                        >
                            <v-form
                                    ref="form"
                                    v-model="valid"
                            >
                                <v-text-field
                                        v-model="host.name"
                                        label="Friendly name"
                                        :rules="[ v => !!v || 'host name is required',]"
                                        required
                                />

                                <v-combobox
                                        v-model="selected"
                                        :items="meshList"
                                        label="Join this mesh"
                                        :rules="[ v => !!v || 'Mesh is required', ]"
                                        single
                                        persistent-hint
                                        required
                                />
                                <v-combobox
                                        v-model="host.address"
                                        chips
                                        hint="Write IPv4 or IPv6 CIDR and hit enter"
                                        label="Addresses"
                                        multiple
                                        dark
                                >
                                    <template v-slot:selection="{ attrs, item, select, selected }">
                                        <v-chip
                                                v-bind="attrs"
                                                :input-value="selected"
                                                close
                                                @click="select"
                                                @click:close="host.address.splice(host.address.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-combobox
                                    v-model="host.dns"
                                    chips
                                    hint="Write IP address(es) and hit enter or leave empty.  If not empty, be sure to include your local resolver."
                                    label="DNS servers for this host"
                                    multiple
                                    dark
                                >
                                    <template v-slot:selection="{ attrs, item, select, selected }">
                                        <v-chip
                                                v-bind="attrs"
                                                :input-value="selected"
                                                close
                                                @click="select"
                                                @click:close="mesh.default.splice(mesh.default.dns.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>

                                <v-combobox
                                        v-model="host.allowedIPs"
                                        chips
                                        hint="Write IPv4 or IPv6 CIDR and hit enter"
                                        label="Allowed IPs"
                                        multiple
                                        dark
                                >
                                    <template v-slot:selection="{ attrs, item, select, selected }">
                                        <v-chip
                                                v-bind="attrs"
                                                :input-value="selected"
                                                close
                                                @click="select"
                                                @click:close="host.allowedIPs.splice(host.allowedIPs.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-text-field
                                        type="number"
                                        v-model="host.mtu"
                                        label="Define global MTU"
                                        hint="Leave at 0 and let us take care of MTU"
                                />
                                <v-text-field
                                        type="number"
                                        v-model="host.persistentKeepalive"
                                        label="Persistent keepalive"
                                        hint="To disable, set to 0.  Recommended value 29 (seconds)"
                                />
                                <v-combobox
                                        v-model="host.tags"
                                        chips
                                        hint="Write tag name and hit enter"
                                        label="Tags"
                                        multiple
                                        dark
                                >
                                    <template v-slot:selection="{ attrs, item, select, selected }">
                                        <v-chip
                                                v-bind="attrs"
                                                :input-value="selected"
                                                close
                                                @click="select"
                                                @click:close="host.tags.splice(host.tags.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                            </v-form>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
            <v-expansion-panels>
            <v-expansion-panel>
                <v-expansion-panel-header dark>Server configuration</v-expansion-panel-header>
                    <v-expansion-panel-content>
                        <div class="d-flex flex-no-wrap justify-space-between">
                            <v-col cols="12">
                                <v-text-field
                                        v-model="host.publicKey"
                                        label="Public key"
                                        disabled
                                />
                                <v-text-field
                                        v-model="host.endpoint"
                                        label="Public endpoint for hosts to connect to"
                                        :rules="[
                            v => !!v || 'Public endpoint for hosts to connect to is required',
                            ]"
                                        required
                                />
                                <v-text-field
                                        v-model="host.listenPort"
                                        type="number"
                                        :rules="[
                            v => !!v || 'Listen port is required',
                            ]"
                                        label="Listen port"
                                        required
                                />
                            </v-col>
                        </div>
                    </v-expansion-panel-content>
                </v-expansion-panel>
            </v-expansion-panels>
                    <v-spacer/>
                    <v-row>
                        <v-col cols="12">

                    <v-btn
                            :disabled="!valid"
                            color="success"
                            @click="update(host)"
                    >
                        Submit
                        <v-icon right dark>mdi-check-outline</v-icon>
                    </v-btn>
                    <v-btn
                            color="primary"
                            @click="dialogUpdate = false"
                    >
                        Cancel
                        <v-icon right dark>mdi-close-circle-outline</v-icon>
                    </v-btn>
                        </v-col>
                    </v-row>


        </v-dialog>
    </v-container>
</template>
<script>
  import { mapActions, mapGetters } from 'vuex'

  export default {
    name: 'Hosts',

    data: () => ({
      listView: true,
      dialogCreate: false,
      dialogUpdate: false,
      host: null,
      mesh: null,
      panel: 1,
      valid: false,
      meshList: [],
      selected: '',
      search: '',
      headers: [
        { text: 'Name', value: 'name', },
        { text: 'Mesh', value: 'meshName', },
        { text: "Endpoint", value: 'endpoint', },
        { text: 'IP addresses', value: 'address', },
        { text: 'Created by', value: 'created', sortable: false, },
        { text: 'Tags', value: 'tags', },
        { text: 'Actions', value: 'action', sortable: false, },

      ],
    }),

    computed:{
      ...mapGetters({
        gethostQrcode: 'host/gethostQrcode',
        gethostConfig: 'host/gethostConfig',
        user: 'auth/user',
        server: 'server/server',
        hosts: 'host/hosts',
        meshes: 'mesh/meshes',
        hostQrcodes: 'host/hostQrcodes',
      }),
    },

    mounted () {
      this.readAllHosts()
      this.readAllMeshes()
    },

    methods: {
      ...mapActions('host', {
        errorhost: 'error',
        readAllHosts: 'readAll',
        creathost: 'create',
        updatehost: 'update',
        deletehost: 'delete',
        emailhost: 'email',
      }),
      ...mapActions('mesh', {
        readAllMeshes: 'readAll',
      }),

      startCreate() {
        this.host = {
          name: "",
          email: this.user.email,
          enable: true,
//          meshName: this.meshes[0].meshName,
//          meshID: this.meshes[0].meshid,
//          allowedIPs: this.meshes[0].default.allowedIPs,
//          address: this.meshes[0].default.address,
//          meshName: this.meshes[0].default.meshName,
//          meshid: this.meshes[0].default.meshid,
          tags: [],
        }
        
        for (let i=0; i<this.meshes.length; i++) {
            this.meshList[i] = this.meshes[i].meshName
        }
        this.dialogCreate = true;
      },

      create(host) {
/*        if (host.allowedIPs.length < 0) {
          this.errorhost('Please provide at least one valid CIDR address for host allowed IPs')
          return;
        }
        for (let i = 0; i < host.allowedIPs.length; i++){
          if (this.$isCidr(host.allowedIPs[i]) === 0) {
            this.errorhost('Invalid CIDR detected, please correct before submitting')
            return
          }
        }*/
        this.dialogCreate = false;
        this.creathost(host)
      },

      remove(host) {
        if(confirm(`Do you really want to delete ${host.name} ?`)){
          this.deletehost(host)
        }
      },

      email(host) {
        if (!host.email){
          this.errorhost('host email is not defined')
          return
        }

        if(confirm(`Do you really want to send email to ${host.email} with all configurations ?`)){
          this.emailhost(host)
        }
      },

      startUpdate(host) {
        this.host = host;
        this.dialogUpdate = true;
        this.selected = this.host.meshName;

       var i=0;
        for (i=0; i<this.meshes.length; i++) {
            this.meshList[i] = this.meshes[i].meshName
        }

      },

      update(host) {

        this.host.listenPort = parseInt(this.host.listenPort, 10);
        this.host.persistentKeepalive = parseInt(this.host.persistentKeepalive, 10);
        this.host.mtu = parseInt(this.host.mtu, 10);
        this.host.meshName = this.selected
//        this.host.meshid = this.server.meshid
//        this.host.meshName = this.server.meshName


        // check allowed IPs
        if (host.allowedIPs.length < 1) {
          this.errorhost('Please provide at least one valid CIDR address for host allowed IPs');
          return;
        }
        for (let i = 0; i < host.allowedIPs.length; i++){
          if (this.$isCidr(host.allowedIPs[i]) === 0) {
            this.errorhost('Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // check address
        if (host.address.length < 1) {
          this.errorhost('Please provide at least one valid CIDR address for host');
          return;
        }
        for (let i = 0; i < host.address.length; i++){
          if (this.$isCidr(host.address[i]) === 0) {
            this.errorhost('Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // all good, submit
        this.dialogUpdate = false;
        this.updatehost(host)
      },

      forceFileDownload(host){
        let config = this.gethostConfig(host.id)
        if (!config) {
          this.errorhost('Failed to download host config');
          return
        }
        const url = window.URL.createObjectURL(new Blob([config]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', host.name.split(' ').join('-') + '.conf') //or any other extension
        document.body.appendChild(link)
        link.click()
      },
    }
  };
</script>
