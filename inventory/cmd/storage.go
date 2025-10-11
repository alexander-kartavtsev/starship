package main

import (
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func GetPartsByFilter(filter *inventoryV1.PartsFilter) map[string]*inventoryV1.Part {
	parts := getAll()
	if filter == nil {
		return parts
	}

	uuids := filter.GetUuids()
	if uuids == nil {
		return parts
	}

	res := map[string]*inventoryV1.Part{}
	for _, guid := range uuids {
		part, ok := parts[guid]
		if ok {
			res[part.Uuid] = part
		}
	}
	return res
}

func getAll() map[string]*inventoryV1.Part {
	parts := make(map[string]*inventoryV1.Part)

	uuids := []string{
		"a0ad507d-2b70-49e4-9378-3d92ebf9e405",
		"905ed12b-3934-45e1-a9af-67f00e00ff3d",
		"21c03a7f-0760-4d10-86a4-3273c025a3c3",
		"ff5c4e99-4c46-4422-8ff7-c8b7162b49c2",
		"b1450d16-5e0c-4685-b7a8-ef8eeb57e255",
		"1a1f5905-af52-4e0a-a295-b6547fac6313",
		"74162091-a328-4232-ae12-d8f5e31cd6b9",
		"e6bc160c-efa9-4f41-b585-c1bcbf2bf4c5",
		"cf3b7cda-b486-4a4c-bdc4-cb440b36c821",
		"dcfc15a9-5bd2-4897-9d2b-bcbda59d113a",
	}

	for _, guid := range uuids {
		part := &inventoryV1.Part{}
		_ = gofakeit.Struct(part)
		part.Uuid = guid
		part.CreatedAt = timestamppb.New(time.Now())
		part.UpdatedAt = timestamppb.New(time.Now())
		parts[guid] = part
	}

	return parts
}
