# FAQ

## Allow custom MAC vendor overrides 
LanLookout uses `arp-scan`, so most of its options are available.

1. Prepare a [mac-vendor.txt](https://manpages.debian.org/testing/arp-scan/mac-vendor.5.en.html) file with additional MACs and put it in the mounted data directory (`/data/lanlookout`).
2. If you are using `IFACES` variable to define interfaces, add path to mac-vendor.txt to `ARP_ARGS`
```yaml
arp_args: --macfile=/data/lanlookout/mac-vendor.txt
```
3. For interfaces defined in `ARP_STRS` add the same directly in the beginning of `ARP_STRS` string
```yaml
arp_strs:
    - --macfile=/data/lanlookout/mac-vendor.txt -gNx 10.144.0.1/24 -I eth0
```
4. **WARNING!** To see an updated vendor, you'll have to delete host and wait for the next scan.