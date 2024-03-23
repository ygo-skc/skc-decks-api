mkdir -p certs

aws secretsmanager get-secret-value --secret-id "/prod/skc/deck-api/ssl" --region us-east-2 \
  | jq -r '.SecretString' \
  | jq -r "with_entries(select(.key | startswith(\"SSL\")))" > certs/base64-certs-json
jq -r ".SSL_PRIVATE_KEY" < certs/base64-certs-json | base64 -d > certs/private.key
jq -r ".SSL_CA_BUNDLE_CRT" < certs/base64-certs-json | base64 -d > certs/ca_bundle.crt
jq -r ".SSL_CERTIFICATE_CRT" < certs/base64-certs-json | base64 -d > certs/certificate.crt

aws secretsmanager get-secret-value --secret-id "/prod/skc/deck-api/db" --region us-east-2 \
  | jq -r '.SecretString'  > certs/base64-certs-json
jq -r ".DB_PEM" < certs/base64-certs-json | base64 -d > certs/skc-deck-api-db.pem

rm certs/base64-certs-json

#############################################
aws secretsmanager get-secret-value --secret-id "/prod/skc/deck-api/env" --region us-east-2 \
  | jq -r '.SecretString' | jq -r "to_entries|map(\"\(.key)=\\\"\(.value|tostring)\\\"\")|.[]" | tee .env .env_docker_local .env_prod

