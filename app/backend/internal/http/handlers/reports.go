package handlers

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ReportsHandler агрегирует отчёты и отдаёт CSV и JSON
type ReportsHandler struct {
	DB *pgxpool.Pool
}

// GET /api/reports/summary
// Возвращает JSON с общей сводкой
func (h *ReportsHandler) SummaryJSON(c *gin.Context) {
	ctx := c.Request.Context()

	// --- читаем фильтры
	var (
		args       []any
		whereParts []string
	)
	whereParts = append(whereParts, "1=1")

	if pid, ok := queryInt(c, "project_id"); ok && pid > 0 {
		args = append(args, pid)
		whereParts = append(whereParts, fmt.Sprintf("d.project_id = $%d", len(args)))
	}
	if from, ok := queryDate(c, "from"); ok {
		args = append(args, from)
		whereParts = append(whereParts, fmt.Sprintf("d.created_at >= $%d", len(args)))
	}
	if to, ok := queryDate(c, "to"); ok {
		args = append(args, to.AddDate(0, 0, 1))
		whereParts = append(whereParts, fmt.Sprintf("d.created_at < $%d", len(args)))
	}
	where := strings.Join(whereParts, " AND ")

	// --- общий счетчик
	var total int
	err := h.DB.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM defects d WHERE %s", where), args...).Scan(&total)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// --- по статусам
	statusRows, err := h.DB.Query(ctx, fmt.Sprintf(`
		SELECT s.name as status, COUNT(*) as count
		FROM defects d
		JOIN statuses s ON s.id = d.status_id
		WHERE %s
		GROUP BY s.name
		ORDER BY count DESC
	`, where), args...)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer statusRows.Close()

	var byStatus []map[string]interface{}
	for statusRows.Next() {
		var status string
		var count int
		if err := statusRows.Scan(&status, &count); err != nil {
			continue
		}
		byStatus = append(byStatus, map[string]interface{}{
			"status": status,
			"count":  count,
		})
	}

	// --- по приоритетам
	priorityRows, err := h.DB.Query(ctx, fmt.Sprintf(`
		SELECT priority, COUNT(*) as count
		FROM defects d
		WHERE %s
		GROUP BY priority
		ORDER BY priority
	`, where), args...)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer priorityRows.Close()

	var byPriority []map[string]interface{}
	for priorityRows.Next() {
		var priority, count int
		if err := priorityRows.Scan(&priority, &count); err != nil {
			continue
		}
		byPriority = append(byPriority, map[string]interface{}{
			"priority": priority,
			"count":    count,
		})
	}

	// --- просроченные
	var overdue int
	err = h.DB.QueryRow(ctx, fmt.Sprintf(`
		SELECT COUNT(*)
		FROM defects d
		WHERE %s AND d.due_date < CURRENT_DATE 
		AND d.status_id NOT IN (SELECT id FROM statuses WHERE name IN ('closed', 'canceled'))
	`, where), args...).Scan(&overdue)
	if err != nil {
		overdue = 0 // не критично
	}

	c.JSON(200, gin.H{
		"total":       total,
		"overdue":     overdue,
		"by_status":   byStatus,
		"by_priority": byPriority,
	})
}

