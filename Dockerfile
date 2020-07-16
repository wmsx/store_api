FROM alpine
ADD html /html
ADD store_api-web /store_api-web
WORKDIR /
ENTRYPOINT [ "/store_api-web" ]
