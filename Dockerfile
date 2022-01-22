FROM moby/buildkit:v0.9.3
WORKDIR /pay
COPY pay README.md /pay/
ENV PATH=/pay:$PATH
ENTRYPOINT [ "/bhojpur/pay" ]