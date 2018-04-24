FROM golang

ADD  . /go/src/github.com/DanielDanteDosSantosViana/gorabbit
WORKDIR /go/src/github.com/DanielDanteDosSantosViana/gorabbit

#RUN curl https://glide.sh/get | sh
#RUN glide install github.com/DanielDanteDosSantosViana/gorabbit
RUN go install github.com/DanielDanteDosSantosViana/gorabbit/cmd/gorabbitd

ENV PORT_ENV=2222

ENTRYPOINT /go/bin/gorabbitd

EXPOSE 2222