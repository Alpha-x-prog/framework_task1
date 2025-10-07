package handlers

import (
	"encoding/csv"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportsHandler struct{ DB *pgxpool.Pool }

func (h *ReportsHandler) SummaryCSV(c *gin.Context) {
	rows, err := h.DB.Query(c, `
SELECT p.name as project, s.name as status, COUNT(*) as cnt
FROM defects d
JOIN projects p ON p.id=d.project_id
JOIN statuses s ON s.id=d.status_id
GROUP BY p.name, s.name
ORDER BY p.name, s.name`)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=summary.csv")
	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"project", "status", "count"})
	for rows.Next() {
		var proj, st string
		var cnt int
		if err := rows.Scan(&proj, &st, &cnt); err != nil {
			continue
		}
		_ = w.Write([]string{proj, st, strconv.Itoa(cnt)})
	}
	w.Flush()
}
