# bogus

## Name

*bogus* - return NXDOMAIN directly if the resovled IP is in the bogus list.

## Description

*bogus* takes a list of IP address and returns NXDOMAIN for any resolved address under them instead
 of return the bogus IP address.
 
Put *bogus* beside *rewrite* in `plugin.cfg`.

## Syntax

~~~ txt
bogus BOGUS_IPs......
~~~

## Examples

~~~ corefile
. {
    bogus 10.255.10.1 10.255.10.2 10.255.20.1 10.255.20.2
    whoami
}
~~~

# Bugs

The list of IP address is just a slice that is traversed, meaning this plugin will get slow when a lof of
IP address are to be bogus.

# add me

```
#!/bin/bash
china=`curl -sSL https://github.com/felixonmars/dnsmasq-china-list/raw/master/accelerated-domains.china.conf | while read line; do awk -F '/' '{print $2}' | grep -v '#' ; done |  paste -sd " " -`
apple=`curl -sSL https://github.com/felixonmars/dnsmasq-china-list/raw/master/apple.china.conf | while read line; do awk -F '/' '{print $2}' | grep -v '#' ; done |  paste -sd " " -`
google=`curl -sSL https://github.com/felixonmars/dnsmasq-china-list/raw/master/google.china.conf | while read line; do awk -F '/' '{print $2}' | grep -v '#' ; done |  paste -sd " " -`
bogus=`curl -sSL https://github.com/felixonmars/dnsmasq-china-list/raw/master/bogus-nxdomain.china.conf | grep "=" | while read line; do awk -F '=' '{print $2}' | grep -v '#' ; done |  paste -sd " " -`
cat>Corefile<<EOF
. {
    ads {
        default-lists
        blacklist https://raw.githubusercontent.com/privacy-protection-tools/anti-AD/master/anti-ad-domains.txt
        whitelist https://files.krnl.eu/whitelist.txt
        log
        auto-update-interval 24h
        list-store ads-cache
    }
    hosts {
        fallthrough
    }
    forward . 208.67.222.222:443 208.67.222.222:5353 208.67.220.220:443 208.67.220.220:5353 127.0.0.1:5301 127.0.0.1:5302 127.0.0.1:5303  {
    except $china $apple $google cdn.jsdelivr.net
    }
    proxy . 192.168.1.1
    bogus $bogus
    log
    cache
    redisc {
        endpoint 127.0.0.1:6379
    }
    health
    reload
}
.:5301 {
    bind 127.0.0.1
    forward . tls://9.9.9.9 tls://9.9.9.10 {
        tls_servername dns.quad9.net
    }
    cache
}
.:5302 {
    bind 127.0.0.1
    forward . tls://1.1.1.1 tls://1.0.0.1 {
        tls_servername cloudflare-dns.com
    }
    cache
}
.:5303 {
    bind 127.0.0.1
    forward . tls://8.8.8.8 tls://8.8.4.4 {
        tls_servername dns.google
    }
    cache
}
EOF
```
