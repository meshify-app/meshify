<template>
    <v-container>
        <v-card>
            <v-card-title>
                Services & Subscriptions
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
                        @click="startCreateService"
                >
                    Add Service Host
                    <v-icon right dark>mdi-weather-cloudy</v-icon>
                </v-btn>
            </v-card-title>
        </v-card>
        <v-card>
            <v-card-title>
                Subscriptions
                <v-spacer></v-spacer>
            </v-card-title>
            <v-data-table
                    v-if="listView"
                    :headers="headers"
                    :items="subscriptions"
                    :search="search"
                     @click:row="startUpdate"
            >
                <template slot="no-data">
                    Welcome to Meshify!  Order a subscription on the <a href="https://meshify.app">Meshify website</a> to get started.
                </template>

                <template v-slot:item.issued="{ item }">
                    {{ item.issued | formatDate }}
                </template>
                <template v-slot:item.expires="{ item }">
                    {{ item.expires | formatDate }}
                </template>
                <template v-slot:item.lastUpdated="{ item }">
                    <v-row>
                        <p>At {{ item.lastUpdated | formatDate }} by {{ item.updatedBy }}</p>
                    </v-row>
                </template>
                <template v-slot:item.action="{ item }">
                    <v-row>
<!--
                        <v-icon
                                class="pr-1 pl-1"
                                @click.stop="startUpdate(item)"
                        >
                            mdi-account-switch
                        </v-icon>
-->
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
                            v-for="(subscription, i) in subscriptions "
                            :key="i"
                            sm12 lg6
                    >
                        <v-card
                                :color="subscription.enable ? '#1F7087' : 'warning'"
                                class="mx-auto"
                                raised
                                shaped
                        >
                            <v-list-item>
                                <v-list-item-content>
                                    <v-list-item-title class="headline">{{ subscription.subscriptionName }}</v-list-item-title>
                                    <v-list-item-subtitle>Created: {{ subscription.created | formatDate }} by {{ subscription.createdBy }}</v-list-item-subtitle>
                                    <v-list-item-subtitle>Updated: {{ subscription.updated | formatDate }} by {{ subscription.updatedBy }}</v-list-item-subtitle>
                                </v-list-item-content>

                                <v-list-item-avatar
                                        tile
                                        size="150"
                                >
                                </v-list-item-avatar>
                            </v-list-item>

                            <v-card-text class="text--primary">
                                <v-chip
                                        v-for="(ip, i) in subscription.default.address"
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
                                        v-for="(tag, i) in subscription.default.tags"
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
                                                @click.stop="startUpdate(subscription)"
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
                                                @click="remove(subscription)"
                                                v-on="on"
                                        >
                                            <span class="d-none d-lg-flex">Delete</span>
                                            <v-icon right dark>mdi-trash-can-outline</v-icon>
                                        </v-btn>
                                    </template>
                                    <span>Delete</span>
                                </v-tooltip>
                            </v-card-actions>
                        </v-card>
                    </v-col>
                </v-row>
            </v-card-text>
        </v-card>
        <v-card>
            <v-card-title>
                Service Hosts
                <v-spacer></v-spacer>
            </v-card-title>
            <v-data-table
                    v-if="listView"
                    :headers="bottom_headers"
                    :items="services"
                    :search="search"
                    @click:row="startUpdateService"
            >
                <template slot="no-data">
                    Creating a service host requires a subscription.  Order a subscription on the <a href="https://meshify.app">Meshify website</a> to get started.
                </template>
                <template v-slot:item.created="{ item }">
                    {{ item.created | formatDate }}
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
                            v-for="(service, i) in services "
                            :key="i"
                            sm12 lg6
                    >
                        <v-card
                                :color="member.enable ? '#1F7087' : 'warning'"
                                class="mx-auto"
                                raised
                                shaped
                        >
                            <v-list-item>
                                <v-list-item-content>
                                    <v-list-item-title class="headline">{{ service.name }}</v-list-item-title>
                                    <v-list-item-subtitle>{{ service.description }}</v-list-item-subtitle>
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
                                                @click.stop="startUpdateMember(member)"
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
                                                @click="removeService(service)"
                                                v-on="on"
                                        >
                                            <span class="d-none d-lg-flex">Delete</span>
                                            <v-icon right dark>mdi-trash-can-outline</v-icon>
                                        </v-btn>
                                    </template>
                                    <span>Delete</span>
                                </v-tooltip>
                            </v-card-actions>
                        </v-card>
                    </v-col>
                </v-row>
            </v-card-text>
        </v-card>
        <v-dialog
                v-if="subscriptions"
                v-model="dialogCreateService"
                max-width="550"
        >
            <v-card>
                <v-card-title class="headline">Create New Service Host</v-card-title>
                <v-card-text>
                    <v-row>
                        <v-col
                                cols="12"
                        >
                            <v-form
                                    ref="form"
                                    v-model="valid"
                            >
                                <v-select return-object
                                        v-model="serverList.selected"
                                        :items="serverList.items"
                                        item-text = "text"
                                        item-value = "value"
                                        label="Pick region"
                                        :rules="[ v => !!v || 'Server is required', ]"
                                        single
                                        persistent-hint
                                        required
                                />
                                <v-select return-object
                                        v-model="meshList.selected"
                                        :items="meshList.items"
                                        item-text = "text"
                                        item-value = "value"
                                        label="To this mesh"
                                        :rules="[ v => !!v || 'Mesh is required', ]"
                                        single
                                        persistent-hint
                                        required
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
                            @click="create()"
                    >
                        Submit
                        <v-icon right dark>mdi-check-outline</v-icon>
                    </v-btn>
                    <v-btn
                            color="primary"
                            @click="dialogCreateService = false"
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
                    <v-card-actions>
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
                    </v-card-actions>
            </v-card>
        </v-dialog>
        <v-dialog
                v-if="member"
                v-model="dialogMember"
                max-width="550"
        >
            <v-card>
                <v-card-title class="headline">Edit Member</v-card-title>
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
                                        v-model="member.accountName"
                                        label="Account Name"
                                        :rules="[ v => !!v || 'Account name is required',]"
                                        required
                                />
                                <v-text-field
                                        v-model="member.email"
                                        label="Email Address"
                                        :rules="[ v => !!v || 'Email address is required',]"
                                        required
                                />
                                <v-text-field
                                        v-model="member.name"
                                        label="Name"
                                        :rules="[ v => !!v || 'Name is required',]"
                                        required
                                />
                                <v-select return-object
                                        v-model="meshList.selected"
                                        :items="meshList.items"
                                        item-text = "text"
                                        item-value = "value"
                                        label="To this mesh"
                                        :rules="[ v => !!v || 'Mesh is required', ]"
                                        single
                                        persistent-hint
                                        required
                                />
                                <v-select
                                    :items="roles"
                                    v-model="member.role"
                                    label="Role"
                                ></v-select>
                                <v-select
                                    :items="statuses"
                                    v-model="member.status"
                                    label="Status"
                                ></v-select>
                            </v-form>
                        </v-col>
                    </v-row>
                </v-card-text>
                    <v-card-actions>
                        <v-btn
                                :disabled="!valid"
                                color="success"
                                @click="updateMember(member)"
                        >
                            Submit
                            <v-icon right dark>mdi-check-outline</v-icon>
                        </v-btn>
                        <v-btn
                                color="primary"
                                @click="dialogMember = false"
                        >
                            Cancel
                            <v-icon right dark>mdi-close-circle-outline</v-icon>
                        </v-btn>
                    </v-card-actions>
            </v-card>
        </v-dialog>
    </v-container>
