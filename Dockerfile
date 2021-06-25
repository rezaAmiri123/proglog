FROM scratch
#RUN CGO_ENABLED=0 go build -o /go/bin/proglog ./cmd/proglog
#RUN CGO_ENABLED=0 go build -o proglog cmd/proglog/main.go
#RUN strip proglog
COPY proglog /bin/proglog
COPY grpc_health_probe /bin/grpc_health_probe
ENTRYPOINT ["/bin/proglog"]
