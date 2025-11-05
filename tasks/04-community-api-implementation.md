# CommunityAPI Implementation Guide

## Overview

The CommunityAPI provides comprehensive community demographics and statistics with **2 endpoints** covering 6 major data categories.

## Endpoints

Base path: `/communityapi/v2.0.0/` and `/neighborhood/`

| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/communityapi/v2.0.0/area/full` | Full community demographics | AreaId |
| `/neighborhood/community` | Community profile data | geoIdV4, AreaId |

## Data Categories (6 total)

| Category | Data Included |
|----------|---------------|
| **Demographics** | Population, households, age distribution, race, ethnicity, household composition |
| **Climate** | Average temperature, precipitation, seasonal weather patterns, climate zones |
| **Economics** | Median income, employment rates, job growth, industry breakdown, poverty levels |
| **Education** | School performance, graduation rates, education attainment levels, literacy |
| **Housing** | Home values, rental rates, vacancy rates, housing types, owner vs. renter ratios |
| **Transportation** | Commute times, transportation modes, walkability scores, public transit access |

## Key Parameters

- **AreaId**: Legacy area identifier (numeric)
- **geoIdV4**: ATTOM's new geographic identifier (SHA-256 hash)

## Implementation Checklist

### Phase 1: Basic Endpoints
- [ ] Implement `GetFullCommunityData()` with AreaId
- [ ] Implement `GetNeighborhoodCommunity()` with geoIdV4
- [ ] Add parameter validation for AreaId and geoIdV4

### Phase 2: Data Category Models
- [ ] Create Demographics model with population stats
- [ ] Create Climate model with weather data
- [ ] Create Economics model with income/employment data
- [ ] Create Education model with school stats
- [ ] Create Housing model with property data
- [ ] Create Transportation model with commute data

### Phase 3: Helper Methods
- [ ] Add demographic summary helpers
- [ ] Add economic trend helpers
- [ ] Add education rating helpers
- [ ] Add housing affordability helpers

## Model Design

### Core Models

```go
package models

// FullCommunityData contains all 6 data categories
type FullCommunityData struct {
    Status         *ResponseStatus  `json:"status,omitempty"`
    GeoIDV4        *string          `json:"geoIdV4,omitempty"`
    AreaID         *string          `json:"areaId,omitempty"`
    Demographics   *Demographics    `json:"demographics,omitempty"`
    Climate        *Climate         `json:"climate,omitempty"`
    Economics      *Economics       `json:"economics,omitempty"`
    Education      *Education       `json:"education,omitempty"`
    Housing        *Housing         `json:"housing,omitempty"`
    Transportation *Transportation  `json:"transportation,omitempty"`
}

