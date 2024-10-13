FROM ubuntu:20.04 as base
MAINTAINER Chen Chen <lukas.c.chen@gmail.com>

RUN apt-get update
RUN apt-get -f install
RUN dpkg --configure -a
RUN apt-get -y install python3-pip
RUN python3 -m pip install pillow
WORKDIR /app
COPY *.py ./ 
COPY test.jpg ./
CMD ["python3", "./server.py"]

