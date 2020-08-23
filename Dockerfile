FROM golang:1.14 AS builder

WORKDIR /workdir

# Assuming the source code is collocated to this Dockerfile, copy the whole
# directory into the container that is building the Docker image.
COPY . .
RUN make status static

# Create a "nobody" non-root user for the next image by crafting an /etc/passwd
# file that the next image can copy in. This is necessary since the next image
# is based on scratch, which doesn't have adduser, cat, echo, or even sh.
RUN \
    mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

FROM scratch
# Copy the static binary from the builder stage
COPY --from=builder /workdir/kubedrainer /usr/local/bin/kubedrainer

# Copy the certs from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the "nobody" non-root user as a security best practice.
COPY --from=builder /user/group /user/passwd /etc/

# Run as non-root by default
USER nobody:nobody

ENTRYPOINT ["kubedrainer", "serve"]
