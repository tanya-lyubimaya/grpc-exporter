package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/tanya.lyubimaya/grpc-exporter/server"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MetricsServiceServer struct {
	server.UnimplementedExporterServer
	startTime time.Time
}

func NewServer() *MetricsServiceServer {
	return &MetricsServiceServer{
		startTime: time.Now(),
	}
}

func (s *MetricsServiceServer) CollectMetrics(context.Context, *emptypb.Empty) (*server.MetricsResponse, error) {
	cpuUsage, err := getCurrentCPUUsage()
	if err != nil {
		fmt.Printf("error getting CPU usage: %v\n", err)
	}
	memoryUsage, err := getCurrentMemoryUsage()
	if err != nil {
		fmt.Printf("error getting memory usage: %v\n", err)
	}
	uptime := calculateUptime(s.startTime)

	response := &server.MetricsResponse{
		CpuUsage:    &cpuUsage,
		MemoryUsage: &memoryUsage,
		Uptime:      &uptime,
	}

	return response, nil
}

func (s *MetricsServiceServer) StreamMetrics(_ *emptypb.Empty, stream server.Exporter_StreamMetricsServer) error {
	for {
		cpuUsage, err := getCurrentCPUUsage()
		if err != nil {
			return err
		}
		memoryUsage, err := getCurrentMemoryUsage()
		if err != nil {
			return err
		}
		uptime := calculateUptime(s.startTime)

		response := &server.MetricsResponse{
			CpuUsage:    &cpuUsage,
			MemoryUsage: &memoryUsage,
			Uptime:      &uptime,
		}

		if err := stream.Send(response); err != nil {
			return err
		}

		time.Sleep(5 * time.Second)
	}
}

func getCurrentCPUUsage() (server.Metric, error) {
	percentage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return server.Metric{}, err
	}

	return server.Metric{
		Name:  "cpu_usage",
		Help:  "Current CPU usage in percentage.",
		Value: percentage[0],
	}, nil
}

func getCurrentMemoryUsage() (server.Metric, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return server.Metric{}, err
	}

	return server.Metric{
		Name:  "memory_usage",
		Help:  "Current memory usage in MB.",
		Value: float64(memInfo.Used),
	}, nil
}

func calculateUptime(startTime time.Time) server.Metric {
	uptime := time.Since(startTime).Seconds()
	return server.Metric{
		Name:  "uptime_seconds",
		Help:  "Uptime of the service in seconds.",
		Value: uptime,
	}
}
