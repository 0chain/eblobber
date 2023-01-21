#!/bin/bash

MIGRATION_ROOT=$HOME/.s3migration
ACCESS_KEY=0chainaccesskey
SECRET_KEY=0chainsecretkey
ALLOCATION=0chainallocation
BUCKET=0chainbucket
# ACCESS_KEY=AKIA4MPQDEZ4OACO45FD
# SECRET_KEY=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
# ALLOCATION=1467bdcedb20e62928eb18960c501c1c50ff35c927b86b34a8f0b00eb29a6d06
# BUCKET=cloud-mig
CONCURRENCY=1
DELETE_SOURCE=false
ENCRYPT=false
REGION=us-east-2
# BLOCK_WORKER_URL=0chainblockworkerurl
BLOCK_WORKER_URL=https://helm.0chain.net/dns
sudo apt update
sudo apt install -y unzip curl containerd docker.io

echo "download docker-compose"
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version

sudo curl -L "https://s3-mig-binaries.s3.us-east-2.amazonaws.com/s3mgrt" -o /usr/local/bin/s3mgrt
chmod +x /usr/local/bin/s3mgrt
CONFIG_DIR=$HOME/.zcn
mkdir -p ${MIGRATION_ROOT}

cat <<EOF >${CONFIG_DIR}/Caddyfile2
blimp76ghf.devnet-0chain.net:9012 {
	route {
		reverse_proxy s3mgrt:8080
	}
}
EOF

cat <<EOF >${CONFIG_DIR}/migration-compose.yml
version: '3.8'
services:
  caddy:
    image: caddy:latest
    ports:
      - 9012:9012
    volumes:
      - ${CONFIG_DIR}/Caddyfile2:/etc/caddy/Caddyfile
      - ${CONFIG_DIR}/caddy/site:/srv
      - ${CONFIG_DIR}/caddy/caddy_data:/data
      - ${CONFIG_DIR}/caddy/caddy_config:/config
    restart: "always"
    
  s3mgrt:
    image: bmanu199/s3mgrt:latest
    restart: always
    volumes:
      - ${MIGRATION_ROOT}:/migrate
EOF

/usr/local/bin/docker-compose -f ${CONFIG_DIR}/migration-compose.yml up -d

#  --concurrency ${CONCURRENCY} --delete-source ${DELETE_SOURCE} --encrypt ${ENCRYPT} --resume true   --skip 1

cd ${MIGRATION_ROOT}
/usr/local/bin/s3mgrt migrate --access-key ${ACCESS_KEY} --secret-key ${SECRET_KEY} --allocation ${ALLOCATION} --bucket ${BUCKET} --region ${REGION}