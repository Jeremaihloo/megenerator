###################Start Script#################  
#!/bin/bash  
# build.sh to build go application for multi platform
# created by http://github.com/jeremaihloo/go-build-sh at {{ .CreateAt }}

if [ ! -d "dist" ]; then
    mkdir dist
fi

windows(){
    # windows
    export GOARCH=386
    export GOOS=windows
    go build -o {{ .Name }}.exe
    if [ ! -d "dist/windows" ]; then
        mkdir dist/windows
    fi
    mv {{ .Name }}.exe dist/windows/
    mv conf dist/windows/
    mv views dist/windows/
}

linux(){
    # linux
    export GOARCH=amd64
    export GOOS=linux
    go build -o {{ .Name }}
    if [ ! -d "dist/linux" ]; then
        mkdir dist/linux
    fi
    mv {{ .Name }} dist/linux/
    mv conf dist/linux/
    mv views dist/linux/
}

darwin(){
    # darwin
    export GOARCH=amd64
    export GOOS=darwin
    go build -o {{ .Name }}
    if [ ! -d "dist/darwin" ]; then
        mkdir dist/darwin
    fi
    mv {{ .Name }} dist/darwin/
    mv conf dist/darwin/
    mv views dist/darwin/
}

all(){
    windows
    linux
    darwin
}

case "$1" in
all)
    all    
    ;;
windows)
    windows
    ;;
linux)
    linux
    ;;
darwin)
    darwin
    ;;
*)
    echo $"Usage: $0 {all|windows|darwin|linux}"
    exit 1
esac

exit 0
#####################End Script##################