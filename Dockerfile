# # Copyright
#
# Copyright (c) 2017-2020 Fiducia DLT Limited
# All Rights Reserved.
#
# NOTICE:  All information contained herein is, and remains
# the property of Fiducia DLT Limited
# The intellectual and technical concepts contained
# herein are proprietary to Fiducia DLT Limited
# and may be covered by U.S. and Foreign Patents,
# patents in process, and are protected by trade secret or copyright law.
# Dissemination of this information or reproduction of this material
# is strictly forbidden unless prior written permission is obtained
# from Fiducia DLT Limited
#
# Written by
# Anton Zhukov <anton@papyrus.global>

FROM alpine

RUN apk add --no-cache \
        libc6-compat \
        ca-certificates && \
    rm -rf /var/cache/apk/*

COPY docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["spacetrouble"]

COPY ./bin/* /usr/local/bin/