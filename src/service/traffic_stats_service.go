// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package service
package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"gopkg.in/yaml.v2"

	"GIDS/common/logger"
	"GIDS/dao"
	"GIDS/models/db"
	"GIDS/models/req"
)

const (
	batchSize = 1000

	getOnline           = "getOnline"
	getOnlineOfModel    = "getOnlineOfModel"
	getOnlineOfInstance = "getOnlineOfInstance"
	getTrafficOfApp     = "getTrafficOfApp"
	getTraffic          = "getTraffic"
)

// TrafficStatsService 接口定义
type TrafficStatsService interface {
	ExportSessionList(month string, csvFilePath string) error
	ExportMediaList(month string, csvFilePath string) error
	ExportControlList(month string, csvFilePath string) error

	GetOnline(startTime, endTime string) []Res
	GetOnlineOfModel(startTime, endTime string) []Res
	GetOnlineOfInstance(startTime, endTime string) []Res
	GetTrafficOfApp(startTime, endTime string) []Res
	GetTraffic(startTime, endTime string) []Res

	BatchInsertStats(tag string, items []json.RawMessage) error
	HandleSessionStats(session *db.SessionStats) error
}

var _ TrafficStatsService = &TrafficStatsServiceImpl{}

// TrafficStatsServiceImpl 实现 TrafficStatsService 接口
type TrafficStatsServiceImpl struct {
	sd        *dao.SessionStatsDao
	md        *dao.MediaTrafficStatsDao
	cd        *dao.ControlTrafficStatsDao
	sqlConfig *SQLConfig // 添加配置字段
}

// NewTrafficStatsService 创建 TrafficStatsService 实例
func NewTrafficStatsService(configPath string) *TrafficStatsServiceImpl {
	service := &TrafficStatsServiceImpl{
		sd: dao.NewSessionLogDao(),
		md: dao.NewMediaTrafficStatsDao(),
		cd: dao.NewControlTrafficStatsDao(),
	}

	if "" != configPath {
		sqlConfig, err := LoadSQLConfig(configPath)
		if err != nil {
			logger.Errorf("load sql config failed: %v", err)
			return service
		}
		service.sqlConfig = sqlConfig
		logger.Infof("load sql config success: %v", service.sqlConfig)
	}
	return service
}

type SQLConfig struct {
	Queries map[string]struct {
		SQL    string   `yaml:"sql"`
		Params []string `yaml:"params"`
	} `yaml:"queries"`
}

