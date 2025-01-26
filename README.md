# toonai.go

# start app

```bash
go run main.go
```

# install dependencies

```bash
go get github.com/gorilla/mux
```

# run app

```bash
go run main.go
```

#######

# install air

```bash
go install github.com/air-verse/air@latest
```

# install

```bash
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

# starting without alias

```bash
$(go env GOPATH)/bin/air
```

# setting up alias in .bashrc

```bash
cd ~/.bashrc
vi .bashrc
alias air='$(go env GOPATH)/bin/air'
:wq
```
