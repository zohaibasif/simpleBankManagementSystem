package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/zohaibAsif/simple_bank_management_system/db/sqlc"
)

type getEntryRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getEntry(ctx *gin.Context) {
	uParams := getEntryRequest{}

	if err := ctx.ShouldBindUri(&uParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := s.store.GetEntry(ctx, uParams.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

type listEntryQuerryParams struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listEntries(ctx *gin.Context) {
	uParams := getEntryRequest{}
	qParams := listEntryQuerryParams{}

	if err := ctx.ShouldBindUri(&uParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&qParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListEntriesParams{
		AccountID: uParams.Id,
		Limit:     qParams.PageSize,
		Offset:    (qParams.PageId - 1) * qParams.PageSize,
	}

	entries, err := s.store.ListEntries(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}
