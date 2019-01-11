## deep_healthcheck
負荷分散装置からのヘルスチェックに応答する用途のAPI
GETされると、複数URLとpostgresへヘルスチェックを行う
OKの時はhttp応答コード200、NGの時はhttp応答コード500 で応答。

### オプション
```
コンフィグパス指定
--config-path deep_healthcheck.yml

LISTENポート番号
--listen-address 1234

ヘルスチェックURL
--healthcheck-url /health/check
```

### systemd登録
```
cat <<'EOT' > /etc/systemd/system/deep_healthcheck.service
[Unit]
Description=deep_healthcheck service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/deep_healthcheck/deep_healthcheck
Restart=always

[Install]
WantedBy=multi-user.target
EOT
```

### rsyslog
```
cat <<EOT >/etc/rsyslog.d/deep_healthcheck.conf
:programname, isequal, "deep_healthcheck" /usr/local/deep_healthcheck/log/deep_healthcheck.log
& stop
EOT
systemctl restart rsyslog.service
```
### logrotate
```
cat <<EOT >/etc/logrotate.d/deep_healthcheck
/data/s01/deep_healthcheck/log/deep_healthcheck.log
{
  daily
  rotate 31
  dateext
  compress
  notifempty
  missingok
}
EOT
```
