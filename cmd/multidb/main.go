package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/connor-davis/zingfibre-core/env"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/radius"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goccy/go-json"
)

func main() {
	context := context.Background()

	zingConnection, err := sql.Open("mysql", string(env.ZING_DSN))

	if err != nil {
		log.Printf("🔥 Failed to connect to Zing database: %s", err.Error())

		return
	}

	defer zingConnection.Close()

	log.Print("✅ Connected to Zingfibre Zing database successfully")

	radiusConnection, err := sql.Open("mysql", string(env.RADIUS_DSN))

	if err != nil {
		log.Printf("🔥 Failed to connect to Radius database: %s", err.Error())

		return
	}

	defer radiusConnection.Close()

	log.Print("✅ Connected to Zingfibre Radius database successfully")

	zingQueries := zing.New(zingConnection)
	radiusQueries := radius.New(radiusConnection)

	start := time.Now()

	expiringCustomersZing, err := zingQueries.GetReportsExpiringCustomersNoPagination(context)

	if err != nil {
		log.Printf("🔥 Failed to retrieve expiring customers from Zing: %s", err.Error())
		return
	}

	expiringCustomersRadius, err := radiusQueries.GetReportsExpiringCustomers(context)

	if err != nil {
		log.Printf("🔥 Failed to retrieve expiring customers from Radius: %s", err.Error())
		return
	}

	log.Printf("✅ Retrieved %d expiring customers from Zing and %d from Radius", len(expiringCustomersZing), len(expiringCustomersRadius))

	// products, err := zingQueries.GetProducts(context)

	// if err != nil {
	// 	log.Printf("🔥 Failed to retrieve products: %s", err.Error())
	// 	return
	// }

	// addresses, err := zingQueries.GetAddresses(context)

	// if err != nil {
	// 	log.Printf("🔥 Failed to retrieve addresses: %s", err.Error())
	// 	return
	// }

	linqExpiringCustomersZing := linq.From(expiringCustomersZing)
	linqExpiringCustomersRadius := linq.From(expiringCustomersRadius)

	joinedData := []interface{}{}

	linqExpiringCustomersZing.
		Join(
			linqExpiringCustomersRadius,
			func(i interface{}) interface{} {
				return strings.ToLower(i.(zing.GetReportsExpiringCustomersNoPaginationRow).RadiusUsername.String)
			},
			func(o interface{}) interface{} {
				return strings.ToLower(o.(radius.GetReportsExpiringCustomersRow).Username)
			},
			func(i interface{}, o interface{}) interface{} {
				return system.ReportExpiringCustomer{
					FullName:             i.(zing.GetReportsExpiringCustomersNoPaginationRow).FullName,
					Email:                i.(zing.GetReportsExpiringCustomersNoPaginationRow).Email.String,
					PhoneNumber:          i.(zing.GetReportsExpiringCustomersNoPaginationRow).PhoneNumber.String,
					RadiusUsername:       i.(zing.GetReportsExpiringCustomersNoPaginationRow).RadiusUsername.String,
					LastPurchaseDuration: i.(zing.GetReportsExpiringCustomersNoPaginationRow).LastPurchaseDuration.String,
					LastPurchaseSpeed:    i.(zing.GetReportsExpiringCustomersNoPaginationRow).LastPurchaseSpeed.String,
					Expiration:           o.(radius.GetReportsExpiringCustomersRow).Expiration.Time.Format(time.DateOnly),
					Address:              i.(zing.GetReportsExpiringCustomersNoPaginationRow).Address.String,
				}
			},
		).
		OrderBy(func(i interface{}) interface{} {
			return strings.ToLower(i.(system.ReportExpiringCustomer).FullName)
		}).
		ToSlice(&joinedData)

	prettyJoinedData, err := json.MarshalIndent(joinedData, "", "  ")

	if err != nil {
		log.Printf("🔥 Failed to marshal joined data: %s", err.Error())
		return
	}

	elapsed := time.Since(start)

	log.Printf("%v", string(prettyJoinedData))
	log.Printf("✅ Joined data retrieved successfully in %v", elapsed)
}
