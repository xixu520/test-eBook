# Build stage
FROM node:20-alpine AS builder

# Set NPM registry for speed in China
RUN npm config set registry https://registry.npmmirror.com

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy the rest of the source code
COPY . .

# Build the Vite project
RUN npm run build

# Production stage
FROM nginx:stable-alpine

# Copy dist from builder
COPY --from=builder /app/dist /usr/share/nginx/html
# Copy custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
