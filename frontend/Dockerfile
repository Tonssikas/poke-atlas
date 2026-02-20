FROM node:24-alpine
WORKDIR /app

COPY build ./build

EXPOSE 3000
ENV API_ADDRESS=backend:8080
CMD [ "node", "build" ]