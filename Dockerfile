FROM centurylink/ca-certs

ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ADD gostatic /

EXPOSE 80

ENTRYPOINT ["/gostatic"]
