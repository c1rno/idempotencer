LIB_ZMQ=v4.3.2
LIB_CZMQ=v4.2.0
LIB_RD_KAFKA=v1.4.0

docker-install-dependencies:
	apt-get update
	apt-get install -y --no-install-recommends \
		ca-certificates \
		gnupg \
		wget
	echo 'deb http://download.opensuse.org/repositories/network:/messaging:/zeromq:/release-stable/Debian_Next/ /' \
		> /etc/apt/sources.list.d/network:messaging:zeromq:release-stable.list
	wget -nv https://download.opensuse.org/repositories/network:messaging:zeromq:release-stable/Debian_Next/Release.key \
		-O Release.key
	apt-key add - < Release.key
	rm Release.key
	apt-get update
	apt-get install -y --no-install-recommends \
		git \
		libsodium-dev \
		build-essential \
		libtool \
		pkg-config \
		autotools-dev \
		automake \
		autoconf \
		cmake \
		uuid-dev \
		libpcre3-dev
	rm -rf /var/cache/apt/*

docker-install-rd-kafka:
	git clone https://github.com/edenhill/librdkafka.git --depth 1 --branch $(LIB_RD_KAFKA)
	cd librdkafka && ./configure --install-deps
	$(MAKE) -C librdkafka
	$(MAKE) -C librdkafka install
	rm -rf librdkafka

docker-install-libzmq:
	git clone https://github.com/zeromq/libzmq.git --depth 1 --branch $(LIB_ZMQ)
	cd libzmq && ./autogen.sh
	cd libzmq && ./configure --disable-shared --enable-static --with-libsodium
	$(MAKE) -C libzmq install
	cd libzmq && ldconfig
	rm -rf libzmq

docker-install-czmq:
	git clone https://github.com/zeromq/czmq.git --depth 1 --branch $(LIB_CZMQ)
	cd czmq && ./autogen.sh
	cd czmq && ./configure --disable-shared --enable-static
	$(MAKE) -C czmq install
	cd czmq && ldconfig
	rm -rf czmq