</template>
<script>
  import { mapActions, mapGetters } from 'vuex'

  export default {
    name: 'Services',

    data: () => ({
      listView: true,
      dialogCreateService: false,
      dialogUpdate: false,
      dialogMember : false,
      inDelete: false,
      meshList: {},
      serverList: {},
      server: null,
      roles : ["Owner", "Admin", "User"],
      statuses : ["Active", "Pending", "Suspended", "Hidden"],
      user: null,
      member: null,
      account: null,
      subscription: null,
      service: null,
      panel: 1,
      valid: false,
      search: '',
      headers: [
        { text: 'Name', value: 'name', },
        { text: "Description", value: 'description', },
        { text: 'Issued', value: 'issued', },
        { text: 'Expires', value: 'expires', },
        { text: 'Status', value: 'status', },
        { text: 'Actions', value: 'action', sortable: false, },

      ],
      bottom_headers: [
        { text: 'Name', value: 'name', },
        { text: "Description", value: 'description', },
        { text: 'Created', value: 'created', },
        { text: 'Status', value: 'status', },
        { text: 'Actions', value: 'action', sortable: false, },

      ],
    }),

    computed:{
      ...mapGetters({
        authuser: 'auth/user',
        accounts: 'account/accounts',
        members: 'account/users',
        subscriptions: 'subscription/subscriptions',
        services: 'service/services',
        servers: 'server/servers',
        meshes: 'mesh/meshes',
      }),
    },

    mounted () {
      this.readAllMeshes()
      this.readSubscriptions(this.authuser.email)
      this.readServices(this.authuser.email)
      this.readServers()

    },

    watch: {
      // whenever accounts changes, this function will run
        //      subscriptions(newSubscriptions, oldSubscriptions) {
          //      this.readSubscriptions(this.authuser.email);
      //}
    },

    methods: {
        ...mapActions('account', {
            readAllAccounts: 'readAll',
            readUsers: 'readUsers',
            createAccount: 'create',
            updateAccount: 'update',
            emailUser: 'email',
        }),

        ...mapActions('subscription', {
            readSubscriptions: 'read',
        }),

        ...mapActions('service', {
            readServices: 'read',
            createService: 'create',
            updateService: 'update',
            deleteService: 'delete',
        }),

        ...mapActions('server', {
            readServers: 'read',
        }),

        ...mapActions('mesh', {
            readAllMeshes: 'readAll',
        }),

      startCreateService() {
        this.dialogCreateService = true;
        this.service = {
          name: "",
          email: this.authuser.email,
        }
        this.meshList = { selected: { "text": "",  "value": ""},
                          items: [] }

        var selected = 0;
        this.meshList.items[0] = { "text": "New Mesh", "value": ""}
        for (let i=0; i<this.meshes.length; i++) {
            this.meshList.items[i+1]= { "text": this.meshes[i].meshName, "value": this.meshes[i].id }
        }

        this.meshList.selected = this.meshList.items[selected];

        this.serverList = { selected: { "text": "",  "value": ""},
                          items: [] }
        for (let i=0; i<this.servers.length; i++) {
            this.serverList.items[i]= { "text": this.servers[i].description, "value": this.servers[i].name }
        }
      },

      create() {

        for (let i=0; i<this.serverList.items.length; i++) {
            if (this.serverList.items[i].value == this.serverList.selected.value) {
                this.server = this.servers[i];
            }
        }

        var range = this.server.portMax - this.server.portMin + 1;
        var port = this.server.portMin + Math.floor(Math.random() * range);

        this.service.servicePort = port;
        this.service.relayHost = {}
        this.service.relayHost.meshName = this.serverList.selected.value;
        this.service.relayHost.current = {}
        this.service.relayHost.current.endpoint = this.server.ipAddress + ":" + port;
        this.service.relayHost.current.listenPort = port;
        this.service.description = this.server.description
        this.service.name = this.server.name
        this.service.serviceGroup = this.server.serviceGroup
        this.service.apiKey = this.server.serviceApiKey

        this.service.serviceType = "relay"

        if (this.service.relayHost.meshName != "") {
            this.service.relayHost.meshId = this.meshList.selected.value;
        }
        else {
            this.service.relayHost.meshId = "";
        }

        this.createService(this.service);

        this.dialogCreateService = false;

      },

      remove(item) {
        this.inDelete = true;
        if (confirm(`Do you really want to delete ${item.name} ?`)){
          this.deleteService(item)
        }
        this.readAllAccounts(this.authuser.email)
        this.readAllMeshes()
        this.readServices(this.authuser.email)

      },

      email(toAddress, mesh) {
        this.dialogCreateService = false;
        if (!toAddress) {
          this.errorUser('email address is not defined')
          return
        }

        for (let i=0; i<this.accounts.length; i++) {
            if (this.accounts[i].id == this.accounts[i].parent) {
                this.emailUser(toAddress + "/" + this.accounts[i].id);
                break;        
            }
        }

      },

      startUpdate(user) {
        if (this.inDelete == true ) {
            this.inDelete = false;
            return
        }
        this.user = user;
        this.dialogUpdate = true;
      },

      update(user) {

        this.dialogUpdate = false;
        this.updateAccount(user)
      },

      updateMember(member) {

        this.dialogMember = false;
        this.member.meshName = this.meshList.selected.text;
        this.member.meshId = this.meshList.selected.value;

        this.updateAccount(member)
      },      

      startUpdateMember(member) {
        if (this.inDelete == true ) {
            this.inDelete = false;
            return
        }

        var selected = 0;
        this.meshList.items = [];
        this.meshList.items[0] = { "text": "All Meshes", "value": ""};
        for (let i=0; i<this.meshes.length; i++) {
            this.meshList.items[i+1]= { "text": this.meshes[i].meshName, "value": this.meshes[i].id }
            if (this.meshes[i].id == member.meshId) {
                selected = i+1;
            }
        }
        this.meshList.selected = this.meshList.items[selected];

        this.member = member;
        this.dialogMember = true;

      },

      startUpdateSubscription(subscription) {
        if (this.inDelete == true ) {
            this.inDelete = false;
            return
        }
        this.subscription = subscription;
        this.dialogSubscription = true;
      },

      updateSubscription(subscription) {

        this.dialogSubscription = false;
        this.updateSubscription(subscription)
      },
    
        startUpdateService(service) {
        if (this.inDelete == true ) {
                this.inDelete = false;
                return
        }
        this.service = service;
        this.dialogService = true;
        },

        updateService(service) {    
            this.dialogService = false;
            this.updateService(service)
        },
        
    }
  };
</script>