func LoadSQLConfig(configPath string) (*SQLConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config SQLConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

type Res struct {
	Obj string  `orm:"obj"`
	Cnt float64 `orm:"cnt"`
}

func (s *TrafficStatsServiceImpl) executeQuery(queryName, startTime, endTime string) []Res {
	// 检查配置是否存在
	if s.sqlConfig == nil {
		logger.Warnf("query %s not any sql config exist", queryName)
		return nil
	}

	// 检查查询配置是否存在
	queryConfig, exists := s.sqlConfig.Queries[queryName]
	if !exists {
		logger.Warnf("query %s not found in config ", queryName)
		return nil
	}

	var result []Res

	// 根据配置中的参数列表构建参数
	var args []interface{}
	for _, param := range queryConfig.Params {
		if param == "startTime" {
			args = append(args, startTime)
		} else if param == "endTime" {
			args = append(args, endTime)
		}
	}

	err := s.sd.QueryMulti(&result, queryConfig.SQL, args...)
	if err != nil {
		logger.Errorf("execute query %s failed, err: %v", queryName, err)
		return nil
	}
	logger.Infof("query success of %s, result: %v", queryName, result)
	return result
}

// GetOnline 获取在线用户数
func (s *TrafficStatsServiceImpl) GetOnline(startTime, endTime string) []Res {
	return s.executeQuery(getOnline, startTime, endTime)
}

// GetOnlineOfModel 获取各机型的在线用户数
func (s *TrafficStatsServiceImpl) GetOnlineOfModel(startTime, endTime string) []Res {
	return s.executeQuery(getOnlineOfModel, startTime, endTime)
}

// GetOnlineOfInstance 获取各个实例的在线用户数
func (s *TrafficStatsServiceImpl) GetOnlineOfInstance(startTime, endTime string) []Res {
	return s.executeQuery(getOnlineOfInstance, startTime, endTime)
}

// GetTrafficOfApp 获取个应用流量
func (s *TrafficStatsServiceImpl) GetTrafficOfApp(startTime, endTime string) []Res {
	return s.executeQuery(getTrafficOfApp, startTime, endTime)
}

// GetTraffic 获取站点/总流量
func (s *TrafficStatsServiceImpl) GetTraffic(startTime, endTime string) []Res {
	return s.executeQuery(getTraffic, startTime, endTime)
}

// ExportSessionList 导出会话数据到 CSV 文件
func (s *TrafficStatsServiceImpl) ExportSessionList(month string, csvFilePath string) error {
	logger.Infof("Begin to export session list")
	return QueryStatsDataAndWriteCSV[db.SessionStats](month, batchSize, func(list *[]db.SessionStats, opt dao.QueryOption) error {
		return s.sd.List(list, opt)
	}, csvFilePath, []string{"SessionID", "AppType", "StartedAt", "FinishedAt"}, func(item interface{}) []string {
		session := item.(db.SessionStats)
		return []string{
			session.SessionID,
			strconv.Itoa(session.AppType),
			session.StartedAt,
			session.FinishedAt,
		}
	})
}

// ExportMediaList 导出媒体流数据到 CSV 文件
func (s *TrafficStatsServiceImpl) ExportMediaList(month string, csvFilePath string) error {
	logger.Infof("Begin to export media list")
	return QueryStatsDataAndWriteCSV[db.MediaTrafficStats](month, batchSize, func(list *[]db.MediaTrafficStats, opt dao.QueryOption) error {
		return s.md.List(list, opt)
	}, csvFilePath, []string{"SessionID", "AppType", "StartedAt", "FinishedAt", "OutBytes", "AccessType"}, func(item interface{}) []string {
		media := item.(db.MediaTrafficStats)
		return []string{
			media.SessionID,
			strconv.Itoa(media.AppType),
			media.StartedAt,
			media.FinishedAt,
			strconv.FormatInt(media.OutBytes, 10),
			strconv.Itoa(media.AccessType),
		}
	})
}

// ExportControlList 导出控制流数据到 CSV 文件
func (s *TrafficStatsServiceImpl) ExportControlList(month string, csvFilePath string) error {
	logger.Infof("Begin to export control list")
	return QueryStatsDataAndWriteCSV[db.ControlTrafficStats](month, batchSize, func(list *[]db.ControlTrafficStats, opt dao.QueryOption) error {
		return s.cd.List(list, opt)
	}, csvFilePath, []string{"SessionID", "AppType", "StartedAt", "FinishedAt", "OutBytes", "AccessType"}, func(item interface{}) []string {
		control := item.(db.ControlTrafficStats)
		return []string{
			control.SessionID,
			strconv.Itoa(control.AppType),
			control.StartedAt,
			control.FinishedAt,
			strconv.FormatInt(control.OutBytes, 10),
			strconv.Itoa(control.AccessType),
		}
	})
}

// QueryStatsDataAndWriteCSV 查询批量数据并写入CSV文件
func QueryStatsDataAndWriteCSV[T any](month string, batchSize int, listFunc func(*[]T, dao.QueryOption) error, csvFilePath string, headers []string, rowFunc func(interface{}) []string) error {
	logger.Infof("Begin to query stats data and write CSV file, month: %s, batchSize: %d, csvFilePath: %s", month, batchSize, csvFilePath)
	file, err := os.Create(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err = writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	offset := 0
	for {
		opt := dao.NewQueryOption()
		opt.Limit(batchSize, offset, "id")
		opt.Filter("started_at__istartswith", month)
		// 分批查询和处理数据
		var batch []T
		if err = listFunc(&batch, *opt); err != nil {
			return fmt.Errorf("get list failed: %v", err)
		}
		if len(batch) == 0 {
			break
		}
		for _, item := range batch {
			if err = writer.Write(rowFunc(item)); err != nil {
				return fmt.Errorf("failed to write CSV row: %v", err)
			}
		}
		offset += batchSize
	}

	return nil
}

// 插入或者更新Session数据
func (s *TrafficStatsServiceImpl) HandleSessionStats(session *db.SessionStats) error {
	// 查询记录是否存在
	exist, err := s.sd.Exist(session.TcpUniqueId)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if !exist {
		// 插入新记录
		logger.Infof("not exist, insert new data")
		if err := s.sd.Insert(session); err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}
	} else {
		// 更新记录
		logger.Infof(" data exist, update data")
		if err := s.sd.UpdatebySession(session); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}
	}

	return nil
}

// 批量插入数据
func (s *TrafficStatsServiceImpl) BatchInsertStats(tag string, items []json.RawMessage) error {
	// 根据 tag 处理不同类型的批量插入
	logger.Infof("Starting BatchInsertStats for tag: %s, items count: %d", tag, len(items))
	switch tag {
	case "media traffic stats":
		return s.batchInsertMediaStats(items)
	case "control traffic stats":
		return s.batchInsertControlStats(items)
	default:
		logger.Errorf("Unsupported tag: %s", tag)
		return fmt.Errorf("unsupported tag: %s", tag)
	}
}

// 媒体流量统计
func (s *TrafficStatsServiceImpl) batchInsertMediaStats(items []json.RawMessage) error {
	return batchInsertStats(
		s.md.BaseInterface.(*dao.BaseDao),
		items,
		func() *db.MediaTrafficStats { return &db.MediaTrafficStats{} },
	)
}

