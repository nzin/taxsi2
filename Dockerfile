######################################
# Prepare go_builder
######################################
FROM golang:1.21 as go_builder
WORKDIR /go/src/github.com/nzin/taxsi2
ADD . .
RUN make build

######################################
# Copy from builder to debian image
######################################
FROM debian:bullseye-slim
RUN mkdir /app
WORKDIR /app

ENV HOST=0.0.0.0
ENV PORT=18000

COPY --from=go_builder /go/src/github.com/nzin/taxsi2/taxsi2 ./taxsi2

RUN useradd --uid 1000 --gid 0 taxsi2 && \
    chown taxsi2:root /app && \
    chmod g=u /app
USER 1000:0

EXPOSE 18000
CMD ./taxsi2
