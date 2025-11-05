# ParcelTilesAPI Implementation Guide

## Overview

The ParcelTilesAPI provides raster tile visualization of parcel boundaries with **1 endpoint** returning PNG tiles compatible with popular mapping libraries.

## Endpoint

| Endpoint | Description | Parameters |
|----------|-------------|------------|
| `/parceltiles/{z}/{x}/{y}.png` | Parcel boundary raster tiles | apikey, z (zoom), x (column), y (row) |

## Tile Specifications

### Zoom Levels
- **14**: Neighborhood view (widest area)
- **15**: Block-level view
- **16**: Multi-property view
- **17**: Property-level view
- **18**: Individual property view (most detailed)

### Tile Format
- **Format**: PNG (image/png)
- **Coordinate System**: Web Mercator (EPSG:3857)
- **Tile Scheme**: Standard XYZ tile scheme
- **Size**: 256x256 pixels per tile

### URL Pattern
```
https://api.gateway.attomdata.com/parceltiles/{z}/{x}/{y}.png?apikey=YOUR_API_KEY
```

## Implementation Checklist

### Phase 1: Basic Tile Fetching
- [ ] Implement `GetParcelTile()` to fetch raw PNG bytes
- [ ] Add zoom level validation (14-18 only)
- [ ] Add tile coordinate validation
- [ ] Handle PNG binary response (not JSON!)

### Phase 2: URL Generation
- [ ] Implement `GetTileURL()` for URL generation
- [ ] Add URL template support for mapping libraries
- [ ] Create tile URL builders for different libraries

### Phase 3: Coordinate Conversion
- [ ] Add lat/lon to tile coordinate conversion
- [ ] Add tile coordinate to lat/lon conversion
- [ ] Add bounding box to tile range calculation
- [ ] Create tile grid calculation helpers

### Phase 4: Integration Helpers
- [ ] Create Leaflet integration helpers
- [ ] Create Mapbox GL JS integration helpers
- [ ] Create Google Maps integration helpers
- [ ] Add tile caching recommendations

## Model Design

### Core Types

```go
package models

// TileParams represents tile request parameters
type TileParams struct {
    Z int  // Zoom level (14-18)
    X int  // Tile column
    Y int  // Tile row
}

// Validate checks if tile parameters are valid
func (t *TileParams) Validate() error {
    if t.Z < 14 || t.Z > 18 {
        return fmt.Errorf("zoom level must be between 14 and 18, got %d", t.Z)
    }
    
    maxTile := 1 << uint(t.Z)  // 2^z
    
    if t.X < 0 || t.X >= maxTile {
        return fmt.Errorf("tile X must be between 0 and %d, got %d", maxTile-1, t.X)
    }
    
    if t.Y < 0 || t.Y >= maxTile {
        return fmt.Errorf("tile Y must be between 0 and %d, got %d", maxTile-1, t.Y)
    }
    
    return nil
}

// TileCoordinates represents a geographic coordinate
type TileCoordinates struct {
    Latitude  float64
    Longitude float64
    Zoom      int
}

// ToTileXY converts lat/lon to tile coordinates
func (tc *TileCoordinates) ToTileXY() (x, y int) {
    n := math.Pow(2, float64(tc.Zoom))
    
    x = int(math.Floor((tc.Longitude + 180.0) / 360.0 * n))
    
    latRad := tc.Latitude * math.Pi / 180.0
    y = int(math.Floor((1.0 - math.Log(math.Tan(latRad)+(1/math.Cos(latRad)))/math.Pi) / 2.0 * n))
    
    return x, y
}

// TileBounds represents the geographic bounds of a tile
type TileBounds struct {
    North float64
    South float64
    East  float64
    West  float64
}

// GetTileBounds calculates the geographic bounds of a tile
func GetTileBounds(x, y, z int) TileBounds {
    n := math.Pow(2, float64(z))
    
    west := float64(x)/n*360.0 - 180.0
    east := float64(x+1)/n*360.0 - 180.0
    
    north := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y)/n))) * 180.0 / math.Pi
    south := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y+1)/n))) * 180.0 / math.Pi
    
    return TileBounds{
        North: north,
        South: south,
        East:  east,
        West:  west,
    }
}

// BoundingBox represents a geographic bounding box
type BoundingBox struct {
    MinLat float64
    MaxLat float64
    MinLon float64
    MaxLon float64
}

// GetTileRange calculates the range of tiles covering a bounding box
func (bb *BoundingBox) GetTileRange(zoom int) (minX, minY, maxX, maxY int) {
    nwCoord := TileCoordinates{Latitude: bb.MaxLat, Longitude: bb.MinLon, Zoom: zoom}
    seCoord := TileCoordinates{Latitude: bb.MinLat, Longitude: bb.MaxLon, Zoom: zoom}
    
    minX, minY = nwCoord.ToTileXY()
    maxX, maxY = seCoord.ToTileXY()
    
    return minX, minY, maxX, maxY
}
```

## Service Implementation

