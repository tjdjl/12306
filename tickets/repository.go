package tickets

import (
	"context"
	"database/sql"

	"12306.com/12306/common"
	"github.com/jinzhu/gorm"
)

type ITicketRepository interface {
	FindTripStartNoAndEndNo(startCity, endCity, date string) ([]TripStartNoAndEndNo, error)
	FindHighSpeedTripStartNoAndEndNo(startCity, endCity, date string) ([]TripStartNoAndEndNo, error)
	FindTripSegment(tripID, startStationNo, endStationNo uint, catogory string) ([]TripSegment, error)
	FindTrain(tripID uint) (Train, error)
	FindStartStaionDetail(tripID, station_no uint) (StaionDetail, error)
	FindEndStaionDetail(tripID, station_no uint) (StaionDetail, error)
}

type ITicketRepositoryTX interface {
	FindTripSegment(tripID, startStationNo, endStationNo uint, catogory string) ([]TripSegment, error)
	FindValidOrder(orderID, userID uint) (Order, error)
	UpdateTripSegment(seats TripSegment) error
	UpdateOrderStatus(order *Order, s string) error
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
		Isolation: sql.LevelReadUncommitted,
	})}
}

// 根据startCity, endCity, date查询TripStartNoAndEndNo
func (c TicketRepository) FindTripStartNoAndEndNo(startCity, endCity, date string) ([]TripStartNoAndEndNo, error) {
	var models []TripStartNoAndEndNo
	sql := "SELECT A.trip_id AS trip_id ,A.sequence AS start_station_no,B.sequence AS end_station_no 	FROM (SELECT trip_id,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? ) A, 	(SELECT trip_id,sequence FROM trip_station WHERE station_name =?) B WHERE A.trip_id = B.trip_id AND A.sequence < B.sequence "
	err := c.DB.Raw(sql, startCity, date, endCity).Find(&models).Error
	return models, err
}

// 根据startCity, endCity, date查询TripStartNoAndEndNo
func (c TicketRepository) FindHighSpeedTripStartNoAndEndNo(startCity, endCity, date string) ([]TripStartNoAndEndNo, error) {
	var models []TripStartNoAndEndNo
	sql := "SELECT A.trip_id AS trip_id ,A.sequence AS start_station_no,B.sequence AS end_station_no 	FROM (SELECT trip_id,sequence FROM trip_station WHERE station_name =? AND date(start_time)=? AND catogory ='highSpeed') A, 	(SELECT trip_id,sequence FROM trip_station WHERE station_name =?) B WHERE A.trip_id = B.trip_id AND A.sequence < B.sequence "
	err := c.DB.Raw(sql, startCity, date, endCity).Find(&models).Error
	return models, err
}
func (c TicketRepository) FindTripSegment(tripID, startStationNo, endStationNo uint, catogory string) ([]TripSegment, error) {
	var models []TripSegment
	sql := "SELECT * FROM trip_segment WHERE trip_id = ? AND segment_no between ? AND ? AND seat_catogory = ?"
	err := c.DB.Raw(sql, tripID, startStationNo, endStationNo-1, catogory).Find(&models).Error
	return models, err
}
func (c TicketRepository) FindTrain(tripID uint) (Train, error) {
	var model Train
	sql := "SELECT trains.catogory,trains.train_number,trains.length FROM trains,trips WHERE trips.id = ? AND trips.train_id = trains.id"
	err := c.DB.Raw(sql, tripID).Find(&model).Error
	return model, err
}
func (c TicketRepository) FindStartStaionDetail(tripID, station_no uint) (StaionDetail, error) {
	var model StaionDetail
	sql := "SELECT station_name,start_time AS station_time FROM trip_station WHERE trip_id = ? AND SEQUENCE = ?"
	err := c.DB.Raw(sql, tripID, station_no).Find(&model).Error
	return model, err
}

func (c TicketRepository) FindEndStaionDetail(tripID, station_no uint) (StaionDetail, error) {
	var model StaionDetail
	sql := "SELECT station_name,arrive_time AS station_time FROM trip_station WHERE trip_id = ? AND SEQUENCE = ?"
	err := c.DB.Raw(sql, tripID, station_no).Find(&model).Error
	return model, err
}
func (c TicketRepositoryTX) FindTripSegment(tripID, startStationNo, endStationNo uint, catogory string) ([]TripSegment, error) {
	var models []TripSegment
	sql := "SELECT * FROM trip_segment WHERE trip_id = ? AND segment_no between ? AND ? AND seat_catogory = ?"
	err := c.TX.Set("gorm:query_option", "FOR UPDATE").Raw(sql, tripID, startStationNo, endStationNo-1, catogory).Find(&models).Error
	if err != nil {
		c.TX.Rollback()
	}
	return models, err
}
func (c TicketRepositoryTX) FindValidOrder(orderID, userID uint) (Order, error) {
	var model Order
	sql := "SELECT * FROM orders WHERE id =  ? AND user_id = ? AND  status != '已退票'"
	err := c.TX.Set("gorm:query_option", "FOR UPDATE").Raw(sql, orderID, userID).Find(&model).Error
	if err != nil {
		c.TX.Rollback()
	}
	return model, err
}
func (c TicketRepositoryTX) UpdateTripSegment(seats TripSegment) error {
	err := c.TX.Save(seats).Error
	return err
}
func (c TicketRepositoryTX) UpdateOrderStatus(order *Order, s string) error {
	err := c.TX.Model(&order).Update("status", s).Error
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
