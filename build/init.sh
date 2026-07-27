#!/bin/bash

chown paas:paas /optemg/csplog/0/gids
chown paas:paas /opt/csplog/0/gids
chown paas:paas -R /opt/csp/gids
if ! grep -q "nameserver 8.8.8.8" /etc/resolv.conf; then
    echo "nameserver 8.8.8.8" >> /etc/resolv.conf
fi