FROM scratch
COPY server /
ENTRYPOINT ["/server"]
