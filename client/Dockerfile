# Development stage
FROM node:22-alpine AS development

WORKDIR /app

# Copy package files and install dependencies
COPY client/package*.json ./
RUN npm install

# Copy the rest of the code
COPY client/ ./

# Expose port for development server
EXPOSE 5173

# Command for development with hot reload
CMD ["npm", "run", "dev"]

# Production stage
FROM node:22-alpine AS build

WORKDIR /app

# Copy package files and install dependencies
COPY client/package*.json ./
RUN npm install

# Copy the rest of the code
COPY client/ ./

# Build the application
RUN npm run build

# Serve stage
FROM nginx:alpine AS production

# Copy built assets from build stage
COPY --from=build /app/dist /usr/share/nginx/html

# Copy nginx config if you have a custom one
# COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]