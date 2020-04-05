apt-get update && apt-get install -y curl
mkdir downloads
cd downloads &&

cat > urls << EOF
https://dist.ipfs.io/go-ipfs/v0.4.23/go-ipfs_v0.4.23_linux-amd64.tar.gz
https://dist.ipfs.io/ipfs-cluster-ctl/v0.12.1/ipfs-cluster-ctl_v0.12.1_linux-amd64.tar.gz
https://dist.ipfs.io/ipfs-cluster-follow/v0.12.1/ipfs-cluster-follow_v0.12.1_linux-amd64.tar.gz
https://dist.ipfs.io/ipfs-cluster-service/v0.12.1/ipfs-cluster-service_v0.12.1_linux-amd64.tar.gz
EOF

for url in $(cat urls);do curl -LO ${url};done

for f in $(ls *.tar.gz);do tar zxvf ${f} --strip-components=1;done

for b in ipfs ipfs-cluster-ctl ipfs-cluster-follow ipfs-cluster-follow;do
  cp ${b} /usr/local/bin/
done

mkdir /mnt/ipfs-data
mount /dev/sda /mnt/ipfs-data
useradd -m \
  --home-dir /mnt/ipfs-data/ipfs \
  --shell /bin/bash \
  ipfs

sudo -u ipfs ipfs init  --profile server

cat <<EOT >> /etc/systemd/system/ipfs.service
[Unit]
Description=IPFS daemon
After=network.target
[Service]
Type=simple
User=ipfs
ExecStart=/usr/local/bin/ipfs daemon
[Install]
WantedBy=multiuser.target
EOT

systemctl enable ipfs
systemctl start ipfs

## !! Mount doesn't persist reboots