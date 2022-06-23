# FROM python:3-alpine
# FROM alpine:3.11
FROM ubuntu:18.04
MAINTAINER endless pancake4u <endless.pnacake4u@gmail.com> 

ARG BUILD_DATE
ARG VCS_REF

# Set labels (see https://microbadger.com/labels)
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/EndlessPancake/python3-alpine-traceroute-docker"

# set for *.html
# ENV FILE_STORE="emptyDir"

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# RUN apk --update --no-cache add tcptraceroute python3 pythn3-dev \
RUN apt-get -y update \
    && apt-get -y upgrade \
    && apt-get install -y locales curl wget tcpdump graphviz python3-dev python3-distutils \
    && curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py \
    && python3 get-pip.py \
    && pip install -U pip \
    && mkdir /code \
    && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8

COPY requirements.txt /usr/src/app/
RUN pip install --no-cache-dir -r requirements.txt

COPY . /usr/src/app

# CMD python main.py $FILE_DIR
# CMD ["/bin/sh"] 
CMD ["python3", "main.py"] 