```go
package parcel

import (
    "context"
    "fmt"
    "io"
    "github.com/my-eq/go-attom/pkg/client"
    "github.com/my-eq/go-attom/pkg/models"
)

type Service struct {
    client *client.Client
}

func NewService(c *client.Client) *Service {
    return &Service{client: c}
}

// GetParcelTile fetches a parcel tile as PNG bytes
func (s *Service) GetParcelTile(ctx context.Context, z, x, y int) ([]byte, error) {
    params := &models.TileParams{Z: z, X: x, Y: y}
    
    if err := params.Validate(); err != nil {
        return nil, err
    }
    
    url := fmt.Sprintf("/parceltiles/%d/%d/%d.png", z, x, y)
    
    resp, err := s.client.DoRequest(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    // Read PNG bytes (not JSON!)
    return io.ReadAll(resp.Body)
}

// GetTileURL generates the full URL for a tile
func (s *Service) GetTileURL(z, x, y int) string {
    return fmt.Sprintf("%s/parceltiles/%d/%d/%d.png?apikey=%s",
        s.client.BaseURL, z, x, y, s.client.APIKey)
}

// GetTileURLTemplate returns a URL template for mapping libraries
func (s *Service) GetTileURLTemplate() string {
    return fmt.Sprintf("%s/parceltiles/{z}/{x}/{y}.png?apikey=%s",
        s.client.BaseURL, s.client.APIKey)
}
```

## Library Integration Examples

### Leaflet Integration

```javascript
// Add ATTOM parcel tiles to Leaflet map
const attomTileLayer = L.tileLayer(
    'https://api.gateway.attomdata.com/parceltiles/{z}/{x}/{y}.png?apikey=YOUR_API_KEY',
    {
        minZoom: 14,
        maxZoom: 18,
        attribution: '&copy; ATTOM Data Solutions'
    }
);

const map = L.map('map').setView([34.0522, -118.2437], 16);
attomTileLayer.addTo(map);
```

### Mapbox GL JS Integration

```javascript
// Add ATTOM parcel tiles to Mapbox GL JS
map.addSource('attom-parcels', {
    'type': 'raster',
    'tiles': [
        'https://api.gateway.attomdata.com/parceltiles/{z}/{x}/{y}.png?apikey=YOUR_API_KEY'
    ],
    'tileSize': 256,
    'minzoom': 14,
    'maxzoom': 18
});

map.addLayer({
    'id': 'attom-parcel-layer',
    'type': 'raster',
    'source': 'attom-parcels',
    'paint': {
        'raster-opacity': 0.7
    }
});
```

### Google Maps Integration

```javascript
// Add ATTOM parcel tiles to Google Maps
const attomTileLayer = new google.maps.ImageMapType({
    getTileUrl: function(coord, zoom) {
        if (zoom < 14 || zoom > 18) return null;
        return 'https://api.gateway.attomdata.com/parceltiles/' + 
               zoom + '/' + coord.x + '/' + coord.y + '.png?apikey=YOUR_API_KEY';
    },
    tileSize: new google.maps.Size(256, 256),
    maxZoom: 18,
    minZoom: 14,
    name: 'ATTOM Parcels'
});

map.overlayMapTypes.push(attomTileLayer);
```

### OpenLayers Integration

```javascript
// Add ATTOM parcel tiles to OpenLayers
const attomSource = new ol.source.XYZ({
    url: 'https://api.gateway.attomdata.com/parceltiles/{z}/{x}/{y}.png?apikey=YOUR_API_KEY',
    minZoom: 14,
    maxZoom: 18
});

const attomLayer = new ol.layer.Tile({
    source: attomSource,
    opacity: 0.7
});

map.addLayer(attomLayer);
```

## Go Usage Examples

### Example 1: Fetch single tile

```go
package main

import (
    "context"
    "os"
    "github.com/my-eq/go-attom/pkg/client"
    "github.com/my-eq/go-attom/pkg/parcel"
)

func main() {
    c := client.NewClient("YOUR_API_KEY")
    svc := parcel.NewService(c)
    ctx := context.Background()
    
    // Los Angeles area at zoom 16
    tileBytes, err := svc.GetParcelTile(ctx, 16, 11274, 26168)
    if err != nil {
        panic(err)
    }
    
    // Save to file
    os.WriteFile("parcel_tile.png", tileBytes, 0644)
}
```

### Example 2: Generate tile URL

```go
// Get tile URL for embedding in web app
url := svc.GetTileURL(16, 11274, 26168)
fmt.Println(url)
// Output: https://api.gateway.attomdata.com/parceltiles/16/11274/26168.png?apikey=YOUR_API_KEY

// Get URL template for Leaflet/Mapbox
template := svc.GetTileURLTemplate()
fmt.Println(template)
// Output: https://api.gateway.attomdata.com/parceltiles/{z}/{x}/{y}.png?apikey=YOUR_API_KEY
```

### Example 3: Convert lat/lon to tile coordinates

