FROM ubuntu:16.04

# system basics
RUN apt-get update && \
  apt-get -y install software-properties-common && \
  add-apt-repository -y ppa:gophers/archive && \
  apt-get update && \
  apt-get -y --no-install-recommends install \
    build-essential \
    curl \
    python3 \
    python3-dev \
    python3-setuptools \
    python3-pip \
    libffi-dev \
    nginx \
    golang-1.11-go \
    git && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /
ENV GOPATH /go
ENV PATH ${PATH}:/usr/lib/go-1.11/bin:/mkv

COPY requirements.txt mkv/requirements.txt
RUN pip3 install --no-cache-dir -r mkv/requirements.txt

COPY *.go rtest.sh bringup.sh master volume test.py mkv/
WORKDIR /mkv
