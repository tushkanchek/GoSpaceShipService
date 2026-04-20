package order_consumer

import (
	"context"
	"errors"
	"math"
	"time"

	"assembly/internal/config"
	"assembly/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"platform/pkg/kafka/consumer"
	"platform/pkg/logger"
)

var BuildTime = 10 * time.Second

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}
	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUuid.String()),
		zap.String("order_uuid", event.OrderUuid.String()),
		zap.String("user_uuid", event.UserUuid.String()),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUuid.String()),
	)

	// Produce OrderAssembled event and 10 seconds wait to simulate assembly time
	start := time.Now()

	timer := time.NewTimer(BuildTime)
	defer timer.Stop()

	select {
	case <-timer.C:

		assembledEvent := model.OrderAssembledEvent{
			EventUuid:    event.EventUuid,
			OrderUuid:    event.OrderUuid,
			UserUuid:     event.UserUuid,
			BuildTimeSec: int32(time.Since(start).Seconds()),
		}

		// Produce with retry
		err = s.produceWithRetry(ctx, assembledEvent)
		if err != nil {
			logger.Error(ctx, "Failed to produce OrderAssembled event after retries",
				zap.String("order_uuid", event.OrderUuid.String()),
				zap.Error(err),
			)
			return err
		}

		logger.Info(ctx, "Order assembled event produced successfully",
			zap.String("order_uuid", event.OrderUuid.String()),
		)

	case <-ctx.Done():
		logger.Info(ctx, "stopping order handler due to context cancellation")
		return ctx.Err()
	}

	return nil
}

func (s *service) produceWithRetry(ctx context.Context, event model.OrderAssembledEvent) error {
	retryConfig := config.AppConfig.Retry
	maxRetries := retryConfig.MaxRetries()
	initialDelay := retryConfig.InitialDelay()
	maxDelay := retryConfig.MaxDelay()
	backoffMultiplier := retryConfig.BackoffMultiplier()

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := s.orderAssembledProducer.ProduceOrderAssembled(ctx, event)

		// Для тестирования - форс ошибка если установлен флаг
		if err == nil && config.AppConfig.Testing.FailProduceForTesting() {
			err = errors.New("forced test error: production failed for testing")
		}

		if err == nil {
			if attempt > 0 {
				logger.Info(ctx, "Producer succeeded after retry",
					zap.Int("attempt", attempt+1),
					zap.String("order_uuid", event.OrderUuid.String()),
				)
			}
			return nil
		}

		lastErr = err
		if attempt < maxRetries-1 {
			delay := calculateBackoffDelay(attempt, initialDelay, maxDelay, backoffMultiplier)
			logger.Warn(ctx, "Producer failed, retrying",
				zap.Int("attempt", attempt+1),
				zap.Int("max_attempts", maxRetries),
				zap.Duration("next_retry_in", delay),
				zap.String("order_uuid", event.OrderUuid.String()),
				zap.Error(err),
			)

			select {
			case <-time.After(delay):
				// Continue to next retry
			case <-ctx.Done():
				logger.Info(ctx, "retry cancelled due to context cancellation")
				return ctx.Err()
			}
		}
	}

	logger.Error(ctx, "Producer failed after all retries",
		zap.Int("max_attempts", maxRetries),
		zap.String("order_uuid", event.OrderUuid.String()),
		zap.Error(lastErr),
	)

	// Отправляем в DLQ
	dlqErr := s.sendToDLQ(ctx, event, lastErr.Error())
	if dlqErr != nil {
		logger.Error(ctx, "Failed to send failed event to DLQ",
			zap.String("order_uuid", event.OrderUuid.String()),
			zap.Error(dlqErr),
		)
		// Если не удалось отправить в DLQ, пробиваем ошибку
		return dlqErr
	}

	// DLQ отправлен успешно, пропускаем обработку
	return nil
}

func calculateBackoffDelay(attempt int, initialDelay, maxDelay time.Duration, multiplier float64) time.Duration {
	delay := time.Duration(float64(initialDelay) * math.Pow(multiplier, float64(attempt)))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}

func (s *service) sendToDLQ(ctx context.Context, event model.OrderAssembledEvent, errorMsg string) error {
	// Создаём DLQ событие со структурированными данными
	dlqEvent := model.DeadLetterEvent{
		EventUuid:       uuid.New(),
		FailedEventType: "OrderAssembled",
		FailedEventData: map[string]interface{}{
			"event_uuid":     event.EventUuid.String(),
			"order_uuid":     event.OrderUuid.String(),
			"user_uuid":      event.UserUuid.String(),
			"build_time_sec": event.BuildTimeSec,
		},
		ErrorMessage: errorMsg,
		Timestamp:    time.Now().Unix(),
	}

	// Кодируем DLQ событие
	data, err := s.dlqEncoder.Encode(dlqEvent)
	if err != nil {
		logger.Error(ctx, "Failed to encode DLQ event", zap.Error(err))
		return err
	}

	// Отправляем в DLQ
	err = s.dlqProducer.Send(ctx, []byte(event.OrderUuid.String()), data)
	if err != nil {
		logger.Error(ctx, "Failed to send event to DLQ",
			zap.String("order_uuid", event.OrderUuid.String()),
			zap.Error(err),
		)
		return err
	}

	logger.Info(ctx, "Event sent to DLQ successfully",
		zap.String("order_uuid", event.OrderUuid.String()),
		zap.String("dlq_event_uuid", dlqEvent.EventUuid.String()),
	)

	return nil
}
