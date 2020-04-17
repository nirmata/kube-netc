FROM golang:latest
ADD main /
EXPOSE 2112
CMD ["/main"]