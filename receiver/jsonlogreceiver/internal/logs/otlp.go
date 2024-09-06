// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logs // // import "github.com/justinmason/otel-collector/receiver/jsonlogreceiver/internal/logs"

import (
	"context"
	"log"
	"strings"
	"time"

	"encoding/json"

	p "github.com/justinmason/otel-collector/receiver/jsonlogreceiver/proto"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

const dataFormatProtobuf = "protobuf"

// Receiver is the type used to handle logs from OpenTelemetry exporters.
type Receiver struct {
	p.UnimplementedRemoteWriteServiceServer
	nextConsumer consumer.Logs
	logMap       map[string]string
	obsreport    *receiverhelper.ObsReport
}

// New creates a new Receiver reference.
func New(nextConsumer consumer.Logs, obsreport *receiverhelper.ObsReport, logTypes map[string][]string) *Receiver {
	logMap := make(map[string]string)
	for k, a := range logTypes {
		for _, v := range a {
			logMap[v] = k
		}
	}
	return &Receiver{
		nextConsumer: nextConsumer,
		logMap:       logMap,
		obsreport:    obsreport,
	}
}

func (r *Receiver) Write(ctx context.Context, req *p.WriteRequest) (*p.WriteResponse, error) {

	logs := plog.NewLogs()
	rlog := logs.ResourceLogs().AppendEmpty()
	scope := rlog.ScopeLogs().AppendEmpty()

	logLines := strings.Split(req.JsonData, "\n")
	for _, line := range logLines {

		var jsonLog map[string]interface{}
		err := json.Unmarshal([]byte(line), &jsonLog)
		if err != nil {
			//return &p.WriteResponse{}, err
			log.Printf("error unmarshaling json: %v\n'%s'", err, line)
			continue
		}

		data := jsonLog["data"].(map[string]interface{})

		lr := scope.LogRecords().AppendEmpty()

		for k1, v1 := range r.logMap {
			if v1 == "timestamp" {
				if value, ok := data[string(k1)]; ok {
					ts, err := time.Parse(time.RFC3339, value.(string))
					if err == nil {
						lr.SetTimestamp(pcommon.NewTimestampFromTime(ts))
					}
					break
				}
			}
		}
		lr.Body().SetStr(line)

		resAttribs := lr.Attributes()
		resAttribs.PutStr("label.Name 3", "label.Value 3")
		resAttribs.PutStr("label.Name 1", "label.Value 1")
		resAttribs.PutStr("label.Name 5", "label.Value 5")
		resAttribs.PutStr("label.Name 2", "label.Value 2")
		resAttribs.PutStr("label.Name 4", "label.Value 4")

	}
	ctx = r.obsreport.StartLogsOp(ctx)
	err := r.nextConsumer.ConsumeLogs(ctx, logs)
	r.obsreport.EndLogsOp(ctx, dataFormatProtobuf, 1, err)

	return &p.WriteResponse{}, err
}
