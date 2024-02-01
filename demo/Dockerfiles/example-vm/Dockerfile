FROM debian:bookworm

# Install base required packages
RUN apt-get update && apt-get -yy install openssh-server

RUN mkdir -p /run/sshd
RUN mkdir -p /root/.ssh

CMD ["/usr/sbin/sshd", "-D", "-o", "ListenAddress=0.0.0.0"]
