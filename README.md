harp - HTTP ip scanner and ARP-cache parser

Fire of HTTP HEAD requests to a set of ips and parse (arp -a) result for
MAC IPs.


## Quick start

    $ go install github.com/gregoryv/harp/cmd/harp@latest
    $ harp -h
    
    Usage: harp [IP]
    
    Examples
    
      $ harp 192.168.1.3
      $ harp 192.168.1.3-9
      $ harp 192.168.1.*
    
    without IP harp shows the arp -a cache.

	
End result is always a list of MAC IP values from the arp cache.  There
is no guarantee that the cache has been populated after a scan, you
might need to run harp without -ip a few times.
