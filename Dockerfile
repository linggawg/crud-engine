FROM golang:1.18

WORKDIR /usr/src/app

# Update package
# RUN apk add --update --no-cache --virtual .build-dev build-base git
RUN apt-get update

COPY . .

RUN make install \
  && make build

# Expose port
EXPOSE 9000

# Run application
CMD ["make", "start"]