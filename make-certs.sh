#!/bin/sh
set -euo pipefail

CERT_EXPIRY=788400h # 90 years = 788,400 hours

generate_certificate_authority() {
  echo '{"CN":"CA","key":{"algo":"ecdsa","size":256}}' \
    | cfssl gencert -initca - \
    | cfssljson -bare ca -
  cat > ca-config.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "$CERT_EXPIRY",
      "usages": [
        "signing",
        "key encipherment",
        "server auth",
        "client auth"
      ]
    }
  }
}
EOF
}

generate_server_cert() {
  local addresses="$1"
  local name=server
  echo '{"CN":"'$name'","hosts":[""],"key":{"algo":"ecdsa","size":256}}' \
    | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem \
      -hostname="$addresses" - \
    | cfssljson -bare $name
}

generate_client_cert() {
  local name="$1"
  echo '{"CN":"'$name'","hosts":[""],"key":{"algo":"ecdsa","size":256}}' \
    | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem \
      -hostname="" - \
    | cfssljson -bare "$name"
}

main() {
  mkdir -p certs && pushd certs
  generate_certificate_authority
  generate_server_cert "localhost,127.0.0.1"
  generate_client_cert "service-1234@accounts.example.com"
  generate_client_cert "service-3456@accounts.example.com"
  popd

  mkdir -p untrusted-certs && pushd untrusted-certs
  generate_certificate_authority
  generate_client_cert "service-4567@accounts.example.com"
  popd

  cp ./untrusted-certs/service-* ./certs
  rm **/*.csr
  rm certs/*.json
  rm certs/ca-key.pem
  rm -Rf untrusted-certs
}

main "$@"
