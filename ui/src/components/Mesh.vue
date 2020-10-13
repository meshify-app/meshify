<template>
    <v-container>
        <v-card>
            <v-card-title>
                Meshes
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
                    Create new mesh
                    <v-icon right dark>mdi-vpn</v-icon>
                </v-btn>
            </v-card-title>
            <v-data-table
                    v-if="listView"
                    :headers="headers"
                    :items="meshes"
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
                            v-for="(mesh, i) in meshes "
                            :key="i"
                            sm12 lg6
                    >
                        <v-card
                                :color="mesh.enable ? '#1F7087' : 'warning'"
                                class="mx-auto"
                                raised
                                shaped
                        >
                            <v-list-item>
                                <v-list-item-content>
                                    <v-list-item-title class="headline">{{ mesh.name }}</v-list-item-title>
                                    <v-list-item-subtitle>{{ mesh.email }}</v-list-item-subtitle>
                                    <v-list-item-subtitle>Created: {{ mesh.created | formatDate }} by {{ mesh.createdBy }}</v-list-item-subtitle>
                                    <v-list-item-subtitle>Updated: {{ mesh.updated | formatDate }} by {{ mesh.updatedBy }}</v-list-item-subtitle>
                                </v-list-item-content>

                                <v-list-item-avatar
                                        tile
                                        size="150"
                                >
                                </v-list-item-avatar>
                            </v-list-item>

                            <v-card-text class="text--primary">
                                <v-chip
                                        v-for="(ip, i) in mesh.address"
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
                                        v-for="(tag, i) in mesh.tags"
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
                                                v-on:click="forceFileDownload(mesh)"
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
                                                @click.stop="startUpdate(mesh)"
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
                                                @click="remove(mesh)"
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
                                                @click="email(mesh)"
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
                                                v-model="mesh.enable"
                                                v-on:change="update(mesh)"
                                        />
                                    </template>
                                    <span> {{mesh.enable ? 'Disable' : 'Enable'}} this mesh</span>
                                </v-tooltip>

                            </v-card-actions>
                        </v-card>
                    </v-col>
                </v-row>
            </v-card-text>
        </v-card>
        <v-dialog
                v-if="mesh"
                v-model="dialogCreate"
                max-width="550"
        >
            <v-card>
                <v-card-title class="headline">Create new mesh</v-card-title>
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
                                        v-model="mesh.name"
                                        label="Mesh friendly name"
                                        :rules="[ v => !!v || 'Mesh name is required', ]"
                                        required
                                />
                                <v-text-field
                                        v-model="mesh.email"
                                        label="Mesh email"
                                        :rules="[ v => (/.+@.+\..+/.test(v) || v === '') || 'E-mail must be valid',]"
                                />
                                <v-select
                                        v-model="mesh.address"
                                        :items="server.address"
                                        label="Mesh IP will be chosen from these networks"
                                        :rules="[ v => !!v || 'Network is required', ]"
                                        multiple
                                        chips
                                        persistent-hint
                                        required
                                />
                                <v-combobox
                                        v-model="mesh.allowedIPs"
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
                                                @click:close="mesh.allowedIPs.splice(mesh.allowedIPs.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-combobox
                                        v-model="mesh.tags"
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
                                                @click:close="mesh.tags.splice(mesh.tags.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-switch
                                        v-model="mesh.enable"
                                        color="success"
                                        inset
                                        :label="mesh.enable ? 'Enable mesh after creation': 'Disable mesh after creation'"
                                />
                                <v-switch
                                        v-model="mesh.ignorePersistentKeepalive"
                                        color="red"
                                        inset
                                        :label="'Ignore global persistent keepalive: ' + (mesh.ignorePersistentKeepalive ? 'Yes': 'NO')"
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
                            @click="create(mesh)"
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
                v-if="mesh"
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
                                        v-model="mesh.name"
                                        label="Friendly name"
                                        :rules="[ v => !!v || 'Mesh name is required',]"
                                        required
                                />
                                <v-combobox
                                        v-model="mesh.address"
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
                                                @click:close="mesh.address.splice(mesh.address.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-combobox
                                    v-model="mesh.dns"
                                    chips
                                    hint="Write IP address(es) and hit enter or leave empty.  If not empty, be sure to include your local resolver."
                                    label="DNS servers for this mesh"
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
                                        v-model="mesh.allowedIPs"
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
                                                @click:close="mesh.allowedIPs.splice(mesh.allowedIPs.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-text-field
                                        type="number"
                                        v-model="mesh.mtu"
                                        label="Define global MTU"
                                        hint="Leave at 0 and let us take care of MTU"
                                />
                                <v-text-field
                                        type="number"
                                        v-model="mesh.persistentKeepalive"
                                        label="Persistent keepalive"
                                        hint="To disable, set to 0.  Recommended value 29 (seconds)"
                                />
                                <v-combobox
                                        v-model="mesh.tags"
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
                                                @click:close="mesh.tags.splice(mesh.tags.indexOf(item), 1)"
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
                                        v-model="mesh.publicKey"
                                        label="Public key"
                                        disabled
                                />
                                <v-text-field
                                        v-model="mesh.endpoint"
                                        label="Public endpoint for clients to connect to"
                                        :rules="[
                            v => !!v || 'Public endpoint for clients to connect to is required',
                            ]"
                                        required
                                />
                                <v-text-field
                                        v-model="mesh.listenPort"
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
                            @click="update(mesh)"
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
    name: 'Meshes',

    data: () => ({
      listView: true,
      dialogCreate: false,
      dialogUpdate: false,
      mesh: null,
      panel: 1,
      valid: false,
      search: '',
      headers: [
        { text: 'Name', value: 'name', },
//        { text: 'Email', value: 'email', },
//        { text: "Endpoint", value: 'endpoint', },
        { text: 'IP address pool', value: 'address', },
        { text: 'Created by', value: 'created', sortable: false, },
        { text: 'Tags', value: 'tags', },
        { text: 'Actions', value: 'action', sortable: false, },

      ],
    }),

    computed:{
      ...mapGetters({
        user: 'auth/user',
        server: 'server/server',
        meshes: 'mesh/mesh',
      }),
    },

    mounted () {
      this.readAllMeshes()
      this.readServer()
    },

    methods: {
      ...mapActions('mesh', {
        errorMesh: 'error',
        readAllMeshes: 'readAll',
        createMesh: 'create',
        updateMesh: 'update',
        deleteMesh: 'delete',
        emailMesh: 'email',
      }),
      ...mapActions('server', {
        readServer: 'read',
      }),

      startCreate() {
        this.mesh = {
          name: "",
          email: this.user.email,
          enable: true,
          allowedIPs: this.server.address,
          address: this.server.address,
          meshName: this.server.meshName,
          meshid: this.server.meshid,
          tags: [],
        }
        this.dialogCreate = true;
      },

      create(mesh) {
        if (mesh.allowedIPs.length < 0) {
          this.errorMesh('Please provide at least one valid CIDR address for mesh allowed IPs')
          return;
        }
        for (let i = 0; i < mesh.allowedIPs.length; i++){
          if (this.$isCidr(mesh.allowedIPs[i]) === 0) {
            this.errorMesh('Invalid CIDR detected, please correct before submitting')
            return
          }
        }
        this.dialogCreate = false;
        this.createMesh(mesh)
      },

      remove(mesh) {
        if(confirm(`Do you really want to delete ${mesh.name} ?`)){
          this.deleteMesh(mesh)
        }
      },

      email(mesh) {
        if (!mesh.email){
          this.errorMesh('Mesh email is not defined')
          return
        }

        if(confirm(`Do you really want to send email to ${mesh.email} with all configurations ?`)){
          this.emailMesh(mesh)
        }
      },

      startUpdate(mesh) {
        this.mesh = mesh;
        this.dialogUpdate = true;
      },

      update(mesh) {

        this.mesh.listenPort = parseInt(this.mesh.listenPort, 10);
        this.mesh.persistentKeepalive = parseInt(this.mesh.persistentKeepalive, 10);
        this.mesh.mtu = parseInt(this.mesh.mtu, 10);
//        this.mesh.meshid = this.server.meshid
//        this.mesh.meshName = this.server.meshName


        // check allowed IPs
        if (mesh.allowedIPs.length < 1) {
          this.errorMesh('Please provide at least one valid CIDR address for mesh allowed IPs');
          return;
        }
        for (let i = 0; i < mesh.allowedIPs.length; i++){
          if (this.$isCidr(mesh.allowedIPs[i]) === 0) {
            this.errorMesh('Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // check address
        if (mesh.address.length < 1) {
          this.errorMesh('Please provide at least one valid CIDR address for mesh');
          return;
        }
        for (let i = 0; i < mesh.address.length; i++){
          if (this.$isCidr(mesh.address[i]) === 0) {
            this.errorMesh('Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // all good, submit
        this.dialogUpdate = false;
        this.updateMesh(mesh)
      },

      forceFileDownload(mesh){
        let config = this.getMeshConfig(mesh.meshid)
        if (!config) {
          this.errorMesh('Failed to download mesh config');
          return
        }
        const url = window.URL.createObjectURL(new Blob([config]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', mesh.name.split(' ').join('-') + '.conf') //or any other extension
        document.body.appendChild(link)
        link.click()
      },
    }
  };
</script>
