FROM fsouza/fake-gcs-server:1.34.1
RUN apk add --no-cache bash curl
COPY healthcheck.sh /bin/healthcheck.sh
HEALTHCHECK --interval=1s CMD /bin/healthcheck.sh
ENTRYPOINT ["/bin/fake-gcs-server", "-scheme=http", "-port=80", "-backend=memory", "-public-host=cloud-storage"]
