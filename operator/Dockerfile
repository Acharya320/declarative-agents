# Copyright (c) 2026 Nokia
# SPDX-License-Identifier: BSD-3-Clause
# The published operator image (applications srd009 R4.1). The build context
# is the operator module directory. The kind platform builds the same runtime
# layer from the checkout's binary (magefiles/kindrig/operator.go).
ARG GO_IMAGE=docker.io/library/golang:1.26-alpine@sha256:8ac98ca534ac3f51e1f420a1dd2c15e74c75cfa0f23f3ad27eb5d7236c349a0c
ARG RUNTIME_IMAGE=docker.io/library/alpine:3.22@sha256:5291449c3df73caf6ed85e649dec1b9e818b39a5d8c871e97afc13e9cd5e8fa8

FROM ${GO_IMAGE} AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/agent-operator ./cmd/agent-operator

FROM ${RUNTIME_IMAGE} AS runtime
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/agent-operator /usr/local/bin/agent-operator
USER 65532:65532
ENTRYPOINT ["agent-operator"]
