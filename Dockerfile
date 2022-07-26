FROM gcr.io/distroless/static
ARG TARGETARCH
COPY ./_output/${TARGETARCH}/aurora-tracker /bin/

LABEL maintainers="nilekh"

ENTRYPOINT ["aurora-tracker"]
