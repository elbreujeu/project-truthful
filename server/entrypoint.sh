#!/bin/bash

# Execute the script to generate keys
bash generate_rsa_keys.sh

# Execute the command passed as arguments
exec "$@"
