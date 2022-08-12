FROM golang:1.16.3-alpine as backend

ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on

WORKDIR /app
COPY ./filestore-api/go.mod go.mod
COPY ./filestore-api/go.sum go.sum
RUN go mod download
COPY ./filestore-api .
RUN go build .

FROM node:16.13.0-alpine as frontend
WORKDIR /app
COPY ./filestore-web/ .
RUN npm install
RUN npm run build

FROM alpine:latest
RUN apk --no-cache add curl nginx openssl
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait
WORKDIR /app

COPY --from=backend /app/summercamp-filestore ./
COPY --from=backend /app/.env.production ./

COPY --from=frontend /app/dist /usr/share/nginx/html
COPY ./docker-script /.docker-script

WORKDIR /.docker-script
EXPOSE 80 8080
#ENTRYPOINT ./summercamp-filestore -e production
CMD ["sh", "/.docker-script/run.sh"]
#CMD ["http-server", "dist", "-p", "8181"]

#CMD ["./summercamp-filestore"]
