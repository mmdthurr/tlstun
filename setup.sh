#!/bin/bash

sudo cp tlstun /usr/bin 
sudo cp tls.cert / 
sudo cp tls.key / 

echo enter addr 
read addr 

echo enter passwd 
read passwd

sudo cat << EOF > /etc/systemd/system/tlstun.service
[Unit]
Description=tls tunnel
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=tlstun -addr $addr  -passwd $passwd
Restart=always
RestartSec=60s
User=root

[Install]
WantedBy=default.target
EOF 


sudo systemctl daemon-reload
sudo systemctl start tlstun.service 