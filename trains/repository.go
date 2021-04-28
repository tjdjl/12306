package trains

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"12306.com/12306/common"
	"github.com/jinzhu/gorm"
)

type ITicketRepository interface {
	FindCityID(cityName string) (uint, error)
	FindTrainStaions(trainID string) ([]TrainStaion, error)
	FindTrainStationPairs(id string, date string, isToday, isFast bool) ([]TrainStaionPair, error)
	FindTripSegments(tripID string, startStationNo, endStationNo uint) ([]TripSegment, error)
}

type ITicketRepositoryTX interface {
	FindTripSegments(tripID string, startStationNo, endStationNo uint, catogory string) ([]TripSegment, error)
	FindValidOrder(orderID, userID uint) (Order, error)
	UpdateTripSegment(seats []TripSegment) error
	UpdateOrderStatus(order *Order, s string) error
	UpdateOrder(order *Order, newMap map[string]interface{}) error
	CreateOrder(order *Order) error
	Rollback()
	Commit() error
}
type TicketRepository struct {
	DB *gorm.DB
}

type TicketRepositoryTX struct {
	TX *gorm.DB
}

func NewTicketRepository() ITicketRepository {
	return TicketRepository{DB: common.GetDB()}
}

func NewTicketRepositoryTX() ITicketRepositoryTX {
	return TicketRepositoryTX{TX: common.GetDB().BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})}
}

//FindCityID 城市名转id
func (c TicketRepository) FindCityID(cityName string) (uint, error) {
	var res []uint
	err := c.DB.Table("cities").Where("name = ?", cityName).Pluck("ID", &res).Error
	if len(res) == 0 {
		return 0, errors.New("不存在该城市")
	}
	return res[0], err
}

//FindTrainStaions
func (c TicketRepository) FindTrainStaions(trainID string) ([]TrainStaion, error) {
	var res []TrainStaion
	err := c.DB.Where("train_id = ?", trainID).Find(&res).Error
	return res, err
}

//FindTrainStationPairs 找到TrainStationPairs
func (c TicketRepository) FindTrainStationPairs(id string, date string, isToday, isFast bool) ([]TrainStaionPair, error) {
	var models []TrainStaionPair
	var sql string
	if isToday == false {
		sql = `
		SELECT train_station_pair.start_station_no,train_station_pair.end_station_no,
		train_station_pair.start_time,train_station_pair.end_time,
		train_station_pair.start_station_name,train_station_pair.end_station_name,
		train_station_pair.start_station_no,train_station_pair.end_station_no,
		trains.id AS train_id,trains.catogory AS train_type,
		trains.station_nums AS train_staion_nums
		FROM train_station_pair,trains,trips
		WHERE train_station_pair.id Like ?
				AND trains.id  = train_station_pair.train_id
				AND trips.id = CONCAT_WS("-",trains.id,?)
		`
	}
	//如果是查找今天的票，查找的是今天时间，过滤掉过期的车次；
	if isToday == true {
		sql = `
		SELECT train_station_pair.start_station_no,train_station_pair.end_station_no,
		train_station_pair.today_start_time,train_station_pair.today_end_time,
		train_station_pair.start_station_name,train_station_pair.end_station_name,
		train_station_pair.start_station_no,train_station_pair.end_station_no,
		trains.id AS train_id,trains.catogory AS train_type,
		trains.station_nums AS train_staion_nums
		FROM train_station_pair,trains,trips
		WHERE train_station_pair.id Like ?
				AND trains.id  = train_station_pair.train_id
				AND trips.id = CONCAT_WS("-",trains.id,?)
		`
		sql = sql + "AND start_time > time(now())"
	}
	//如果是查找快车的票，加上条件
	if isFast == true {
		sql = sql + "AND trains.is_fast  = 1"
	}
	err := c.DB.Raw(sql, id, date).Find(&models).Error
	return models, err
}

func (c TicketRepository) FindTripSegments(tripID string, startStationNo, endStationNo uint) ([]TripSegment, error) {
	var models []TripSegment
	sql := "SELECT * FROM trip_segment WHERE trip_id = ? AND segment_no between ? AND ?"
	err := c.DB.Raw(sql, tripID, startStationNo, endStationNo-1).Find(&models).Error
	return models, err
}

func (c TicketRepositoryTX) FindTripSegments(tripID string, startStationNo, endStationNo uint, catogory string) ([]TripSegment, error) {
	var models []TripSegment
	sql := "SELECT * FROM trip_segment WHERE trip_id = ? AND segment_no between ? AND ? AND seat_catogory = ?"
	err := c.TX.Set("gorm:query_option", "FOR UPDATE").Raw(sql, tripID, startStationNo, endStationNo-1, catogory).Find(&models).Error
	return models, err
}
func (c TicketRepositoryTX) FindValidOrder(orderID, userID uint) (Order, error) {
	var model Order
	sql := "SELECT * FROM orders WHERE id =  ? AND user_id = ? AND  status != '已退票'"
	err := c.TX.Set("gorm:query_option", "FOR UPDATE").Raw(sql, orderID, userID).Find(&model).Error
	return model, err
}
func (c TicketRepositoryTX) UpdateTripSegment(seats []TripSegment) error {
	var err error
	fmt.Println(len(seats))
	for i := 0; i < len(seats); i++ {
		err = c.TX.Model(&seats[i]).Update("seat_bytes", seats[i].SeatBytes).Error
		if err != nil {
			return err
		}

	}
	return err
}
func (c TicketRepositoryTX) UpdateOrderStatus(order *Order, s string) error {
	err := c.TX.Model(&order).Update("status", s).Error
	return err
}
func (c TicketRepositoryTX) UpdateOrder(order *Order, newMap map[string]interface{}) error {
	err := c.TX.Model(&order).Update(newMap).Error
	return err
}
func (c TicketRepositoryTX) CreateOrder(order *Order) error {
	err := c.TX.Create(&order).Error
	return err
}
func (c TicketRepositoryTX) Rollback() {
	c.TX.Rollback()
}
func (c TicketRepositoryTX) Commit() error {
	err := c.TX.Commit().Error
	return err
}
