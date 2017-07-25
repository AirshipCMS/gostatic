FROM centurylink/ca-certs

COPY gostatic /
ENTRYPOINT ["/gostatic"]
