package main

type InfluxConfig struct {
	Measurement string		`json:"measurement"`
	Host 		string		`json:"host"`
	DB 			string		`json:"db"`
	Precision	string		`json:"precision"`
	Retention 	string		`json:"retention"`
}