// GET /api/reports/trends
// Возвращает временные тренды
func (h *ReportsHandler) TrendsJSON(c *gin.Context) {
	ctx := c.Request.Context()

	groupBy := c.DefaultQuery("group", "week")
	if groupBy != "day" && groupBy != "week" && groupBy != "month" {
		groupBy = "week"
	}

	// --- читаем фильтры
	var (
		args       []any
		whereParts []string
	)
	whereParts = append(whereParts, "1=1")

	if pid, ok := queryInt(c, "project_id"); ok && pid > 0 {
		args = append(args, pid)
		whereParts = append(whereParts, fmt.Sprintf("d.project_id = $%d", len(args)))
	}
	if from, ok := queryDate(c, "from"); ok {
		args = append(args, from)
		whereParts = append(whereParts, fmt.Sprintf("d.created_at >= $%d", len(args)))
	}
	if to, ok := queryDate(c, "to"); ok {
		args = append(args, to.AddDate(0, 0, 1))
		whereParts = append(whereParts, fmt.Sprintf("d.created_at < $%d", len(args)))
	}
	where := strings.Join(whereParts, " AND ")

	// --- группировка по времени
	var dateExpr string
	switch groupBy {
	case "day":
		dateExpr = "DATE(d.created_at)"
	case "week":
		dateExpr = "DATE_TRUNC('week', d.created_at)"
	case "month":
		dateExpr = "DATE_TRUNC('month', d.created_at)"
	}

	rows, err := h.DB.Query(ctx, fmt.Sprintf(`
		SELECT %s as bucket,
			COUNT(CASE WHEN s.name = 'new' THEN 1 END) as new,
			COUNT(CASE WHEN s.name = 'in_work' THEN 1 END) as in_work,
			COUNT(CASE WHEN s.name = 'review' THEN 1 END) as review,
			COUNT(CASE WHEN s.name = 'closed' THEN 1 END) as closed,
			COUNT(CASE WHEN s.name = 'canceled' THEN 1 END) as canceled
		FROM defects d
		JOIN statuses s ON s.id = d.status_id
		WHERE %s
		GROUP BY %s
		ORDER BY bucket
	`, dateExpr, where, dateExpr), args...)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var series []map[string]interface{}
	for rows.Next() {
		var bucket time.Time
		var new, inWork, review, closed, canceled int
		if err := rows.Scan(&bucket, &new, &inWork, &review, &closed, &canceled); err != nil {
			continue
		}
		series = append(series, map[string]interface{}{
			"bucket":   bucket.Format("2006-01-02"),
			"new":      new,
			"in_work":  inWork,
			"review":   review,
			"closed":   closed,
			"canceled": canceled,
		})
	}

	c.JSON(200, gin.H{
		"series": series,
	})
}

