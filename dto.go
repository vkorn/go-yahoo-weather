package weather

// Weather defines Yahoo current weather info and forecast.
type Weather struct {
	Location    *Location       `json:"location"`
	Observation *Observation    `json:"current_observation"`
	Forecasts   []*ForecastInfo `json:"forecasts"`
}

// Observation defines current weather conditions.
type Observation struct {
	Wind       *WindInfo       `json:"wind"`
	Atmosphere *AtmosphereInfo `json:"atmosphere"`
	Astronomy  *AstronomyInfo  `json:"astronomy"`
	Condition  *ConditionInfo  `json:"condition"`
}

// Location defines parsed location from the request.
type Location struct {
	WoeID    int     `json:"woeid"`
	City     string  `json:"city"`
	Region   string  `json:"region"`
	Country  string  `json:"country"`
	Lat      float32 `json:"lat"`
	Lon      float32 `json:"long"`
	Timezone string  `json:"timezone_id"`
}

// WindInfo defines current wind conditions.
type WindInfo struct {
	Chill     int     `json:"chill"`
	Direction int     `json:"direction"`
	Speed     float32 `json:"speed"`
}

// PressureStates describes barometric pressure state.
type PressureStates int

const (
	PressureSteady  PressureStates = iota
	PressureRising                 = 1
	PressureFalling                = 2
)

// AtmosphereInfo defines current humidity, visibility and pressure.
type AtmosphereInfo struct {
	Humidity   int            `json:"humidity"`
	Visibility float32        `json:"visibility"`
	Pressure   float32        `json:"pressure"`
	State      PressureStates `json:"rising"`
}

// AstronomyInfo defines current astronomical conditions.
type AstronomyInfo struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}

// ConditionCodes describes weather condition enum.
type ConditionCodes int

const (
	ConditionTornado                   ConditionCodes = iota
	ConditionTropicalStorm                            = 1
	ConditionHurricane                                = 2
	ConditionSevereThunderstorms                      = 3
	ConditionThunderstorms                            = 4
	ConditionMixedRainAndSnow                         = 5
	ConditionMixedRainAndSleet                        = 6
	ConditionMixedSnowAndSleet                        = 7
	ConditionFreezingDrizzle                          = 8
	ConditionDrizzle                                  = 9
	ConditionFreezingRain                             = 10
	ConditionShowers                                  = 11
	ConditionRain                                     = 12
	ConditionSnowFlurries                             = 13
	ConditionLightSnowShowers                         = 14
	ConditionBlowingSnow                              = 15
	ConditionSnow                                     = 16
	ConditionHail                                     = 17
	ConditionSleet                                    = 18
	ConditionDust                                     = 19
	ConditionFoggy                                    = 20
	ConditionHaze                                     = 21
	ConditionSmoky                                    = 22
	ConditionBlustery                                 = 23
	ConditionWindy                                    = 24
	ConditionCold                                     = 25
	ConditionCloudy                                   = 26
	ConditionMostlyCloudyNight                        = 27
	ConditionMostlyCloudyDay                          = 28
	ConditionPartlyCloudyNight                        = 29
	ConditionPartlyCloudyDay                          = 30
	ConditionClearNight                               = 31
	ConditionSunny                                    = 32
	ConditionFairNight                                = 33
	ConditionFairDay                                  = 34
	ConditionMixedRainAndHail                         = 35
	ConditionHot                                      = 36
	ConditionIsolatedThunderstorms                    = 37
	ConditionScatteredThunderstorms                   = 38
	ConditionScatteredShowersDay                      = 39
	ConditionHeavyRain                                = 40
	ConditionScatteredSnowShowersDay                  = 41
	ConditionHeavySnow                                = 42
	ConditionBlizzard                                 = 43
	ConditionNA                                       = 44
	ConditionScatteredShowersNight                    = 45
	ConditionScatteredSnowShowersNight                = 46
	ConditionScatteredThundershowers                  = 47
)

// ConditionInfo defines current temperature and condition.
type ConditionInfo struct {
	Text        string         `json:"text"`
	Code        ConditionCodes `json:"code"`
	Temperature int            `json:"temperature"`
}

// ForecastInfo defines forecast information.
type ForecastInfo struct {
	Day  string         `json:"day"`
	Date int64          `json:"date"`
	Low  int            `json:"low"`
	High int            `json:"high"`
	Text string         `json:"text"`
	Code ConditionCodes `json:"code"`
}
