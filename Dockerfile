FROM gcr.io/distroless/static-debian13:nonroot

COPY sqlc-wizard /sqlc-wizard

USER 65532:65532

ENTRYPOINT ["/sqlc-wizard"]
CMD ["--help"]
