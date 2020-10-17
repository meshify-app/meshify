<template>
    <v-container>
        <v-card>
            <v-card-title>
                Users
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
                    Invite new user
                    <v-icon right dark>mdi-vpn</v-icon>
                </v-btn>
            </v-card-title>
            <v-data-table
                    v-if="listView"
                    :headers="headers"
                    :items="users"
                    :search="search"
            >
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
                    </v-row>
                </template>

            </v-data-table>
            <v-card-text v-else>
                <v-row>
                    <v-col
                            v-for="(user, i) in users "
                            :key="i"
                            sm12 lg6
                    >
                        <v-card
                                :color="user.enable ? '#1F7087' : 'warning'"
                                class="mx-auto"
                                raised
                                shaped
                        >
                            <v-list-item>
                                <v-list-item-content>
                                    <v-list-item-title class="headline">{{ user.name }}</v-list-item-title>
                                    <v-list-item-subtitle>{{ user.email }}</v-list-item-subtitle>
                                </v-list-item-content>

                                <v-list-item-avatar
                                        tile
                                        size="150"
                                >
                                </v-list-item-avatar>
                            </v-list-item>
                            <v-card-actions>
                                <v-tooltip bottom>
                                    <template v-slot:activator="{ on }">
                                        <v-btn
                                                text
                                                v-on:click="forceFileDownload(user)"
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
                                                @click.stop="startUpdate(user)"
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
                                                @click="remove(user)"
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
                                                @click="email(user)"
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
                                                v-model="user.enable"
                                                v-on:change="update(user)"
                                        />
                                    </template>
                                    <span> {{user.enable ? 'Disable' : 'Enable'}} this user</span>
                                </v-tooltip>

                            </v-card-actions>
                        </v-card>
                    </v-col>
                </v-row>
            </v-card-text>
        </v-card>
        <v-dialog
                v-if="user"
                v-model="dialogCreate"
                max-width="550"
        >
            <v-card>
                <v-card-title class="headline">Invite new user</v-card-title>
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
                                        v-model="user.name"
                                        label="User friendly name"
                                        :rules="[ v => !!v || 'User name is required', ]"
                                        required
                                />
                                <v-text-field
                                        v-model="user.email"
                                        label="User email"
                                        :rules="[ v => (/.+@.+\..+/.test(v) || v === '') || 'E-mail must be valid',]"
                                />
                                <v-select
                                        v-model="user.address"
                                        :items="server.address"
                                        label="User IP will be chosen from these networks"
                                        :rules="[ v => !!v || 'Network is required', ]"
                                        multiple
                                        chips
                                        persistent-hint
                                        required
                                />
                                <v-combobox
                                        v-model="user.allowedIPs"
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
                                                @click:close="user.allowedIPs.splice(user.allowedIPs.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-combobox
                                        v-model="user.tags"
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
                                                @click:close="user.tags.splice(user.tags.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-switch
                                        v-model="user.enable"
                                        color="success"
                                        inset
                                        :label="user.enable ? 'Enable user after creation': 'Disable user after creation'"
                                />
                                <v-switch
                                        v-model="user.ignorePersistentKeepalive"
                                        color="red"
                                        inset
                                        :label="'Ignore global persistent keepalive: ' + (user.ignorePersistentKeepalive ? 'Yes': 'NO')"
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
                            @click="create(user)"
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
                v-if="user"
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
                                        v-model="user.name"
                                        label="Friendly name"
                                        :rules="[ v => !!v || 'User name is required',]"
                                        required
                                />
                                <v-combobox
                                        v-model="user.address"
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
                                                @click:close="user.address.splice(user.address.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-combobox
                                    v-model="user.dns"
                                    chips
                                    hint="Write IP address(es) and hit enter or leave empty.  If not empty, be sure to include your local resolver."
                                    label="DNS servers for this user"
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
                                        v-model="user.allowedIPs"
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
                                                @click:close="user.allowedIPs.splice(user.allowedIPs.indexOf(item), 1)"
                                        >
                                            <strong>{{ item }}</strong>&nbsp;
                                        </v-chip>
                                    </template>
                                </v-combobox>
                                <v-text-field
                                        type="number"
                                        v-model="user.mtu"
                                        label="Define global MTU"
                                        hint="Leave at 0 and let us take care of MTU"
                                />
                                <v-text-field
                                        type="number"
                                        v-model="user.persistentKeepalive"
                                        label="Persistent keepalive"
                                        hint="To disable, set to 0.  Recommended value 29 (seconds)"
                                />
                                <v-combobox
                                        v-model="user.tags"
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
                                                @click:close="user.tags.splice(user.tags.indexOf(item), 1)"
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
                                        v-model="user.publicKey"
                                        label="Public key"
                                        disabled
                                />
                                <v-text-field
                                        v-model="user.endpoint"
                                        label="Public endpoint for clients to connect to"
                                        :rules="[
                            v => !!v || 'Public endpoint for clients to connect to is required',
                            ]"
                                        required
                                />
                                <v-text-field
                                        v-model="user.listenPort"
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
                            @click="update(user)"
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
    name: 'Users',

    data: () => ({
      listView: true,
      dialogCreate: false,
      dialogUpdate: false,
      user: null,
      panel: 1,
      valid: false,
      search: '',
      headers: [
        { text: 'Email', value: 'email', },
        { text: 'Name', value: 'name', },
//        { text: "Endpoint", value: 'endpoint', },
//        { text: 'IP address pool', value: 'address', },
        { text: 'Created', value: 'created_at', sortable: false, },
//        { text: 'Tags', value: 'tags', },
        { text: 'Actions', value: 'action', sortable: false, },

      ],
    }),

    computed:{
      ...mapGetters({
        authuser: 'auth/user',
        server: 'server/server',
        users: 'user/users',
      }),
    },

    mounted () {
      this.readAllUsers()
      this.readServer()
    },

    methods: {
      ...mapActions('user', {
        errorUser: 'error',
        readAllUsers: 'readAll',
        createUser: 'create',
        updateUser: 'update',
        deleteUser: 'delete',
        emailUser: 'email',
      }),
      ...mapActions('server', {
        readServer: 'read',
      }),

      startCreate() {
        this.user = {
          name: "",
          email: this.user.email,
          enable: true,
          allowedIPs: this.server.address,
          address: this.server.address,
          userName: this.server.userName,
          userid: this.server.userid,
          tags: [],
        }
        this.dialogCreate = true;
      },

      create(user) {
        if (user.allowedIPs.length < 0) {
          this.errorUser('Please provide at least one valid CIDR address for user allowed IPs')
          return;
        }
        for (let i = 0; i < user.allowedIPs.length; i++){
          if (this.$isCidr(user.allowedIPs[i]) === 0) {
            this.errorUser('Invalid CIDR detected, please correct before submitting')
            return
          }
        }
        this.dialogCreate = false;
        this.createUser(user)
      },

      remove(user) {
        if(confirm(`Do you really want to delete ${user.name} ?`)){
          this.deleteUser(user)
        }
      },

      email(user) {
        if (!user.email){
          this.errorUser('User email is not defined')
          return
        }

        if(confirm(`Do you really want to send email to ${user.email} with all configurations ?`)){
          this.emailUser(user)
        }
      },

      startUpdate(user) {
        this.user = user;
        this.dialogUpdate = true;
      },

      update(user) {

        this.user.listenPort = parseInt(this.user.listenPort, 10);
        this.user.persistentKeepalive = parseInt(this.user.persistentKeepalive, 10);
        this.user.mtu = parseInt(this.user.mtu, 10);
//        this.user.userid = this.server.userid
//        this.user.userName = this.server.userName


        // check allowed IPs
        if (user.allowedIPs.length < 1) {
          this.errorUser('Please provide at least one valid CIDR address for user allowed IPs');
          return;
        }
        for (let i = 0; i < user.allowedIPs.length; i++){
          if (this.$isCidr(user.allowedIPs[i]) === 0) {
            this.errorUser('Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // check address
        if (user.address.length < 1) {
          this.errorUser('Please provide at least one valid CIDR address for user');
          return;
        }
        for (let i = 0; i < user.address.length; i++){
          if (this.$isCidr(user.address[i]) === 0) {
            this.errorUser('Invalid CIDR detected, please correct before submitting');
            return
          }
        }
        // all good, submit
        this.dialogUpdate = false;
        this.updateUser(user)
      },

      forceFileDownload(user){
        let config = this.getUserConfig(user.userid)
        if (!config) {
          this.errorUser('Failed to download user config');
          return
        }
        const url = window.URL.createObjectURL(new Blob([config]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', user.name.split(' ').join('-') + '.conf') //or any other extension
        document.body.appendChild(link)
        link.click()
      },
    }
  };
</script>
