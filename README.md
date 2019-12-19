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