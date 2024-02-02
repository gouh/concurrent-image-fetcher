# Etapa de compilación para el backend de Go
FROM golang:1.20 AS go-builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
COPY ./.env /root/.env
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-binary ./cmd/server/main.go

# Etapa de compilación para el frontend de React
FROM node:alpine as react-builder
WORKDIR /app
COPY front/package.json front/package-lock.json ./
RUN npm ci
COPY front/ ./
RUN npm run build

# Etapa final para ejecutar la aplicación
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copia el binario compilado de Go desde la etapa de compilación de Go
COPY --from=go-builder /root/.env /root/.env
COPY --from=go-builder /go-binary .
# Copia los archivos estáticos compilados de React desde la etapa de compilación de React
COPY --from=react-builder /app/build ./web
# Asegúrate de que tu servidor Go sirve los archivos desde /root/public

# Expone el puerto en el que tu servidor Go escucha
EXPOSE 8080
CMD ["./go-binary"]
