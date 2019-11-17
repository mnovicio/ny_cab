package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	pbdata "mnovicio.com/nycab/protocol/objects"
)

var (
	sqlDBOnce     sync.Once
	sqlDBInstance *MySQLDBContext
)

// Cache synchronized cache for cab trips
type Cache struct {
	sync.RWMutex
	m *pbdata.CabTripsPerDay
}

// MySQLDBContext is an MySQL DB Context with simple caching support
type MySQLDBContext struct {
	db    *sql.DB
	cache *Cache
}

// CabTripsPerDay is used for unmarhalling row bytes from query
type CabTripsPerDay struct {
	CabID      string `json:"cab_id"`
	PickUpDate string `json:"pickup_date"`
	TripCount  uint32 `json:"total_trip_count"`
}

// GetSQLDBContextInstance returns single instance of SQL DB context
func GetSQLDBContextInstance(db *sql.DB) *MySQLDBContext {
	sqlDBOnce.Do(func() {
		sqlDBInstance = &MySQLDBContext{
			db: db,
			cache: &Cache{
				m: &pbdata.CabTripsPerDay{
					CabTrips: make(map[string]*pbdata.TripsPerDay),
				},
			},
		}
	})
	return sqlDBInstance
}

// GetTripCountsForCabsByPickupDate returns the total number of trips the cab has made based on pickup_datetime column with time ignored
// cabIDs: list of cab IDs to search
// pickupDate: pickup date in 'YYYY-MM-DD' format
// ignoreCache: true - ignores cache and make query to DB. uses cached data otherwise.
func (m *MySQLDBContext) GetTripCountsForCabsByPickupDate(cabIDs []string, pickupDate string, ignoreCache bool) (*pbdata.CabTripsPerDay, error) {
	cabTripsPerDay := &pbdata.CabTripsPerDay{
		CabTrips: make(map[string]*pbdata.TripsPerDay),
	}

	notInCache := []string{}
	if ignoreCache {
		// if ignore cache, search everthing from db
		notInCache = append(notInCache, cabIDs...)
	} else {
		// else, check cache if cabID with pickup date exists
		m.cache.Lock()
		defer m.cache.Unlock()

		for _, cabID := range cabIDs {
			cachedDataFound := false
			cachedCabTripsPerDay, cachedCabTripFound := m.cache.m.CabTrips[cabID]
			if cachedCabTripFound {
				cachedTripCountOnDate, cachedTripCountOnDateFound := cachedCabTripsPerDay.TripsPerDay[pickupDate]
				if cachedTripCountOnDateFound {
					cachedDataFound = true
					log.Println(fmt.Sprintf("Found in cache [cab_id='%s', pickup_date='%s', count='%d']", cabID, pickupDate, cachedTripCountOnDate))
					cabTripsPerDay.CabTrips[cabID] = &pbdata.TripsPerDay{
						TripsPerDay: map[string]uint32{
							pickupDate: cachedTripCountOnDate,
						},
					}
				}
			}

			// cached data not found, add into list to be be queried from DB
			if !cachedDataFound {
				notInCache = append(notInCache, cabID)
				_cabTripsPerDay := CabTripsPerDay{
					CabID:      cabID,
					PickUpDate: pickupDate,
					TripCount:  0,
				}
				m.addTripCountToSet(cabTripsPerDay, _cabTripsPerDay.CabID, _cabTripsPerDay.PickUpDate, _cabTripsPerDay.TripCount)
				m.addTripCountToSet(m.cache.m, _cabTripsPerDay.CabID, _cabTripsPerDay.PickUpDate, _cabTripsPerDay.TripCount)
			}
		}
	}

	if len(notInCache) > 0 {
		log.Println("fetching data from db for ff cabIDs: ", notInCache)
		query := fmt.Sprintf("SELECT medallion AS cab_id, DATE(pickup_datetime) AS pickup_date, COUNT(pickup_datetime) AS total_trip_cnt FROM cab_trip_data WHERE medallion IN ('%s') AND DATE(pickup_datetime) = '%s' GROUP BY medallion, pickup_date",
			strings.Join(notInCache, "', '"), pickupDate)

		log.Printf("running query: [%s]", query)
		results, err := m.db.Query(query)
		if err != nil {
			panic(err.Error())
		}

		for results.Next() {
			var _cabTripsPerDay CabTripsPerDay
			// for each row, scan the result into our tag composite object
			err = results.Scan(&_cabTripsPerDay.CabID, &_cabTripsPerDay.PickUpDate, &_cabTripsPerDay.TripCount)
			if err != nil {
				panic(err.Error())
			}

			// format date to 'YYYY-MM-DD'
			t, _ := time.Parse(time.RFC3339, _cabTripsPerDay.PickUpDate)
			_cabTripsPerDay.PickUpDate = t.Format("2006-01-02")

			m.addTripCountToSet(cabTripsPerDay, _cabTripsPerDay.CabID, _cabTripsPerDay.PickUpDate, _cabTripsPerDay.TripCount)
			m.addTripCountToSet(m.cache.m, _cabTripsPerDay.CabID, _cabTripsPerDay.PickUpDate, _cabTripsPerDay.TripCount)
		}
	}

	return cabTripsPerDay, nil
}

