# Upload - DEV
doppler secrets upload -p skc-deck-api -c dev .env
doppler secrets upload -p skc-deck-api -c dev_docker .env_docker_local

cat certs/skc-suggestion-engine-db.pem | doppler secrets set -p skc-deck-api -c dev MONGODB_X509
cat certs/skc-suggestion-engine-db.pem | doppler secrets set -p skc-deck-api -c dev_docker MONGODB_X509

# upload tls certs
cat certs/certificate.crt | doppler secrets set -p skc-deck-api -c dev SSL_CERTIFICATE
cat certs/certificate.crt | doppler secrets set -p skc-deck-api -c dev_docker SSL_CERTIFICATE

cat certs/private.key | doppler secrets set -p skc-deck-api -c dev SSL_PRIVATE_KEY
cat certs/private.key | doppler secrets set -p skc-deck-api -c dev_docker SSL_PRIVATE_KEY

cat certs/ca_bundle.crt | doppler secrets set -p skc-deck-api -c dev SSL_CA_BUNDLE
cat certs/ca_bundle.crt | doppler secrets set -p skc-deck-api -c dev_docker SSL_CA_BUNDLE

#######################
# Upload - Prod
doppler secrets upload -p skc-deck-api -c prod .env_prod
cat certs/skc-suggestion-engine-db.pem | doppler secrets set -p skc-deck-api -c prod MONGODB_X509

# upload tls certs
cat certs/certificate.crt | doppler secrets set -p skc-deck-api -c prod SSL_CERTIFICATE
cat certs/private.key | doppler secrets set -p skc-deck-api -c prod SSL_PRIVATE_KEY
cat certs/ca_bundle.crt | doppler secrets set -p skc-deck-api -c prod SSL_CA_BUNDLE