```go
import "github.com/my-eq/go-attom/pkg/models"

// Los Angeles coordinates
coord := models.TileCoordinates{
    Latitude:  34.0522,
    Longitude: -118.2437,
    Zoom:      16,
}

x, y := coord.ToTileXY()
fmt.Printf("Tile coordinates: %d, %d at zoom %d\n", x, y, coord.Zoom)
// Output: Tile coordinates: 11274, 26168 at zoom 16

// Fetch the tile
tileBytes, _ := svc.GetParcelTile(ctx, coord.Zoom, x, y)
```

### Example 4: Download all tiles in bounding box

```go
import "github.com/my-eq/go-attom/pkg/models"

// Define bounding box (e.g., downtown Los Angeles)
bbox := models.BoundingBox{
    MinLat: 34.0400,
    MaxLat: 34.0600,
    MinLon: -118.2600,
    MaxLon: -118.2300,
}

// Get tile range at zoom 16
zoom := 16
minX, minY, maxX, maxY := bbox.GetTileRange(zoom)

fmt.Printf("Tile range: X[%d-%d], Y[%d-%d]\n", minX, maxX, minY, maxY)

// Download all tiles in range
for x := minX; x <= maxX; x++ {
    for y := minY; y <= maxY; y++ {
        tileBytes, err := svc.GetParcelTile(ctx, zoom, x, y)
        if err != nil {
            fmt.Printf("Error fetching tile %d/%d/%d: %v\n", zoom, x, y, err)
            continue
        }
        
        filename := fmt.Sprintf("tiles/%d_%d_%d.png", zoom, x, y)
        os.WriteFile(filename, tileBytes, 0644)
    }
}
```

## Testing Strategy

1. **Zoom Level Tests**: Test all zoom levels 14-18
2. **Coordinate Validation**: Test boundary conditions (x=0, y=0, max values)
3. **Conversion Tests**: Verify lat/lon ↔ tile coordinate conversions
4. **Binary Response**: Verify PNG bytes are valid image data
5. **Error Handling**: Test invalid zoom levels, out-of-bounds coordinates

## Common Pitfalls

1. **Binary vs JSON**: Tile endpoint returns PNG bytes, NOT JSON!
2. **Zoom Range**: Only zoom levels 14-18 are supported (not 1-20 like some tile services)
3. **Coordinate Order**: Tile coordinates are (x, y) not (lat, lon)
4. **Y-Axis Direction**: Tile Y increases downward (not upward like lat)
5. **API Key in URL**: Must include apikey as query parameter
6. **Tile Size**: Always 256x256 pixels (not configurable)
7. **CORS**: If using in browser, ensure proper CORS headers

## Caching Recommendations

```go
// Example tile cache implementation
type TileCache struct {
    cache map[string][]byte
    mu    sync.RWMutex
}

func (tc *TileCache) Get(z, x, y int) ([]byte, bool) {
    key := fmt.Sprintf("%d/%d/%d", z, x, y)
    tc.mu.RLock()
    defer tc.mu.RUnlock()
    
    data, ok := tc.cache[key]
    return data, ok
}

func (tc *TileCache) Set(z, x, y int, data []byte) {
    key := fmt.Sprintf("%d/%d/%d", z, x, y)
    tc.mu.Lock()
    defer tc.mu.Unlock()
    
    tc.cache[key] = data
}

// Use with parcel service
func GetCachedTile(ctx context.Context, svc *parcel.Service, cache *TileCache, z, x, y int) ([]byte, error) {
    if data, ok := cache.Get(z, x, y); ok {
        return data, nil
    }
    
    data, err := svc.GetParcelTile(ctx, z, x, y)
    if err != nil {
        return nil, err
    }
    
    cache.Set(z, x, y, data)
    return data, nil
}
```

## Performance Considerations

1. **Rate Limiting**: Tile requests count against API rate limits
2. **Bandwidth**: PNG tiles can be 10-50KB each - bulk downloads add up
3. **Caching**: Implement aggressive caching - tiles don't change frequently
4. **Parallel Requests**: Use goroutines for bulk downloads, but respect rate limits
5. **CDN**: Consider proxying tiles through a CDN for web applications

## Web Server Example

```go
// Serve tiles through a web server with caching
package main

import (
    "context"
    "fmt"
    "net/http"
    "strconv"
    "github.com/my-eq/go-attom/pkg/client"
    "github.com/my-eq/go-attom/pkg/parcel"
)

func main() {
    c := client.NewClient("YOUR_API_KEY")
    svc := parcel.NewService(c)
    cache := &TileCache{cache: make(map[string][]byte)}
    
    http.HandleFunc("/tiles/", func(w http.ResponseWriter, r *http.Request) {
        // Parse URL: /tiles/16/11274/26168.png
        var z, x, y int
        fmt.Sscanf(r.URL.Path, "/tiles/%d/%d/%d.png", &z, &x, &y)
        
        // Get tile (with caching)
        ctx := r.Context()
        data, err := GetCachedTile(ctx, svc, cache, z, x, y)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // Serve PNG
        w.Header().Set("Content-Type", "image/png")
        w.Header().Set("Cache-Control", "public, max-age=86400")  // Cache 24 hours
        w.Write(data)
    })
    
    http.ListenAndServe(":8080", nil)
}
```
