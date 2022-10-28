# STEP 1 build executable binary
FROM alpine:latest as builder
RUN apk add go 
#git
WORKDIR /linx/
COPY ./ ./
#RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o "/go/bin/linx" .
RUN mkdir -p ./data

# STEP 2 build a small image
FROM scratch
COPY --from=builder "/go/bin/linx" "/linx"
COPY --from=builder "/linx/data" "/data"
#EXPOSE 8080
ENTRYPOINT ["/linx"]