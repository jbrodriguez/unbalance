#!/bin/bash

usermod -u 99 nobody
usermod -g 100 nobody
usermod -d /home nobody
chown -R nobody:users /home


## Disable SSH
rm -rf /etc/service/sshd
rm /etc/my_init.d/00_regen_ssh_host_keys.sh

## Run unBALANCE
mkdir /etc/service/unbalance
cat <<'EOT' > /etc/service/unbalance/run
#!/bin/bash

cp /root/mdcmd /usr/bin
chown -R nobody:users /usr/local/share/unbalance /config /usr/bin/unbalance /usr/bin/diskmv /usr/bin/mdcmd
chmod +x /usr/bin/unbalance /usr/bin/diskmv
chown -R nobody:users /etc/ssmtp/ssmtp.conf

if [[ -d /log ]]; then
	UNBALANCE_LOGFILEPATH=/log GIN_MODE=release UNBALANCE_DOCKER=y /sbin/setuser nobody /usr/bin/unbalance -c /config
else
	GIN_MODE=release UNBALANCE_DOCKER=y /sbin/setuser nobody /usr/bin/unbalance -c /config	
fi

EOT

chmod -R +x /etc/service/ /etc/my_init.d/

# Dependencies
# apt-get update
# apt-get install -y \
# 		rsync \
# 		wget \
# 		ssmtp

# wget --no-check-certificate https://github.com/jbrodriguez/unbalance/releases/download/0.5.1/unbalance-0.5.1-linux-amd64.tar.gz -O - | tar -xzf - -C /tmp

# ls -al /tmp
# ls -al /tmp/unbalance

mv /tmp/unbalance/unbalance /usr/bin
mv /tmp/unbalance/diskmv /usr/bin
mv /tmp/unbalance /usr/local/share/

## Clean up APT when done.
# apt-get clean -y
rm -rf /var/lib/apt/lists/* /var/cache/* /tmp/* /var/tmp/*