// NeighborhoodCommunity response
type NeighborhoodCommunity struct {
    Status         *ResponseStatus  `json:"status,omitempty"`
    GeoIDV4        *string          `json:"geoIdV4,omitempty"`
    Demographics   *Demographics    `json:"demographics,omitempty"`
    Climate        *Climate         `json:"climate,omitempty"`
    Economics      *Economics       `json:"economics,omitempty"`
    Education      *Education       `json:"education,omitempty"`
    Housing        *Housing         `json:"housing,omitempty"`
    Transportation *Transportation  `json:"transportation,omitempty"`
}
```

### Demographics Model

```go
// Demographics contains population and household statistics
type Demographics struct {
    Population              *int     `json:"population,omitempty"`
    PopulationDensity       *float64 `json:"populationDensity,omitempty"`
    TotalHouseholds         *int     `json:"totalHouseholds,omitempty"`
    AverageHouseholdSize    *float64 `json:"averageHouseholdSize,omitempty"`
    MedianAge               *float64 `json:"medianAge,omitempty"`
    
    // Age distribution (percentages)
    AgeUnder18              *float64 `json:"ageUnder18,omitempty"`
    Age18To34               *float64 `json:"age18To34,omitempty"`
    Age35To54               *float64 `json:"age35To54,omitempty"`
    Age55To74               *float64 `json:"age55To74,omitempty"`
    Age75Plus               *float64 `json:"age75Plus,omitempty"`
    
    // Race/ethnicity (percentages)
    RaceWhite               *float64 `json:"raceWhite,omitempty"`
    RaceBlack               *float64 `json:"raceBlack,omitempty"`
    RaceAsian               *float64 `json:"raceAsian,omitempty"`
    RaceHispanic            *float64 `json:"raceHispanic,omitempty"`
    RaceOther               *float64 `json:"raceOther,omitempty"`
    
    // Household composition
    MarriedCouples          *float64 `json:"marriedCouples,omitempty"`
    SingleParents           *float64 `json:"singleParents,omitempty"`
    NonFamilyHouseholds     *float64 `json:"nonFamilyHouseholds,omitempty"`
}
```

### Climate Model

```go
// Climate contains weather and climate statistics
type Climate struct {
    AverageTemperature      *float64 `json:"averageTemperature,omitempty"`       // Fahrenheit
    AverageTempJanuary      *float64 `json:"averageTempJanuary,omitempty"`
    AverageTempJuly         *float64 `json:"averageTempJuly,omitempty"`
    AnnualPrecipitation     *float64 `json:"annualPrecipitation,omitempty"`      // Inches
    AnnualSnowfall          *float64 `json:"annualSnowfall,omitempty"`           // Inches
    SunnyDays               *int     `json:"sunnyDays,omitempty"`
    ClimateZone             *string  `json:"climateZone,omitempty"`
    HumidityAverage         *float64 `json:"humidityAverage,omitempty"`          // Percentage
}
```

### Economics Model

```go
// Economics contains income and employment statistics
type Economics struct {
    MedianHouseholdIncome   *int     `json:"medianHouseholdIncome,omitempty"`    // USD
    MedianFamilyIncome      *int     `json:"medianFamilyIncome,omitempty"`
    PerCapitaIncome         *int     `json:"perCapitaIncome,omitempty"`
    PovertyRate             *float64 `json:"povertyRate,omitempty"`              // Percentage
    UnemploymentRate        *float64 `json:"unemploymentRate,omitempty"`         // Percentage
    
    // Income distribution (percentages)
    IncomeUnder25K          *float64 `json:"incomeUnder25K,omitempty"`
    Income25KTo50K          *float64 `json:"income25KTo50K,omitempty"`
    Income50KTo75K          *float64 `json:"income50KTo75K,omitempty"`
    Income75KTo100K         *float64 `json:"income75KTo100K,omitempty"`
    Income100KTo150K        *float64 `json:"income100KTo150K,omitempty"`
    Income150KPlus          *float64 `json:"income150KPlus,omitempty"`
    
    // Employment by industry (percentages)
    IndustryManufacturing   *float64 `json:"industryManufacturing,omitempty"`
    IndustryRetail          *float64 `json:"industryRetail,omitempty"`
    IndustryHealthcare      *float64 `json:"industryHealthcare,omitempty"`
    IndustryEducation       *float64 `json:"industryEducation,omitempty"`
    IndustryTechnology      *float64 `json:"industryTechnology,omitempty"`
    IndustryFinance         *float64 `json:"industryFinance,omitempty"`
}
```

### Education Model

```go
// Education contains school and education statistics
type Education struct {
    // Educational attainment (percentages)
    HighSchoolGradRate      *float64 `json:"highSchoolGradRate,omitempty"`
    BachelorsRate           *float64 `json:"bachelorsRate,omitempty"`
    GraduateDegreeRate      *float64 `json:"graduateDegreeRate,omitempty"`
    
    // School enrollment
    EnrollmentPreschool     *int     `json:"enrollmentPreschool,omitempty"`
    EnrollmentElementary    *int     `json:"enrollmentElementary,omitempty"`
    EnrollmentHighSchool    *int     `json:"enrollmentHighSchool,omitempty"`
    EnrollmentCollege       *int     `json:"enrollmentCollege,omitempty"`
    
    // School performance
    AverageTestScores       *float64 `json:"averageTestScores,omitempty"`
    StudentTeacherRatio     *float64 `json:"studentTeacherRatio,omitempty"`
    SchoolRating            *float64 `json:"schoolRating,omitempty"`          // 1-10 scale
}
```

### Housing Model

```go
// Housing contains property and housing statistics
type Housing struct {
    MedianHomeValue         *int     `json:"medianHomeValue,omitempty"`        // USD
    MedianRent              *int     `json:"medianRent,omitempty"`             // USD per month
    HomeownershipRate       *float64 `json:"homeownershipRate,omitempty"`      // Percentage
    VacancyRate             *float64 `json:"vacancyRate,omitempty"`            // Percentage
    
    // Housing types (percentages)
    SingleFamilyHomes       *float64 `json:"singleFamilyHomes,omitempty"`
    MultiUnitHomes          *float64 `json:"multiUnitHomes,omitempty"`
    MobileHomes             *float64 `json:"mobileHomes,omitempty"`
    
    // Housing age
    MedianYearBuilt         *int     `json:"medianYearBuilt,omitempty"`
    HomesBuiltBefore1980    *float64 `json:"homesBuiltBefore1980,omitempty"`
    HomesBuiltAfter2000     *float64 `json:"homesBuiltAfter2000,omitempty"`
    
    // Affordability
    PriceToIncomeRatio      *float64 `json:"priceToIncomeRatio,omitempty"`
    RentToIncomeRatio       *float64 `json:"rentToIncomeRatio,omitempty"`
}
```

### Transportation Model

```go
// Transportation contains commute and transit statistics
type Transportation struct {
    AverageCommuteTime      *float64 `json:"averageCommuteTime,omitempty"`     // Minutes
    
    // Commute modes (percentages)
    CommuteDriveAlone       *float64 `json:"commuteDriveAlone,omitempty"`
    CommuteCarpool          *float64 `json:"commuteCarpool,omitempty"`
    CommutePublicTransit    *float64 `json:"commutePublicTransit,omitempty"`
    CommuteBike             *float64 `json:"commuteBike,omitempty"`
    CommuteWalk             *float64 `json:"commuteWalk,omitempty"`
    CommuteWorkFromHome     *float64 `json:"commuteWorkFromHome,omitempty"`
    
    // Accessibility scores (0-100)
    WalkabilityScore        *int     `json:"walkabilityScore,omitempty"`
    TransitScore            *int     `json:"transitScore,omitempty"`
    BikeScore               *int     `json:"bikeScore,omitempty"`
    
    // Vehicle ownership
    VehiclesPerHousehold    *float64 `json:"vehiclesPerHousehold,omitempty"`
    NoVehicleHouseholds     *float64 `json:"noVehicleHouseholds,omitempty"`    // Percentage
}
```

### Parameter Struct

```go
// CommunityParams for community API requests
type CommunityParams struct {
    GeoIDV4 *string `json:"-"`
    AreaID  *string `json:"-"`
}
```

## Testing Strategy

1. **Data Completeness**: Verify all 6 categories are present
2. **Percentage Validation**: Ensure percentages sum to 100%
3. **Null Handling**: Test with areas having sparse data
4. **Numeric Ranges**: Validate income/population/temperature ranges
5. **GeoID Formats**: Test both legacy AreaId and new geoIdV4

## Common Pitfalls

1. **Missing Data**: Community data is VERY sparse - most fields will be nil
2. **Percentage Format**: Some APIs return decimals (0.25), others percentages (25.0)
3. **Temperature Units**: Always Fahrenheit, not Celsius
4. **Income Values**: USD integers, not float currency values
5. **Nested Nulls**: Check both parent struct and child fields for nil

## Example Usage

### Example 1: Get full community data

```go
package main

