package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	pbsvc "mnovicio.com/nycab/protocol/rpc"

	persistence "mnovicio.com/nycab/server/data/persistence"
)

var (
	serviceSyncOnce sync.Once
	serviceInstance *NYCabServiceImpl
)

// NYCabServiceImpl implements NYCabService
type NYCabServiceImpl struct {
	dbContext *persistence.MySQLDBContext
}

// GetServiceInstance returns single instance of NYCabServiceImpl
func GetServiceInstance(db *sql.DB) *NYCabServiceImpl {
	serviceSyncOnce.Do(func() {
		serviceInstance = &NYCabServiceImpl{
			dbContext: persistence.GetSQLDBContextInstance(db),
		}
	})

	return serviceInstance
}

// GetTripCountsForCabIDsV1 returns the total number of trips the cab has made based on pickup_datetime column with time ignored
func (s *NYCabServiceImpl) GetTripCountsForCabIDsV1(ctx context.Context, in *pbsvc.GetTripCountsForCabIDsRequestV1) (*pbsvc.GetTripCountsForCabIDsResponseV1, error) {
	log.Println("GetTripCountsForCabIDsV1: request = ", in)
	// check date format
	_, err := time.Parse("2006-01-02", in.PickupDate)
	if err != nil {
		errString := fmt.Sprintf("wrong file format for [%s], expecting 'YYYY-MM-DD'", in.PickupDate)
		log.Println(errString)
		return &pbsvc.GetTripCountsForCabIDsResponseV1{
			Error: fmt.Sprintf("%s. Error: %s", errString, err.Error()),
		}, nil
	}

	cabTrips, err := s.dbContext.GetTripCountsForCabsByPickupDate(in.CabIds, in.PickupDate, in.IgnoreCache)
	if err != nil {
		return &pbsvc.GetTripCountsForCabIDsResponseV1{}, err
	}

	return &pbsvc.GetTripCountsForCabIDsResponseV1{
		CabTripsPerDay: cabTrips,
	}, nil
}

// GetAllCabTripCountPerDayV1 returns number of trips per day on record for each cab
func (s *NYCabServiceImpl) GetAllCabTripCountPerDayV1(ctx context.Context, in *pbsvc.GetAllCabTripsRequestV1) (*pbsvc.GetAllCabTripsResponseV1, error) {
	log.Println("GetAllCabTripCountPerDayV1: request = ", in)
	cabTrips, err := s.dbContext.GetAllCabTrips(in.IgnoreCache)
	if err != nil {
		return &pbsvc.GetAllCabTripsResponseV1{}, err
	}

	return &pbsvc.GetAllCabTripsResponseV1{
		CabTripsPerDay: cabTrips,
	}, nil
}

// ClearCacheV1 clears the cache
func (s *NYCabServiceImpl) ClearCacheV1(ctx context.Context, in *pbsvc.ClearCacheRequestV1) (*pbsvc.ClearCacheResponseV1, error) {
	cleared, err := s.dbContext.ClearCache()
	return &pbsvc.ClearCacheResponseV1{
		CacheCleared: cleared,
	}, err
}
