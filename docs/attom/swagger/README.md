# ATTOM API Swagger Specifications

This directory contains the OpenAPI/Swagger specifications for various ATTOM Data API endpoints.

## Downloaded Files

All files have been successfully downloaded and pretty-printed from the ATTOM API swagger specification repository.

| Output File | Source URL | Description | Size |
|-------------|------------|-------------|------|
| `propertyapi_allevents.pretty.json` | `https://api.developer.attomdata.com/swagger/spec/allevents_extended_by_address.json` | All property events combined | 2.9 KB |
| `propertyapi_assessmenthistory.pretty.json` | `https://api.developer.attomdata.com/swagger/spec/propertyapi_assessment_history_address.json` | Historical assessment records | 3.0 KB |
| `propertyapi_salestrend.pretty.json` | `https://api.developer.attomdata.com/swagger/spec/sales_trend_by_year.json` | Geographic sales trends by year | 3.7 KB |
| `propertyapi_school.pretty.json` | `https://api.developer.attomdata.com/swagger/spec/propertyapi_schools.json` | School information and search | 9.7 KB |
| `propertyapi_valuation.pretty.json` | `https://api.developer.attomdata.com/swagger/spec/propertyapi-valuationv1.json` | Property valuation data | 23 KB |

## File Mapping Notes

The original issue requested files with specific names, but the actual ATTOM API uses slightly different naming conventions. Here's the mapping:

- **allevents**: Source file is `allevents_extended_by_address.json` (not `propertyapi_allevents.json`)
- **assessmenthistory**: Source file is `propertyapi_assessment_history_address.json` (address-based version)
- **salestrend**: Source file is `sales_trend_by_year.json` (yearly trend version; monthly and quarterly versions also available)
- **school**: Source file is `propertyapi_schools.json` (plural form)
- **valuation**: Source file is `propertyapi-valuationv1.json` (with v1 suffix and hyphen)

## Other Available Swagger Specs

The ATTOM API provides many other swagger specifications. To see the complete list, access:
```
https://api.developer.attomdata.com/swagger/spec/
```

This returns a JSON index of all available specifications.

## Verification

All files have been verified to:
- Be valid JSON (parseable by `python3 -m json.tool`)
- Contain actual content (not empty)
- Be properly formatted (pretty-printed)
- Follow the OpenAPI 2.0 (Swagger) specification format

## Usage

These swagger specifications can be used to:
- Generate client code using tools like Swagger Codegen or OpenAPI Generator
- Understand API endpoints, parameters, and response schemas
- Import into API testing tools like Postman
- Reference during implementation of the go-attom client library
