#!/bin/bash -e

PASS_FILE="pass.txt"
YML_FILE="bens.yml"
PUB_KEY="pub.key"
PRI_KEY="pri.key"

if [ -e "$PASS_FILE" ]; then
    echo already initialized in this directory >&2
    exit 1
fi

OPENSSL=$(which openssl)
BASE64=$(which base64)

origMask=$(umask)
umask 077

echo generating pass file...
$OPENSSL rand -base64 -out "$PASS_FILE" 33

echo generating private key...
$OPENSSL genrsa -aes128 -out "$PRI_KEY" -passout "file:$PASS_FILE" 2048

umask $origMask

echo computing public key from private key...
$OPENSSL rsa -passin "file:$PASS_FILE" -in $PRI_KEY \
    -outform pem -out "$PUB_KEY" -pubout

val=$(echo bar | \
    $OPENSSL rsautl -encrypt -inkey "$PUB_KEY" -keyform pem -oaep -pubin | \
    $BASE64)

echo writing yaml template file...
cat > "$YML_FILE" <<EOF
version: 1
environment:
  - name: FOO
    encryptedValue: $val
EOF
