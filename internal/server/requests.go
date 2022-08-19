//nolint:tagliatelle
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var (
	errReaderFail      = errors.New("failed to read request body")
	errInvalidPayload  = errors.New("request has invalid format")
	errRequiredMissing = errors.New("some required params not filled")
)

type bannerRequest struct {
	BannerID int `json:"banner_id"`
	SlotID   int `json:"slot_id"`
}

type clickRequest struct {
	BannerID int `json:"banner_id"`
	SlotID   int `json:"slot_id"`
	GroupID  int `json:"group_id"`
}

type showRequest struct {
	SlotID  int `json:"slot_id"`
	GroupID int `json:"group_id"`
}

func parseBannerRequest(r *http.Request) (bannerRequest, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return bannerRequest{}, errReaderFail
	}
	request := bannerRequest{}
	err = json.Unmarshal(data, &request)
	if err != nil {
		return bannerRequest{}, errInvalidPayload
	}
	if request.SlotID == 0 {
		return bannerRequest{}, fmt.Errorf("%w: slot_id", errRequiredMissing)
	}
	if request.BannerID == 0 {
		return bannerRequest{}, fmt.Errorf("%w: banner_id", errRequiredMissing)
	}
	return request, nil
}

func parseClickRequest(r *http.Request) (clickRequest, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return clickRequest{}, errReaderFail
	}
	request := clickRequest{}
	err = json.Unmarshal(data, &request)
	if err != nil {
		return clickRequest{}, errInvalidPayload
	}
	if request.SlotID == 0 {
		return clickRequest{}, fmt.Errorf("%w: slot_id", errRequiredMissing)
	}
	if request.BannerID == 0 {
		return clickRequest{}, fmt.Errorf("%w: banner_id", errRequiredMissing)
	}
	if request.GroupID == 0 {
		return clickRequest{}, fmt.Errorf("%w: group_id", errRequiredMissing)
	}
	return request, nil
}

func parseShowRequest(r *http.Request) (showRequest, error) {
	slotID := r.URL.Query().Get("slot_id")
	if slotID == "" {
		return showRequest{}, fmt.Errorf("%w: slot_id", errRequiredMissing)
	}
	iSlotID, err := strconv.Atoi(slotID)
	if err != nil {
		return showRequest{}, fmt.Errorf("%w: slot_id is not int", errInvalidPayload)
	}
	groupID := r.URL.Query().Get("slot_id")
	if groupID == "" {
		return showRequest{}, fmt.Errorf("%w: group_id", errRequiredMissing)
	}
	iGroupID, err := strconv.Atoi(groupID)
	if err != nil {
		return showRequest{}, fmt.Errorf("%w: group_id is not int", errInvalidPayload)
	}
	return showRequest{
		SlotID:  iSlotID,
		GroupID: iGroupID,
	}, nil
}
