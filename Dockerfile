FROM opensuse
COPY src/weightServer /

ENTRYPOINT [ "/weightServer" ]
