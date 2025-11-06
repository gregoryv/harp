harp - HTTP ip scanner and ARP-cache parser

Fire of HTTP HEAD requests to a set of ips and parse (arp -a) result for
MAC IPs.


## Quick start

  $ go install github.com/gregoryv/harp/cmd/harp@latest



Single

    $ harp -ip 192.168.1.3

Range

	$ harp -ip 192.168.1.4-29

All

	$ harp -ip 192.168.1.*

None

    $ harp
	
End result is always a list of MAC IP values from the arp cache.  There
is no guarantee that the cache has been populated after a scan, you
might need to run harp without -ip a few times.
