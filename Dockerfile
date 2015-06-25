FROM scratch
COPY artifacts/server /server
ENTRYPOINT ["/server"]
