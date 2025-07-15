FROM node:24 AS build

WORKDIR /app

COPY app/frontend/package.json ./
COPY app/frontend/package-lock.json ./
RUN npm install
COPY app/frontend/. ./
RUN npm run build

FROM golang:1.24.5-bookworm AS prod

WORKDIR /app
ENV GOCACHE=/app/.cache
ARG TARGETARCH
ENV GOARCH=$TARGETARCH

RUN chown -R 1001:1001 /app
USER 1001
COPY --chown=1001:1001 app/go.mod app/*.go /app/
COPY --from=build --chown=1001:1001 /app/dist /app/frontend/dist

RUN GOOS=linux GOARCH=$TARGETARCH go build -o pomodoro
CMD ["./pomodoro"]
