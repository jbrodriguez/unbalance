#!/bin/bash

cp /root/mdcmd /usr/bin
chown -R nobody:users /usr/local/share/unbalance /config /usr/bin/unbalance /usr/bin/diskmv /usr/bin/mdcmd
chmod +x /usr/bin/unbalance /usr/bin/diskmv

if [[ -d /log ]]; then
	UNBALANCE_LOGFILEPATH=/log UNBALANCE_DOCKER=y /sbin/setuser nobody /usr/bin/unbalance -c /config
else
	UNBALANCE_DOCKER=y /sbin/setuser nobody /usr/bin/unbalance -c /config	
fi