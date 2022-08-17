package http

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/hierynomus/p1-monitor/internal/logging"
)

var _ http.Handler = (*TelegramHandler)(nil)

type TelegramHandler struct {
	ctx            context.Context
	mutex          sync.RWMutex
	LastUpdateTime time.Time
	Telegram       string
}

func NewTelegramHandler(ctx context.Context) *TelegramHandler {
	return &TelegramHandler{
		ctx: ctx,
	}
}

func (h *TelegramHandler) UpdateTelegram(telegram string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Telegram = telegram
	h.LastUpdateTime = time.Now()
}

func (h *TelegramHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	logger := logging.LoggerFor(h.ctx, "telegram-handler")
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Last-Modified", h.LastUpdateTime.Format(http.TimeFormat))
	if _, err := w.Write([]byte(h.Telegram)); err != nil {
		logger.Error().Err(err).Msg("Failed to write telegram")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
