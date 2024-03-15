FROM golang:1.21.3-alpine3.18

WORKDIR gosocket
COPY gosocket.bin .
COPY . .
RUN chmod +x gosocket.bin

EXPOSE 8081
CMD ["./gosocket.bin"]