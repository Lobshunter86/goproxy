# goproxy
A secure proxy base on SOCKS5 and QUIC.

Client side listen for new connections, forward them to remote server via QUIC. 
Server act as a SOCKS5 server but handles quic connection instead of TCP connetcion. <br>
Also, server use client certificate for client authentication.

# Requirements
Both client and server side need certificate & private key.

# MISC
The official QUIC impletement only support CUBIC congestion control algo now. 
I replace official module to edit default cwnd parameters for better performance, you should do it on your own or just use official implementaion.


# TODO
- Add test code
- implement http proxy from scratch then change Handler interface