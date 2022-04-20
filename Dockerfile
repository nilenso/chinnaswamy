FROM golang:1.18-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /out/chinnaswamy cmd/chinnaswamy.go

FROM scratch as bin
COPY --from=build /out/chinnaswamy /
CMD ["/chinnaswamy"]