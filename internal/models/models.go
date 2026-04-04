package models

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Coordinates struct {
	XMin float64 `json:"x_min"`
	XMax float64 `json:"x_max" validate:"gtfield=XMin"`
	YMin float64 `json:"y_min"`
	YMax float64 `json:"y_max" validate:"gtfield=YMin"`
}

type SimulationConfig struct {
	UECnt            int         `json:"ue_cnt" validate:"required,min=1,max=1000"`
	UEIDs            []int       `json:"ue_ids" validate:"required"`
	UEMovePattern    string      `json:"ue_move_pattern" validate:"required,oneof=RandomWalk Static ConstantVelocity"`
	UECoords         Coordinates `json:"ue_coords" validate:"required"`
	UETrafficPattern string      `json:"ue_traffic_pattern" validate:"required,oneof=Poisson CBR FTP"`
	UEPause          float64     `json:"ue_pause,omitempty" validate:"min=0"`
	BSScheduler      string      `json:"bs_scheduler" validate:"required,oneof=BestCQI RoundRobin ProportionalFair"`
	BSCoords         Point       `json:"bs_coords" validate:"required"`
	BSBwMHz          int         `json:"bs_bw_mhz" validate:"required,oneof=5 10 15 20"`
	SimPacketRate    int         `json:"sim_packet_rate" validate:"required,min=1"`
}

type Response struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Data   interface{} `json:"data,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
}
