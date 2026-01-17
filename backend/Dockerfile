FROM debian:bookworm-slim

WORKDIR /app

# Install SQLite libraries if needed
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0

COPY ./main .

# Ensure executable
RUN chmod +x ./main

EXPOSE 8080
CMD ["./main"]