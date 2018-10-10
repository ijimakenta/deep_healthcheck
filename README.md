# deep_healthcheck
負荷分散装置からのヘルスチェックに応答するAPI


# CentOS7
・systemd
cat <<'EOT' > /etc/systemd/system/deep_healthcheck.service
[Unit]
Description=deep_healthcheck service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/deep_healthcheck/deep_healthcheck \
      --config-path /usr/local/deep_healthcheck/deep_healthcheck.yml
Restart=always

[Install]
WantedBy=multi-user.target
EOT


・ログ
cat <<EOT >/etc/rsyslog.d/deep_healthcheck.conf
:programname, isequal, "deep_healthcheck" /var/log/deep_healthcheck.log
& stop
EOT

・ログローテーション
sed -i '1i\/var\/log\/deep_healthcheck\.log' /etc/logrotate.d/syslog
