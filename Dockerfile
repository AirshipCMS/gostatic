FROM centurylink/ca-certs

ADD gostatic /

EXPOSE 80

ENTRYPOINT ["/gostatic"]
