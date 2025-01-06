# Build OpenRCT2
FROM ubuntu:24.04 AS build-env
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
 && apt-get install --no-install-recommends -y git cmake pkg-config ninja-build clang nlohmann-json3-dev libcurl4-openssl-dev libcrypto++-dev libfontconfig1-dev libfreetype6-dev libpng-dev libzip-dev libsdl2-dev libspeexdsp-dev libssl-dev libicu-dev libflac-dev libvorbis-dev \
 && rm -rf /var/lib/apt/lists/*

ARG OPENRCT2_REF=v0.4.17
WORKDIR /openrct2
RUN git -c http.sslVerify=false clone --depth 1 -b $OPENRCT2_REF https://github.com/OpenRCT2/OpenRCT2 . \
 && mkdir build \
 && cd build \
 && cmake .. -G Ninja -DCMAKE_BUILD_TYPE=release -DCMAKE_INSTALL_PREFIX=/openrct2-install/usr -DENABLE_HEADERS_CHECK=OFF \
 && ninja -k0 install \
 && rm /openrct2-install/usr/lib/libopenrct2.a

# Build runtime image
FROM ubuntu:24.04
# Install OpenRCT2
COPY --from=build-env /openrct2-install /openrct2-install
RUN apt-get update \
 && apt-get install --no-install-recommends -y rsync ca-certificates libpng16-16 libzip4 libcurl4 libfreetype6 libfontconfig1 libicu74 libsdl2-2.0-0 libgl1 libspeexdsp1 libvorbisfile3 x11vnc \
 && rm -rf /var/lib/apt/lists/* \
 && rsync -a /openrct2-install/* / \
 && rm -rf /openrct2-install \
 && openrct2 --version

# Set up ordinary user
RUN useradd -m openrct2
USER openrct2
WORKDIR /home/openrct2

EXPOSE 11753
EXPOSE 5901

# Test run and scan
RUN openrct2 --version \
 && openrct2 scan-objects

COPY entrypoint.sh /bin/entrypoint.sh
COPY lib/socket.js /home/openrct2/.config/OpenRCT2/plugin/socket.js

# Done
ENTRYPOINT ["/bin/entrypoint.sh"]