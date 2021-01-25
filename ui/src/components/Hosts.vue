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
                    <v-icon right dark>mdi-network-outline</v-icon>
                </v-btn>
            </v-card-title>
            <v-data-table
                    v-if="listView"
                    :headers="headers"
                    :items="hosts"
                    :search="search"
            >
                <template v-slot:item.name="{ item }">
                        {{ item.name }}
                    </v-chip>
                </template>
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
                                        label="Host friendly name"
                                        :rules="[ v => !!v || 'host name is required', ]"
                                        required
                                />
                                <v-select return-object
                                        v-model="meshList.selected"
                                        :items="meshList.items"
                                        item-text = "text"
                                        item-value = "value"
                                        label="Join this mesh"
                                        :rules="[ v => !!v || 'Mesh is required', ]"
                                        single
                                        persistent-hint
                                        required
                                />
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

                                <v-select return-object
                                        v-model="meshList.selected"
                                        :items="meshList.items"
                                        item-text = "text"
                                        item-value = "value"
                                        label="Join this mesh"
                                        :rules="[ v => !!v || 'Mesh is required', ]"
                                        single
                                        persistent-hint
                                        required
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
                <v-expansion-panel-header dark>Advanced Configuration</v-expansion-panel-header>
                    <v-expansion-panel-content>
                        <div class="d-flex flex-no-wrap justify-space-between">
                            <v-col cols="12">
                                <v-combobox
                                        v-model="host.current.address"
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
                                                @click:close="host.current.address.splice(host.current.address.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-combobox
                                        v-model="host.current.allowedIPs"
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
                                                @click:close="host.current.allowedIPs.splice(host.current.allowedIPs.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-combobox
                                    v-model="host.current.dns"
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
                                                @click:close="host.current.splice(host.current.dns.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-text-field
                                        v-model="host.current.publicKey"
                                        label="Public key"
                                        disabled
                                />
                                <v-text-field
                                        v-model="host.id"
                                        label="Host ID"
                                        disabled
                                />
                                <v-text-field
                                        v-model="host.current.endpoint"
                                        label="Public endpoint for hosts to connect to"
                                />
                                <v-text-field
                                        v-model="host.current.listenPort"
                                        type="number"
                                        label="Listen port"
                                />
                                <v-text-field
                                        type="number"
                                        v-model="host.current.mtu"
                                        label="Define global MTU"
                                        hint="Leave at 0 and let us take care of MTU"
                                />
                                <v-text-field
                                        type="number"
                                        v-model="host.current.persistentKeepalive"
                                        label="Persistent keepalive"
                                        hint="To disable, set to 0.  Recommended value 29 (seconds)"
                                />
                                <v-switch
                                        v-model="host.current.subnetRouting"
                                        color="success"
                                        inset
                                        label="Enable subnet routing"
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
      meshList: {},
      selected: '',
      search: '',
      headers: [
        { text: 'Name', value: 'name', },
        { text: 'Mesh', value: 'meshName', },
        { text: 'IP addresses', value: 'current.address', },
//        { text: 'ID', value:'id', },
        { text: "Endpoint", value: 'current.endpoint', },
//        { text: 'Created by', value: 'created', sortable: false, },
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
        accounts: 'account/accounts',
        hosts: 'host/hosts',
        meshes: 'mesh/meshes',
        hostQrcodes: 'host/hostQrcodes',
      }),
    },

    mounted () {
      this.readAllAccounts(this.user.email)
      this.readAllHosts()
      this.readAllMeshes()
    },

    methods: {
      ...mapActions('host', {
        errorhost: 'error',
        readAllHosts: 'readAll',
        createhost: 'create',
        updatehost: 'update',
        deletehost: 'delete',
        emailhost: 'email',
      }),
      ...mapActions('mesh', {
        readAllMeshes: 'readAll',
      }),
      ...mapActions('account', {
          readAllAccounts: 'readAll',
      }),


      startCreate() {
        this.host = {
          name: "",
          email: this.user.email,
          enable: true,
//          meshName: this.meshes[0].meshName,
//          meshID: this.meshes[0].id,
//          allowedIPs: this.meshes[0].default.allowedIPs,
//          address: this.meshes[0].default.address,
//          meshName: this.meshes[0].default.meshName,
//          id: this.meshes[0].default.id,
          tags: [],
          current: {},
        }
        
        this.meshList = { selected: { "text": "",  "value": ""},
                          items: [] }

        var selected = 0;
        for (let i=0; i<this.meshes.length; i++) {
            this.meshList.items[i]= { "text": this.meshes[i].meshName, "value": this.meshes[i].id }
            if (this.meshList.items[i].text == this.host.meshName) {
                selected = i
            }
        }

        this.meshList.selected = this.meshList.items[selected];
        this.dialogCreate = true;
      },

      create(host) {

        this.host.meshName = this.meshList.selected.text
        this.host.meshid = this.meshList.selected.value
        this.host.accountid = this.accounts[0].id
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

        this.meshList = { selected: { "text": this.host.meshName,  "value": this.host.meshid },
                          items: [] }

        var selected = 0;
        for (let i=0; i<this.meshes.length; i++) {
            this.meshList.items[i]= { "text": this.meshes[i].meshName, "value": this.meshes[i].id }
            if (this.meshList.items[i].text == this.host.meshName) {
                selected = i
            }
        }

        this.meshList.selected = this.meshList.items[selected];

        this.dialogUpdate = true;

      },

      reconcile(host, template) {

      },

      update(host) {

        this.host = host

        this.host.accountid = this.accounts[0].id;

        this.host.current.listenPort = parseInt(this.host.current.listenPort, 10);
        this.host.current.persistentKeepalive = parseInt(this.host.current.persistentKeepalive, 10);
        this.host.current.mtu = parseInt(this.host.current.mtu, 10);

        var changed = false;
        if (this.host.meshid != this.meshList.selected.value) {
            this.host.meshName = this.meshList.selected.text
            this.host.meshid = this.meshList.selected.value
            changed = true;
        }

        if (changed) {
            for (let i=0; i<this.meshes.length; i++) {
                if ( this.host.meshid == this.meshList.items[i].value ) {
                    var template = this.meshes[i]
                    this.host.current.address = []

                    if (host.current.listenPort == host.default.listenPort) {
                        host.current.listenPort = template.listenPort
                    }

                    if (host.current.allowedIPs == host.default.allowedIPs) {
                        host.current.allowedIPs = template.allowedIPs
                    } else {
                        host.current.allowedIPs = []
                    }

                    if (host.current.mtu == host.default.mtu) {
                        host.current.mtu = template.mtu
                    }

                    if (host.current.dns == host.default.dns) {
                        host.current.dns = template.dns
                    } else {
                        host.current.dns = []
                    }

                    if (host.current.persistentKeepalive == host.default.persistentKeepalive) {
                        host.current.persistentKeepalive = template.persistentKeepalive
                    }
                    this.host.default = this.meshes[i].default
                } 
            }
        }

/*
        // check allowed IPs
        if (host.current.allowedIPs.length < 1) {
          this.errorhost('Please provide at least one valid CIDR address for host allowed IPs');
          return;
        }
        for (let i = 0; i < host.current.allowedIPs.length; i++){
          if (this.$isCidr(host.current.allowedIPs[i]) === 0) {
            this.errorhost('Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // check address
        if (host.current.address.length < 1) {
          this.errorhost('Please provide at least one valid CIDR address for host');
          return;
        }
        for (let i = 0; i < host.current.address.length; i++){
          if (this.$isCidr(host.current.address[i]) === 0) {
            this.errorhost('Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        */
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