// GetAllCabTrips returns number of trips per day on record for each cab
// ignoreCache: true - ignores cache and make query to DB. uses cached data otherwise.
func (m *MySQLDBContext) GetAllCabTrips(ignoreCache bool) (*pbdata.CabTripsPerDay, error) {
	cabTripsPerDay := &pbdata.CabTripsPerDay{
		CabTrips: make(map[string]*pbdata.TripsPerDay),
	}

	m.cache.Lock()
	defer m.cache.Unlock()

	// if ignore cache or cach is empty, hit the db
	if ignoreCache || len(m.cache.m.CabTrips) == 0 {
		log.Printf("getting data from db")
		query := "select medallion as cab_id, DATE(pickup_datetime) as pickup_date, count(pickup_datetime) as total_trip_cnt from cab_trip_data group by medallion, pickup_date"
		log.Printf("running query: [%s]", query)
		results, err := m.db.Query(query)
		if err != nil {
			panic(err.Error())
		}

		for results.Next() {
			var _cabTripsPerDay CabTripsPerDay
			// for each row, scan the result into our tag composite object
			err = results.Scan(&_cabTripsPerDay.CabID, &_cabTripsPerDay.PickUpDate, &_cabTripsPerDay.TripCount)
			if err != nil {
				panic(err.Error())
			}

			// format date to 'YYYY-MM-DD'
			t, _ := time.Parse(time.RFC3339, _cabTripsPerDay.PickUpDate)
			_cabTripsPerDay.PickUpDate = t.Format("2006-01-02")

			m.addTripCountToSet(cabTripsPerDay, _cabTripsPerDay.CabID, _cabTripsPerDay.PickUpDate, _cabTripsPerDay.TripCount)
			m.addTripCountToSet(m.cache.m, _cabTripsPerDay.CabID, _cabTripsPerDay.PickUpDate, _cabTripsPerDay.TripCount)
		}
	} else {
		log.Printf("returning cached data")
		for cabID, tripsPerDay := range m.cache.m.CabTrips {
			copyTripsPerDay := &pbdata.TripsPerDay{
				TripsPerDay: make(map[string]uint32),
			}

			for pickupDate, cnt := range tripsPerDay.TripsPerDay {
				copyTripsPerDay.TripsPerDay[pickupDate] = cnt
			}
			cabTripsPerDay.CabTrips[cabID] = copyTripsPerDay
		}
	}

	return cabTripsPerDay, nil
}

func (m *MySQLDBContext) addTripCountToSet(set *pbdata.CabTripsPerDay, cabID, pickUpDate string, tripCount uint32) {
	tripsPerDay := set.CabTrips[cabID]
	if tripsPerDay == nil {
		tripsPerDay = &pbdata.TripsPerDay{
			TripsPerDay: make(map[string]uint32),
		}

		set.CabTrips[cabID] = tripsPerDay
	}

	tripsPerDay.TripsPerDay[pickUpDate] = tripCount
}

// ClearCache clears the cache
func (m *MySQLDBContext) ClearCache() (bool, error) {
	m.cache.Lock()
	defer m.cache.Unlock()

	log.Printf("clearing cache")
	m.cache.m = &pbdata.CabTripsPerDay{
		CabTrips: make(map[string]*pbdata.TripsPerDay),
	}
	log.Printf("cache cleared")

	return true, nil
}
