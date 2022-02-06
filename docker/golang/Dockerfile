FROM golang:latest

# コンテナ作業ディレクトリの変更
WORKDIR /go/src/kcjweo/go
# ホストOSの ./go の中身を作業ディレクトリに追加
ADD ./go .

#GOPATHの設定
RUN go mod init github.com/kcjweo/gowithpostgres

# パッケージのインストール
RUN go get github.com/lib/pq
RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq

EXPOSE 8080 8080