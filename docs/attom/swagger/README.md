# ATTOM API Swagger Specifications

This directory contains the latest OpenAPI/Swagger specifications for various ATTOM Data API endpoints, updated to reflect the current API versions as of November 2025.

## Latest API Versions Downloaded

All files have been downloaded from the ATTOM API swagger specification repository and updated to the latest versions.

| API | Version | Base Path | File | Description |
|-----|---------|-----------|------|-------------|
| Property API | v4.0.0 | `/v4/property` | `propertyapi_v4.pretty.json` | Main property endpoints |
| Property Detail | v4.0.0 | `/v4/property` | `propertyapi_propertydetailv4.pretty.json` | Detailed property information |
| Transaction/All Events | v4.0.0 | `/v4` | `propertyapi_transactionv4.pretty.json` | Property transactions and events |
| All Events Extended | v1.2 | `/propertyapi/v1.0.0` | `allevents_extended_v4.pretty.json` | Extended all events by address |
| Area API | v4.0.0 | `/v4/area` | `areaapi_v4.pretty.json` | Geographic area data |
| POI API | v4.0.0 | `/v4/neighborhood` | `poiapi_v4.pretty.json` | Points of interest |
| Community API | v4.0.0 | `/v4` | `communityapi_v4.pretty.json` | Community data |

## Legacy Files (Deprecated)

The following files contain older API versions and are kept for reference only. New development should use the latest versions above.

- `propertyapi_allevents.pretty.json` (superseded by `allevents_extended_v4.pretty.json`)
- `propertyapi_area.pretty.json` (superseded by `areaapi_v4.pretty.json`)
- `propertyapi_assessment.pretty.json` (superseded by `propertyapi_v4.pretty.json`)
- `propertyapi_assessmenthistory.pretty.json` (superseded by `propertyapi_v4.pretty.json`)
- `propertyapi_avm.pretty.json` (superseded by `propertyapi_v4.pretty.json`)
- `propertyapi_community.pretty.json` (superseded by `communityapi_v4.pretty.json`)
- `propertyapi_poi.pretty.json` (superseded by `poiapi_v4.pretty.json`)
- `propertyapi_property.pretty.json` (superseded by `propertyapi_v4.pretty.json`)
- `propertyapi_sale.pretty.json` (superseded by `propertyapi_v4.pretty.json`)
- `propertyapi_saleshistory.pretty.json` (superseded by `propertyapi_v4.pretty.json`)
- `propertyapi_valuation.pretty.json` (superseded by `propertyapi_v4.pretty.json`)

## Download Source

All current files were downloaded from:
```
https://api.developer.attomdata.com/swagger/spec/
```

Specific URLs used:
- Property API v4: `https://api.developer.attomdata.com/swagger/spec/propertyapi_v4.json`
- Property Detail v4: `https://api.developer.attomdata.com/swagger/spec/propertyapi_propertydetailv4.json`
- Transaction v4: `https://api.developer.attomdata.com/swagger/spec/propertyapi_transactionv4.json`
- All Events Extended: `https://api.developer.attomdata.com/swagger/spec/allevents_extended_by_address.json`
- Area API v4: `https://api.developer.attomdata.com/swagger/spec/AreaAPI.json`
- POI API v4: `https://api.developer.attomdata.com/swagger/spec/POIAPIV4.json`
- Community API v4: `https://api.developer.attomdata.com/swagger/spec/communityv4_api.json`

## Verification

All files have been verified to:
- Be valid JSON (parseable by `python3 -m json.tool`)
- Contain actual content (not empty)
- Be properly formatted (pretty-printed)
- Follow the OpenAPI 2.0 (Swagger) specification format
- Reflect the current API versions as documented on the ATTOM developer portal

## Usage

These swagger specifications can be used to:
- Generate client code using tools like Swagger Codegen or OpenAPI Generator
- Understand API endpoints, parameters, and response schemas
- Import into API testing tools like Postman
- Reference during implementation of the go-attom client library
- Ensure implementation matches the latest ATTOM API versions

## Version History

- **November 2025**: Updated all specifications to latest v4.0.0 versions where available
- **Previous**: Mixed v1.0.0 and v2.0.0 versions
