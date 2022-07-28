FROM debian:bullseye

# Install base required packages
RUN apt-get update && apt-get -yy install openssh-server

RUN mkdir /run/sshd
RUN mkdir /root/.ssh

CMD ["/usr/sbin/sshd", "-D", "-o", "ListenAddress=0.0.0.0"]
