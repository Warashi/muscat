package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/atotto/clipboard"
	"github.com/skratchdot/open-golang/open"

	"github.com/Warashi/muscat/pb"
)

type Muscat struct {
	pb.UnimplementedMuscatServer
}

func (m *Muscat) Open(ctx context.Context, request *pb.OpenRequest) (*pb.OpenResponse, error) {
	if err := open.Run(request.Uri); err != nil {
		return nil, fmt.Errorf("open.Run: %w", err)
	}
	return new(pb.OpenResponse), nil
}

func (m *Muscat) Copy(stream pb.Muscat_CopyServer) error {
	var buf bytes.Buffer

	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("stream.Recv: %w", err)
		}
		n, err := buf.Write(msg.GetBody())
		if err != nil {
			return fmt.Errorf("buf.Write: %w", err)
		}
		if n != len(msg.GetBody()) {
			return fmt.Errorf("could not write whole body")
		}
	}
	if err := clipboard.WriteAll(buf.String()); err != nil {
		return fmt.Errorf("clipboard.WriteAll: %w", err)
	}
	return nil
}

func (m *Muscat) Paste(_ *pb.PasteRequest, stream pb.Muscat_PasteServer) error {
	body, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("clipboard.ReadAll: %w", err)
	}
	for i := 0; i < len(body); i += 1024 {
		start, end := i, i+1024
		if len(body) <= end {
			end = len(body)
		}
		if err := stream.Send(&pb.PasteResponse{Body: []byte(body[start:end])}); err != nil {
			return fmt.Errorf("stream.Send: %w", err)
		}
	}
	return nil
}
