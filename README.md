# COS_SP
This is a simplistic implementation of onion routing for study purposes developed in golang.
The project will eventually consist of a client api ,a peer protocol , and a server protocol developed in Go.
The aim of this project is to imitate hidden services which allow both clients and servers to exchange messages while being fully anonymized.
Ideally this implementation would include a DHT to store service descriptors however due to time constraints  a NodeJS  app will act as a registry by using a redis server will be used instead.