// 控制流量统计
func (s *TrafficStatsServiceImpl) batchInsertControlStats(items []json.RawMessage) error {
	return batchInsertStats(
		s.cd.BaseInterface.(*dao.BaseDao),
		items,
		func() *db.ControlTrafficStats { return &db.ControlTrafficStats{} },
	)
}

// 泛型批量插入函数
func batchInsertStats[T req.IRequest](
	dao *dao.BaseDao,
	items []json.RawMessage,
	newModel func() T,
) error {
	var statsList []T

	// 反序列化和验证
	for i, raw := range items {
		model := newModel()
		if err := json.Unmarshal(raw, model); err != nil {
			logger.Errorf("Failed to unmarshal stat item (index %d): %v", i, err)
			return fmt.Errorf("failed to unmarshal stat item (index %d): %w", i, err)
		}

		if err := model.Validate(); err != nil {
			logger.Errorf("Validation failed for stat (index %d): %v", i, err)
			return fmt.Errorf("validation failed for stat (index %d): %w", i, err)
		}

		statsList = append(statsList, model)
	}

	logger.Infof("Total valid records: %d", len(statsList))

	// 使用 baseDao 的 DoTxWithCtx 方法
	return dao.DoTxWithCtx(context.Background(), func(ctx context.Context, txOrm orm.TxOrmer) error {
		batchSize := 1000
		total := len(statsList)

		for i := 0; i < total; i += batchSize {
			end := i + batchSize
			if end > total {
				end = total
			}

			batch := statsList[i:end]
			logger.Infof("Inserting batch %d-%d", i, end)

			// 调用 DAO 的 InsertMultiWithOrm 方法
			_, err := dao.InsertMultiWithOrm(ctx, txOrm, batch)
			if err != nil {
				logger.Errorf("Batch insert failed (range %d-%d): %v", i, end, err)
				return fmt.Errorf("batch insert failed (range %d-%d): %w", i, end, err)
			}
			logger.Infof("Successfully inserted batch %d-%d", i, end)
		}

		logger.Infof("Transaction completed, total records inserted: %d", total)
		return nil
	})
}

// CleanOldStats 清理过期的统计数据
func (s *TrafficStatsServiceImpl) CleanOldStats(months int) error {
	// 计算过期时间点（当前时间 - 指定月数）
	cutoffTime := time.Now().AddDate(0, -months, 0).Format(time.RFC3339)
	logger.Infof("Starting data cleanup for data older than: %s", cutoffTime)

	// 清理会话统计表
	if err := s.cleanSessionStats(cutoffTime); err != nil {
		return fmt.Errorf("failed to clean session stats: %w", err)
	}

	// 清理媒体流量统计表
	if err := s.cleanMediaStats(cutoffTime); err != nil {
		return fmt.Errorf("failed to clean media stats: %w", err)
	}

	// 清理控制流量统计表
	if err := s.cleanControlStats(cutoffTime); err != nil {
		return fmt.Errorf("failed to clean control stats: %w", err)
	}

	logger.Infof("Data cleanup completed for data older than: %s", cutoffTime)
	return nil
}

// 清理会话统计表
func (s *TrafficStatsServiceImpl) cleanSessionStats(cutoffTime string) error {
	logger.Infof("Cleaning session stats before %s", cutoffTime)
	return s.sd.DoTxWithCtx(context.Background(), func(ctx context.Context, txOrm orm.TxOrmer) error {
		_, err := txOrm.QueryTable("t_session_stats").Filter("started_at__lt", cutoffTime).Delete()
		if err != nil {
			return fmt.Errorf("failed to delete session stats: %w", err)
		}
		return nil
	})
}

// 清理媒体流量统计表
func (s *TrafficStatsServiceImpl) cleanMediaStats(cutoffTime string) error {
	logger.Infof("Cleaning media stats before %s", cutoffTime)
	return s.md.DoTxWithCtx(context.Background(), func(ctx context.Context, txOrm orm.TxOrmer) error {
		_, err := txOrm.QueryTable("t_media_traffic_stats").Filter("started_at__lt", cutoffTime).Delete()
		if err != nil {
			return fmt.Errorf("failed to delete media stats: %w", err)
		}
		return nil
	})
}

// 清理控制流量统计表
func (s *TrafficStatsServiceImpl) cleanControlStats(cutoffTime string) error {
	logger.Infof("Cleaning control stats before %s", cutoffTime)
	return s.cd.DoTxWithCtx(context.Background(), func(ctx context.Context, txOrm orm.TxOrmer) error {
		_, err := txOrm.QueryTable("t_control_traffic_stats").Filter("started_at__lt", cutoffTime).Delete()
		if err != nil {
			return fmt.Errorf("failed to delete control stats: %w", err)
		}
		return nil
	})
}
