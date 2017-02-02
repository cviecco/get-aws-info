# This is how we want to name the binary output
BINARY=get_aws_info

# These are the values we want to pass for Version and BuildTime
VERSION=0.1.0
#BUILD_TIME=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
#LDFLAGS=-ldflags "-X github.com/ariejan/roll/core.Version=${VERSION} -X github.com/ariejan/roll/core.BuildTime=${BUILD_TIME}"

all:
	go build -ldflags "-X main.Version=${VERSION}"  -o ${BINARY} main.go

clean:
	rm -f ${BINARY}
	rm -f *.gz

tar:
	mkdir ${BINARY}-${VERSION}
	cp main.go LICENSE get_aws_info.spec ${BINARY}-${VERSION}/
	tar -cvzf ${BINARY}-${VERSION}.tar.gz ${BINARY}-${VERSION}/
	rm -rf ${BINARY}-${VERSION}/	
