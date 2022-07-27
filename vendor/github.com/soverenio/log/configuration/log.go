// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package configuration

// Log holds configuration for logging
type Log struct {
	Level               string `insconfig:"info| Default level for logger"`
	Adapter             string `insconfig:"bilog| Logging adapter - zerolog (json, text) or bilog (json, text, pbuf)"`
	Formatter           string `insconfig:"json| Log output format - e.g. json or text"`
	OutputType          string `insconfig:"stderr| Log output type - e.g. stdout, stderr, syslog"`
	OutputParallelLimit string `insconfig:"| Write-parallel limit for the output"`
	OutputParams        string `insconfig:"| Parameter for output - depends on OutputType"`
	BufferSize          int    `insconfig:"0| Number of regular log events that can be buffered, =0 to disable"`
	LLBufferSize        int    `insconfig:"0| Number of low-latency log events that can be buffered, =-1 to disable, =0 - default size"`
}

// NewLog creates new default configuration for logging
func NewLog() Log {
	return Log{
		Level:      "info",
		Adapter:    "bilog",
		Formatter:  "json",
		OutputType: "stderr",
		BufferSize: 0,
	}
}
