FROM golang:1.18

WORKDIR /
COPY . .

# Set the entrypoint script
COPY /entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Specify the entrypoint script as the entrypoint
ENTRYPOINT ["/entrypoint.sh"]

# Specify the command to run your Go program
CMD ["go", "run", "main.go"]
