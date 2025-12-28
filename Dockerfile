# GoReleaser Dockerfile
# The binary is pre-built by GoReleaser and copied into the image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary (GoReleaser will inject this)
COPY tools .

# Create config directory
RUN mkdir -p /root/.config/tools

# Set the binary as entrypoint
ENTRYPOINT ["./tools"]

# Default command shows help
CMD ["--help"]
