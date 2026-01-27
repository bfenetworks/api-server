# Copyright 2026 The BFE Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
FROM --platform=${BUILDPLATFORM} golang:1.22.2-alpine3.19 AS build

ARG TARGETARCH
ARG TARGETOS
ARG ALPINE_MIRROR=

WORKDIR /src

RUN set -ex; \
  if [ -n "${ALPINE_MIRROR}" ]; then \
    echo "${ALPINE_MIRROR}/v3.19/main" > /etc/apk/repositories; \
    echo "${ALPINE_MIRROR}/v3.19/community" >> /etc/apk/repositories; \
  fi; \
  apk add --no-cache git ca-certificates curl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOTOOLCHAIN=local
RUN set -ex; \
  GOOS=${TARGETOS:-linux} \
  GOARCH=${TARGETARCH:-$(go env GOARCH)} \
  go build -trimpath -ldflags "-s -w" -o /out/api-server ./

# Dashboard download stage
FROM --platform=${BUILDPLATFORM} alpine:3.19 AS dashboard

ARG ALPINE_MIRROR=
ARG DASHBOARD_VERSION=v0.0.2

RUN set -ex; \
  if [ -n "${ALPINE_MIRROR}" ]; then \
    echo "${ALPINE_MIRROR}/v3.19/main" > /etc/apk/repositories; \
    echo "${ALPINE_MIRROR}/v3.19/community" >> /etc/apk/repositories; \
  fi; \
  apk add --no-cache curl tar; \
  curl -fsSL -o /tmp/dashboard.tar.gz \
    "https://github.com/bfenetworks/dashboard/releases/download/${DASHBOARD_VERSION}/bfe_dashboard_${DASHBOARD_VERSION#v}.tar.gz"; \
  mkdir -p /dashboard; \
  tar -C /dashboard -zxf /tmp/dashboard.tar.gz; \
  rm /tmp/dashboard.tar.gz

# Final runtime stage
FROM alpine:3.19

ARG ALPINE_MIRROR=

RUN set -ex; \
  if [ -n "${ALPINE_MIRROR}" ]; then \
    echo "${ALPINE_MIRROR}/v3.19/main" > /etc/apk/repositories; \
    echo "${ALPINE_MIRROR}/v3.19/community" >> /etc/apk/repositories; \
  fi; \
  apk add --no-cache ca-certificates tzdata; \
  addgroup -S app; \
  adduser -S -G app -u 10001 app

RUN mkdir -p /home/work/api-server/static

WORKDIR /home/work/api-server

COPY --from=build /out/api-server ./api-server
COPY --from=dashboard /dashboard/bfe_dashboard_* /dashboard_tmp
RUN cp -r /dashboard_tmp/* ./static/ && rm -rf /dashboard_tmp
COPY conf ./conf

RUN mkdir -p ./log \
  && chown -R app:app /home/work/api-server

USER app

EXPOSE 8183 8284

ENTRYPOINT ["./api-server"]
CMD ["-c","./conf/","-sc","api_server.toml","-l","./log"]
