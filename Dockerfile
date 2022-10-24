FROM node:16 AS nodebuilder

WORKDIR /app

COPY web .

RUN yarn && yarn build

FROM golang AS gobuilder

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential libjpeg-dev libpng-dev libgif-dev libwebp-dev \
    libfontconfig1-dev libfreetype6-dev libgomp1 libexpat1-dev

ENV IMAGEMAGICK_VERSION=7.1.0-47

RUN wget -q https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz \
    && tar xf ${IMAGEMAGICK_VERSION}.tar.gz \
    && cd ImageMagick-${IMAGEMAGICK_VERSION} \
    && ./configure --enable-shared=no --with-zip=no && make && make install \
    && ldconfig /usr/local/lib

WORKDIR /app

COPY . .
COPY --from=nodebuilder /app/dist /app/web/dist

RUN export CGO_CFLAGS_ALLOW="-fopenmp" && \
    export CGO_CFLAGS="`pkg-config --cflags MagickWand MagickCore`" && \
    export CGO_LDFLAGS="`pkg-config --libs MagickWand MagickCore` \
    -ljpeg -lpng -lwebp -lwebpmux -lwebpdemux \
    -lbrotlienc -lbrotlidec -lbrotlicommon \
    -lfontconfig -lfreetype -lgomp -lexpat -luuid -lz -lm -ldl" && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go build -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick && \
    export KIMG_TAG="`git describe "--abbrev=0" "--tags"`" && \
    go build -tags netgo -ldflags "-linkmode 'external' -extldflags '-static' -w -s -X 'main.KimgVersion=${KIMG_TAG}'" -o kimg main/kimg.go

FROM scratch
COPY --from=gobuilder /app/kimg .
COPY --from=gobuilder /app/kimg.yaml .
COPY --from=gobuilder /app/fonts/arial.ttf .
COPY --from=gobuilder /app/LICENSE .
EXPOSE 80
ENTRYPOINT ["/kimg"]

