FROM golang AS gobuilder

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential libjpeg-dev libpng-dev libgif-dev libwebp-dev \
    libfontconfig1-dev libfreetype6-dev libgomp1 libexpat1-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/qizhidata.com/zhoukk/kimg

ENV IMAGEMAGICK_VERSION=7.0.8-40

RUN wget -q https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz \
    && tar xf ${IMAGEMAGICK_VERSION}.tar.gz \
    && cd ImageMagick-${IMAGEMAGICK_VERSION} \
    && ./configure && make && make install \
    && ldconfig /usr/local/lib \
    && cd - && rm -rf ImageMagick-*

COPY . .

RUN export CGO_CFLAGS_ALLOW='-fopenmp' && \
    export CGO_CFLAGS="`pkg-config --cflags MagickWand MagickCore`" && \
    export CGO_LDFLAGS="`pkg-config --libs MagickWand MagickCore` \
    -ljpeg -lpng -lwebpmux -lwebp -lfontconfig -lfreetype -lgomp -lexpat -lz -lm -ldl" && \
    go get -d && go install -tags no_pkgconfig -v gopkg.in/gographics/imagick.v3/imagick && \
    export KIMG_TAG="`git describe "--abbrev=0" "--tags"`" && \
    go build -ldflags "-linkmode 'external' -extldflags '-static' -w -s -X 'main.KimgVersion=${KIMG_TAG#*release-}'" -o kimg

FROM node AS nodebuilder

WORKDIR /app

COPY . .

RUN cd web && yarn && yarn build

FROM scratch
COPY --from=gobuilder /go/src/qizhidata.com/zhoukk/kimg/kimg .
COPY --from=gobuilder /go/src/qizhidata.com/zhoukk/kimg/kimg.ini .
COPY --from=nodebuilder /app/www www
EXPOSE 80
ENTRYPOINT ["/kimg"]