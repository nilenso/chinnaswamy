FROM golang:1.18-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go build -o /out/chinnaswamy cmd/chinnaswamy.go

FROM scratch as bin
COPY --from=build /out/chinnaswamy /
CMD ["/chinnaswamy"]