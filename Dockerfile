FROM --platform=$BUILDPLATFORM docker.io/library/node:16.17-alpine AS frontend-deps

WORKDIR /frontend/

RUN yarn global add @quasar/cli

COPY ./frontend/package.json ./frontend/yarn.lock /frontend/

RUN yarn

FROM --platform=$BUILDPLATFORM frontend-deps AS build-frontend

WORKDIR /frontend/

COPY ./frontend/ /frontend/

RUN quasar build -m spa

FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.19.1-alpine as build-backend

WORKDIR /backend/

RUN go version

COPY ./backend/go.mod ./backend/go.sum /backend/
RUN go mod download

COPY ./backend/ /backend/

ARG TARGETOS TARGETARCH

RUN chmod +x scripts/Build.sh
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH scripts/Build.sh

COPY --from=build-frontend /frontend/dist/spa/ /backend/frontend-dist/public/original/

FROM scratch

WORKDIR /app/

COPY --from=build-backend /backend/SimpleIPLocation /backend/frontend-dist /app/

CMD ["/app/SimpleIPLocation"]


