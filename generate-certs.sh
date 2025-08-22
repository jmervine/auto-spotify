#!/bin/bash

# Generate SSL certificates for localhost HTTPS server
echo "Generating SSL certificates for localhost..."

# Create certificates directory
mkdir -p certs

# Generate private key
openssl genrsa -out certs/server.key 2048

# Generate certificate signing request
openssl req -new -key certs/server.key -out certs/server.csr -subj "/C=US/ST=Local/L=Local/O=Auto-Spotify/CN=127.0.0.1"

# Create config file for SAN (Subject Alternative Names)
cat > certs/server.conf << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = US
ST = Local
L = Local
O = Auto-Spotify
CN = 127.0.0.1

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = 127.0.0.1
IP.1 = 127.0.0.1
EOF

# Generate self-signed certificate
openssl x509 -req -in certs/server.csr -signkey certs/server.key -out certs/server.crt -days 365 -extensions v3_req -extfile certs/server.conf

# Clean up CSR file
rm certs/server.csr certs/server.conf

echo "✅ Certificates generated in certs/ directory"
echo "   - certs/server.key (private key)"
echo "   - certs/server.crt (certificate)"
echo ""
echo "Now update your Spotify app to use: https://127.0.0.1:8080/callback"
