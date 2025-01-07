package app

import (
	"net/http"

	"github.com/inconshreveable/log15" //nolint:logging // TODO move all logging to sourcegraph/log

	"github.com/khulnasoft/khulnasoft/internal/auth"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/usagestats"
)

func usageStatsArchiveHandler(db database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 🚨SECURITY: Only site admins may get this archive.
		if err := auth.CheckCurrentUserIsSiteAdmin(r.Context(), db); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\"SourcegraphUsersUsageArchive.zip\"")

		archive, err := usagestats.GetArchive(r.Context(), db)
		if err != nil {
			log15.Error("usagestats.WriteArchive", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(archive)
	}
}
