FROM node:slim

RUN echo 'deb http://emdebian.org/tools/debian/ jessie main' > /etc/apt/sources.list.d/crosstools.list
RUN dpkg --add-architecture armhf
RUN apt-get update
RUN apt-get install -y --force-yes git build-essential crossbuild-essential-armhf

ADD https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz /golang/go1.7.1.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf /golang/go1.7.1.linux-amd64.tar.gz
ENV GOPATH=/app
ENV PATH=$PATH:/usr/local/go/bin:/app/bin

RUN npm -g install bower gulp

COPY package.json /app/src/github.com/mobyos/mobyos-admin-app/package.json
WORKDIR /app/src/github.com/mobyos/mobyos-admin-app
RUN npm install

COPY bower.json /app/src/github.com/mobyos/mobyos-admin-app/bower.json
RUN bower install --allow-root

COPY server /app/src/github.com/mobyos/mobyos-admin-app/server
WORKDIR /app/src/github.com/mobyos/mobyos-admin-app/server
RUN go get ./...
RUN CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build -v -o server 

WORKDIR /app/src/github.com/mobyos/mobyos-admin-app/server/utils
RUN ./run.sh

WORKDIR /app/src/github.com/mobyos/mobyos-admin-app
COPY gulpfile.js /app/src/github.com/mobyos/mobyos-admin-app/gulpfile.js
COPY src /app/src/github.com/mobyos/mobyos-admin-app/src
RUN gulp prepare

RUN mkdir -p /dst/server
RUN cp server/server /dst/server
RUN cp -R server/db /dst/server/db
RUN cp -R src /dst/src

COPY Dockerfile.run /dst/Dockerfile

WORKDIR /dst
CMD tar -cf - . 
