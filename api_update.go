package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func UpdateHandler(ctx *Context) {
	request, err := NewUpdateRequest(ctx.Request)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.HasData() {
		ctx.Response.Header().Set("Content-Type", "application/json")
		ctx.Response.WriteHeader(http.StatusOK)
		json.NewEncoder(ctx.Response).Encode(ctx.Server.Manager.LatestRelease)
		return
	}

	releaseType := ctx.Server.Config.GetReleaseTypeByIdentifier(request.releaseType)
	if releaseType == nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	targetItem := ctx.Server.Manager.LatestRelease.GetItemByFilename(releaseType.Filename)
	if targetItem == nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	if targetItem.Checksum == request.releaseChecksum {
		ctx.Response.WriteHeader(http.StatusOK)
		return
	}

	// Redirect to new release file
	targetUrl := fmt.Sprintf("/releases/%s/%s", ctx.Server.Manager.LatestRelease.TagName, targetItem.Filename)
	http.Redirect(ctx.Response, ctx.Request, targetUrl, http.StatusFound)
}

type UpdateRequest struct {
	releaseType     string
	releaseChecksum string
}

func (request *UpdateRequest) HasData() bool {
	return request.releaseType != "" && request.releaseChecksum != ""
}

func NewUpdateRequest(request *http.Request) (*UpdateRequest, error) {
	query := request.URL.Query()
	releaseType := query.Get("type")
	releaseChecksum := query.Get("checksum")

	// Either both are empty or both are given
	isEmpty := releaseType == "" && releaseChecksum == ""
	isGiven := releaseType != "" && releaseChecksum != ""

	if isEmpty || isGiven {
		return &UpdateRequest{
			releaseType:     releaseType,
			releaseChecksum: releaseChecksum,
		}, nil
	}

	return nil, fmt.Errorf("invalid update request")
}
