FROM alpine
RUN apk add --no-cache tini
COPY gsoc2 /bin/gsoc2
ENTRYPOINT ["/sbin/tini", "--", "/bin/gsoc2"]