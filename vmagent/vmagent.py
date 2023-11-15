import os
import threading

import grpc
import sys

import exporter_pb2_grpc
import google.protobuf.empty_pb2

sys.path.append("server")


def collect_metrics(request, context):
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = exporter_pb2_grpc.ExporterStub(channel)
        response = stub.CollectMetrics(request, context)
        print("[UNARY] received metrics response:", response)
    return response


class ExporterServicer(exporter_pb2_grpc.ExporterServicer):
    pass


def call_collect_metrics():
    ExporterServicer()
    empty = google.protobuf.empty_pb2.Empty()
    collect_metrics(empty, None)
    threading.Timer(5.0, call_collect_metrics).start()


def run_server_streaming_client():
    channel = grpc.insecure_channel('localhost:50051')  # Replace with your server address
    stub = exporter_pb2_grpc.ExporterStub(channel)

    request = google.protobuf.empty_pb2.Empty()

    stream = stub.StreamMetrics(request)

    try:
        for response in stream:
            print("[STREAMING] received metrics response:", response)
    except grpc.RpcError as e:
        print(f"error while reading stream: {e}")


if __name__ == '__main__':
    mode = os.environ['MODE']
    if mode == 'STREAM':
        run_server_streaming_client()
    elif mode == 'UNARY':
        call_collect_metrics()
    else:
        print("unknown mode provided")

