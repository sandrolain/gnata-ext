# extnet — Network Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extnet`

---

### `$ipVersion(str)`

Returns `4` for IPv4, `6` for IPv6, or `-1` if *str* is not a valid IP address.

```jsonata
$ipVersion("192.168.1.1")   /* → 4 */
$ipVersion("::1")           /* → 6 */
$ipVersion("not-an-ip")     /* → -1 */
```

---

### `$isPrivateIP(str)`

Returns `true` if *str* is a private, loopback, or link-local address:

- RFC 1918: `10/8`, `172.16/12`, `192.168/16`
- Loopback: `127/8`, `::1`
- Link-local: `169.254/16`, `fe80::/10`

```jsonata
$isPrivateIP("192.168.1.1")   /* → true */
$isPrivateIP("8.8.8.8")       /* → false */
$isPrivateIP("::1")           /* → true */
```

---

### `$ipToInt(str)`

Converts an IPv4 address string to its 32-bit unsigned integer representation (returned as `float64`). Only IPv4 is supported.

```jsonata
$ipToInt("0.0.0.0")           /* → 0 */
$ipToInt("255.255.255.255")   /* → 4294967295 */
$ipToInt("192.168.1.1")       /* → 3232235777 */
```

---

### `$intToIP(n)`

Converts a 32-bit unsigned integer (as `float64`) back to an IPv4 address string.

```jsonata
$intToIP(0)           /* → "0.0.0.0" */
$intToIP(3232235777)  /* → "192.168.1.1" */
```

---

### `$ipInCIDR(ip, cidr)`

Returns `true` if *ip* is contained within the *cidr* block.

```jsonata
$ipInCIDR("192.168.1.100", "192.168.1.0/24")   /* → true */
$ipInCIDR("10.0.0.1",      "192.168.1.0/24")   /* → false */
```

---

### `$expandCIDR(cidr)`

Returns a network-info object for the CIDR block.

**IPv4** — returns `{network, broadcast, first, last, count}`:

```jsonata
$expandCIDR("192.168.1.0/24")
/* → {
     "network":   "192.168.1.0",
     "broadcast": "192.168.1.255",
     "first":     "192.168.1.1",
     "last":      "192.168.1.254",
     "count":     256
   } */
```

**IPv6** — returns `{network, first, last}` (no broadcast, count omitted due to scale):

```jsonata
$expandCIDR("2001:db8::/32")
/* → {
     "network": "2001:db8::",
     "first":   "2001:db8::",
     "last":    "2001:db8:ffff:ffff:ffff:ffff:ffff:ffff"
   } */
```

---
