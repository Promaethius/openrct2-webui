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
 && apt-get install --no-install-recommends -y rsync ca-certificates libpng16-16 libzip4 libcurl4 libfreetype6 libfontconfig1 libicu74 libsdl2-2.0-0 libgl1 libspeexdsp1 libvorbisfile3 libc6 libgcc-s1 gcc-14-base libcurl4t64 libbrotli1 libgssapi-krb5-2 libcom-err2 libk5crypto3 libkrb5support0 libkrb5-3 libkeyutils1 libssl3t64 libidn2-0 libunistring5 libldap2 libgnutls30t64 libgmp10 libhogweed6t64 libnettle8t64 libp11-kit0 libffi8 libtasn1-6 libsasl2-2 libsasl2-modules-db libdb5.3t64 libnghttp2-14 libpsl5t64 librtmp1 libgnutls30 libhogweed6 libnettle8 zlib1g libssh-4 libssl3 libzstd1 libdiscord-rpc3 libstdc++6 libduktape207 libflac12t64 libogg0 libfontconfig1 fontconfig-config fonts-croscore fonts-dejavu-core fonts-dejavu-mono fonts-freefont-otf fonts-freefont-ttf fonts-liberation fonts-noto-core fonts-noto-mono fonts-texgyre fonts-urw-base35 xfonts-utils libfontenc1 libfreetype6 libbz2-1.0 libpng16-16t64 x11-common lsb-base sysvinit-utils xfonts-encodings libexpat1 libgzstream0 libicu74 libsdl2-2.0-0 libasound2t64 libasound2-data libdecor-0-0 libwayland-client0 libdrm2 libdrm-common libgbm1 libwayland-server0 libxcb-randr0 libxcb1 libxau6 libxdmcp6 libbsd0 libmd0 libpulse0 libapparmor1 libasyncns0 libdbus-1-3 libsystemd0 libcap2 libgcrypt20 libgpg-error0 liblz4-1 liblzma5 libsndfile1 libmp3lame0 libmpg123-0t64 libopus0 libvorbis0a libvorbisenc2 libx11-6 libx11-data libx11-xcb1 libsamplerate0 libwayland-cursor0 libwayland-egl1 libxcursor1 libxfixes3 libxrender1 libxext6 libxi6 libxkbcommon0 xkb-data libxrandr2 libxss1 libsdl2-ttf-2.0-0 libharfbuzz0b libglib2.0-0 libgraphite2-3 libspeexdsp1 libvorbisfile3 libzip4t64 x11vnc \
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