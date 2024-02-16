#!/bin/bash

# Execute the script to generate keys
./generate_rsa_keys.sh

# Execute the command passed as arguments
exec "$@"