// GET /api/reports/summary.csv
// Query: ?project_id=&from=&to=   (from/to в формате YYYY-MM-DD)
func (h *ReportsHandler) SummaryCSV(c *gin.Context) {
	ctx := c.Request.Context()

	// --- читаем фильтры
	var (
		args       []any
		whereParts []string
	)
	whereParts = append(whereParts, "1=1")

	if pid, ok := queryInt(c, "project_id"); ok && pid > 0 {
		args = append(args, pid)
		whereParts = append(whereParts, fmt.Sprintf("d.project_id = $%d", len(args)))
	}
	if from, ok := queryDate(c, "from"); ok {
		args = append(args, from)
		whereParts = append(whereParts, fmt.Sprintf("d.created_at >= $%d", len(args)))
	}
	if to, ok := queryDate(c, "to"); ok {
		// convention: to — не включительно следующего дня
		args = append(args, to.AddDate(0, 0, 1))
		whereParts = append(whereParts, fmt.Sprintf("d.created_at < $%d", len(args)))
	}
	where := strings.Join(whereParts, " AND ")

	// --- основной запрос по проектам
	// one row = one project
	// считаем total, статусы, приоритеты, просрочку и среднее время закрытия
	sql := `
SELECT
  p.id AS project_id,
  p.name AS project_name,

  COUNT(d.*)                              AS total,

  COUNT(d.*) FILTER (WHERE s.name = 'new')       AS new_cnt,
  COUNT(d.*) FILTER (WHERE s.name = 'in_work')   AS in_work_cnt,
  COUNT(d.*) FILTER (WHERE s.name = 'review')    AS review_cnt,
  COUNT(d.*) FILTER (WHERE s.name = 'closed')    AS closed_cnt,
  COUNT(d.*) FILTER (WHERE s.name = 'canceled')  AS canceled_cnt,

  COUNT(d.*) FILTER (
    WHERE d.due_date IS NOT NULL
      AND d.due_date < CURRENT_DATE
      AND s.name NOT IN ('closed','canceled')
  ) AS overdue_cnt,

  COUNT(d.*) FILTER (WHERE d.priority = 1) AS p1,
  COUNT(d.*) FILTER (WHERE d.priority = 2) AS p2,
  COUNT(d.*) FILTER (WHERE d.priority = 3) AS p3,
  COUNT(d.*) FILTER (WHERE d.priority = 4) AS p4,
  COUNT(d.*) FILTER (WHERE d.priority = 5) AS p5,

  COALESCE(
    AVG( EXTRACT(EPOCH FROM (d.updated_at - d.created_at)) / 86400.0 )
      FILTER (WHERE s.name = 'closed'),
    0
  ) AS avg_close_days

FROM projects p
LEFT JOIN defects d ON d.project_id = p.id
LEFT JOIN statuses s ON s.id = d.status_id
WHERE ` + where + `
GROUP BY p.id, p.name
ORDER BY p.id;
`

	rows, err := h.DB.Query(ctx, sql, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed: " + err.Error()})
		return
	}
	defer rows.Close()

	type row struct {
		ProjectID    int
		ProjectName  string
		Total        int
		New          int
		InWork       int
		Review       int
		Closed       int
		Canceled     int
		Overdue      int
		P1           int
		P2           int
		P3           int
		P4           int
		P5           int
		AvgCloseDays float64
	}
	var data []row

	for rows.Next() {
		var r row
		if err := rows.Scan(
			&r.ProjectID, &r.ProjectName,
			&r.Total,
			&r.New, &r.InWork, &r.Review, &r.Closed, &r.Canceled,
			&r.Overdue,
			&r.P1, &r.P2, &r.P3, &r.P4, &r.P5,
			&r.AvgCloseDays,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed: " + err.Error()})
			return
		}
		data = append(data, r)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rows err: " + err.Error()})
		return
	}

	// --- готовим CSV
	filename := "summary.csv"
	if d, ok := queryDate(c, "from"); ok {
		filename = "summary_from_" + d.Format("2006-01-02") + ".csv"
	}
	// заголовки ответа
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)

	// Добавим BOM, чтобы Excel на Windows открыл кириллицу нормально
	if _, err := c.Writer.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		// игнорируем, всё равно попробуем писать CSV дальше
	}

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	// Шапка
	_ = w.Write([]string{
		"project_id", "project_name",
		"total", "new", "in_work", "review", "closed", "canceled",
		"overdue",
		"p1", "p2", "p3", "p4", "p5",
		"avg_close_days",
	})

	// Строки
	for _, r := range data {
		_ = w.Write([]string{
			strconv.Itoa(r.ProjectID),
			r.ProjectName,
			strconv.Itoa(r.Total),
			strconv.Itoa(r.New),
			strconv.Itoa(r.InWork),
			strconv.Itoa(r.Review),
			strconv.Itoa(r.Closed),
			strconv.Itoa(r.Canceled),
			strconv.Itoa(r.Overdue),
			strconv.Itoa(r.P1),
			strconv.Itoa(r.P2),
			strconv.Itoa(r.P3),
			strconv.Itoa(r.P4),
			strconv.Itoa(r.P5),
			fmt.Sprintf("%.1f", r.AvgCloseDays),
		})
	}
}

// ----- вспомогательные парсеры -----

// queryInt читает число из query (?name=123)
func queryInt(c *gin.Context, name string) (int, bool) {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return 0, false
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return v, true
}

// queryDate читает дату YYYY-MM-DD из query (?from=2025-10-01)
func queryDate(c *gin.Context, name string) (time.Time, bool) {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return time.Time{}, false
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// (иногда полезно дернуть контекст напрямую)
func ctx(c *gin.Context) context.Context { return c.Request.Context() }
