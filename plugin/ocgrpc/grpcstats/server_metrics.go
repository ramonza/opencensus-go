// Copyright 2017, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package grpcstats

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// The following variables are measures and views made available for gRPC clients.
// Server needs to use a ServerStatsHandler in order to enable collection.
var (
	ServerErrorCount, _        = stats.Int64("grpc.io/server/error_count", "RPC Errors", stats.UnitNone)
	ServerServerElapsedTime, _ = stats.Float64("grpc.io/server/server_elapsed_time", "Server elapsed time in msecs", stats.UnitMilliseconds)
	ServerRequestBytes, _      = stats.Int64("grpc.io/server/request_bytes", "Request bytes", stats.UnitBytes)
	ServerResponseBytes, _     = stats.Int64("grpc.io/server/response_bytes", "Response bytes", stats.UnitBytes)
	ServerStartedCount, _      = stats.Int64("grpc.io/server/started_count", "Number of server RPCs (streams) started", stats.UnitNone)
	ServerFinishedCount, _     = stats.Int64("grpc.io/server/finished_count", "Number of server RPCs (streams) finished", stats.UnitNone)
	ServerRequestCount, _      = stats.Int64("grpc.io/server/request_count", "Number of server RPC request messages", stats.UnitNone)
	ServerResponseCount, _     = stats.Int64("grpc.io/server/response_count", "Number of server RPC response messages", stats.UnitNone)
)

// TODO(acetechnologist): This is temporary and will need to be replaced by a
// mechanism to load these defaults from a common repository/config shared by
// all supported languages. Likely a serialized protobuf of these defaults.

// Predefined server views
var (
	RPCServerErrorCountView, _ = view.New(
		"grpc.io/server/error_count",
		"RPC Errors",
		[]tag.Key{KeyMethod, KeyStatus},
		ServerErrorCount,
		view.CountAggregation{})
	RPCServerServerElapsedTimeView, _ = view.New(
		"grpc.io/server/server_elapsed_time",
		"Server elapsed time in msecs",
		[]tag.Key{KeyMethod},
		ServerServerElapsedTime,
		DefaultMillisecondsDistribution)
	RPCServerRequestBytesView, _ = view.New(
		"grpc.io/server/request_bytes",
		"Request bytes",
		[]tag.Key{KeyMethod},
		ServerRequestBytes,
		DefaultBytesDistribution)
	RPCServerResponseBytesView, _ = view.New(
		"grpc.io/server/response_bytes",
		"Response bytes",
		[]tag.Key{KeyMethod},
		ServerResponseBytes,
		DefaultBytesDistribution)
	RPCServerRequestCountView, _ = view.New(
		"grpc.io/server/request_count",
		"Count of request messages per server RPC",
		[]tag.Key{KeyMethod},
		ServerRequestCount,
		DefaultMessageCountDistribution)
	RPCServerResponseCountView, _ = view.New(
		"grpc.io/server/response_count",
		"Count of response messages per server RPC",
		[]tag.Key{KeyMethod},
		ServerResponseCount,
		DefaultMessageCountDistribution)

	DefaultServerViews = []*view.View{
		RPCServerErrorCountView,
		RPCServerServerElapsedTimeView,
		RPCServerRequestBytesView,
		RPCServerResponseBytesView,
		RPCServerRequestCountView,
		RPCServerResponseCountView,
	}
)
// TODO(jbd): Add roundtrip_latency, uncompressed_request_bytes, uncompressed_response_bytes, request_count, response_count.
