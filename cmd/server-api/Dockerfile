FROM registry.ronaksoftware.com/base/docker/ubuntu:18.04
MAINTAINER Ehsan N. Moosa <ehsan@ronaksoft.com>

# Add contents from root directory to the docker
ADD ./_build/ /ronak/bin

# Set working directory
WORKDIR /ronak/bin

# Entry point
ENTRYPOINT ["/ronak/bin/server-api"]