// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package jsonlogreceiver // import "go.opentelemetry.io/collector/receiver/otlpreceiver"

import (
	"fmt"
	"io"
	"mime"
	"net/http"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"

	"github.com/justinmason/otel-collector/receiver/jsonlogreceiver/internal/errors"
	"github.com/justinmason/otel-collector/receiver/jsonlogreceiver/internal/logs"
)

// Pre-computed status with code=Internal to be used in case of a marshaling error.
var fallbackMsg = []byte(`{"code": 13, "message": "failed to marshal error message"}`)

const fallbackContentType = "application/json"

func handleLogs(resp http.ResponseWriter, req *http.Request, logsReceiver *logs.Receiver) {
	enc, ok := readContentType(resp, req)
	if !ok {
		return
	}

	body, ok := readAndCloseBody(resp, req, enc)
	if !ok {
		return
	}

	logReq, err := enc.unmarshalWriteRequest(body)
	if err != nil {
		writeError(resp, enc, err, http.StatusBadRequest)
		return
	}

	_, err = logsReceiver.Write(req.Context(), &logReq)
	if err != nil {
		writeError(resp, enc, err, http.StatusInternalServerError)
		return
	}

	writeResponse(resp, enc.contentType(), http.StatusOK, nil)
}
func readContentType(resp http.ResponseWriter, req *http.Request) (encoder, bool) {
	if req.Method != http.MethodPost {
		handleUnmatchedMethod(resp)
		return nil, false
	}

	switch getMimeTypeFromContentType(req.Header.Get("Content-Type")) {
	case pbContentType:
		return pbEncoder, true
	case jsonContentType:
		return jsEncoder, true
	default:
		handleUnmatchedContentType(resp)
		return nil, false
	}
}

func readAndCloseBody(resp http.ResponseWriter, req *http.Request, enc encoder) ([]byte, bool) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeError(resp, enc, err, http.StatusBadRequest)
		return nil, false
	}
	if err = req.Body.Close(); err != nil {
		writeError(resp, enc, err, http.StatusBadRequest)
		return nil, false
	}
	return body, true
}

// writeError encodes the HTTP error inside a rpc.Status message as required by the OTLP protocol.
func writeError(w http.ResponseWriter, encoder encoder, err error, statusCode int) {
	s, ok := status.FromError(err)
	if ok {
		statusCode = errors.GetHTTPStatusCodeFromStatus(s)
	} else {
		s = errors.NewStatusFromMsgAndHTTPCode(err.Error(), statusCode)
	}
	writeStatusResponse(w, encoder, statusCode, s.Proto())
}

// errorHandler encodes the HTTP error message inside a rpc.Status message as required
// by the OTLP protocol.
func errorHandler(w http.ResponseWriter, r *http.Request, errMsg string, statusCode int) {
	s := errors.NewStatusFromMsgAndHTTPCode(errMsg, statusCode)
	switch getMimeTypeFromContentType(r.Header.Get("Content-Type")) {
	case pbContentType:
		writeStatusResponse(w, pbEncoder, statusCode, s.Proto())
		return
	case jsonContentType:
		writeStatusResponse(w, jsEncoder, statusCode, s.Proto())
		return
	}
	writeResponse(w, fallbackContentType, http.StatusInternalServerError, fallbackMsg)
}

func writeStatusResponse(w http.ResponseWriter, enc encoder, statusCode int, rsp *spb.Status) {
	msg, err := enc.marshalStatus(rsp)
	if err != nil {
		writeResponse(w, fallbackContentType, http.StatusInternalServerError, fallbackMsg)
		return
	}

	writeResponse(w, enc.contentType(), statusCode, msg)
}

func writeResponse(w http.ResponseWriter, contentType string, statusCode int, msg []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	// Nothing we can do with the error if we cannot write to the response.
	_, _ = w.Write(msg)
}

func getMimeTypeFromContentType(contentType string) string {
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return ""
	}
	return mediatype
}

func handleUnmatchedMethod(resp http.ResponseWriter) {
	status := http.StatusMethodNotAllowed
	writeResponse(resp, "text/plain", status, []byte(fmt.Sprintf("%v method not allowed, supported: [POST]", status)))
}

func handleUnmatchedContentType(resp http.ResponseWriter) {
	status := http.StatusUnsupportedMediaType
	writeResponse(resp, "text/plain", status, []byte(fmt.Sprintf("%v unsupported media type, supported: [%s, %s]", status, jsonContentType, pbContentType)))
}
