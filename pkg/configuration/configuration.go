package configuration

import filehandler "work-mini-project/pkg/fileHandler"

type CustomerConfig struct {
	FilePath string `json:"filePath"`
}

type CompanyConfig struct {
	GridX int `json:"gridX"`
	GridY int `json:"gridY"`
}

type UsersConfig struct {
	FilePath string `json:"filePath"`
}

type GridLimitsConfig struct {
	MinX int `json:"minX"`
	MaxX int `json:"maxX"`
	MinY int `json:"minY"`
	MaxY int `json:"maxY"`
}

type VehiclesConfig struct {
	Lorry struct {
		Speed                 int `json:"speed"`
		TrafficDelayTime      int `json:"trafficDelayTime"`
		TrafficDelayFrequency int `json:"trafficDelayFrequency"`
	} `json:"lorry"`
	CanalBoat struct {
		Speed int `json:"speed"`
	} `json:"canalBoat"`
	Helicopter struct {
		Speed        int `json:"speed"`
		InitialDelay int `json:"initialDelay"`
	} `json:"helicopter"`
}

type Config struct {
	Customers  CustomerConfig   `json:"customers"`
	Company    CompanyConfig    `json:"company"`
	Users      UsersConfig      `json:"users"`
	GridLimits GridLimitsConfig `json:"gridLimits"`
	Vehicles   VehiclesConfig   `json:"vehicles"`
}

func LoadConfig() (*Config, error) {
	config, err := filehandler.ReadFile[Config]("./config.json")
	if err != nil {
		return nil, err
	}

	return config, err
}
