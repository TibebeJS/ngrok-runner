!#/usr/bin/bash
mkdir -p /usr/local/ngrok-runner
go build
cp ./ngrok-runner /usr/local/ngrok-runner/ngrok-runner
cp ./ngrok-runner.service /etc/systemd/system/ngrok-runner.service
systemctl daemon-reload
systemctl start ngrok-runner