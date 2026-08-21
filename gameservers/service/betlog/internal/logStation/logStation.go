package logStation

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"os"
	"serve/comm/mq"
	"serve/comm/slotsmongo"
	"sync"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

type LogService struct {
	buffer              []*slotsmongo.DocBetLog
	bufferMutex         sync.Mutex
	clickhouseDB        *sql.DB
	flushInterval       time.Duration
	ctx                 context.Context
	BatchSize           int
	logSignal           chan *slotsmongo.DocBetLog // 用于接收新日志的信号
	ReTryTimes          int
	minibuffer          []*slotsmongo.DocBetLogAviator
	minilogSignal       chan *slotsmongo.DocBetLogAviator // 用于接收新日志的信号
	miniUpdatebuffer    []*slotsmongo.DocBetLogAviator
	miniUpdatelogSignal chan *slotsmongo.DocBetLogAviator // 用于接收新日志的信号
	miniflushInterval   time.Duration
}

func NewLogService(ctx context.Context, clickHouseAddr string, BatchSize, reTryTimes int, flushInterval time.Duration, miniflushInterval time.Duration) (*LogService, error) {
	conn, err := sql.Open("clickhouse", clickHouseAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	ls := &LogService{
		buffer:              make([]*slotsmongo.DocBetLog, 0),
		BatchSize:           BatchSize,
		clickhouseDB:        conn,
		flushInterval:       flushInterval,
		ctx:                 ctx,
		logSignal:           make(chan *slotsmongo.DocBetLog),
		ReTryTimes:          reTryTimes,
		minibuffer:          make([]*slotsmongo.DocBetLogAviator, 0),
		minilogSignal:       make(chan *slotsmongo.DocBetLogAviator),
		miniUpdatebuffer:    make([]*slotsmongo.DocBetLogAviator, 0),
		miniUpdatelogSignal: make(chan *slotsmongo.DocBetLogAviator),
		miniflushInterval:   miniflushInterval,
	}

	return ls, nil
}

func (ls *LogService) Start(c chan os.Signal) {
	go ls.subscribeLogs()
	ls.processAndFlushLogs(c) // 启动合并后的处理 Goroutine
}

func (ls *LogService) Close() {
	close(ls.logSignal)
	if ls.clickhouseDB != nil {
		ls.clickhouseDB.Close()
	}
}

func (ls *LogService) subscribeLogs() {
	_, err := mq.JsonNC.Subscribe("/games/betlog", func(logEntry *slotsmongo.DocBetLog) {
		fmt.Println(logEntry)
		ls.logSignal <- logEntry // 将接收到的日志发送到 logSignal
	})
	if err != nil {
		slog.Error("Failed to subscribe to subject: %v", err)
	}
	_, err = mq.JsonNC.Subscribe("/center/transfer_log", func(logEntry *slotsmongo.DocBetLog) {
		fmt.Println(logEntry)
		ls.logSignal <- logEntry // 将接收到的日志发送到 logSignal
	})
	if err != nil {
		slog.Error("Failed to subscribe to subject: %v", err)
	}
	_, err = mq.JsonNC.Subscribe("/games/minibetlog", func(minilogEntry *slotsmongo.DocBetLogAviator) {
		fmt.Println(minilogEntry)
		ls.minilogSignal <- minilogEntry // 将接收到的日志发送到 logSignal
	})
	if err != nil {
		slog.Error("Failed to subscribe to subject: %v", err)
	}
	_, err = mq.JsonNC.Subscribe("/games/minibetlogUpdate", func(minilogupdateEntry *slotsmongo.DocBetLogAviator) {
		fmt.Println(minilogupdateEntry)
		ls.miniUpdatelogSignal <- minilogupdateEntry // 将接收到的日志发送到 logSignal
	})
	if err != nil {
		slog.Error("Failed to subscribe to subject: %v", err)
	}
}

func (ls *LogService) taskReTry(batchToFlush []*slotsmongo.DocBetLog) {
	for i := 0; i <= ls.ReTryTimes; i++ {
		slog.Error("Failed to flush batch to ClickHouse, retry attempt: %d/%d, batch: %v", i+1, ls.ReTryTimes+1, batchToFlush)
		if err := ls.flushBatch(batchToFlush); err == nil {
			slog.Info("Successfully flushed batch after %d retries.", i+1)
			return
		}

		if i < ls.ReTryTimes {
			// 指数退避策略，可以根据需要调整 base 和 maxDelay
			baseDelay := 1 * time.Second
			maxDelay := 60 * time.Second
			delay := baseDelay * time.Duration(math.Pow(2, float64(i)))
			if delay > maxDelay {
				delay = maxDelay
			}
			slog.Warn("Retrying in %v...", delay)
			time.Sleep(delay)
		}
	}
	slog.Error("Failed to flush batch after all %d retries.", ls.ReTryTimes+1)
}

func (ls *LogService) processAndFlushLogs(c chan os.Signal) {
	ticker := time.NewTicker(ls.flushInterval)
	defer ticker.Stop()
	miniticker := time.NewTicker(ls.miniflushInterval)
	defer miniticker.Stop()

	for {
		select {
		case logEntry := <-ls.logSignal:
			ls.bufferMutex.Lock()
			ls.buffer = append(ls.buffer, logEntry)
			if len(ls.buffer) >= ls.BatchSize {
				// 达到阈值，触发批量更新
				batchToFlush := make([]*slotsmongo.DocBetLog, len(ls.buffer))
				copy(batchToFlush, ls.buffer)
				ls.buffer = ls.buffer[:0]
				ls.bufferMutex.Unlock()
				fmt.Println("logSignal Flushing betlogs...")
				if err := ls.flushBatch(batchToFlush); err != nil {
					go ls.taskReTry(batchToFlush)
					slog.Error("Failed to flush batch due to threshold: %v", err)
				}
			} else {
				ls.bufferMutex.Unlock()
			}
		case <-ticker.C:
			// 定时器到期，触发强制更新
			ls.bufferMutex.Lock()
			if len(ls.buffer) > 0 {
				batchToFlush := make([]*slotsmongo.DocBetLog, len(ls.buffer))
				copy(batchToFlush, ls.buffer)
				ls.buffer = ls.buffer[:0]
				ls.bufferMutex.Unlock()
				fmt.Println("ticker Flushing betlogs...")
				if err := ls.flushBatch(batchToFlush); err != nil {
					go ls.taskReTry(batchToFlush)
					slog.Error("Failed to flush batch due to ticker: %v", err)
				}
			} else if len(ls.minibuffer) > 0 {
				batchToFlush := make([]*slotsmongo.DocBetLogAviator, len(ls.minibuffer))
				copy(batchToFlush, ls.minibuffer)
				ls.minibuffer = ls.minibuffer[:0]
				ls.bufferMutex.Unlock()
				slog.Info("ticker 准备插入投注记录")
				if err := ls.flushBatchAviator(batchToFlush); err != nil {
					go ls.taskReTryAviator(batchToFlush)
					slog.Error(fmt.Sprintf("Failed to flush batch due to ticker: %v", err))
				}
			} else if len(ls.miniUpdatebuffer) > 0 {
				// 达到阈值，触发批量更新
				batchToFlush := make([]*slotsmongo.DocBetLogAviator, len(ls.miniUpdatebuffer))
				copy(batchToFlush, ls.miniUpdatebuffer)
				ls.miniUpdatebuffer = ls.miniUpdatebuffer[:0]
				ls.bufferMutex.Unlock()
				slog.Info("ticker 准备更新投注记录")
				if err := ls.flushBatchAviatorUpdate(batchToFlush); err != nil {
					go ls.taskReTryAviatorUpdate(batchToFlush)
					slog.Error(fmt.Sprintf("Failed to flush batch due to threshold: %v", err))
				}
			} else {
				ls.bufferMutex.Unlock()
			}
		case <-miniticker.C:
			// 定时器到期，触发强制更新
			ls.bufferMutex.Lock()
			if len(ls.minibuffer) > 0 {
				batchToFlush := make([]*slotsmongo.DocBetLogAviator, len(ls.minibuffer))
				copy(batchToFlush, ls.minibuffer)
				ls.minibuffer = ls.minibuffer[:0]
				ls.bufferMutex.Unlock()
				slog.Info("ticker 准备插入投注记录")
				if err := ls.flushBatchAviator(batchToFlush); err != nil {
					go ls.taskReTryAviator(batchToFlush)
					slog.Error(fmt.Sprintf("Failed to flush batch due to ticker: %v", err))
				}
			} else if len(ls.miniUpdatebuffer) > 0 {
				// 达到阈值，触发批量更新
				batchToFlush := make([]*slotsmongo.DocBetLogAviator, len(ls.miniUpdatebuffer))
				copy(batchToFlush, ls.miniUpdatebuffer)
				ls.miniUpdatebuffer = ls.miniUpdatebuffer[:0]
				ls.bufferMutex.Unlock()
				slog.Info("ticker 准备更新投注记录")
				if err := ls.flushBatchAviatorUpdate(batchToFlush); err != nil {
					go ls.taskReTryAviatorUpdate(batchToFlush)
					slog.Error(fmt.Sprintf("Failed to flush batch due to threshold: %v", err))
				}
			} else {
				ls.bufferMutex.Unlock()
			}
		case <-c:
			slog.Info("LogServer processing stopped.")
			// 在退出前刷新剩余的日志
			ls.bufferMutex.Lock()
			if len(ls.buffer) > 0 {
				batchToFlush := make([]*slotsmongo.DocBetLog, len(ls.buffer))
				copy(batchToFlush, ls.buffer)
				ls.buffer = ls.buffer[:0]
				ls.bufferMutex.Unlock()
				if err := ls.flushBatch(batchToFlush); err != nil {
					ls.taskReTry(batchToFlush)
					slog.Error("Failed to flush remaining batch on shutdown: %v", err)
				}
			} else if len(ls.minibuffer) > 0 {
				batchToFlush := make([]*slotsmongo.DocBetLogAviator, len(ls.minibuffer))
				copy(batchToFlush, ls.minibuffer)
				ls.minibuffer = ls.minibuffer[:0]
				ls.bufferMutex.Unlock()
				fmt.Println("ticker Flushing betlogs...")
				if err := ls.flushBatchAviator(batchToFlush); err != nil {
					go ls.taskReTryAviator(batchToFlush)
					slog.Error("Failed to flush batch due to ticker: %v", err)
				}
			} else if len(ls.miniUpdatebuffer) > 0 {
				// 达到阈值，触发批量更新
				batchToFlush := make([]*slotsmongo.DocBetLogAviator, len(ls.miniUpdatebuffer))
				copy(batchToFlush, ls.miniUpdatebuffer)
				ls.miniUpdatebuffer = ls.miniUpdatebuffer[:0]
				ls.bufferMutex.Unlock()
				fmt.Println("logSignal Flushing betlogs...")
				if err := ls.flushBatchAviatorUpdate(batchToFlush); err != nil {
					go ls.taskReTryAviatorUpdate(batchToFlush)
					slog.Error("Failed to flush batch due to threshold: %v", err)
				}
			} else {
				ls.bufferMutex.Unlock()
			}
			return

		case minilogEntry := <-ls.minilogSignal:
			ls.bufferMutex.Lock()
			ls.minibuffer = append(ls.minibuffer, minilogEntry)
			if len(ls.minibuffer) >= ls.BatchSize {
				// 达到阈值，触发批量更新
				batchToFlush := make([]*slotsmongo.DocBetLogAviator, len(ls.minibuffer))
				copy(batchToFlush, ls.minibuffer)
				ls.minibuffer = ls.minibuffer[:0]
				ls.bufferMutex.Unlock()
				fmt.Println("minilogSignal 准备插入历史记录")
				if err := ls.flushBatchAviator(batchToFlush); err != nil {
					go ls.taskReTryAviator(batchToFlush)
					slog.Error(fmt.Sprintf("Failed to flush batch due to threshold: %v", err))
				}
			} else {
				ls.bufferMutex.Unlock()
			}
		case minilogUpdateEntry := <-ls.miniUpdatelogSignal:
			ls.bufferMutex.Lock()
			ls.miniUpdatebuffer = append(ls.miniUpdatebuffer, minilogUpdateEntry)
			if len(ls.miniUpdatebuffer) >= ls.BatchSize {
				// 达到阈值，触发批量更新
				batchToFlush := make([]*slotsmongo.DocBetLogAviator, len(ls.miniUpdatebuffer))
				copy(batchToFlush, ls.miniUpdatebuffer)
				ls.miniUpdatebuffer = ls.miniUpdatebuffer[:0]
				ls.bufferMutex.Unlock()
				fmt.Println("miniUpdatelogSignal 准备插入历史记录")
				if err := ls.flushBatchAviatorUpdate(batchToFlush); err != nil {
					go ls.taskReTryAviatorUpdate(batchToFlush)
					slog.Error("Failed to flush batch due to threshold: %v", err)
				}
			} else {
				ls.bufferMutex.Unlock()
			}

		}
	}
}

// flushBatch 将一批日志写入 ClickHouse
func (ls *LogService) flushBatch(logEntries []*slotsmongo.DocBetLog) error {
	if len(logEntries) == 0 {
		return nil
	}

	tx, err := ls.clickhouseDB.BeginTx(ls.ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rberr := tx.Rollback(); rberr != nil && rberr != sql.ErrTxDone {
			slog.Error("Failed to rollback transaction: %v", rberr)
		}
	}()

	stmt, err := tx.PrepareContext(ls.ctx, "INSERT INTO gamelogs (id, Pid, GameID, TotalWinLoss, PGBetID, UserID, RoundID, Win, Balance, TransferAmount, Comment, Completed, WinLose, SpinDetailsJson, GameType, Bet, InsertTime, AppID, Grade, HitBigReward, LogType, ManufacturerName, CurrencyKey, UserName) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	ls.buildBatchInsertQuery(stmt, logEntries)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// buildBatchInsertQuery 构建批量插入 SQL 语句和参数
func (ls *LogService) buildBatchInsertQuery(stmt *sql.Stmt, logEntries []*slotsmongo.DocBetLog) {
	for _, entry := range logEntries {
		completedInt8 := int8(0) // 默认值 0 (false)
		if entry.Completed {
			completedInt8 = 1 // 如果为 true，设置为 1
		}
		_, err := stmt.ExecContext(ls.ctx,
			entry.ID,
			entry.Pid,              // Pid
			entry.GameID,           // GameID
			entry.TotalWinLoss,     // TotalWinLoss
			entry.PGBetID,          // PGBetID
			entry.UserID,           // UserID
			entry.RoundID,          // RoundID
			entry.Win,              // Win
			entry.Balance,          // Balance
			entry.TransferAmount,   // TransferAmount
			entry.Comment,          // Comment
			completedInt8,          // Completed (1 = true) 不支持bool值
			entry.WinLose,          // WinLose
			entry.SpinDetailsJson,  // SpinDetailsJson
			entry.GameType,         // GameType
			entry.Bet,              // Bet
			time.Now(),             // InsertTime
			entry.AppID,            // AppID
			entry.Grade,            // Grade
			entry.HitBigReward,     // HitBigReward
			entry.LogType,          // LogType
			entry.ManufacturerName, // ManufacturerName
			entry.CurrencyKey,      // CurrencyKey
			entry.UserName,         // UserName
		)
		if err != nil {
			fmt.Errorf("failed to execute statement: %w", err)
			return
		}
	}

}

// flushBatch 将一批日志写入 ClickHouse
func (ls *LogService) flushBatchAviator(logEntries []*slotsmongo.DocBetLogAviator) error {
	if len(logEntries) == 0 {
		return nil
	}

	tx, err := ls.clickhouseDB.BeginTx(ls.ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rberr := tx.Rollback(); rberr != nil && rberr != sql.ErrTxDone {
			slog.Error("Failed to rollback transaction: %v", rberr)
		}
	}()

	stmt, err := tx.PrepareContext(ls.ctx, "INSERT INTO aviatorgamelogs (id, Pid, AppID, Uid, UserName, CurrencyKey, RoomID, Bet,  Win, Balance, CashOutDate, Payout, Profit, GameID, RoundBetId, RoundID, MaxMultiplier, InsertTime, RoundMaxMultiplier, LogType, FinishType, ManufacturerName, BetId) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	ls.buildBatchInsertQueryAviator(stmt, logEntries)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// buildBatchInsertQuery 构建批量插入 SQL 语句和参数
func (ls *LogService) buildBatchInsertQueryAviator(stmt *sql.Stmt, logEntries []*slotsmongo.DocBetLogAviator) {
	for _, entry := range logEntries {
		_, err := stmt.ExecContext(ls.ctx,
			entry.ID,
			entry.Pid,         // Pid
			entry.AppID,       // AppID
			entry.Uid,         // UserID
			entry.UserName,    // UserName
			entry.CurrencyKey, // CurrencyKey
			entry.RoomId,      // RoomId
			entry.Bet,         // Bet
			entry.Win,         // Win
			entry.Balance,     // Balance
			entry.CashOutDate,
			entry.Payout,
			entry.Profit,
			entry.GameID,
			entry.RoundBetId,
			entry.RoundId, // RoundID
			entry.MaxMultiplier,
			time.Now(), // InsertTime
			entry.RoundMaxMultiplier,
			entry.LogType,
			entry.FinishType,
			entry.ManufacturerName,
			entry.BetId,
		)
		if err != nil {
			slog.Error("buildBatchInsertQueryAviator err", "stmt.ExecContext err", err)
			fmt.Errorf("failed to execute statement: %w", err)
			return
		}
	}

}

func (ls *LogService) taskReTryAviator(batchToFlush []*slotsmongo.DocBetLogAviator) {
	for i := 0; i <= ls.ReTryTimes; i++ {
		slog.Error(fmt.Sprintf("Failed to flush batch to ClickHouse, retry attempt: %d/%d, batch: %v", i+1, ls.ReTryTimes+1, batchToFlush))
		if err := ls.flushBatchAviator(batchToFlush); err == nil {
			slog.Info("Successfully flushed batch after %d retries.", i+1)
			return
		}

		if i < ls.ReTryTimes {
			// 指数退避策略，可以根据需要调整 base 和 maxDelay
			baseDelay := 1 * time.Second
			maxDelay := 60 * time.Second
			delay := baseDelay * time.Duration(math.Pow(2, float64(i)))
			if delay > maxDelay {
				delay = maxDelay
			}
			slog.Warn("Retrying in %v...", delay)
			time.Sleep(delay)
		}
	}
	slog.Error(fmt.Sprintf("Failed to flush batch after all %d retries.", ls.ReTryTimes+1))
}

// flushBatch 将一批日志写入 ClickHouse
func (ls *LogService) flushBatchAviatorUpdate(logEntries []*slotsmongo.DocBetLogAviator) error {
	if len(logEntries) == 0 {
		return nil
	}
	slog.Info("我要更新")
	tx, err := ls.clickhouseDB.BeginTx(ls.ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rberr := tx.Rollback(); rberr != nil && rberr != sql.ErrTxDone {
			slog.Error("Failed to rollback transaction: %v", rberr)
		}
	}()
	sqlStr := `Win = ?, Balance = ?, CashOutDate = ?, Payout = ?, Profit = ?, MaxMultiplier = ?, FinishType = ?`
	query := fmt.Sprintf(`ALTER TABLE aviatorgamelogs UPDATE %v WHERE id = ?;`, sqlStr)
	stmt, err := tx.PrepareContext(ls.ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	ls.buildBatchUpdateAviator(stmt, logEntries)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// buildBatchInsertQuery 构建批量插入 SQL 语句和参数
func (ls *LogService) buildBatchUpdateAviator(stmt *sql.Stmt, logEntries []*slotsmongo.DocBetLogAviator) {
	for _, entry := range logEntries {
		_, err := stmt.ExecContext(ls.ctx,
			entry.Win,     // Win
			entry.Balance, // Balance
			entry.CashOutDate,
			entry.Payout,
			entry.Profit,
			entry.MaxMultiplier,
			entry.FinishType,
			entry.ID,
		)
		if err != nil {
			slog.Error("buildBatchUpdateAviator err", "stmt.ExecContext err", err)
			fmt.Errorf("failed to execute statement: %w", err)
			return
		}
	}

}

func (ls *LogService) taskReTryAviatorUpdate(batchToFlush []*slotsmongo.DocBetLogAviator) {
	for i := 0; i <= ls.ReTryTimes; i++ {
		slog.Error("Failed to flush batch to ClickHouse, retry attempt: %d/%d, batch: %v", i+1, ls.ReTryTimes+1, batchToFlush)
		if err := ls.flushBatchAviatorUpdate(batchToFlush); err == nil {
			slog.Info("Successfully flushed batch after %d retries.", i+1)
			return
		}

		if i < ls.ReTryTimes {
			// 指数退避策略，可以根据需要调整 base 和 maxDelay
			baseDelay := 1 * time.Second
			maxDelay := 60 * time.Second
			delay := baseDelay * time.Duration(math.Pow(2, float64(i)))
			if delay > maxDelay {
				delay = maxDelay
			}
			slog.Warn("Retrying in %v...", delay)
			time.Sleep(delay)
		}
	}
	slog.Error(fmt.Sprintf("Failed to flush batch after all %d retries.", ls.ReTryTimes+1))
}
