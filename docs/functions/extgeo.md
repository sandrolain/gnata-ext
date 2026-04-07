# extgeo — Geographic Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extgeo`

All calculations use the WGS-84 mean Earth radius **6 371 km**. Coordinates are always **decimal degrees**. No external dependencies — pure `math` stdlib.

---

### `$haversine(lat1, lon1, lat2, lon2)`

Returns the great-circle distance between two points in **kilometres**.

```jsonata
/* London → Paris ≈ 340 km */
$haversine(51.5074, -0.1278, 48.8566, 2.3522)
/* → ~340.6 */

/* Same point */
$haversine(48.0, 2.0, 48.0, 2.0)   /* → 0 */
```

---

### `$bearing(lat1, lon1, lat2, lon2)`

Returns the initial bearing from point 1 to point 2 in **degrees**, clockwise from north (0–360).

```jsonata
$bearing(0, 0, 0, 1)    /* → ~90  (east) */
$bearing(0, 0, 1, 0)    /* → ~0   (north) */
```

---

### `$geoFormat(lat, lon [, format])`

Formats a coordinate pair as a string.

- `"decimal"` (default): `"48.8566, 2.3522"` — 4 decimal places
- `"dms"`: degrees/minutes/seconds with cardinal directions, e.g. `"48°51'23.76\"N 2°21'7.92\"E"`

```jsonata
$geoFormat(48.8566, 2.3522)           /* → "48.8566, 2.3522" */
$geoFormat(48.8566, 2.3522, "dms")    /* → "48°51'23.76\"N 2°21'7.92\"E" */
```

---

### `$geoParse(str)`

Parses a `"lat, lon"` decimal string and returns `{lat, lon}`.

```jsonata
$geoParse("48.8566, 2.3522")
/* → {"lat": 48.8566, "lon": 2.3522} */
```

---

### `$inBoundingBox(lat, lon, minLat, minLon, maxLat, maxLon)`

Returns `true` if the point (`lat`, `lon`) lies within the axis-aligned bounding box.

```jsonata
/* Paris inside Europe bbox */
$inBoundingBox(48.8566, 2.3522, 36.0, -10.0, 71.0, 40.0)   /* → true */
```

---

### `$geoDistance(point, points)`

Computes the haversine distance from *point* to each element of *points* and returns an array of distances in **kilometres**.

- *point*: `{"lat": float, "lon": float}`
- *points*: array of `{"lat": float, "lon": float}`

```jsonata
$geoDistance(
  {"lat": 51.5074, "lon": -0.1278},
  [
    {"lat": 48.8566, "lon": 2.3522},
    {"lat": 52.5200, "lon": 13.4050}
  ]
)
/* → [~340.6, ~930.9] */
```

---
