# Use an official Node.js runtime as the base image
FROM node:latest

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json files
COPY package*.json ./

# Install dependencies during build (not copy)
RUN npm install

# Copy the rest of the application code
COPY . .

# Specify the command to run the React app
CMD [ "npm", "start" ]

# Expose the port for the React app to listen on
EXPOSE 3000
