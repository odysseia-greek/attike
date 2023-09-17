# build stage
FROM node:lts-alpine as build-stage
WORKDIR /app
COPY package*.json ./
ENV GRAPHQL_URL=http://euripides:8080/graphql
RUN npm install
COPY . .
RUN npm run build

# production stage
FROM nginx:stable-alpine as runtime
RUN rm -rf /etc/nginx/conf.d
COPY conf /etc/nginx
COPY --from=build-stage /app/public /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]