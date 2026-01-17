FROM debian:bookworm-slim

WORKDIR /app

# Install dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    libsqlite3-0 \
 && rm -rf /var/lib/apt/lists/*

COPY ./main .
COPY ./.prod.env ./.env

# Ensure executable
RUN chmod +x ./main

EXPOSE 8080
CMD ["./main"]