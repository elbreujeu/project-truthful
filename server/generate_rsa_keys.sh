#!/bin/bash

KEY_DIR="/cert"
PRIVATE_KEY="$KEY_DIR/id_rsa"
PUBLIC_KEY="$KEY_DIR/id_rsa.pub"

if [ ! -f "$PRIVATE_KEY" ] || [ ! -f "$PUBLIC_KEY" ]; then
    echo "RSA keys not found. Generating..."
    mkdir -p "$KEY_DIR"
    openssl genrsa -out "$PRIVATE_KEY" 4096
    openssl rsa -in "$PRIVATE_KEY" -pubout -out "$PUBLIC_KEY"
    echo "RSA keys generated successfully."
else
    echo "RSA keys already exist."
fi