import (
    "context"
    "fmt"
    "github.com/my-eq/go-attom/pkg/client"
    "github.com/my-eq/go-attom/pkg/community"
)

func main() {
    c := client.NewClient("YOUR_API_KEY")
    svc := community.NewService(c)
    ctx := context.Background()
    
    // Get full community data by area ID
    data, err := svc.GetFullCommunityData(ctx, "12345")
    if err != nil {
        panic(err)
    }
    
    // Safely access demographics
    if data.Demographics != nil {
        if data.Demographics.Population != nil {
            fmt.Printf("Population: %d\n", *data.Demographics.Population)
        }
        if data.Demographics.MedianAge != nil {
            fmt.Printf("Median Age: %.1f\n", *data.Demographics.MedianAge)
        }
    }
    
    // Safely access economics
    if data.Economics != nil && data.Economics.MedianHouseholdIncome != nil {
        fmt.Printf("Median Income: $%d\n", *data.Economics.MedianHouseholdIncome)
    }
}
```

### Example 2: Check housing affordability

```go
params := &models.CommunityParams{
    GeoIDV4: stringPtr("abc123..."),
}

community, err := svc.GetNeighborhoodCommunity(ctx, params)
if err != nil {
    panic(err)
}

// Calculate affordability
if community.Housing != nil && 
   community.Housing.MedianHomeValue != nil &&
   community.Economics != nil &&
   community.Economics.MedianHouseholdIncome != nil {
    
    priceToIncome := float64(*community.Housing.MedianHomeValue) / 
                     float64(*community.Economics.MedianHouseholdIncome)
    
    fmt.Printf("Price-to-Income Ratio: %.2f\n", priceToIncome)
    
    if priceToIncome > 5.0 {
        fmt.Println("Housing is expensive for this area")
    }
}
```

### Example 3: Education summary

```go
if data.Education != nil {
    fmt.Println("Education Statistics:")
    
    if data.Education.HighSchoolGradRate != nil {
        fmt.Printf("  High School Graduation: %.1f%%\n", 
                  *data.Education.HighSchoolGradRate)
    }
    
    if data.Education.BachelorsRate != nil {
        fmt.Printf("  Bachelor's Degree: %.1f%%\n", 
                  *data.Education.BachelorsRate)
    }
    
    if data.Education.SchoolRating != nil {
        fmt.Printf("  School Rating: %.1f/10\n", 
                  *data.Education.SchoolRating)
    }
}
```

## Helper Functions

```go
package community

// GetAffordabilityIndex calculates housing affordability
func GetAffordabilityIndex(data *models.FullCommunityData) *float64 {
    if data.Housing == nil || data.Economics == nil {
        return nil
    }
    
    homeValue := data.Housing.MedianHomeValue
    income := data.Economics.MedianHouseholdIncome
    
    if homeValue == nil || income == nil || *income == 0 {
        return nil
    }
    
    ratio := float64(*homeValue) / float64(*income)
    return &ratio
}

// GetDemographicSummary creates a human-readable summary
func GetDemographicSummary(demo *models.Demographics) string {
    if demo == nil {
        return "No demographic data available"
    }
    
    summary := ""
    if demo.Population != nil {
        summary += fmt.Sprintf("Population: %d\n", *demo.Population)
    }
    if demo.MedianAge != nil {
        summary += fmt.Sprintf("Median Age: %.1f years\n", *demo.MedianAge)
    }
    if demo.AverageHouseholdSize != nil {
        summary += fmt.Sprintf("Avg Household Size: %.1f\n", 
                              *demo.AverageHouseholdSize)
    }
    
    return summary
}
```
