# Multi-stage build; final image runs as a non-root UID (Choreo requires 10000-20000).
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o /bin/delay-endpoint .

FROM gcr.io/distroless/static-debian12
COPY --from=build /bin/delay-endpoint /bin/delay-endpoint
# Choreo mandates a non-root user in the 10000-20000 range.
USER 10001
EXPOSE 8080
ENTRYPOINT ["/bin/delay-endpoint"]
