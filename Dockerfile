FROM opensuse
COPY weightServer /

ENTRYPOINT [ "/weightServer" ]
