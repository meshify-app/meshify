# Meshify

<h1><img src="./ui/src/assets/meshify.png" alt="A WireGuard control plane"></h1>

A control plane for [WireGuard](https://wireguard.com).

## Requirements

* OIDC compliant OAuth2 implementation
* MongoDB
* Mail Server credentials for sending outgoing email
* golang
* nginx
* NodeJS / Vue 2

![Screenshot](meshify-architecture.png)

## Features

 * Self-hosted and web based management of wireguard networks
 * Mesh define the configuration of the hosts in the network
 * Invite people to network with email
 * Authenticate them with OAuth2
 * Generation of configuration files on demand
 * User authentication (Oauth2 OIDC)
 * Fully configure all aspects of your VPN
 * Manage hosts remotely
 * Simple
 * Lightweight
 * Secure



![Screenshot](meshify-screenshot.png)

## Running


### Directly

Install dependencies

Sample NGINX Config:

```
server {

        server_name meshifyvpn.com;

        root /usr/share/meshify/ui/dist; index index.html; location / {
            try_files $uri $uri/ /index.html;
       }

    location /api/ {
        # app2 reverse proxy settings follow
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header Host localhost;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_pass http://127.0.0.1:8080;
    }


    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/meshifyvpn.com/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/meshifyvpn.com/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

}
server {
    if ($host = meshifyvpn.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot



        server_name meshifyvpn.com;
    listen 80;
    return 404; # managed by Certbot


}
```

Example `.env` file:

```
# IP address to listen to
SERVER=0.0.0.0
# port to bind
PORT=8080
# Gin framework release mode
GIN_MODE=release

# SMTP settings to send email to clients
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=
SMTP_FROM=Meshify <info@meshifyvpn.com>

# MONGO settings
MONGODB_CONNECTION_STRING=mongodb://127.0.0.1:27017

# example with google
#OAUTH2_PROVIDER_NAME=google
#OAUTH2_PROVIDER=
#OAUTH2_CLIENT_ID=
#OAUTH2_CLIENT_SECRET=
#OAUTH2_REDIRECT_URL=

# example with github
#OAUTH2_PROVIDER_NAME=github
#OAUTH2_PROVIDER=https://github.com
#OAUTH2_CLIENT_ID=
#OAUTH2_CLIENT_SECRET=
#OAUTH2_REDIRECT_URL=

#OAUTH2_PROVIDER_NAME=oauth2oidc
#OAUTH2_PROVIDER=https://auth.meshifyvpn.com/
#OAUTH2_PROVIDER_URL=meshifyvpn.us.auth0.com
#OAUTH2_CLIENT_ID=
#OAUTH2_CLIENT_ID_WINDOWS=
#OAUTH2_CLIENT_SECRET=
#OAUTH2_REDIRECT_URL=https://dev.meshifyvpn.com

OAUTH2_PROVIDER_NAME=microsoft
OAUTH2_PROVIDER=https://login.microsoftonline.com/.../v2.0
OAUTH2_CLIENT_ID=
OAUTH2_CLIENT_ID_WINDOWS=
OAUTH2_CLIENT_SECRET=
OAUTH2_REDIRECT_URL=https://meshifyvpn.com
OAUTH2_TENET=...

# set provider name to fake to disable auth, also the default
OAUTH2_PROVIDER_NAME=microsoft
```

Create a systemd service for the API:

```
cat  /lib/systemd/system/meshify-api.service
[Unit]
Description=Meshify API
ConditionPathExists=/usr/share/meshify/cmd/meshify
After=network.target

[Service]
Type=simple
User=root
Group=root
LimitNOFILE=1024000

Restart=on-failure
RestartSec=10
#startLimitIntervalSec=60

WorkingDirectory=/usr/share/meshify/
ExecStart=/usr/share/meshify/cmd/meshify/meshify

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/meshify
ExecStartPre=/bin/chown syslog:adm /var/log/meshify
ExecStartPre=/bin/chmod 755 /var/log/meshify
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=meshify

[Install]
WantedBy=multi-user.target
```

Build the API
```
cd /usr/share/meshify/cmd/meshify
go build
```

Enable the service:

```
sudo systemctl enable meshify-api
sudo systemctl start meshify-api
```

Install NodeJS using NVM
```
nvm use lts-latest
```
cd ui
npm Install
npm run build

With the given nginx config, you should now be able to use your website.  Don't forget
to get a cert using certbot

## Need Help

mailto:support@meshify.app

## License
* Released under MIT License
