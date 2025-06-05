cp -p ./bserver /usr/local/bin/bserver
cp ./bserver.service /etc/systemd/system

systemctl enable bserver
systemctl start bserver