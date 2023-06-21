FROM alpine:latest
WORKDIR /root
COPY ./rbq_anonymous_bot ./rbq_anonymous_bot
COPY ./config.json ./config.json
RUN chmod +x ./rbq_anonymous_bot
ENTRYPOINT ["./rbq_anonymous_bot"]
HEALTHCHECK --interval=60s --timeout=10s CMD cat 'healthcheck.lock' || exit 1
