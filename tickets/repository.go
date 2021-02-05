package tickets

import (
	"context"
	"database/sql"

	"12306.com/12306/common"
	"github.com/jinzhu/gorm"
)

type ITicketRepository interface {
	FindTripStationPair(startCity, endCity, date string) ([]TripStaionPair, error)
	FindFastTripStationPair(startCity, endCity, date string) ([]TripStaionPair, error)
	FindTripSegment(tripID, startStationNo, endStationNo uint, catogory string) ([]TripSegment, error)
}

type ITicketRepositoryTX interface {
	FindTripSegment(tripID, startStationNo, endStationNo uint, catogory string) ([]TripSegment, error)
	FindValidOrder(orderID, userID uint) (Order, error)
	UpdateTripSegment(seats []TripSegment) error
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
func (c TicketRepository) FindTripStationPair(startCity, endCity, date string) ([]TripStaionPair, error) {
	var models []TripStaionPair
	sql := `
			SELECT train_station_pair.start_station_no,train_station_pair.end_station_no,
				CONCAT(trips.start_date,' ',train_station_pair.start_time)  AS start_time,
				CONCAT(trips.start_date,' ',train_station_pair.end_time)  AS end_time,
				train_station_pair.start_station_name,train_station_pair.end_station_name,
				train_station_pair.start_station_no,train_station_pair.end_station_no,
				trips.id AS trip_id,trips.start_date ,
				trains.train_number,trains.catogory AS train_type,
				trains.station_nums AS train_staion_nums，
			FROM train_station_pair,trips,trains
			WHERE train_station_pair.start_station_name =? AND train_station_pair.end_station_name =?
				AND trips.start_date = ? 
				AND CONCAT(trips.start_date,' ',train_station_pair.start_time)  > NOW()
				AND train_station_pair.train_id = trips.train_id 
				AND train_station_pair.train_id = trains.id  		
	`
	err := c.DB.Raw(sql, startCity, endCity, date).Find(&models).Error
	return models, err
}

// 根据startCity, endCity, date查询TripStartNoAndEndNo
func (c TicketRepository) FindFastTripStationPair(startCity, endCity, date string) ([]TripStaionPair, error) {
	var models []TripStaionPair
	sql := `
			SELECT train_station_pair.start_station_no,train_station_pair.end_station_no,
				CONCAT(trips.start_date,' ',train_station_pair.start_time)  AS start_time,
				CONCAT(trips.start_date,' ',train_station_pair.end_time)  AS end_time,
				train_station_pair.start_station_name,train_station_pair.end_station_name,
				train_station_pair.start_station_no,train_station_pair.end_station_no,
				trips.id AS trip_id,trips.start_date ,
				trains.train_number,trains.catogory AS train_type,
				trains.station_nums AS train_staion_nums，
			FROM train_station_pair,trips,trains
			WHERE train_station_pair.start_station_name =? AND train_station_pair.end_station_name =?
				AND trips.start_date = ? 
				AND CONCAT(trips.start_date,' ',train_station_pair.start_time)  > NOW()
				AND train_station_pair.train_id = trips.train_id 
				AND train_station_pair.train_id = trains.id  	
				AND trains.is_fast = 1	
	`
	err := c.DB.Raw(sql, startCity, endCity, date).Find(&models).Error
	return models, err
}

func (c TicketRepository) FindTripSegment(tripID, startStationNo, endStationNo uint, catogory string) ([]TripSegment, error) {
	var models []TripSegment
	sql := "SELECT * FROM trip_segment WHERE trip_id = ? AND segment_no between ? AND ? AND seat_catogory = ?"
	err := c.DB.Raw(sql, tripID, startStationNo, endStationNo-1, catogory).Find(&models).Error
	return models, err
}

func (c TicketRepositoryTX) FindTripSegment(tripID, startStationNo, endStationNo uint, catogory string) ([]TripSegment, error) {
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
	for i := 0; i < len(seats); i++ {
		err = c.TX.Save(seats[i]).Error
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
