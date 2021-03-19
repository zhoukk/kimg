FROM golang AS gobuilder

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential libjpeg-dev libpng-dev libgif-dev libwebp-dev \
    libfontconfig1-dev libfreetype6-dev libgomp1 libexpat1-dev && \
    rm -rf /var/lib/apt/lists/*

ENV IMAGEMAGICK_VERSION=7.0.8-59

RUN wget -q https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz \
    && tar xf ${IMAGEMAGICK_VERSION}.tar.gz \
    && cd ImageMagick-${IMAGEMAGICK_VERSION} \
    && ./configure && make && make install \
    && ldconfig /usr/local/lib \
    && cd - && rm -rf ImageMagick-*

RUN go get github.com/zhoukk/kimg

WORKDIR /go/src/github.com/zhoukk/kimg

RUN export CGO_CFLAGS_ALLOW="-fopenmp" && \
    export CGO_CFLAGS="`pkg-config --cflags MagickWand MagickCore`" && \
    export CGO_LDFLAGS="`pkg-config --libs MagickWand MagickCore` \
    -ljpeg -lpng -lwebpmux -lwebp -lfontconfig -lfreetype -lgomp -lexpat -luuid -lz -lm -ldl" && \
    go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct && \
    go install -tags no_pkgconfig -v gopkg.in/gographics/imagick.v3/imagick@latest && \
    export KIMG_TAG="`git describe "--abbrev=0" "--tags"`" && \
    go build -tags netgo -ldflags "-linkmode 'external' -extldflags '-static' -w -s -X 'main.KimgVersion=${KIMG_TAG}'" -o kimg main/kimg.go

FROM node AS nodebuilder

WORKDIR /web

COPY --from=gobuilder /go/src/github.com/zhoukk/kimg/web .

RUN yarn && yarn build

FROM scratch
COPY --from=gobuilder /go/src/github.com/zhoukk/kimg/kimg .
COPY --from=gobuilder /go/src/github.com/zhoukk/kimg/kimg.yaml .
COPY --from=gobuilder /go/src/github.com/zhoukk/kimg/LICENSE .
COPY --from=nodebuilder /www www
EXPOSE 80
ENTRYPOINT ["/kimg"]
