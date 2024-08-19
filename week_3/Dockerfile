FROM ubuntu:20.04
COPY webook /app/webook
WORKDIR /app
RUN apt-get update && apt-get install -y curl
CMD ["/app/webook"]