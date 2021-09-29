#!/usr/bin/bash
mkdir -p /usr/share/ngrok-runner
mkdir -p /var/log/ngrok-runner
chown syslog:syslog /var/log/ngrok-runner

mkdir -p /etc/ngrok-runner
cp ./sample.config.yaml /etc/ngrok-runner/config.yaml

cat <<EOT >> /etc/rsyslog.d/ngrok-runner.conf
if $programname == 'ngrok-runner-service' then /var/log/ngrok-runner/output.log
& stop
EOT

cp ./ngrok-runner /usr/share/ngrok-runner/ngrok-runner
cp ./ngrok /usr/share/ngrok-runner/ngrok
cp ./start_ngrok.sh /usr/share/ngrok-runner/start_ngrok.sh
chmod +x /usr/share/ngrok-runner/start_ngrok.sh
rm /etc/systemd/system/ngrok-runner.service
systemctl daemon-reload
cp ./ngrok-runner.service /etc/systemd/system/ngrok-runner.service
systemctl daemon-reload
systemctl enable ngrok-runner
echo "
[+] Successfully installed ngrok-runner.

Use the following commands to:
    start the service:
        - $ systemctl start ngrok-runner

    stop the service:
        - $ systemctl stop ngrok-runner
    
    enable auto-start on system boot:
        - $ systemctl enable ngrok-runner
"